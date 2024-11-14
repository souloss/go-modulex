package gin

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/souloss/go-clean-arch/pkg/errors"
	"github.com/souloss/go-clean-arch/pkg/logger"
	pkgreflect "github.com/souloss/go-clean-arch/pkg/reflect"
)

// response 默认响应结构体
type response struct {
	Code    errors.ErrorCode `json:"code"`
	Message string           `json:"msg"`
	Data    any              `json:"data"`
}

// DefaultBindding 默认的绑定优化器
var DefaultBinding = func(c *gin.Context, obj any) error {
	return AutoBind(c, obj)
}

// DefaultResponseWrap 默认的统一的响应处理
var DefaultJSONResponseWrap = func(c *gin.Context, data any, err error) {
	if err != nil {
		logger.L().Error("API Error, err: %v", err)
		e, ok := err.(*errors.Error)
		if !ok {
			e = errors.Wrap(err)
		}
		c.JSON(e.GetCode().ToHTTPStatusCode(), response{
			Code:    e.GetCode(),
			Message: e.ToMessage(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, response{
		Code:    errors.Success,
		Message: errors.Success.String(),
		Data:    data,
	})
}

// JSONWithDataConvert
// HandlerWithData[T any] func(ctx context.Context) (T, error)
type JSONWithDataConvert struct{}

func (j *JSONWithDataConvert) Convert(controller any) gin.HandlerFunc {
	if !j.Check(controller) {
		panic("controller convert failed")
	}
	return func(ctx *gin.Context) {
		values := reflect.ValueOf(controller).Call([]reflect.Value{
			reflect.ValueOf(ctx.Request.Context()),
		})
		data := values[0].Interface()
		err, _ := values[1].Interface().(error)
		DefaultJSONResponseWrap(ctx, data, err)
	}
}

func (j *JSONWithDataConvert) Check(controller any) bool {
	signatureSpec := pkgreflect.NewSignatureSpec(1, 2,
		[]pkgreflect.TypeChecker{
			pkgreflect.IsContext,
		},
		[]pkgreflect.TypeChecker{
			pkgreflect.AsInterface,
			pkgreflect.IsError,
		})
	if err := signatureSpec.Validate(controller); err != nil {
		return false
	}
	return true
}

// JSONWithReqAndDataConvert
// HandlerWithReqAndData[Req any, Resp any] func(ctx context.Context, req *Req) (Resp, error)
type JSONWithReqAndDataConvert struct{}

func (j *JSONWithReqAndDataConvert) Convert(controller any) gin.HandlerFunc {
	if !j.Check(controller) {
		panic("controller convert failed")
	}
	controllerType := reflect.TypeOf(controller)
	return func(ctx *gin.Context) {
		reqType := controllerType.In(1).Elem()
		req := reflect.New(reqType).Interface()
		if err := DefaultBinding(ctx, req); err != nil {
			logger.L().Warnf("params auto bind failed, err: %v", err)
			DefaultJSONResponseWrap(ctx, nil, errors.New(errors.ErrInvalidInput))
			return
		}
		values := reflect.ValueOf(controller).Call([]reflect.Value{
			reflect.ValueOf(ctx.Request.Context()),
			reflect.ValueOf(req),
		})
		data := values[0].Interface()
		err, _ := values[1].Interface().(error)
		DefaultJSONResponseWrap(ctx, data, err)
	}
}

func (j *JSONWithReqAndDataConvert) Check(controller any) bool {
	signatureSpec := pkgreflect.NewSignatureSpec(2, 2,
		[]pkgreflect.TypeChecker{
			pkgreflect.IsContext,
			pkgreflect.IsStructPtr,
		},
		[]pkgreflect.TypeChecker{
			pkgreflect.AsInterface,
			pkgreflect.IsError,
		})
	if err := signatureSpec.Validate(controller); err != nil {
		return false
	}
	return true
}

// JSONWithReqConvert
// HandlerWithReq[Req any, Resp any] func(ctx context.Context, req *Req) (error)
type JSONWithReqConvert struct{}

func (j *JSONWithReqConvert) Convert(controller any) gin.HandlerFunc {
	if !j.Check(controller) {
		panic("controller convert failed")
	}
	controllerType := reflect.TypeOf(controller)
	return func(ctx *gin.Context) {
		reqType := controllerType.In(1).Elem()
		req := reflect.New(reqType).Interface()
		if err := DefaultBinding(ctx, req); err != nil {
			logger.L().Warnf("params auto bind failed, err: %v", err)
			DefaultJSONResponseWrap(ctx, nil, errors.New(errors.ErrInvalidInput))
			return
		}
		values := reflect.ValueOf(controller).Call([]reflect.Value{
			reflect.ValueOf(ctx.Request.Context()),
			reflect.ValueOf(req),
		})
		err, _ := values[0].Interface().(error)
		DefaultJSONResponseWrap(ctx, nil, err)
	}
}

func (j *JSONWithReqConvert) Check(controller any) bool {
	signatureSpec := pkgreflect.NewSignatureSpec(2, 1,
		[]pkgreflect.TypeChecker{
			pkgreflect.IsContext,
			pkgreflect.IsStructPtr,
		},
		[]pkgreflect.TypeChecker{
			pkgreflect.IsError,
		})
	if err := signatureSpec.Validate(controller); err != nil {
		return false
	}
	return true
}
