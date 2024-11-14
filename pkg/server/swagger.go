package server

// SwaggerGenerator 定义 Swagger 文档生成器接口
type SwaggerGenerator interface {
	// 解析 Handler 获取请求和响应信息
	ParseHandler(handler interface{}) (*RequestInfo, *ResponseInfo, error)

	// 生成 Swagger YAML 文档
	GenerateYAML(groups []RouteGroup) ([]byte, error)

	// 生成 Swagger JSON 文档
	GenerateJSON(groups []RouteGroup) ([]byte, error)
}

// SwaggerInfo 定义单个接口的 Swagger 文档信息
type SwaggerInfo struct {
	Tags        []string // API 标签分类
	Summary     string   // API 简要说明
	Description string   // API 详细说明
	Deprecated  bool     // 是否废弃
	Security    []string // 安全认证方式
}

// RequestInfo 定义请求参数信息
type RequestInfo struct {
	ContentType string      // 请求内容类型
	Params      []ParamInfo // 参数信息列表
}

// ResponseInfo 定义响应信息
type ResponseInfo struct {
	StatusCode  int         // HTTP 状态码
	ContentType string      // 响应内容类型
	Schema      interface{} // 响应结构体示例
}

// ParamInfo 定义参数详细信息
type ParamInfo struct {
	Name        string // 参数名称
	In          string // 参数位置(path/query/header/body)
	Required    bool   // 是否必须
	Type        string // 参数类型
	Description string // 参数说明
}
