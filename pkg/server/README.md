## 整洁架构服务器

### 参数绑定
框架支持参数自动绑定到用户的控制器的请求结构体上，但对函数签名有一定限制：
```go
// 第一个参数必须是 ctx
// 第二个参数必须是结构体指针
// 第一个返回值可以是任意类型
// 第二个返回值必须是 error
fun(ctx context.Context, req *Struct) (interface{}, err)

// 第一个参数必须是 ctx
// 第一个返回值可以是任意类型
// 第二个返回值必须是 error
fun(ctx context.Context) (interface{}, err)
```

其中 HTTP 请求参数到结构体的绑定规则为：
1. Body 类绑定由 ContentType 决定，同时可以使用 Tag 修改绑定的参数名，比如:
```go
// 同时支持 ContentType: application/json, applicatnion/yaml, application/xml, application/toml
type UserReq struct{
    ID      int    `json:"id" yaml:"id" xml:"id" toml:"id"`
	Name    string `json:"name" yaml:"name" xml:"name" toml:"name"`
}
```
>**注**：Get 请求不支持从 Body 中解析数据

2. 表单类绑定在 Get 请求的场景下仅从查询参数进行绑定，其它类请求不仅支持从查询参数绑定还会对 Body 进行 `application/x-www-form-urlencoded` 式的解析并绑定，也支持通过 Tag 修改绑定的参数名：
```go
// Get场景下，会从查询参数绑定 id 和 name
// 其它场景也会从查询参数绑定，同时也会从 `application/x-www-form-urlencoded` 类型的 Body 中进行解析和绑定
type UserReq struct{
    ID      int    `form:"id"`
	Name    string `form:"name"`
}
```
3. 路径参数绑定，需要使用 Tag: uri 标识，它会根据动态路由参数进行绑定：
```go

```
4. 请求头参数绑定，需要使用 Tag: header 标识，它会根据请求头进行绑定：
```go

```
5. Context绑定，需要使用 Tag: context 标识，它会根据动态路由参数进行绑定：
```go

```
6. Cookie 绑定，需要使用 Tag: cookie 标识，它会根据动态路由参数进行绑定：
```go

```



如何支持 path 通配符绑定？

支持 context、cookie 的绑定

支持 websocket handler

支持 Controller 和 RestController
视图控制器支持 redirect:，forward:

文件上传下载支持

支持 Web 配置注册，比如替换JSON，validate实现，拦截器，验证器添加 cors 支持，等

异步任务支持(控制器返回 callable，封装 ctx 的复制，协程池)

支持  opentelemetry 集成