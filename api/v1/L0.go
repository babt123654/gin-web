package v1

import (
	"gin-web/pkg/request"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
	"github.com/piupuer/go-helper/pkg/tracing"
	"time"
)

// FindL0
// @Security Bearer
// @Accept json
// @Produce json
// @Success 201 {object} resp.Resp "success"
// @Tags L0
// @Description FindL0
// @Param params query request.L0 true "params"
// @Router /L0/list [GET]
//func FindL0(c *gin.Context) {
//	ctx := tracing.RealCtx(c)
//	_, span := tracer.Start(ctx, tracing.Name(tracing.Rest, "FindL0"))
//	defer span.End()
//	var r request.L0
//	req.ShouldBind(c, &r)
//	user := GetCurrentUser(c)
//	r.L0Id = user.Id
//	my := service.New(c)
//	list := my.FindL0(&r)
//	resp.SuccessWithPageData(list, &[]response.L0{}, r.Page)
//}

// CreateL0
// @Security Bearer
// @Accept json
// @Produce json
// @Success 201 {object} resp.Resp "success"
// @Tags L0
// @Description CreateL0
// @Param params body request.CreateL0 true "params"
// @Router /L0/create [POST]
func CreateL0(c *gin.Context) {
	ctx := tracing.RealCtx(c)
	_, span := tracer.Start(ctx, tracing.Name(tracing.Rest, "CreateL0"))
	defer span.End()
	//var r request.CreateL0
	//req.ShouldBind(c, &r)
	//req.Validate(c, r, r.FieldTrans())
	path := "F:\\L0\\L0"
	my := service.New(c)
	err := my.CreateL0(path)
	resp.CheckErr(err)
	resp.Success()
}

// UpdateL0ById
// @Security Bearer
// @Accept json
// @Produce json
// @Success 201 {object} resp.Resp "success"
// @Tags L0
// @Description UpdateL0ById
// @Param id path uint true "id"
// @Param params body request.UpdateL0 true "params"
// @Router /L0/update/{id} [PATCH]
func UpdateL0ById(c *gin.Context) {
	ctx := tracing.RealCtx(c)
	_, span := tracer.Start(ctx, tracing.Name(tracing.Rest, "UpdateL0ById"))
	defer span.End()
	var r request.UpdateL0
	req.ShouldBind(c, &r)
	id := req.UintId(c)
	my := service.New(c)
	user := GetCurrentUser(c)
	err := my.UpdateL0ById(id, r, user)
	resp.CheckErr(err)
	resp.Success()
}

// BatchDeleteL0ByIds
// @Security Bearer
// @Accept json
// @Produce json
// @Success 201 {object} resp.Resp "success"
// @Tags L0
// @Description BatchDeleteL0ByIds
// @Param ids body req.Ids true "ids"
// @Router /L0/delete/batch [DELETE]
func BatchDeleteL0ByIds(c *gin.Context) {
	ctx := tracing.RealCtx(c)
	_, span := tracer.Start(ctx, tracing.Name(tracing.Rest, "BatchDeleteL0ByIds"))
	defer span.End()
	var r req.Ids
	req.ShouldBind(c, &r)
	my := service.New(c)
	user := GetCurrentUser(c)
	err := my.DeleteL0ByIds(r.Uints(), user)
	resp.CheckErr(err)
	resp.Success()
}
func OminMonitor(c *gin.Context) {
	ctx := tracing.RealCtx(c)
	_, span := tracer.Start(ctx, tracing.Name(tracing.Rest, "BatchDeleteL0ByIds"))
	defer span.End()
	var r req.Ids
	req.ShouldBind(c, &r)
	q := service.New(c)
	ticker := time.NewTicker(5 * time.Minute)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				// 每隔5分钟调用一次q.GetOminiUsdtOp1(nil)函数
				q.GetOminiUsdtOp1(nil)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	// 假设你的程序需要运行一段时间后才退出
	// 在这里可以添加退出条件，例如等待某个事件发生或接收到某个信号
	// 例如，使用time.Sleep(time.Hour)来模拟程序的运行时间
	time.Sleep(time.Hour * 9999)

	// 关闭定时器，并退出循环
	close(quit)
	//resp.CheckErr(err)
	resp.Success()
}
