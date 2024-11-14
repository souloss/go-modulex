package gin

import (
	"github.com/gin-gonic/gin"
)

// CtrlToGinHandlerConverter gin Handler转换器
type CtrlToGinHandlerConverter interface {

	// Convert 对 Controller 进行转换
	//  @param controller 用户控制器
	//  @return gin.HandlerFunc gin处理器
	Convert(controller any) gin.HandlerFunc

	// Check 检查 controller 是否符合转换条件
	//  @param controller
	//  @return bool
	Check(controller any) bool
}

// ginHandlerAdapt gin Handler 适配器
// 一个适配器可以注册多个转换器，它会对每个非 gin.HandlerFunc 的控制器进行 Check，符合条件则转换成 gin.HandlerFunc
type ginHandlerAdapt struct {
	converter []CtrlToGinHandlerConverter
}

func NewHandlerAdapt() *ginHandlerAdapt {
	return &ginHandlerAdapt{
		converter: []CtrlToGinHandlerConverter{},
	}
}

func (g *ginHandlerAdapt) RegisterConverter(c CtrlToGinHandlerConverter) {
	g.converter = append(g.converter, c)
}
