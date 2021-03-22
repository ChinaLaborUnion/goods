package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"grpc-demo/views/order"
)

func RegisterOrderRouters(app *iris.Application){
	OrderRouter := app.Party("_order")

	OrderRouter.Post("", hero.Handler(order.CreateOrder))
	OrderRouter.Get("/{oid:int}", hero.Handler(order.GetOrder))
	OrderRouter.Put("/{oid:int}", hero.Handler(order.PutOrder))
	OrderRouter.Delete("/{ooid:int}/child_order/{oid:int}", hero.Handler(order.DeleteOrder)) //删除子订单
	OrderRouter.Delete("/{ooid:int}", hero.Handler(order.DeleteAllOrder)) //删除总订单下所有子订单
	OrderRouter.Get("/list", hero.Handler(order.ListOrder))
	OrderRouter.Post("/_mget", hero.Handler(order.MgetOrder))

	//批量获取订单明细
	OrderRouter.Post("/detail/_mget", hero.Handler(order.MgetOrderDetail))

	//取消订单
	OrderRouter.Post("/{oid:int}/cancel", hero.Handler(order.CancelOrder))

	//确认收货
	OrderRouter.Get("/{ooid:int}/child_order/{oid:int}/check_get", hero.Handler(order.CheckGet))

	//获取价钱
	OrderRouter.Post("/get_price",hero.Handler(order.GetPrice))

	//获取个数
	OrderRouter.Get("/list/sum",hero.Handler(order.ListSum))

	OrderRouter.Put("/status/{ooid:int}/child_order/{oid:int}", hero.Handler(order.PutOrderStatus)) //仅供测试，修改订单状态

}
