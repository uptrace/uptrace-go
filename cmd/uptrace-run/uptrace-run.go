package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/label"
)

const outputLimit = 1024

var (
	cmdFlag     = flag.String("cmd", "", "command to run")
	timeoutFlag = flag.Duration("timeout", time.Hour, "command timeout")
)

var tracer = global.Tracer("github.com/uptrace/uptrace-go")

func main() {
	flag.Parse()

	var args []string

	if strings.IndexByte(*cmdFlag, ' ') >= 0 {
		args = []string{"sh", "-c", *cmdFlag}
	} else {
		args = []string{*cmdFlag}
	}

	ctx := context.Background()
	if timeoutFlag != nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *timeoutFlag)
		defer cancel()
	}

	upclient := setupUptrace()

	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	ctx, span := tracer.Start(ctx, *cmdFlag)
	defer span.End()

	span.SetAttributes(
		label.String("process.command_line", *cmdFlag),
	)

	stdout := NewTailWriter(make([]byte, outputLimit))
	stderr := NewTailWriter(make([]byte, outputLimit))

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdout = io.MultiWriter(os.Stdout, stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr, stderr)

	if err := cmd.Start(); err != nil {
		span.RecordError(ctx, err)
		log.Print(err)
		return
	}

	span.SetAttributes(label.Int("process.pid", cmd.Process.Pid))

	err := cmd.Wait()

	if stdout.Len() > 0 {
		span.SetAttributes(label.String("process.stdout", stdout.Text()))
	}
	if stderr.Len() > 0 {
		span.SetAttributes(label.String("process.stderr", stderr.Text()))
	}

	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			span.SetAttributes(
				label.Int("process.exit_code", err.ExitCode()),
			)
		}

		span.RecordError(ctx, err)
		return
	}
}

func setupUptrace() *uptrace.Client {
	hostname, _ := os.Hostname()
	upclient := uptrace.NewClient(&uptrace.Config{
		Resource: map[string]interface{}{
			"host.name": hostname,
		},
	})
	return upclient
}

//------------------------------------------------------------------------------

type TailWriter struct {
	b []byte
	i int
}

func NewTailWriter(b []byte) *TailWriter {
	return &TailWriter{
		b: b,
	}
}

func (w *TailWriter) Len() int {
	return w.i
}

func (w *TailWriter) Text() string {
	return string(w.b[:w.i])
}

func (w *TailWriter) Write(b []byte) (int, error) {
	written := len(b)

	if len(b) >= len(w.b) {
		w.i = 0
		b = b[len(b)-len(w.b):]
	} else if len(b) > w.available() {
		i := len(w.b) - len(b)
		copy(w.b, w.b[w.i-i:])
		w.i = i
	}

	copy(w.b[w.i:], b)
	w.i += len(b)

	return written, nil
}

func (w *TailWriter) available() int {
	return len(w.b) - w.i
}
