package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"grpc-demo/views/goods_info"
)

func RegisterGoodsRouters(app *iris.Application){
	goodsRouter := app.Party("goods")

	goodsRouter.Post("", hero.Handler(goods_info.CreateGoods))
	goodsRouter.Get("/list", hero.Handler(goods_info.ListGoods))
	goodsRouter.Put("/{gid:int}", hero.Handler(goods_info.PutGoods))
	goodsRouter.Delete("/{gid:int}", hero.Handler(goods_info.DeleteGoods))
	goodsRouter.Post("/_mget", hero.Handler(goods_info.MgetGoods))
	//goodsRouter.Post("/little_mget", hero.Handler(goods_info.LittleMgetGoods))



	goodsSpecificationRouter := app.Party("goods/specification")

	goodsSpecificationRouter.Post("/{gid:int}", hero.Handler(goods_info.CreateGoodsSpecification))
	goodsSpecificationRouter.Put("/{gid:int}", hero.Handler(goods_info.PutGoodsSpecification))
	goodsSpecificationRouter.Delete("/{gid:int}", hero.Handler(goods_info.DeleteGoodsSpecification))
}