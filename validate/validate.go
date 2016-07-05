package validate

import "sync"

// Validator
type Validator interface {
	Validate(value string) bool
}

// Validator创建器
//  source:验证规则
type ValidatorCreator func(source string) (Validator, error)

var (
	mu       sync.Mutex                          //互斥锁
	creators = make(map[string]ValidatorCreator) //Validator创建器映射
)

// NewValidator 创建一个新的Validator
//  expire:最大过期时间(秒)
func NewValidator(kind string, source string) (Validator, error) {
	var creator, ok = creators[kind]
	if !ok {
		return nil, ErrorInvalidValidatorKind.Format(kind).Error()
	}
	return creator(source)
}

// Register 注册Validator创建器
func Register(kind string, creator ValidatorCreator) error {
	if creator == nil {
		panic(ErrorInvalidValidatorCreator)
	}
	mu.Lock()
	defer mu.Unlock()
	creators[kind] = creator
	return nil
}
