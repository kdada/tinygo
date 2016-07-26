//  Package validator 实现了验证器的基本接口
package validator

import "sync"

// 验证器
type Validator interface {
	Validate(str string) bool
}

// 创建器
//  suorce: 用于创建验证器的信息
type ValidatorCreator func(source string) (Validator, error)

var (
	mu       sync.Mutex                          //互斥锁
	creators = make(map[string]ValidatorCreator) //创建器映射
)

// NewValidator 创建一个新的Validator
//  kind:类型
func NewValidator(kind string, source string) (Validator, error) {
	var creator, ok = creators[kind]
	if !ok {
		return nil, ErrorInvalidKind.Format(kind).Error()
	}
	return creator(source)
}

// Register 注册创建器
func Register(kind string, creator ValidatorCreator) {
	if creator == nil {
		panic(ErrorInvalidValidatorCreator.Format(kind))
	}
	mu.Lock()
	defer mu.Unlock()
	creators[kind] = creator
}
