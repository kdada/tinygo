package meta

import "reflect"

// 不能明确的类型元数据
type OtherMetadata struct {
	Name string       //名称
	Type reflect.Type //类型
}

// Generate 根据vp提供的值生成相应值
func (this *OtherMetadata) Generate(vp ValueProvider) (interface{}, error) {
	if vp.Contains(this.Name, this.Type) {
		return vp.Value(this.Name, this.Type), nil
	}
	return reflect.New(this.Type).Elem().Interface(), nil
}

// AnalyzeOther 分析其他类型并生成OtherMetadata
func AnalyzeOther(t reflect.Type) (Generator, error) {
	return &OtherMetadata{
		t.Name(),
		t,
	}, nil
}
