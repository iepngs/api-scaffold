package validator

import "github.com/go-playground/validator/v10"

// Validator 封装验证器
type Validator struct {
	v *validator.Validate
}

// NewValidator 初始化验证器
func NewValidator() *Validator {
	return &Validator{
		v: validator.New(),
	}
}

// Validate 验证结构体
func (v *Validator) Validate(i interface{}) error {
	return v.v.Struct(i)
}
