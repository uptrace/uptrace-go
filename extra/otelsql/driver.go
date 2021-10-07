package otelsql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var dbRowsAffected = attribute.Key("db.rows_affected")

// Open is a wrapper over sql.Open that instruments the sql.DB to record executed queries
// using OpenTelemetry API.
func Open(driverName, dsn string, opts ...Option) (*sql.DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}

	dbDriver := db.Driver()
	d := newDriver(dbDriver, opts...)

	if _, ok := dbDriver.(driver.DriverContext); ok {
		connector, err := d.OpenConnector(dsn)
		if err != nil {
			return nil, err
		}
		return sql.OpenDB(connector), nil
	}

	return sql.OpenDB(&dsnConnector{
		driver: d,
		dsn:    dsn,
	}), nil
}

// OpenDB is a wrapper over sql.OpenDB that instruments the sql.DB to record executed queries
// using OpenTelemetry API.
func OpenDB(connector driver.Connector, opts ...Option) *sql.DB {
	cfg := newConfig(opts...)
	return sql.OpenDB(newConnector(connector.Driver(), connector, cfg))
}

type dsnConnector struct {
	driver *otelDriver
	dsn    string
}

func (c *dsnConnector) Connect(ctx context.Context) (driver.Conn, error) {
	var conn driver.Conn
	err := c.driver.cfg.withSpan(ctx, "db.Connect", func(ctx context.Context, span trace.Span) error {
		var err error
		conn, err = c.driver.Open(c.dsn)
		return err
	})
	return conn, err
}

func (c *dsnConnector) Driver() driver.Driver {
	return c.driver
}

//------------------------------------------------------------------------------

type otelDriver struct {
	driver    driver.Driver
	driverCtx driver.DriverContext
	cfg       *config
}

var _ driver.DriverContext = (*otelDriver)(nil)

func newDriver(dr driver.Driver, opts ...Option) *otelDriver {
	driverCtx, _ := dr.(driver.DriverContext)
	d := &otelDriver{
		driver:    dr,
		driverCtx: driverCtx,
		cfg:       newConfig(opts...),
	}
	return d
}

func (d *otelDriver) Open(name string) (driver.Conn, error) {
	conn, err := d.driver.Open(name)
	if err != nil {
		return nil, err
	}
	return newConn(conn, d.cfg), nil
}

func (d *otelDriver) OpenConnector(dsn string) (driver.Connector, error) {
	connector, err := d.driverCtx.OpenConnector(dsn)
	if err != nil {
		return nil, err
	}
	return newConnector(d, connector, d.cfg), nil
}

//------------------------------------------------------------------------------

type connector struct {
	driver.Connector
	driver driver.Driver
	cfg    *config
}

var _ driver.Connector = (*connector)(nil)

func newConnector(d driver.Driver, c driver.Connector, cfg *config) *connector {
	return &connector{
		driver:    d,
		Connector: c,
		cfg:       cfg,
	}
}

func (c *connector) Connect(ctx context.Context) (driver.Conn, error) {
	var conn driver.Conn
	if err := c.cfg.withSpan(ctx, "db.Connect", func(ctx context.Context, span trace.Span) error {
		var err error
		conn, err = c.Connector.Connect(ctx)
		return err
	}); err != nil {
		return nil, err
	}
	return newConn(conn, c.cfg), nil
}

func (c *connector) Driver() driver.Driver {
	return c.driver
}

//------------------------------------------------------------------------------

type otelConn struct {
	driver.Conn

	cfg *config

	ping         pingFunc
	exec         execFunc
	execCtx      execCtxFunc
	query        queryFunc
	queryCtx     queryCtxFunc
	prepareCtx   prepareCtxFunc
	beginTx      beginTxFunc
	resetSession resetSessionFunc
}

var _ driver.Conn = (*otelConn)(nil)

func newConn(conn driver.Conn, cfg *config) *otelConn {
	cn := &otelConn{
		Conn: conn,
		cfg:  cfg,
	}

	cn.ping = cn.createPingFunc(conn)
	cn.exec = cn.createExecFunc(conn)
	cn.execCtx = cn.createExecCtxFunc(conn)
	cn.query = cn.createQueryFunc(conn)
	cn.queryCtx = cn.createQueryCtxFunc(conn)
	cn.prepareCtx = cn.createPrepareCtxFunc(conn)
	cn.beginTx = cn.createBeginTxFunc(conn)
	cn.resetSession = cn.createResetSessionFunc(conn)

	return cn
}

var _ driver.Pinger = (*otelConn)(nil)

func (c *otelConn) Ping(ctx context.Context) error {
	return c.ping(ctx)
}

type pingFunc func(ctx context.Context) error

func (c *otelConn) createPingFunc(conn driver.Conn) pingFunc {
	if pinger, ok := conn.(driver.Pinger); ok {
		return func(ctx context.Context) error {
			return c.cfg.withSpan(ctx, "db.Ping", func(ctx context.Context, span trace.Span) error {
				return pinger.Ping(ctx)
			})
		}
	}
	return func(ctx context.Context) error {
		return driver.ErrSkip
	}
}

//------------------------------------------------------------------------------

var _ driver.Execer = (*otelConn)(nil)

func (c *otelConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	return c.exec(query, args)
}

type execFunc func(query string, args []driver.Value) (driver.Result, error)

func (c *otelConn) createExecFunc(conn driver.Conn) execFunc {
	if execer, ok := conn.(driver.Execer); ok {
		return func(query string, args []driver.Value) (driver.Result, error) {
			return execer.Exec(query, args)
		}
	}
	return func(query string, args []driver.Value) (driver.Result, error) {
		return nil, driver.ErrSkip
	}
}

//------------------------------------------------------------------------------

var _ driver.ExecerContext = (*otelConn)(nil)

func (c *otelConn) ExecContext(
	ctx context.Context, query string, args []driver.NamedValue,
) (driver.Result, error) {
	return c.execCtx(ctx, query, args)
}

type execCtxFunc func(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error)

func (c *otelConn) createExecCtxFunc(conn driver.Conn) execCtxFunc {
	var fn execCtxFunc

	if execer, ok := conn.(driver.ExecerContext); ok {
		fn = execer.ExecContext
	} else {
		fn = func(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
			vArgs, err := namedValueToValue(args)
			if err != nil {
				return nil, err
			}
			return c.exec(query, vArgs)
		}
	}

	return func(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
		var res driver.Result
		if err := c.cfg.withSpan(ctx, "db.Exec", func(ctx context.Context, span trace.Span) error {
			isRecording := span.IsRecording()

			if isRecording {
				span.SetAttributes(semconv.DBStatementKey.String(c.cfg.formatQuery(query)))
			}

			var err error
			res, err = fn(ctx, query, args)
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
		}); err != nil {
			return nil, err
		}
		return res, nil
	}
}

//------------------------------------------------------------------------------

var _ driver.Queryer = (*otelConn)(nil)

func (c *otelConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	return c.query(query, args)
}

type queryFunc func(query string, args []driver.Value) (driver.Rows, error)

func (c *otelConn) createQueryFunc(conn driver.Conn) queryFunc {
	if queryer, ok := c.Conn.(driver.Queryer); ok {
		return func(query string, args []driver.Value) (driver.Rows, error) {
			return queryer.Query(query, args)
		}
	}
	return func(query string, args []driver.Value) (driver.Rows, error) {
		return nil, driver.ErrSkip
	}
}

//------------------------------------------------------------------------------

var _ driver.QueryerContext = (*otelConn)(nil)

func (c *otelConn) QueryContext(
	ctx context.Context, query string, args []driver.NamedValue,
) (driver.Rows, error) {
	return c.queryCtx(ctx, query, args)
}

type queryCtxFunc func(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error)

func (c *otelConn) createQueryCtxFunc(conn driver.Conn) queryCtxFunc {
	var fn queryCtxFunc

	if queryer, ok := c.Conn.(driver.QueryerContext); ok {
		fn = queryer.QueryContext
	} else {
		fn = func(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
			vArgs, err := namedValueToValue(args)
			if err != nil {
				return nil, err
			}
			return c.query(query, vArgs)
		}
	}

	return func(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
		var rows driver.Rows
		err := c.cfg.withSpan(ctx, "db.Query", func(ctx context.Context, span trace.Span) error {
			if span.IsRecording() {
				span.SetAttributes(semconv.DBStatementKey.String(c.cfg.formatQuery(query)))
			}

			var err error
			rows, err = fn(ctx, query, args)
			return err
		})
		return rows, err
	}
}

//------------------------------------------------------------------------------

var _ driver.ConnPrepareContext = (*otelConn)(nil)

func (c *otelConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return c.prepareCtx(ctx, query)
}

type prepareCtxFunc func(ctx context.Context, query string) (driver.Stmt, error)

func (c *otelConn) createPrepareCtxFunc(conn driver.Conn) prepareCtxFunc {
	var fn prepareCtxFunc

	if preparer, ok := c.Conn.(driver.ConnPrepareContext); ok {
		fn = preparer.PrepareContext
	} else {
		fn = func(ctx context.Context, query string) (driver.Stmt, error) {
			return c.Conn.Prepare(query)
		}
	}

	return func(ctx context.Context, query string) (driver.Stmt, error) {
		var stmt driver.Stmt
		if err := c.cfg.withSpan(ctx, "db.Prepare", func(
			ctx context.Context, span trace.Span) error {
			if span.IsRecording() {
				span.SetAttributes(semconv.DBStatementKey.String(c.cfg.formatQuery(query)))
			}

			var err error
			stmt, err = fn(ctx, query)
			return err
		}); err != nil {
			return nil, err
		}
		return newStmt(stmt, query, c.cfg), nil
	}
}

var _ driver.ConnBeginTx = (*otelConn)(nil)

func (c *otelConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return c.beginTx(ctx, opts)
}

type beginTxFunc func(ctx context.Context, opts driver.TxOptions) (driver.Tx, error)

func (c *otelConn) createBeginTxFunc(conn driver.Conn) beginTxFunc {
	var fn beginTxFunc

	if txor, ok := conn.(driver.ConnBeginTx); ok {
		fn = txor.BeginTx
	} else {
		fn = func(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
			return conn.Begin()
		}
	}

	return func(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
		var tx driver.Tx
		if err := c.cfg.withSpan(ctx, "db.Begin", func(
			ctx context.Context, span trace.Span) error {
			var err error
			tx, err = fn(ctx, opts)
			return err
		}); err != nil {
			return nil, err
		}
		return newTx(ctx, tx, c.cfg), nil
	}
}

//------------------------------------------------------------------------------

var _ driver.SessionResetter = (*otelConn)(nil)

func (c *otelConn) ResetSession(ctx context.Context) error {
	return c.resetSession(ctx)
}

type resetSessionFunc func(ctx context.Context) error

func (c *otelConn) createResetSessionFunc(conn driver.Conn) resetSessionFunc {
	if resetter, ok := c.Conn.(driver.SessionResetter); ok {
		return func(ctx context.Context) error {
			return resetter.ResetSession(ctx)
		}
	}
	return func(ctx context.Context) error {
		return driver.ErrSkip
	}
}

//------------------------------------------------------------------------------

func namedValueToValue(named []driver.NamedValue) ([]driver.Value, error) {
	args := make([]driver.Value, len(named))
	for n, param := range named {
		if len(param.Name) > 0 {
			return nil, errors.New("otelsql: driver does not support the use of Named Parameters")
		}
		args[n] = param.Value
	}
	return args, nil
}
