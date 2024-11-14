package gin

import (
	"errors"
	"reflect"

	"github.com/gin-gonic/gin"
)

// 绑定接口
type CtxBinding interface {
	Name() string
	Bind(ctx *gin.Context, obj interface{}) error
}

// contextBinding 实现
type contextBinding struct{}

// cookieBinding 实现
type cookieBinding struct{}

var (
	ContextBinding CtxBinding = contextBinding{}
	CookieBinding  CtxBinding = cookieBinding{}
)

// 绑定字段
func bindFields(ctx *gin.Context, obj interface{}, tag string) error {
	v := reflect.ValueOf(obj).Elem()
	t := v.Type()

	var fields []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get(tag) != "" {
			fields = append(fields, field.Name)
			if err := setFieldValue(ctx, field, v.Field(i), tag); err != nil {
				return err
			}
		}

		// 递归处理嵌套结构体
		if field.Type.Kind() == reflect.Struct {
			if err := bindFields(ctx, v.Field(i).Addr().Interface(), tag); err != nil {
				return err
			}
		}
	}

	return nil
}

// 设置字段值
func setFieldValue(ctx *gin.Context, field reflect.StructField, fieldValue reflect.Value, tag string) error {
	var value interface{}

	if tag == ContextBinding.Name() {
		value = ctx.Value(field.Tag.Get(tag))
	} else if tag == CookieBinding.Name() {
		cookieValue, err := ctx.Cookie(field.Tag.Get(tag))
		if err == nil {
			value = cookieValue
		} else {
			return err
		}
	}

	if value == nil {
		return nil // 如果没有值，跳过
	}

	// 判断类型是否匹配
	if fieldValue.Kind() != reflect.TypeOf(value).Kind() {
		return errors.New("type mismatch for binding field: " + field.Name)
	}

	fieldValue.Set(reflect.ValueOf(value))
	return nil
}

// 绑定函数的实现
func (c contextBinding) Name() string {
	return "context"
}
func (c contextBinding) Bind(ctx *gin.Context, obj interface{}) error {
	return bindFields(ctx, obj, c.Name())
}

func (c cookieBinding) Name() string {
	return "cookie"
}
func (c cookieBinding) Bind(ctx *gin.Context, obj interface{}) error {
	return bindFields(ctx, obj, c.Name())
}
