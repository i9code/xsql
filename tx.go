package xsql

import (
	`github.com/i9code/xlog`
	`xorm.io/xorm`
)

type (
	// Tx 事务控制
	Tx struct {
		engine *xorm.Engine
	}

	txFun  func(tx *Session) (err error)
	txpFun func(tx *Session, params ...interface{}) (err error)
)

// 事务控制
func newTx(engine *xorm.Engine) *Tx {
	return &Tx{
		engine: engine,
	}
}

func (t *Tx) Do(fun txFun, fields ...interface{}) (err error) {
	return t.do(func(tx *Session) error {
		return fun(tx)
	}, fields...)
}

func (t *Tx) Dop(fun txpFun, params []interface{}, fields ...interface{}) (err error) {
	return t.do(func(tx *Session) error {
		return fun(tx, params...)
	}, fields...)
}

func (t *Tx) do(fun func(tx *Session) error, fields ...interface{}) (err error) {
	session := t.engine.NewSession()
	if err = t.begin(session, fields...); nil != err {
		return
	}
	defer t.close(session, fields...)

	if err = fun(&Session{Session: session}); nil != err {
		t.rollback(session, fields...)
	} else {
		t.commit(session, fields...)
	}

	return
}

func (t *Tx) begin(tx *xorm.Session, fields ...interface{}) (err error) {
	if err = tx.Commit(); nil != err {
		t.error(err, "开始数据库事务出错", fields...)
	}

	return
}

func (t *Tx) commit(tx *xorm.Session, fields ...interface{}) {
	if err := tx.Commit(); nil != err {
		t.error(err, "提交数据库事务出错", fields...)
	}
}

func (t *Tx) close(tx *xorm.Session, fields ...interface{}) {
	if err := tx.Close(); nil != err {
		t.error(err, "关闭数据库事务出错", fields...)
	}
}

func (t *Tx) rollback(tx *xorm.Session, fields ...interface{}) {
	if err := tx.Rollback(); nil != err {
		t.error(err, "回退数据库事务出错", fields...)
	}
}

func (t *Tx) error(err error, msg string, fields ...interface{}) {
	xlog.Error(err, msg, fields)
}
