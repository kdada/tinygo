package sql

import (
	"database/sql"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// 数据行
type Rows struct {
	rows    *sql.Rows
	err     error
	columns map[string]int
}

// 正则替换
var reg = regexp.MustCompile(`\B[A-Z]`)

// TransFieldName 转换字段名称,默认在大写字母前添加下划线(首字母除外)
var TransFieldName = func(name string) string {
	return reg.ReplaceAllString(name, "_$0")
}

// parse 解析fields值到value中
func (this *Rows) parse(value reflect.Value, index int, fields []interface{}) error {
	switch value.Kind() {
	case reflect.Bool:
		var b = sql.NullBool{}
		var err = b.Scan(*(fields[index].(*interface{})))
		if err != nil {
			return err
		}
		if b.Valid {
			value.SetBool(b.Bool)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var i = sql.NullInt64{}
		var err = i.Scan(*(fields[index].(*interface{})))
		if err != nil {
			return err
		}
		if i.Valid {
			value.SetInt(i.Int64)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var i = sql.NullInt64{}
		var err = i.Scan(*(fields[index].(*interface{})))
		if err != nil {
			return err
		}
		if i.Valid {
			value.SetUint(uint64(i.Int64))
		}
	case reflect.Float32, reflect.Float64:
		var f = sql.NullFloat64{}
		var err = f.Scan(*(fields[index].(*interface{})))
		if err != nil {
			return err
		}
		if f.Valid {
			value.SetFloat(f.Float64)
		}
	case reflect.String:
		var s = sql.NullString{}
		var err = s.Scan(*(fields[index].(*interface{})))
		if err != nil {
			return err
		}
		if s.Valid {
			value.SetString(s.String)
		}
	case reflect.Struct:
		{
			if value.Type().String() == "time.Time" {
				//时间结构体解析
				var v = *(fields[index].(*interface{}))
				if v != nil && reflect.TypeOf(v).String() == "time.Time" {
					value.Set(reflect.ValueOf(v))
				} else {
					var s = sql.NullString{}
					var err = s.Scan(v)
					if err != nil {
						return err
					}
					if s.Valid {
						result, err := time.ParseInLocation("2006-01-02 15:04:05", s.String, time.Local)
						if err == nil {
							value.Set(reflect.ValueOf(result))
						} else {
							return err
						}
					}
				}
			} else {
				//常规结构体解析
				for i := 0; i < value.NumField(); i++ {
					var fieldValue = value.Field(i)
					var fieldType = value.Type().Field(i)
					if fieldType.Anonymous {
						//匿名组合字段,进行递归解析
						this.parse(fieldValue, 0, fields)
					} else {
						//非匿名字段
						if fieldValue.CanSet() {
							var fieldName = fieldType.Tag.Get("col")
							if fieldName == "-" {
								//如果是-,则忽略当前字段
								continue
							}
							if fieldName == "" {
								//如果为空,则使用字段名
								fieldName = TransFieldName(fieldType.Name)
							}
							var index, ok = this.columns[strings.ToLower(fieldName)]
							if ok {
								this.parse(fieldValue, index, fields)
							}
						}
					}
				}
			}
		}

	}
	return nil
}

// scan 扫描单行数据
func (this *Rows) scan(data reflect.Value) error {
	if this.columns == nil {
		var columns, err = this.rows.Columns()
		if err != nil {
			return this.err
		}
		this.columns = make(map[string]int, len(columns))
		for i, n := range columns {
			this.columns[n] = i
		}
	}
	var fields = make([]interface{}, len(this.columns))
	for i := 0; i < len(fields); i++ {
		var pif interface{}
		fields[i] = &pif
	}
	var err = this.rows.Scan(fields...)
	if err == nil {
		err = this.parse(data, 0, fields)
	}
	return err
}

// Scan 扫描数据行
//  data:将数据行中的数据解析到data中,data可以是 基础类型,time.Time类型,结构体,数组类型 的指针
//  return:(扫描的行数,错误)
func (this *Rows) Scan(data interface{}) (int, error) {
	if this.err == nil {
		// 类型解析
		var d, err = newData(data)
		if err != nil {
			return 0, err
		}
		//行解析
		for this.rows.Next() && d.Next() {
			var n = d.New()
			this.err = this.scan(n)
			if this.err != nil {
				return 0, this.err
			}
			d.SetBack(n)
		}
		this.rows.Close()
		return d.length, nil
	}
	return 0, this.err
}

// Error 返回数据行错误
func (this *Rows) Error() error {
	return this.err
}

// 数据类型
type data struct {
	t       reflect.Type
	v       reflect.Value
	setType reflect.Type
	length  int
	slice   bool
}

// newData 创建一个data
func newData(v interface{}) (*data, error) {
	var d = new(data)
	d.t = reflect.TypeOf(v)
	d.v = reflect.ValueOf(v)
	if d.v.Kind() == reflect.Ptr {
		//取指针指向的值
		d.t = d.t.Elem()
		d.v = d.v.Elem()
		switch d.t.Kind() {
		case reflect.Slice:
			{
				d.slice = true
				d.setType = d.t.Elem()
			}
		default:
			{
				d.setType = d.t
				if d.t.Kind() == reflect.Ptr {
					//如果对象为指针
					d.t = d.t.Elem()
					d.v.Set(reflect.New(d.t))
					d.v = d.v.Elem()
				}
			}
		}
		return d, nil
	}
	return nil, ErrorInvalidParamType.Format(d.t.Name()).Error()
}

// New 获取一个可Set的值
func (this *data) New() reflect.Value {
	this.length++
	if this.slice {
		var v reflect.Value
		if this.setType.Kind() == reflect.Ptr {
			v = reflect.New(this.setType.Elem()).Elem()
		} else {
			v = reflect.New(this.setType).Elem()
		}
		return v
	}
	return this.v
}

// SetBack 将New()的值设置回data
func (this *data) SetBack(value reflect.Value) {
	if this.slice {
		var v = value
		if this.setType.Kind() == reflect.Ptr {
			v = v.Addr()
		}
		this.v.Set(reflect.Append(this.v, v))
	}

}

// Next 能否继续获取
func (this *data) Next() bool {
	if this.slice {
		return true
	}
	return this.length < 1
}
