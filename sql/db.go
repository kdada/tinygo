package sql

import "database/sql"

// 执行结果
type Result struct {
	result sql.Result //执行返回结果
	err    error      //执行错误
}

// LastInsertId 返回Insert的记录的Id
func (this *Result) LastInsertId() (int64, error) {
	return this.result.LastInsertId()
}

// RowsAffected 返回执行过程中影响的行数
func (this *Result) RowsAffected() (int64, error) {
	return this.result.RowsAffected()
}

// Error 返回执行过程中的错误
func (this *Result) Error() error {
	return this.err
}

// 数据库链接
type DB struct {
	db          *sql.DB //原始数据库链接
	transaction *sql.Tx //事务
}

// Query 查询sql
func (this *DB) Query(sqlStr string, params ...interface{}) *Rows {
	var rows *sql.Rows
	var err error
	if this.IsTransaction() {
		rows, err = this.transaction.Query(sqlStr, params...)
	} else {
		rows, err = this.db.Query(sqlStr, params...)
	}
	return &Rows{rows, err, nil}
}

// Exec 执行sql
func (this *DB) Exec(sqlStr string, params ...interface{}) *Result {
	var result sql.Result
	var err error
	if this.IsTransaction() {
		result, err = this.transaction.Exec(sqlStr, params...)
	} else {
		result, err = this.db.Exec(sqlStr, params...)
	}
	return &Result{result, err}
}

// IsTransaction 返回当前是否开启了事务
func (this *DB) IsTransaction() bool {
	return this.transaction != nil
}

// Begin 开始事务
func (this *DB) Begin() error {
	if this.IsTransaction() {
		return ErrorBeginFailed.Error()
	}
	var err error
	this.transaction, err = this.db.Begin()
	if err != nil {
		return err
	}
	return nil
}

// Commit 提交事务
func (this *DB) Commit() error {
	if this.IsTransaction() {
		var err = this.transaction.Commit()
		this.transaction = nil
		return err
	}
	return nil
}

// Rollback 回滚事务
func (this *DB) Rollback() error {
	if this.IsTransaction() {
		var err = this.transaction.Rollback()
		this.transaction = nil
		return err
	}
	return nil
}
