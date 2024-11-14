package reflect

import (
	"context"
	"fmt"
	"reflect"
)

// TypeChecker 定义类型检查函数
type TypeChecker func(reflect.Type) bool

// SignatureSpec 定义函数签名规范
type SignatureSpec struct {
	ExpectedNumIn   int           // 预期输入参数个数
	ExpectedNumOut  int           // 预期输出参数个数
	InTypeCheckers  []TypeChecker // 输入参数类型检查函数
	OutTypeCheckers []TypeChecker // 输出参数类型检查函数
}

// NewSignatureSpec 创建一个新的签名规范
func NewSignatureSpec(numIn, numOut int, inCheckers, outCheckers []TypeChecker) *SignatureSpec {
	return &SignatureSpec{
		ExpectedNumIn:   numIn,
		ExpectedNumOut:  numOut,
		InTypeCheckers:  inCheckers,
		OutTypeCheckers: outCheckers,
	}
}

// Validate 验证函数是否符合签名规范
func (s *SignatureSpec) Validate(f interface{}) error {
	if f == nil {
		return fmt.Errorf("function cannot be nil")
	}

	fType := reflect.TypeOf(f)
	if fType.Kind() != reflect.Func {
		return fmt.Errorf("expected a function, got %T", f)
	}

	if err := s.validateParameterCount(fType); err != nil {
		return err
	}

	if err := s.validateInputTypes(fType); err != nil {
		return err
	}

	if err := s.validateOutputTypes(fType); err != nil {
		return err
	}

	return nil
}

// validateParameterCount 验证参数数量
func (s *SignatureSpec) validateParameterCount(fType reflect.Type) error {
	if s.ExpectedNumIn > 0 && fType.NumIn() != s.ExpectedNumIn {
		return fmt.Errorf("expected %d input parameters, got %d", s.ExpectedNumIn, fType.NumIn())
	}

	if s.ExpectedNumOut > 0 && fType.NumOut() != s.ExpectedNumOut {
		return fmt.Errorf("expected %d output parameters, got %d", s.ExpectedNumOut, fType.NumOut())
	}

	return nil
}

// validateInputTypes 验证输入参数类型
func (s *SignatureSpec) validateInputTypes(fType reflect.Type) error {
	for i := 0; i < fType.NumIn(); i++ {
		if !s.checkType(fType.In(i), s.InTypeCheckers[i]) {
			return fmt.Errorf("input parameter %d (type %v) does not match expected type", i, fType.In(i))
		}
	}
	return nil
}

// validateOutputTypes 验证输出参数类型
func (s *SignatureSpec) validateOutputTypes(fType reflect.Type) error {
	for i := 0; i < fType.NumOut(); i++ {
		if !s.checkType(fType.Out(i), s.OutTypeCheckers[i]) {
			return fmt.Errorf("output parameter %d (type %v) does not match expected type", i, fType.Out(i))
		}
	}
	return nil
}

// checkType 检查类型是否符合检查器的要求
func (s *SignatureSpec) checkType(t reflect.Type, checker TypeChecker) bool {
	if checker(t) {
		return true
	}
	return false
}

// Common Type Checkers
var (
	IsString    = func(t reflect.Type) bool { return t.Kind() == reflect.String }
	IsInt       = func(t reflect.Type) bool { return t.Kind() == reflect.Int }
	IsError     = func(t reflect.Type) bool { return t.Implements(reflect.TypeOf((*error)(nil)).Elem()) }
	IsInterface = func(t reflect.Type) bool { return t.Kind() == reflect.Interface }
	AsInterface = func(t reflect.Type) bool { return t.AssignableTo(reflect.TypeOf((*interface{})(nil)).Elem()) }
	IsContext   = func(t reflect.Type) bool { return t.Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) }
	IsStructPtr = func(t reflect.Type) bool { return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct }
)
