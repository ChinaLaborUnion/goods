package order

import (
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	orderEnums "grpc-demo/enums/order"
	"grpc-demo/models/db"
)

func ListSum(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckLogin()

	//var submitOrders []db.TestChildOrder
	//db.Driver.Where("account_id = ? and order_status = ?",auth.AccountModel().Id,orderEnums.SUBMIT).Find(&submitOrders)
	var count int
	table := db.Driver.Table("test_child_order").Where("account_id = ? and order_status = ?",auth.AccountModel().Id,orderEnums.SUBMIT)
	table.Count(&count).Group("order_id")

	var wfShippedOrders []db.TestChildOrder
	db.Driver.Where("account_id = ? and order_status = ?",auth.AccountModel().Id,orderEnums.WfShipped).Find(&wfShippedOrders)

	var tbShippedOrders []db.TestChildOrder
	db.Driver.Where("account_id = ? and order_status in (?)",auth.AccountModel().Id,[]int16{orderEnums.TbShipped,orderEnums.ToPick}).Find(&tbShippedOrders)

	var receivedDetails []db.TestOrderDetail
	db.Driver.Where("account_id = ? and is_comment = 0",auth.AccountModel().Id).Find(&receivedDetails)


	var afterSalesDetails []db.TestOrderDetail
	db.Driver.Where("account_id = ? and is_after_serve != 0 and is_pass = 1",auth.AccountModel().Id).Find(&afterSalesDetails)

	ctx.JSON(iris.Map{
		"submit":count,
		"wf_shipped":len(wfShippedOrders),
		"tb_shipped":len(tbShippedOrders),
		"received":len(receivedDetails),
		"after_sales":len(afterSalesDetails),
	})
}
