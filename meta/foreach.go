package meta

import "reflect"

// IsStructPtrType 判断是否是结构体指针类型
func IsStructPtrType(instance reflect.Type) bool {
	return instance.Kind() == reflect.Ptr && instance.Elem().Kind() == reflect.Struct
}

// ForeachMethod 遍历instance的所有方法
func ForeachMethod(instance reflect.Type, solve func(method reflect.Method) error) error {
	for i := 0; i < instance.NumMethod(); i++ {
		var err = solve(instance.Method(i))
		if err != nil {
			return err
		}
	}
	return nil
}

// ForeachField 遍历instance的所有字段(包括匿名字段,但不包括私有字段),instance必须是结构体或结构体指针类型
func ForeachField(instance reflect.Type, solve func(field reflect.StructField) error) error {
	if IsStructPtrType(instance) {
		instance = instance.Elem()
	}
	if instance.Kind() == reflect.Struct {
		var preIndex = make([]int, 0, 2)
		return foreachField(preIndex, instance, solve)
	}
	return ErrorParamNotStruct.Format(instance.String()).Error()
}

// foreachField 遍历instance的所有字段(包括匿名字段,但不包括私有字段),instance必须是结构体或结构体指针类型
func foreachField(preIndex []int, instance reflect.Type, solve func(field reflect.StructField) error) error {
	for i := 0; i < instance.NumField(); i++ {
		var field = instance.Field(i)
		if (field.Name[0] >= 'A') && (field.Name[0] <= 'Z') {
			if field.Anonymous {
				preIndex = append(preIndex, i)
				foreachField(preIndex, field.Type, solve)
				preIndex = preIndex[:len(preIndex)-1]
			} else {
				var lenPI = len(preIndex)
				if lenPI > 0 {
					var index = make([]int, lenPI+len(field.Index))
					copy(index, preIndex)
					copy(index[lenPI:], field.Index)
					field.Index = index
				}
				var err = solve(field)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// ForeachParam 遍历instance所有参数,instance必须是函数类型
func ForeachParam(instance reflect.Type, solve func(param reflect.Type) error) error {
	if instance.Kind() == reflect.Func {
		for i := 0; i < instance.NumIn(); i++ {
			var err = solve(instance.In(i))
			if err != nil {
				return err
			}
		}
		return nil
	}
	return ErrorParamNotFunc.Format(instance.String()).Error()
}

// ForeachResult 遍历instance所有返回值,instance必须是函数类型
func ForeachResult(instance reflect.Type, solve func(result reflect.Type) error) error {
	if instance.Kind() == reflect.Func {
		for i := 0; i < instance.NumOut(); i++ {
			var err = solve(instance.Out(i))
			if err != nil {
				return err
			}
		}
		return nil
	}
	return ErrorParamNotFunc.Format(instance.String()).Error()
}
