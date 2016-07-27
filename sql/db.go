package sql

import "database/sql"

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
func (this *DB) Exec(sqlStr string, params ...interface{}) (Result, error) {
	if this.IsTransaction() {
		return this.transaction.Exec(sqlStr, params...)
	} else {
		return this.db.Exec(sqlStr, params...)
	}
}

// IsTransaction 返回当前是否开启了事务
func (this *DB) IsTransaction() bool {
	return this.transaction != nil
}

// Begin 开始事务
func (this *DB) Begin() error {
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
