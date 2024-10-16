package router

import (
	v1 "gin-web/api/v1"
	"github.com/piupuer/go-helper/router"
)

func InitL0Router(r *router.Router) {
	//router1 := r.Casbin("/l0")
	router2 := r.CasbinAndIdempotence("/l0")
	//router1.POST("/list", v1.CreateL0)
	router2.POST("/create", v1.CreateL0)
	//router1.PATCH("/update/:id", v1.UpdateLeaveById)
	//router1.DELETE("/delete/batch", v1.BatchDeleteLeaveByIds)
}
