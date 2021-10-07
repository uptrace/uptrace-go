package otelsql

import (
	"context"
	"database/sql/driver"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type otelStmt struct {
	driver.Stmt

	query string
	cfg   *config

	execCtx  stmtExecCtxFunc
	queryCtx stmtQueryCtxFunc
}

var _ driver.Stmt = (*otelStmt)(nil)

func newStmt(stmt driver.Stmt, query string, cfg *config) *otelStmt {
	s := &otelStmt{
		Stmt:  stmt,
		query: query,
		cfg:   cfg,
	}
	s.execCtx = s.createExecCtxFunc(stmt)
	s.queryCtx = s.createQueryCtxFunc(stmt)
	return s
}

//------------------------------------------------------------------------------

var _ driver.StmtExecContext = (*otelStmt)(nil)

func (stmt *otelStmt) ExecContext(
	ctx context.Context, args []driver.NamedValue,
) (driver.Result, error) {
	return stmt.execCtx(ctx, args)
}

type stmtExecCtxFunc func(ctx context.Context, args []driver.NamedValue) (driver.Result, error)

func (s *otelStmt) createExecCtxFunc(stmt driver.Stmt) stmtExecCtxFunc {
	var fn stmtExecCtxFunc

	if execer, ok := s.Stmt.(driver.StmtExecContext); ok {
		fn = execer.ExecContext
	} else {
		fn = func(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
			vArgs, err := namedValueToValue(args)
			if err != nil {
				return nil, err
			}
			return stmt.Exec(vArgs)
		}
	}

	return func(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
		var res driver.Result
		err := s.cfg.withSpan(ctx, "stmt.Exec", func(ctx context.Context, span trace.Span) error {
			isRecording := span.IsRecording()

			if isRecording {
				span.SetAttributes(semconv.DBStatementKey.String(s.cfg.formatQuery(s.query)))
			}

			var err error
			res, err = fn(ctx, args)
			if err != nil {
				return err
			}

			if isRecording {
				rows, err := res.RowsAffected()
				if err == nil {
					span.SetAttributes(dbRowsAffected.Int64(rows))
				}
			}

			return nil
		})
		return res, err
	}
}

//------------------------------------------------------------------------------

var _ driver.StmtQueryContext = (*otelStmt)(nil)

func (stmt *otelStmt) QueryContext(
	ctx context.Context, args []driver.NamedValue,
) (driver.Rows, error) {
	return stmt.queryCtx(ctx, args)
}

type stmtQueryCtxFunc func(ctx context.Context, args []driver.NamedValue) (driver.Rows, error)

func (s *otelStmt) createQueryCtxFunc(stmt driver.Stmt) stmtQueryCtxFunc {
	var fn stmtQueryCtxFunc

	if queryer, ok := s.Stmt.(driver.StmtQueryContext); ok {
		fn = queryer.QueryContext
	} else {
		fn = func(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
			vArgs, err := namedValueToValue(args)
			if err != nil {
				return nil, err
			}
			return s.Query(vArgs)
		}
	}

	return func(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
		var rows driver.Rows
		err := s.cfg.withSpan(ctx, "stmt.Query", func(ctx context.Context, span trace.Span) error {
			if span.IsRecording() {
				span.SetAttributes(semconv.DBStatementKey.String(s.cfg.formatQuery(s.query)))
			}

			var err error
			rows, err = fn(ctx, args)
			return err
		})
		return rows, err
	}
}
