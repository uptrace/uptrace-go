package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

const outputLimit = 1024

var (
	cmdFlag     = flag.String("cmd", "", "command to run")
	timeoutFlag = flag.Duration("timeout", time.Hour, "command timeout")
)

var tracer = otel.Tracer("github.com/uptrace/uptrace-go")

func main() {
	flag.Usage = usage
	flag.Parse()

	if *cmdFlag == "" {
		usage()
		os.Exit(2)
	}

	var exitCode int
	defer func() {
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()

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

	uptrace.ConfigureOpentelemetry(&uptrace.Config{})
	defer uptrace.Shutdown(ctx)

	ctx, span := tracer.Start(ctx, *cmdFlag)
	defer span.End()

	span.SetAttributes(
		attribute.String("process.command_line", *cmdFlag),
	)

	stdout := NewTailWriter(make([]byte, outputLimit))
	stderr := NewTailWriter(make([]byte, outputLimit))

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdout = io.MultiWriter(os.Stdout, stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr, stderr)

	if err := cmd.Start(); err != nil {
		span.RecordError(err)
		log.Print(err)
		exitCode = 1
		return
	}

	span.SetAttributes(attribute.Int("process.pid", cmd.Process.Pid))

	err := cmd.Wait()

	if stdout.Len() > 0 {
		span.SetAttributes(attribute.String("process.stdout", stdout.Text()))
	}
	if stderr.Len() > 0 {
		span.SetAttributes(attribute.String("process.stderr", stderr.Text()))
	}

	if err != nil {
		span.RecordError(err)
		if err, ok := err.(*exec.ExitError); ok {
			exitCode = err.ExitCode()
			span.SetAttributes(
				attribute.Int("process.exit_code", exitCode),
			)
			return
		}

		log.Print(err)
		exitCode = 1
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `usage: uptrace-run [flags] -cmd="/path/to/executable"`+"\n")
	flag.PrintDefaults()
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
