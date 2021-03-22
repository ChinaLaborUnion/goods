package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	courseTag "grpc-demo/views/course_tag"
)

func RegisterCourseTagRouters(app *iris.Application){

	courseTagRouter := app.Party("study/course/tag")

	courseTagRouter.Post("", hero.Handler(courseTag.CreateTag))
	courseTagRouter.Get("/list", hero.Handler(courseTag.ListTag))
	courseTagRouter.Put("/{tid:int}", hero.Handler(courseTag.PutTag))
	courseTagRouter.Delete("/{tid:int}", hero.Handler(courseTag.DeleteTag))
	courseTagRouter.Post("/_mget", hero.Handler(courseTag.MgetTag))

	courseKindRouter := app.Party("study/course/kind")
	courseKindRouter.Post("", hero.Handler(courseTag.CreateKind))
	courseKindRouter.Get("/list", hero.Handler(courseTag.ListKind))
	courseKindRouter.Put("/{kid:int}", hero.Handler(courseTag.PutKind))
	courseKindRouter.Delete("/{kid:int}", hero.Handler(courseTag.DeleteKind))
	courseKindRouter.Post("/_mget", hero.Handler(courseTag.MgetKind))

}
