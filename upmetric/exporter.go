package upmetric

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/uptrace-go/internal"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

// Config is the configuration to be used when initializing an exporter.
type Config struct {
	// DSN is a data source to connect to uptrace.dev.
	// Example: https://<key>@uptrace.dev/<project_id>
	// The default is to use UPTRACE_DSN environment var.
	DSN string

	// Disabled disables the exporter.
	Disabled bool

	endpoint string
	token    string
}

func (cfg *Config) init() {
	dsn, err := internal.ParseDSN(cfg.DSN)
	if err != nil {
		internal.Logger.Print(err.Error())
		cfg.Disabled = true
		return
	}

	cfg.endpoint = fmt.Sprintf("%s://%s/api/v1/projects/%s/metrics",
		dsn.Scheme, dsn.Host, dsn.ProjectID)
	cfg.token = dsn.Token
}

type Exporter struct {
	cfg *Config

	mmsc      []mmsc
	quantiles []quantile
}

var _ export.Exporter = (*Exporter)(nil)

func NewRawExporter(cfg *Config) *Exporter {
	cfg.init()
	return &Exporter{
		cfg: cfg,
	}
}

// InstallNewPipeline instantiates a NewExportPipeline and registers it globally.
// Typically called as:
//
// 	pipeline := stdout.InstallNewPipeline(stdout.Config{...})
// 	defer pipeline.Stop()
// 	... Done
func InstallNewPipeline(config *Config, options ...push.Option) *push.Controller {
	options = append(options, push.WithPeriod(10*time.Second))
	ctrl := NewExportPipeline(config, options...)
	global.SetMeterProvider(ctrl.Provider())
	return ctrl
}

// NewExportPipeline sets up a complete export pipeline with the recommended setup,
// chaining a NewRawExporter into the recommended selectors and integrators.
func NewExportPipeline(
	config *Config, options ...push.Option,
) *push.Controller {
	exporter := NewRawExporter(config)

	// Not stateful.
	pusher := push.New(
		simple.NewWithInexpensiveDistribution(),
		exporter,
		options...,
	)
	pusher.Start()

	return pusher
}

func (e *Exporter) Export(_ context.Context, checkpointSet export.CheckpointSet) error {
	if e.cfg.Disabled {
		return nil
	}

	if err := e.export(checkpointSet); err != nil {
		return err
	}
	e.flush()

	return nil
}

func (e *Exporter) export(checkpointSet export.CheckpointSet) error {
	return checkpointSet.ForEach(func(record export.Record) error {
		switch agg := record.Aggregator().(type) {
		case aggregator.Quantile:
			return e.exportQuantile(record, agg)
		case aggregator.MinMaxSumCount:
			return e.exportMMSC(record, agg)
		default:
			//log.Printf("unsupported aggregator type: %T", agg)
			return nil
		}
	})
}

func (e *Exporter) exportMMSC(
	record export.Record, agg aggregator.MinMaxSumCount,
) error {
	var expose mmsc

	if err := exportCommon(record, &expose.baseRecord); err != nil {
		return err
	}

	desc := record.Descriptor()
	numKind := desc.NumberKind()

	min, err := agg.Min()
	if err != nil {
		return err
	}
	expose.Min = float32(min.CoerceToFloat64(numKind))

	max, err := agg.Max()
	if err != nil {
		return err
	}
	expose.Max = float32(max.CoerceToFloat64(numKind))

	sum, err := agg.Sum()
	if err != nil {
		return err
	}
	expose.Sum = sum.CoerceToFloat64(numKind)

	count, err := agg.Count()
	if err != nil {
		return err
	}
	expose.Count = count

	e.mmsc = append(e.mmsc, expose)

	return nil
}

var quantiles = []float64{0.5, 0.75, 0.9, 0.95, 0.99}

func (e *Exporter) exportQuantile(
	record export.Record, agg aggregator.Quantile,
) error {
	var expose quantile

	if err := exportCommon(record, &expose.baseRecord); err != nil {
		return err
	}

	desc := record.Descriptor()
	numKind := desc.NumberKind()

	if agg, ok := agg.(aggregator.Count); ok {
		count, err := agg.Count()
		if err != nil {
			return err
		}
		expose.Count = count
	}

	for _, q := range quantiles {
		n, err := agg.Quantile(q)
		if err != nil {
			return err
		}
		expose.Quantiles = append(expose.Quantiles, float32(n.CoerceToFloat64(numKind)))
	}

	e.quantiles = append(e.quantiles, expose)

	return nil
}

func exportCommon(record export.Record, expose *baseRecord) error {
	desc := record.Descriptor()

	expose.Name = desc.Name()
	expose.Description = desc.Description()
	expose.Kind = int8(desc.MetricKind()) // use string?
	expose.Unit = string(desc.Unit())
	expose.Time = time.Now().UnixNano()

	if iter := record.Labels().Iter(); iter.Len() > 0 {
		attrs := record.Resource().Attributes()
		labels := make([]kv.KeyValue, 0, len(attrs)+iter.Len())
		labels = append(labels, attrs...)

		for iter.Next() {
			labels = append(labels, iter.Label())
		}

		expose.Labels = labels
	}

	return nil
}

func (e *Exporter) flush() {
	if len(e.mmsc) == 0 && len(e.quantiles) == 0 {
		return
	}

	go func(mmsc []mmsc, quantiles []quantile) {
		out := make(map[string]interface{})
		if len(mmsc) > 0 {
			out["mmsc"] = mmsc
		}
		if len(quantiles) > 0 {
			out["quantiles"] = quantiles
		}

		if err := e.send(out); err != nil {
			logrus.WithError(err).Error("send failed")
		}
	}(e.mmsc, e.quantiles)

	e.mmsc = nil
	e.quantiles = nil
}

func (e *Exporter) send(out map[string]interface{}) error {
	enc := internal.GetEncoder()
	defer internal.PutEncoder(enc)

	buf, err := enc.EncodeS2(out)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", e.cfg.endpoint, buf)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+e.cfg.token)
	req.Header.Set("Content-Type", "application/msgpack")
	req.Header.Set("Content-Encoding", "s2")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, _ = io.Copy(ioutil.Discard, resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("got %s, wanted 200 OK", resp.Status)
	}

	return nil
}

type baseRecord struct {
	Name        string           `msgpack:"name"`
	Description string           `msgpack:"description"`
	Unit        string           `msgpack:"unit"`
	Kind        int8             `msgpack:"kind"`
	Labels      internal.KVSlice `msgpack:"labels"`

	Time int64 `msgpack:"time"`
}

type mmsc struct {
	baseRecord

	Min   float32 `msgpack:"min"`
	Max   float32 `msgpack:"max"`
	Sum   float64 `msgpack:"sum"`
	Count int64   `msgpack:"count"`
}

func (rec *mmsc) String() string {
	return fmt.Sprintf("name=%s min=%f max=%f sum=%f count=%d",
		rec.Name, rec.Min, rec.Max, rec.Sum, rec.Count)
}

type quantile struct {
	baseRecord

	Count     int64     `msgpack:"count"`
	Quantiles []float32 `msgpack:"quantiles"`
}

func (rec *quantile) String() string {
	return fmt.Sprintf("name=%s count=%d quantiles=%v",
		rec.Name, rec.Count, rec.Quantiles)
}
