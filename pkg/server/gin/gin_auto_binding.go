package gin

import (
	"reflect"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// autoBindingCache 缓存绑定信息，避免重复解析结构体标签
var autoBindingCache sync.Map

// structBindings 缓存结构体的绑定类型
type structBindings struct {
	Bindings    []binding.Binding
	CtxBindings []CtxBinding
	uri         bool
	context     bool
}

// getBindingsFromStruct 获取结构体的绑定类型
//
//	@param t
//	@return *structBindings
//	@return error
func getBindingsFromStruct(t reflect.Type) (*structBindings, error) {
	if v, ok := autoBindingCache.Load(t); ok {
		return v.(*structBindings), nil
	}

	bindings := &structBindings{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		switch {
		// 请求体绑定根据 ContentType 进行动态绑定
		// 查询参数绑定
		case field.Tag.Get("form") != "":
			bindings.Bindings = append(bindings.Bindings, binding.Query)
		// 路径绑定
		case field.Tag.Get("uri") != "":
			bindings.uri = true
		// 请求头绑定
		case field.Tag.Get("header") != "":
			bindings.Bindings = append(bindings.Bindings, binding.Header)
		case field.Tag.Get("cookie") != "":
			bindings.CtxBindings = append(bindings.CtxBindings, CookieBinding)
		case field.Tag.Get("context") != "":
			bindings.CtxBindings = append(bindings.CtxBindings, ContextBinding)
		}
	}
	autoBindingCache.Store(t, bindings)
	return bindings, nil
}

// AutoBind 自动绑定结构体
//
//	@param ctx
//	@param obj
//	@return error
func AutoBind(ctx *gin.Context, obj interface{}) error {
	t := reflect.TypeOf(obj).Elem()
	bindings, err := getBindingsFromStruct(t)
	if err != nil {
		return err
	}
	// 进行请求头，查询参数绑定
	for _, binding := range bindings.Bindings {
		if err := binding.Bind(ctx.Request, obj); err != nil {
			return err
		}
	}
	// 进行路径绑定
	if bindings.uri {
		if err := ctx.ShouldBindUri(obj); err != nil {
			return err
		}
	}
	// 进行Cookie，Context绑定
	for _, binding := range bindings.CtxBindings {
		if err := binding.Bind(ctx, obj); err != nil {
			return err
		}
	}
	// 根据 ContentType 进行动态绑定
	if err := ctx.ShouldBind(obj); err != nil {
		return err
	}

	return nil
}
