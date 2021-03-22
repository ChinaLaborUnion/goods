package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"grpc-demo/views/goods_home"
)

func RegisterGoodsSlideshowRouters(app *iris.Application){
	goodsSlideshowRouter := app.Party("goods/slideshow")

	goodsSlideshowRouter.Post("", hero.Handler(goods_home.CreateGoodsSlideshow))
	goodsSlideshowRouter.Get("/list", hero.Handler(goods_home.ListGoodsSlideshow))
	goodsSlideshowRouter.Put("/{sid:int}", hero.Handler(goods_home.PutGoodsSlideshow))
	goodsSlideshowRouter.Delete("/{sid:int}", hero.Handler(goods_home.DeleteGoodsSlideshow))
	goodsSlideshowRouter.Post("/_mget", hero.Handler(goods_home.MgetGoodsSlideshow))
}
