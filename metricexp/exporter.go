/*
metricexp provides metric exporter for OpenTelemetry.
*/
package metricexp

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"

	"github.com/uptrace/uptrace-go/internal"
)

const (
	mmscAgg      = "mmsc"
	sumAgg       = "sum"
	lastValueAgg = "last-value"
)

type Exporter struct {
	cfg *Config

	kindSelector export.ExportKindSelector

	client   internal.SimpleClient
	endpoint string

	utime   int64
	records []metricRecord
}

var _ export.Exporter = (*Exporter)(nil)

func NewExporter(cfg *Config) (*Exporter, error) {
	cfg.init()

	e := &Exporter{
		cfg:          cfg,
		kindSelector: kindSelector(),
	}

	dsn, err := internal.ParseDSN(cfg.DSN)
	if err != nil {
		return nil, err
	}

	e.client.Client = cfg.HTTPClient
	e.client.Token = dsn.Token
	e.client.MaxRetries = 3

	e.endpoint = fmt.Sprintf("%s://%s/api/v1/meters/%s/metrics",
		dsn.Scheme, dsn.Host, dsn.ProjectID)

	return e, nil
}

func (e *Exporter) send(out interface{}) error {
	ctx := context.Background()

	data, err := internal.EncodeMsgpack(out)
	if err != nil {
		return err
	}

	return e.client.Post(ctx, e.endpoint, data)
}

func (e *Exporter) ExportKindFor(desc *metric.Descriptor, kind aggregation.Kind) export.ExportKind {
	return e.kindSelector.ExportKindFor(desc, kind)
}

func (e *Exporter) Export(_ context.Context, checkpoint export.CheckpointSet) error {
	e.utime = time.Now().UnixNano()

	if err := checkpoint.ForEach(e.kindSelector, func(record export.Record) error {
		switch agg := record.Aggregation().(type) {
		case aggregation.MinMaxSumCount:
			return e.exportMMSC(record, agg)
		case aggregation.Histogram:
			// TODO
			return nil
		case aggregation.Sum:
			return e.exportSum(record, agg)
		case aggregation.LastValue:
			return e.exportLastValue(record, agg)
		default:
			name := record.Descriptor().Name()
			internal.Logger.Printf("%s has unsupported aggregation type: %T", name, agg)
			return nil
		}
	}); err != nil {
		return err
	}

	if len(e.records) == 0 {
		return nil
	}

	if err := e.send(map[string]interface{}{"records": e.records}); err != nil {
		internal.Logger.Printf("send failed: %s", err)
	}
	e.records = nil

	return nil
}

func (e *Exporter) exportMMSC(record export.Record, agg aggregation.MinMaxSumCount) error {
	e.records = append(e.records, metricRecord{})
	out := &e.records[len(e.records)-1]

	if err := e.exportCommon(record, out); err != nil {
		return err
	}

	out.Aggregation = mmscAgg
	desc := record.Descriptor()
	numKind := desc.NumberKind()

	min, err := agg.Min()
	if err != nil {
		return err
	}

	max, err := agg.Max()
	if err != nil {
		return err
	}

	sum, err := agg.Sum()
	if err != nil {
		return err
	}

	count, err := agg.Count()
	if err != nil {
		return err
	}

	out.Data = &mmscData{
		Min:   min.CoerceToFloat64(numKind),
		Max:   max.CoerceToFloat64(numKind),
		Sum:   sum.CoerceToFloat64(numKind),
		Count: count,
	}

	return nil
}

func (e *Exporter) exportSum(record export.Record, agg aggregation.Sum) error {
	e.records = append(e.records, metricRecord{})
	out := &e.records[len(e.records)-1]

	if err := e.exportCommon(record, out); err != nil {
		return err
	}

	out.Aggregation = sumAgg
	desc := record.Descriptor()
	numKind := desc.NumberKind()

	sum, err := agg.Sum()
	if err != nil {
		return err
	}

	out.Data = &sumData{
		Sum: sum.CoerceToFloat64(numKind),
	}

	return nil
}

func (e *Exporter) exportLastValue(
	record export.Record, agg aggregation.LastValue,
) error {
	e.records = append(e.records, metricRecord{})
	out := &e.records[len(e.records)-1]

	if err := e.exportCommon(record, out); err != nil {
		return err
	}

	out.Aggregation = lastValueAgg
	desc := record.Descriptor()
	numKind := desc.NumberKind()

	value, _, err := agg.LastValue()
	if err != nil {
		return err
	}

	out.Data = &lastValueData{
		Value: value.CoerceToFloat64(numKind),
	}

	return nil
}

func (e *Exporter) exportCommon(record export.Record, out *metricRecord) error {
	desc := record.Descriptor()

	out.Name = desc.Name()
	out.Description = desc.Description()
	out.Unit = string(desc.Unit())
	out.Instrument = instrumentKind(desc.InstrumentKind())
	out.Time = e.utime

	out.MeterName = desc.InstrumentationName()
	out.MeterVersion = desc.InstrumentationVersion()

	if res := record.Resource(); res != nil {
		out.Resource = res.Attributes()
	}

	if iter := record.Labels().Iter(); iter.Len() > 0 {
		attrs := make([]attribute.KeyValue, 0, iter.Len())
		for iter.Next() {
			attrs = append(attrs, iter.Attribute())
		}
		out.Attrs = attrs
	}

	return nil
}

func kindSelector() export.ExportKindSelector {
	return export.StatelessExportKindSelector()
}

type metricRecord struct {
	Name        string `msgpack:"name"`
	Description string `msgpack:"description"`
	Unit        string `msgpack:"unit"`
	Aggregation string `msgpack:"aggregation"`
	Instrument  string `msgpack:"instrument"`

	MeterName    string `msgpack:"meterName"`
	MeterVersion string `msgpack:"meterVersion"`

	Resource internal.KeyValueSlice `msgpack:"resource"`
	Attrs    internal.KeyValueSlice `msgpack:"attrs"`

	Data interface{} `msgpack:"data"`

	Time int64 `msgpack:"time"`
}

type mmscData struct {
	Min   float64 `msgpack:"min"`
	Max   float64 `msgpack:"max"`
	Sum   float64 `msgpack:"sum"`
	Count uint64  `msgpack:"count"`
}

type lastValueData struct {
	Value float64 `msgpack:"value"`
}

type sumData struct {
	Sum float64 `msgpack:"sum"`
}

func instrumentKind(kind metric.InstrumentKind) string {
	switch kind {
	case metric.ValueRecorderInstrumentKind:
		return "value-recorder"
	case metric.ValueObserverInstrumentKind:
		return "value-observer"
	case metric.CounterInstrumentKind:
		return "counter"
	case metric.UpDownCounterInstrumentKind:
		return "up-down-counter"
	case metric.SumObserverInstrumentKind:
		return "sum-observer"
	case metric.UpDownSumObserverInstrumentKind:
		return "up-down-sum-observer"
	default:
		return "invalid"
	}
}
