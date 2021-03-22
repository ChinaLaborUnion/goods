package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	courseApply "grpc-demo/views/course_apply"
)

func RegisterCourseApplyRouters(app *iris.Application){

	preApplyRouter := app.Party("study/course")

	preApplyRouter.Post("/{cid:int}/pre_apply", hero.Handler(courseApply.PreApplyCreate))
	preApplyRouter.Post("/pre_apply/_mget", hero.Handler(courseApply.PreApplyMget))
	preApplyRouter.Get("/pre_apply/list", hero.Handler(courseApply.PreApplyList))
	//修改
	preApplyRouter.Put("/pre_apply/{pid:int}", hero.Handler(courseApply.PreApplyPut))
	//撤销
	preApplyRouter.Post("/pre_apply/{pid:int}/cancel", hero.Handler(courseApply.PreApplyCancel))

	applyRouter := app.Party("study/course")

	applyRouter.Post("/pre_apply/{pid:int}/apply", hero.Handler(courseApply.ApplyCreate))
	applyRouter.Post("/apply/_mget", hero.Handler(courseApply.ApplyMget))
	applyRouter.Get("/apply/list", hero.Handler(courseApply.ApplyList))
	//修改
	applyRouter.Put("/apply/{aid:int}", hero.Handler(courseApply.ApplyPut))
	//撤销
	applyRouter.Post("/apply/{aid:int}/cancel", hero.Handler(courseApply.ApplyCancel))

	//todo 直接报名
	applyRouter.Post("/{cid:int}/apply", hero.Handler(courseApply.ApplyDirectCreate))
	//todo excel导入

	applyRouter.Get("/apply/all/list", hero.Handler(courseApply.ApplyAllList))

	//退款
	applyRouter.Post("/apply/{aid:int}/back_money", hero.Handler(courseApply.BackMoney))

}
