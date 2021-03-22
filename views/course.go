package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"grpc-demo/views/course"
)

func RegisterCourseRouters(app *iris.Application){

	courseRouter := app.Party("study/course")

	courseRouter.Post("", hero.Handler(course.Create))
	courseRouter.Get("/list", hero.Handler(course.List))
	courseRouter.Put("/{cid:int}", hero.Handler(course.Put))
	courseRouter.Delete("/{cid:int}", hero.Handler(course.Delete))
	courseRouter.Post("/_mget", hero.Handler(course.Mget))

	//applyRouter := app.Party("study/course/take")
	//applyRouter.Post("/{cid:int}", hero.Handler(course.ApplyCreate))
	//applyRouter.Post("/_mget", hero.Handler(course.ApplyMget))
	//applyRouter.Get("/list", hero.Handler(course.ApplyList))
	//applyRouter.Post("/cancel/{tid:int}", hero.Handler(course.ApplyCancel))


}
