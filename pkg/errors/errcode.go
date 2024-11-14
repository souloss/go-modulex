package errors

import (
	"net/http"
)

//go:generate stringer -type=ErrorCode -linecomment
type ErrorCode int64

func (i ErrorCode) ToHTTPStatusCode() int {
	switch {
	case 1000000 <= i && i < 2000000:
		return http.StatusInternalServerError
	case 2000000 <= i && i <= 4000000:
		return http.StatusBadRequest
	default:
		return http.StatusBadRequest
	}
}

const Success ErrorCode = 0 // 操作成功

const ErrUndefined ErrorCode = 1001000 // {{ .error }}

const (
	ErrSystemInternal ErrorCode = 1001001 // 系统内部错误，例如服务崩溃、未知异常

	ErrSystemDatabase ErrorCode = 1001002 // 数据库操作失败，如连接失败、查询异常等

	ErrSystemCache ErrorCode = 1001003 // 缓存操作失败，比如 Redis 连接失败

	ErrSystemTimeout ErrorCode = 1001004 // 系统请求超时，比如网络延迟或依赖服务未响应

	ErrSystemDependency ErrorCode = 1001005 // 依赖服务错误，通常指外部服务、微服务依赖等
)

const (
	ErrInvalidInput ErrorCode = 2002001 // 输入参数错误，常用于接口参数验证失败

	ErrUnauthorized ErrorCode = 2002002 // 未授权，通常表示用户未登录或认证失败

	ErrForbidden ErrorCode = 2002003 // 权限不足，用户登录了但没有执行此操作的权限

	ErrResourceNotFound ErrorCode = 2002004 // 资源未找到，如查询数据库时 ID 不存在的情况

	ErrResourceAlreadyExists ErrorCode = 2002005 // 资源已存在，通常用于创建已存在的唯一资源

	ErrBusinessLogicFailure ErrorCode = 2002006 // 一般的业务逻辑处理失败，比如订单处理异常
)
