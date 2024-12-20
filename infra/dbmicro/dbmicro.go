package dbmicro

import (
	"context"
	"database/sql"
	"time"

	"github.com/langgeng-jbt/langgengpkg/contextwrap"

	"github.com/langgeng-jbt/langgengpkg/basicdto/trace"
)

type DBmicro struct {
	db *sql.DB
}

type TxMicro struct {
	tx *sql.Tx
}

func NewDBMicro(db *sql.DB) *DBmicro {
	return &DBmicro{
		db: db,
	}
}

func (d *DBmicro) QueryContext(ctx context.Context, query string, args ...interface{}) (context.Context, *sql.Rows, error) {
	start := time.Now()
	trOri := contextwrap.GetTraceFromContext(ctx)

	rows, err := d.db.QueryContext(ctx, query, args...)

	tr := &trace.TraceDatabase{
		Query:   query,
		Elapsed: time.Since(start).String(),
	}

	trProcessed := append(trOri, tr)
	ctx = contextwrap.SetTraceFromContext(ctx, trProcessed)

	if err != nil {
		return ctx, &sql.Rows{}, err
	}

	return ctx, rows, nil
}

func (d *DBmicro) ExecContext(ctx context.Context, query string, args ...interface{}) (context.Context, sql.Result, error) {
	start := time.Now()
	trOri := contextwrap.GetTraceFromContext(ctx)

	rows, err := d.db.ExecContext(ctx, query, args...)

	tr := &trace.TraceDatabase{
		Query:   query,
		Elapsed: time.Since(start).String(),
	}

	trProcessed := append(trOri, tr)
	ctx = contextwrap.SetTraceFromContext(ctx, trProcessed)

	if err != nil {
		return ctx, rows, err
	}

	return ctx, rows, nil
}

func (t *TxMicro) QueryContext(ctx context.Context, query string, args ...interface{}) (context.Context, *sql.Rows, error) {
	start := time.Now()
	trOri := contextwrap.GetTraceFromContext(ctx)

	rows, err := t.tx.QueryContext(ctx, query, args...)

	tr := &trace.TraceDatabase{
		Query:   query,
		Elapsed: time.Since(start).String(),
	}

	trProcessed := append(trOri, tr)
	ctx = contextwrap.SetTraceFromContext(ctx, trProcessed)

	if err != nil {
		return ctx, &sql.Rows{}, err
	}

	return ctx, rows, nil
}

func (t *TxMicro) ExecContext(ctx context.Context, query string, args ...interface{}) (context.Context, sql.Result, error) {
	start := time.Now()
	trOri := contextwrap.GetTraceFromContext(ctx)

	rows, err := t.tx.ExecContext(ctx, query, args...)

	tr := &trace.TraceDatabase{
		Query:   query,
		Elapsed: time.Since(start).String(),
	}

	trProcessed := append(trOri, tr)
	ctx = contextwrap.SetTraceFromContext(ctx, trProcessed)

	if err != nil {
		return ctx, rows, err
	}

	return ctx, rows, nil
}

func (t *TxMicro) Rollback() error {
	err := t.tx.Rollback()
	return err
}

func (t *TxMicro) Commit() error {
	err := t.tx.Commit()
	return err
}

func (d *DBmicro) BeginTx(ctx context.Context, opts *sql.TxOptions) (*TxMicro, error) {
	tx, err := d.db.BeginTx(ctx, opts)
	txalburaq := &TxMicro{tx}
	return txalburaq, err
}

func (d *DBmicro) Begin() (*TxMicro, error) {
	tx, err := d.db.Begin()
	txalburaq := &TxMicro{tx}
	return txalburaq, err
}
