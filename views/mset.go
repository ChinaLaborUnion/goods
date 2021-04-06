package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"grpc-demo/views/mset"
)

func RegisterMsetRouters(app *iris.Application) {
	//icon
	msetRouter := app.Party("mset")

	msetRouter.Post("", hero.Handler(mset.Mset))
}
