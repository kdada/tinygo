package meta

import "reflect"

// 结构体元数据
type StructMetadata struct {
	Name   string           //结构体全名(包含包名)
	Struct reflect.Type     //结构体类型
	Fields []*FieldMetadata //结构体字段元数据
}

// generate 根据vp提供的值生成相应结构体值的指针
func (this *StructMetadata) generate(vc ValueContainer) (interface{}, error) {
	var vp, ok = vc.Contains(this.Name, this.Struct)
	if ok {
		return vp.Value(), nil
	}
	var result = reflect.New(this.Struct)
	for _, fMd := range this.Fields {
		var err = fMd.Set(result, vc)
		if err != nil {
			return nil, err
		}
	}
	return result.Interface(), nil
}

// Generate 根据vp提供的值生成相应值
func (this *StructMetadata) Generate(vc ValueContainer) (interface{}, error) {
	var vp, ok = vc.Contains(this.Name, this.Struct)
	if ok {
		return vp.Value(), nil
	}
	var v, err = this.generate(vc)
	if err != nil {
		return nil, err
	}
	return reflect.ValueOf(v).Elem().Interface(), err
}

// 结构体指针元数据
type StructPtrMetadata struct {
	Name string          //结构体指针全名(包含包名)
	Type reflect.Type    //结构体指针类型
	Ptr  *StructMetadata //结构体元数据
}

// Generate 根据vp提供的值生成相应值
func (this *StructPtrMetadata) Generate(vc ValueContainer) (interface{}, error) {
	var vp, ok = vc.Contains(this.Name, this.Type)
	if ok {
		return vp.Value(), nil
	}
	return this.Ptr.generate(vc)
}

// 全局结构体元数据信息
var globalStructMetadata = make(map[string]*StructMetadata)

// AnalyzeField 分析结构体
//  s:结构体类型,必须是结构体或者结构体指针
func AnalyzeStruct(s reflect.Type) (Generator, error) {
	if s.Kind() != reflect.Struct && !IsStructPtrType(s) {
		return nil, ErrorParamNotStruct.Format(s.String()).Error()
	}
	var st = s
	if IsStructPtrType(s) {
		st = s.Elem()
	}
	var sMd, ok = globalStructMetadata[st.String()]
	if !ok {
		sMd = new(StructMetadata)
		sMd.Name = st.String()
		sMd.Struct = st
		sMd.Fields = make([]*FieldMetadata, 0)
		var err = ForeachField(st, func(field reflect.StructField) error {
			var fMd, e = AnalyzeField(&field)
			if e != nil {
				return e
			}
			sMd.Fields = append(sMd.Fields, fMd)
			return nil
		})
		if err != nil {
			return nil, err
		}
		globalStructMetadata[st.String()] = sMd
	}

	// 当参数是结构体指针时,构造结构体指针元数据
	if s.Kind() == reflect.Ptr {
		return &StructPtrMetadata{
			s.String(),
			s,
			sMd,
		}, nil
	}
	return sMd, nil
}
