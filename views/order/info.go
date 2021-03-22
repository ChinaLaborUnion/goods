package order

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris"
	"grpc-demo/constants"
	authbase "grpc-demo/core/auth"
	"grpc-demo/core/cache"
	specificationEntity "grpc-demo/entity/goods_specification"
	goodsInfoEnum "grpc-demo/enums/goods_info"
	orderEnums "grpc-demo/enums/order"
	accountException "grpc-demo/exceptions/account"
	goodsInfoException "grpc-demo/exceptions/goods_info"
	orderException "grpc-demo/exceptions/order"
	"grpc-demo/models/db"
	"grpc-demo/utils/hash"
	logUtils "grpc-demo/utils/log"
	paramsUtils "grpc-demo/utils/params"
	"strconv"
	"time"
)

var orderDetailField = []string{
	"ID", "ChildOrderID", "GoodsID", "GoodsSpecificationID", "PurchaseQty", "Message", "OrderID",
	"Coupon", "ExpFare", "GoodsAmount", "CreateTime", "UpdateTime","IsComment","IsAfterServe","OrderAmount","IsPass",
}

//创建订单
func CreateOrder(ctx iris.Context, auth authbase.AuthAuthorization) {
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	var addressId int
	//addressId := params.Int("address_id", "地址id")
	goodsList := params.List("goods_list", "所选商品")
	//[{"goods_id":,"goods_specification":,"goods_total":,"message":,"delivery":},{}]
	//TODO 输入格式校验

	//判断地址id是否属于本人
	//var address db.Address
	//if err := db.Driver.Where("account_id = ? and id = ?", auth.AccountModel().Id, addressId).First(&address).Error; err != nil {
	//	panic(orderException.AddressIsNotExsit())
	//}

	expressage := make([]map[string]interface{}, 0)
	sameCity := make([]map[string]interface{}, 0)
	self := make([]map[string]interface{}, 0)

	var totalCoupon, totalExpFare, totalGoodsAmount, totalOrderAmount int
	var childCoupon1, childGoodsAmount1, childExpFare1, childOrderAmount1 int
	var childCoupon2, childGoodsAmount2, childExpFare2, childOrderAmount2 int
	var childCoupon3, childGoodsAmount3, childExpFare3, childOrderAmount3 int

	tx := db.Driver.Begin()

	expFare1 := make([]int, 0)
	expFare2 := make([]int, 0)

	for _, goods := range goodsList {
		gd := goods.(map[string]interface{})
		//判断商品是否存在
		goodsId := gd["goods_id"]
		var g db.Goods
		if err := db.Driver.GetOne("goods", int(goodsId.(float64)), &g); err != nil {
			panic(goodsInfoException.GoodsIsNotExsit())
		}

		//判断商品规格是否存在
		specification := g.GetInfo()["specification"]
		var ok bool
		ok = false

		goodsTotal := int(gd["goods_total"].(float64))
		delivery := int16(gd["delivery"].(float64))
		//TODO 判断此商品是否存在此发货方式

		if g.Sale {
			for index, s := range specification.(specificationEntity.GoodsSpecificationEntity) {
				if s.ID == int(gd["goods_specification"].(float64)) {
					ok = true

					//卡库存
					if s.Total < goodsTotal {
						panic(orderException.TotalNoEnough())
					}

					//减少库存
					var goodsSpecification db.GoodsSpecification
					db.Driver.Where("goods_id = ?", goodsId).First(&goodsSpecification)

					s.Total -= goodsTotal

					specification.(specificationEntity.GoodsSpecificationEntity)[index].Total = s.Total

					if data, err := json.Marshal(specification.(specificationEntity.GoodsSpecificationEntity)); err != nil {
						panic(goodsInfoException.SpecificationMarshalFail())
					} else {
						goodsSpecification.Specification = string(data)

					}
					if err := tx.Save(&goodsSpecification).Error; err != nil {
						tx.Rollback()
						panic(orderException.TotalPutFail())
					}

					//同时减少首页库存显示
					g.Total -= goodsTotal
					if err := tx.Save(&g).Error; err != nil {
						tx.Rollback()
						panic(orderException.TotalPutFail())
					}

					totalCoupon += (s.Price - s.ReducedPrice) * goodsTotal
					totalGoodsAmount += s.Price * goodsTotal

					//订单分单
					if delivery == goodsInfoEnum.GetWayExpressage {

						addressId = params.Int("address_id", "地址id")
						var address db.Address
						if err := db.Driver.Where("account_id = ? and id = ?", auth.AccountModel().Id, addressId).First(&address).Error; err != nil {
							panic(orderException.AddressIsNotExsit())
						}

						//子数值计算
						childCoupon1 += (s.Price - s.ReducedPrice) * goodsTotal
						gd["coupon"] = (s.Price - s.ReducedPrice) * goodsTotal

						childGoodsAmount1 += s.Price * goodsTotal

						gd["exp_fare"] = g.Carriage
						gd["goods_amount"] = s.Price * goodsTotal

						expFare1 = append(expFare1, g.Carriage)

						expressage = append(expressage, gd)
					} else if delivery == goodsInfoEnum.GetWaySameCity {

						addressId = params.Int("address_id", "地址id")
						var address db.Address
						if err := db.Driver.Where("account_id = ? and id = ?", auth.AccountModel().Id, addressId).First(&address).Error; err != nil {
							panic(orderException.AddressIsNotExsit())
						}

						//子数值计算
						childCoupon2 += (s.Price - s.ReducedPrice) * goodsTotal
						childGoodsAmount2 += s.Price * goodsTotal

						gd["coupon"] = (s.Price - s.ReducedPrice) * goodsTotal
						gd["exp_fare"] = g.Carriage
						gd["goods_amount"] = s.Price * goodsTotal

						expFare2 = append(expFare2, g.Carriage)

						sameCity = append(sameCity, gd)
					} else {

						//子数值计算
						childCoupon3 += (s.Price - s.ReducedPrice) * goodsTotal
						childGoodsAmount3 += s.Price * goodsTotal

						gd["coupon"] = (s.Price - s.ReducedPrice) * goodsTotal
						gd["exp_fare"] = g.Carriage
						gd["goods_amount"] = s.Price * goodsTotal

						self = append(self, gd)
					}
					break
				}
			}
		} else {
			for index, s := range specification.(specificationEntity.GoodsSpecificationEntity1) {
				if s.ID == int(gd["goods_specification"].(float64)) {
					ok = true

					//卡库存
					if s.Total < goodsTotal {
						panic(orderException.TotalNoEnough())
					}

					//减少库存
					var goodsSpecification db.GoodsSpecification
					db.Driver.Where("goods_id = ?", goodsId).First(&goodsSpecification)

					s.Total -= goodsTotal

					specification.(specificationEntity.GoodsSpecificationEntity1)[index].Total = s.Total

					if data, err := json.Marshal(specification.(specificationEntity.GoodsSpecificationEntity1)); err != nil {
						panic(goodsInfoException.SpecificationMarshalFail())
					} else {
						goodsSpecification.Specification = string(data)
					}
					if err := tx.Save(&goodsSpecification).Error; err != nil {
						tx.Rollback()
						panic(orderException.TotalPutFail())
					}

					//同时减少首页库存显示
					g.Total -= goodsTotal
					if err := tx.Save(&g).Error; err != nil {
						tx.Rollback()
						panic(orderException.TotalPutFail())
					}

					totalGoodsAmount += s.Price * goodsTotal

					//订单分单
					if delivery == goodsInfoEnum.GetWayExpressage {
						addressId = params.Int("address_id", "地址id")
						var address db.Address
						if err := db.Driver.Where("account_id = ? and id = ?", auth.AccountModel().Id, addressId).First(&address).Error; err != nil {
							panic(orderException.AddressIsNotExsit())
						}

						//子数值计算
						childGoodsAmount1 += s.Price * goodsTotal

						gd["exp_fare"] = g.Carriage
						gd["goods_amount"] = s.Price * goodsTotal

						expFare1 = append(expFare1, g.Carriage)

						expressage = append(expressage, gd)
					} else if delivery == goodsInfoEnum.GetWaySameCity {
						addressId = params.Int("address_id", "地址id")
						var address db.Address
						if err := db.Driver.Where("account_id = ? and id = ?", auth.AccountModel().Id, addressId).First(&address).Error; err != nil {
							panic(orderException.AddressIsNotExsit())
						}

						//子数值计算
						childGoodsAmount2 += s.Price * goodsTotal

						gd["exp_fare"] = g.Carriage
						gd["goods_amount"] = s.Price * goodsTotal

						expFare2 = append(expFare2, g.Carriage)

						sameCity = append(sameCity, gd)
					} else {

						//子数值计算
						childGoodsAmount3 += s.Price * goodsTotal

						gd["exp_fare"] = g.Carriage
						gd["goods_amount"] = s.Price * goodsTotal

						self = append(self, gd)
					}
					break
				}
			}
		}

		if !ok {
			panic(goodsInfoException.SpecificationIsNotExsit())
		}
	}

	if len(expFare1) > 0 {
		childExpFare1 = expFare1[0]
		for _, e1 := range expFare1 {
			if e1 > childExpFare1 {
				childExpFare1 = e1
			}
		}
	}

	if len(expFare2) > 0 {
		childExpFare2 = expFare2[0]
		for _, e2 := range expFare2 {
			if e2 > childExpFare2 {
				childExpFare2 = e2
			}
		}
	}

	childOrderAmount1 += childGoodsAmount1 + childExpFare1 - childCoupon1
	childOrderAmount2 += childGoodsAmount2 + childExpFare2 - childCoupon2
	childOrderAmount3 += childGoodsAmount3 + childExpFare3 - childCoupon3

	totalExpFare += childExpFare1 + childExpFare2 + childExpFare3
	totalOrderAmount += totalGoodsAmount + totalExpFare - totalCoupon

	//判断是否拆单
	flag := 0
	if len(expressage) > 0 {
		flag++
	}
	if len(self) > 0 {
		flag++
	}
	if len(sameCity) > 0 {
		flag++
	}

	//创建总订单
	order := db.TestOrder{
		PayOrNot:         false,
		AccountID:        auth.AccountModel().Id,
		AddressID:        addressId,
		TotalCoupon:      totalCoupon,
		TotalExpFare:     totalExpFare,
		TotalGoodsAmount: totalGoodsAmount,
		TotalOrderAmount: totalOrderAmount,
	}
	if err := tx.Debug().Create(&order).Error; err != nil {
		tx.Rollback()
		panic(orderException.OrderCreateFail())
	}

	order.OrderNum = strconv.FormatInt(order.CreateTime, 10) + "-" + hash.GetRandomString(8)
	if flag > 1 {
		order.Status = orderEnums.BreakStatusBreak
	} else {
		order.Status = orderEnums.BreakStatusNo
	}
	if err := tx.Debug().Save(&order).Error; err != nil {
		tx.Rollback()
		panic(orderException.OrderCreateFail())
	}

	//创建子订单
	var childOrder1 db.TestChildOrder
	if len(expressage) > 0 {
		childOrder1 = db.TestChildOrder{
			AccountID:        auth.AccountModel().Id,
			AddressID:        addressId,
			OrderStatus:      orderEnums.SUBMIT,
			OrderNum:         strconv.FormatInt(order.CreateTime, 10) + "-" + hash.GetRandomString(8),
			OrderID:          order.ID,
			Delivery:         goodsInfoEnum.GetWayExpressage,
			ChildTotalCoupon: childCoupon1,
			ChildGoodsAmount: childGoodsAmount1,
			ChildExpFare:     childExpFare1,
			ChildOrderAmount: childOrderAmount1,
		}
		if err := tx.Debug().Create(&childOrder1).Error; err != nil {
			tx.Rollback()
			panic(orderException.OrderCreateFail())
		}
	}

	var childOrder2 db.TestChildOrder
	if len(sameCity) > 0 {
		childOrder2 = db.TestChildOrder{
			AccountID:        auth.AccountModel().Id,
			AddressID:        addressId,
			OrderStatus:      orderEnums.SUBMIT,
			OrderNum:         strconv.FormatInt(order.CreateTime, 10) + "-" + hash.GetRandomString(8),
			OrderID:          order.ID,
			Delivery:         goodsInfoEnum.GetWaySameCity,
			ChildTotalCoupon: childCoupon2,
			ChildGoodsAmount: childGoodsAmount2,
			ChildExpFare:     childExpFare2,
			ChildOrderAmount: childOrderAmount2,
		}
		if err := tx.Debug().Create(&childOrder2).Error; err != nil {
			tx.Rollback()
			panic(orderException.OrderCreateFail())
		}
	}

	var childOrder3 db.TestChildOrder
	if len(self) > 0 {
		childOrder3 = db.TestChildOrder{
			AccountID:        auth.AccountModel().Id,
			AddressID:        addressId,
			OrderStatus:      orderEnums.SUBMIT,
			OrderNum:         strconv.FormatInt(order.CreateTime, 10) + "-" + hash.GetRandomString(8),
			OrderID:          order.ID,
			Delivery:         goodsInfoEnum.GetWaySelf,
			ChildTotalCoupon: childCoupon3,
			ChildGoodsAmount: childGoodsAmount3,
			ChildExpFare:     childExpFare3,
			ChildOrderAmount: childOrderAmount3,
		}
		if err := tx.Debug().Create(&childOrder3).Error; err != nil {
			tx.Rollback()
			panic(orderException.OrderCreateFail())
		}
	}

	sql := squirrel.Insert("test_order_detail").Columns(
		"account_id","child_order_id", "order_id", "goods_id", "goods_specification_id", "purchase_qty", "message", "coupon", "exp_fare", "goods_amount",
		"create_time", "update_time", "order_amount","is_comment",
	)
	_time := time.Now().Unix()
	//创建订单明细
	if len(expressage) > 0 {
		for _, _expressage := range expressage {
			sql = sql.Values(
				auth.AccountModel().Id,
				childOrder1.ID,
				order.ID,
				_expressage["goods_id"],
				_expressage["goods_specification"],
				_expressage["goods_total"],
				_expressage["message"],
				_expressage["coupon"],
				_expressage["exp_fare"],
				_expressage["goods_amount"],
				_time,
				_time,
				_expressage["goods_amount"].(int)-_expressage["coupon"].(int),
				0,
			)
		}
	}

	if len(sameCity) > 0 {
		for _, _sameCity := range sameCity {
			sql = sql.Values(
				auth.AccountModel().Id,
				childOrder2.ID,
				order.ID,
				_sameCity["goods_id"],
				_sameCity["goods_specification"],
				_sameCity["goods_total"],
				_sameCity["message"],
				_sameCity["coupon"],
				_sameCity["exp_fare"],
				_sameCity["goods_amount"],
				_time,
				_time,
				_sameCity["goods_amount"].(int)-_sameCity["coupon"].(int),
				0,
			)
		}
	}

	if len(self) > 0 {
		for _, _self := range self {
			sql = sql.Values(
				auth.AccountModel().Id,
				childOrder3.ID,
				order.ID,
				_self["goods_id"],
				_self["goods_specification"],
				_self["goods_total"],
				_self["message"],
				_self["coupon"],
				_self["exp_fare"],
				_self["goods_amount"],
				_time,
				_time,
				_self["goods_amount"].(int)-_self["coupon"].(int),
				0,
			)
		}
	}

	if s, args, err := sql.ToSql(); err != nil {
		logUtils.Println(err)
		tx.Rollback()
		panic(orderException.OrderCreateFail())
	} else {
		if err := db.Driver.Exec(s, args...).Error; err != nil {
			logUtils.Println(err)
			tx.Rollback()
			panic(orderException.OrderCreateFail())
		}
	}
	tx.Commit()

	ctx.JSON(iris.Map{
		"id": order.ID,
	})

}

//只能修改总订单备注，地址
func PutOrder(ctx iris.Context, auth authbase.AuthAuthorization, oid int) {
	auth.CheckLogin()

	var order db.TestOrder
	if err := db.Driver.GetOne("test_order", oid, &order); err != nil {
		panic(orderException.OrderIsNotExsit())
	}

	//已取消订单不能修改
	var cancelChildOrder db.TestChildOrder
	if err := db.Driver.Where("order_id = ? and order_status = ?", oid, orderEnums.CANCEL).First(&cancelChildOrder).Error; err == nil {
		panic(orderException.OrderCancelNoPut())
	}

	if order.AccountID != auth.AccountModel().Id {
		panic(accountException.NoPermission())
	}

	if order.PayOrNot {
		panic(orderException.OrderPutStatusFail())
	}

	var childOrder []db.TestChildOrder
	db.Driver.Where("order_id=?", oid).Find(&childOrder)

	tx := db.Driver.Begin()
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	if params.Has("address_id") {
		addressId := params.Int("address_id", "地址id")
		//判断地址id是否属于本人
		var address db.Address
		if err := db.Driver.Where("account_id = ? and id = ?", auth.AccountModel().Id, addressId).First(&address).Error; err != nil {
			panic(orderException.AddressIsNotExsit())
		}

		order.AddressID = addressId
		for _, c := range childOrder {
			c.AddressID = addressId
			if err := tx.Save(&c).Error; err != nil {
				tx.Rollback()
				panic(orderException.OrderPutFail())
			}
		}

		if err := tx.Save(&order).Error; err != nil {
			tx.Rollback()
			panic(orderException.OrderPutFail())
		}

	}

	if params.Has("message") {
		message := params.List("message", "备注")

		for _, v := range message {
			goodsId := v.(map[string]interface{})["goods_id"]
			var orderDetail db.TestOrderDetail
			db.Driver.Where("order_id = ? and goods_id = ?", oid, goodsId).First(&orderDetail)
			orderDetail.Message = v.(map[string]interface{})["message"].(string)
			if err := tx.Save(&orderDetail).Error; err != nil {
				tx.Rollback()
				panic(orderException.OrderPutFail())
			}
		}
	}

	tx.Commit()
	ctx.JSON(iris.Map{
		"id": order.ID,
	})
}

//获取总订单信息
func GetOrder(ctx iris.Context, auth authbase.AuthAuthorization, oid int) {
	//不走缓存
	auth.CheckLogin()

	var order db.TestOrder
	//if err := db.Driver.GetOne("test_order", oid, &order); err != nil {
	//	panic(orderException.OrderIsNotExsit())
	//}
	if err := db.Driver.Where("id = ?", oid).First(&order).Error; err != nil {
		panic(orderException.OrderIsNotExsit())
	}

	if order.AccountID != auth.AccountModel().Id {
		panic(accountException.NoPermission())
	}

	var cOrder db.TestChildOrder
	if err := db.Driver.Where("order_id = ?", oid).First(&cOrder).Error; err == nil {
		if !order.PayOrNot && order.CreateTime+86400 < time.Now().Unix() {
			Recover(oid)
		}
	}

	ctx.JSON(order.GetInfo())
}

//删除子订单
func DeleteOrder(ctx iris.Context, auth authbase.AuthAuthorization, ooid int, oid int) {
	auth.CheckLogin()

	var oo db.TestOrder
	if err := db.Driver.GetOne("test_order", ooid, &oo); err == nil {
		if oo.AccountID != auth.AccountModel().Id {
			panic(accountException.NoPermission())
		}
		var order db.TestChildOrder
		if err := db.Driver.Where("id = ? and order_id = ?", oid, ooid).First(&order).Error; err != nil {
			panic(orderException.ChildOrderIsNotExsit())
		} else {
			if order.OrderStatus != orderEnums.OVER && order.OrderStatus != orderEnums.CANCEL {
				panic(orderException.OrderDeleteStatusFail())
			}

			//订单备份
			tx := db.Driver.Begin()
			o := db.TestChildOrderCopy{
				OrderStatus:      order.OrderStatus,
				OrderNum:         order.OrderNum,
				AccountID:        order.AccountID,
				AddressID:        order.AddressID,
				OrderID:          order.OrderID,
				Delivery:         order.Delivery,
				ChildTotalCoupon: order.ChildTotalCoupon,
				ChildGoodsAmount: order.ChildGoodsAmount,
				DeliveryTime:     order.DeliveryTime,
				GetTime:          order.GetTime,
				TrackingID:       order.TrackingID,
				CreateTime:       order.CreateTime,
				UpdateTime:       order.UpdateTime,
				ChildExpFare:     order.ChildExpFare,
				ChildOrderAmount: order.ChildOrderAmount,
			}
			if err := tx.Debug().Create(&o).Error; err != nil {
				tx.Rollback()

				panic(orderException.OrderDeleteFail())
			}

			if err := tx.Debug().Delete(order).Error; err != nil {
				tx.Rollback()

				panic(orderException.OrderDeleteFail())
			}
			tx.Commit()
		}

	}

	ctx.JSON(iris.Map{
		"id": oid,
	})
}

func DeleteAllOrder(ctx iris.Context, auth authbase.AuthAuthorization, ooid int) {
	auth.CheckLogin()

	var oo db.TestOrder
	if err := db.Driver.GetOne("test_order", ooid, &oo); err == nil {
		if oo.AccountID != auth.AccountModel().Id {
			panic(accountException.NoPermission())
		}

		var o1 []db.TestChildOrder
		db.Driver.Where("order_id = ?", ooid).Find(&o1)

		var o2 []db.TestChildOrder
		db.Driver.Where("order_id = ? and order_status in (?)", ooid, []int16{orderEnums.CANCEL, orderEnums.OVER}).Find(&o2)

		if len(o1) != len(o2) {
			panic(orderException.OrderDeleteStatusFail())
		}

		tx := db.Driver.Begin()

		sql := squirrel.Insert("test_child_order_copy").Columns(
			"order_status", "order_num", "account_id", "address_id", "order_id", "delivery", "child_total_coupon", "child_goods_amount", "child_exp_fare", "child_order_amount", "delivery_time",
			"get_time", "tracking_id", "create_time", "update_time",
		)

		for _, o := range o1 {
			sql = sql.Values(
				o.OrderStatus,
				o.OrderNum,
				o.AccountID,
				o.AddressID,
				o.OrderID,
				o.Delivery,
				o.ChildTotalCoupon,
				o.ChildGoodsAmount,
				o.ChildExpFare,
				o.ChildOrderAmount,
				o.DeliveryTime,
				o.GetTime,
				o.TrackingID,
				o.CreateTime,
				o.UpdateTime,
			)
		}

		if s, args, err := sql.ToSql(); err != nil {
			logUtils.Println(err)
			tx.Rollback()
			panic(orderException.OrderDeleteFail())
		} else {
			if err := db.Driver.Exec(s, args...).Error; err != nil {
				logUtils.Println(err)
				tx.Rollback()
				panic(orderException.OrderDeleteFail())
			}
		}

		tx.Delete(&o1)

		tx.Commit()
	}

	ctx.JSON(iris.Map{
		"status": "success",
	})
}

//获取子订单列表
func ListOrder(ctx iris.Context, auth authbase.AuthAuthorization) {
	auth.CheckLogin()

	var lists []struct {
		Id         int   `json:"id"`
		UpdateTime int64 `json:"update_time"`
		CreateTime int64 `json:"create_time"`
	}
	var count int

	table := db.Driver.Table("test_child_order")

	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)

	if !auth.IsAdmin() {
		table = table.Where("account_id = ?", auth.AccountModel().Id)
	}

	//所属人过滤
	if author := ctx.URLParamIntDefault("author_id", 0); author != 0 && auth.IsAdmin() {
		table = table.Where("account_id = ?", author)
	}

	//状态过滤
	if status := ctx.URLParamIntDefault("status", 0); status != 0 {
		if status == 10{
			table = table.Where("order_status in (?)", []int16{orderEnums.ToPick,orderEnums.TbShipped})
		}else {
			table = table.Where("order_status = ?", status)
		}
	}

	//名字过滤
	if name := ctx.URLParam("name"); len(name) > 0 {
		keyString := fmt.Sprintf("%%%s%%", name)

		var data []struct {
			db.TestOrderDetail
			db.Goods
		}
		db.Driver.Table("test_order_detail, goods").Where("goods.id = test_order_detail.goods_id and goods.name like ?", keyString).Find(&data)
		ids := make([]int, 0)
		for _, v := range data {
			ids = append(ids, v.TestOrderDetail.ChildOrderID)
		}
		table = table.Where("id in (?)", ids)
	}

	if startTime := ctx.URLParamInt64Default("start_time", 0); startTime != 0 {
		endTime := ctx.URLParamInt64Default("end_time", 0)
		table = table.Where("create_time between ? and ?", startTime, endTime)
	}

	table.Count(&count).Order("create_time desc").Offset((page - 1) * limit).Limit(limit).Select("id, update_time, create_time").Find(&lists)
	ctx.JSON(iris.Map{
		"orders": lists,
		"total":  count,
		"limit":  limit,
		"page":   page,
	})
}

//获取子订单信息
func MgetOrder(ctx iris.Context, auth authbase.AuthAuthorization) {
	//不走缓存
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "id列表")

	//tx := db.Driver.Begin()

	data := make([]interface{}, 0, len(ids))
	//orders := db.Driver.GetMany("test_child_order", ids, db.TestChildOrder{})
	var orders []db.TestChildOrder
	db.Driver.Where("id in (?)", ids).Find(&orders)

	o := make([]db.TestChildOrder, 0)

	for _, i := range ids {

		for _, order := range orders {
			//fmt.Println("eeeee")
			//fmt.Println(order.ID)
			if order.ID == int(i.(float64)) {
				//fmt.Println(order.ID)
				o = append(o, order)
				//fmt.Println(order.GetInfo())
				//fmt.Println("11")
			}
		}
	}
	for _, order := range o {

		if order.AccountID != auth.AccountModel().Id && !auth.IsAdmin() {
			continue
		}

		if order.OrderStatus == orderEnums.SUBMIT {
			key := paramsUtils.CacheBuildKey(constants.GlobalOrderTime, "auto_cancel")
			re, err := redis.Int64(cache.Redis.Do(constants.DbNumberOther, "get", key))
			if err != nil || re == 0 {
				re = 86400
			}

			var parent db.TestOrder
			//db.Driver.GetOne("test_order",order.OrderID,&parent)
			db.Driver.Where("id = ?", order.OrderID).First(&parent)

			if !parent.PayOrNot && parent.CreateTime+re < time.Now().Unix() {
				Recover(parent.ID)
			}
		}

		if order.OrderStatus == orderEnums.ToPick{
			key := paramsUtils.CacheBuildKey(constants.GlobalOrderTime, "pick_auto_get")
			re, err := redis.Int64(cache.Redis.Do(constants.DbNumberOther, "get", key))
			if err != nil || re == 0 {
				re = 604800
			}

			if order.CreateTime + re < time.Now().Unix(){
				order.OrderStatus = orderEnums.RECEIVED
				order.GetTime = order.CreateTime + re

				db.Driver.Debug().Save(&order)
			}
		}

		//TODO 已发货8天后改为已完成
		if order.OrderStatus == orderEnums.TbShipped {
			key := paramsUtils.CacheBuildKey(constants.GlobalOrderTime, "tran_auto_get")
			re, err := redis.Int64(cache.Redis.Do(constants.DbNumberOther, "get", key))
			if err != nil || re == 0 {
				re = 604800
			}

			if order.DeliveryTime+re < time.Now().Unix() {
				order.OrderStatus = orderEnums.RECEIVED
				order.GetTime = order.DeliveryTime + re

				db.Driver.Debug().Save(&order)
			}
		}

		// 已收货8天后改为订单结束
		if order.OrderStatus == orderEnums.RECEIVED {
			key := paramsUtils.CacheBuildKey(constants.GlobalOrderTime, "to_over")
			re, err := redis.Int64(cache.Redis.Do(constants.DbNumberOther, "get", key))
			if err != nil || re == 0 {
				re = 604800
			}

			if order.GetTime+re < time.Now().Unix() {
				order.OrderStatus = orderEnums.OVER
				order.OverTime = time.Now().Unix()

				//todo 此子订单下所有订单明细全部默认好评
				var orderDetails []db.TestOrderDetail
				db.Driver.Where("child_order_id = ?", order.ID).Find(&orderDetails)

				sql := squirrel.Insert("comment").Columns(
					"author_id", "good_id", "order_id", "comment_tag", "content", "create_time", "update_time",
				)

				flag := false
				tx := db.Driver.Begin()
				for _, detail := range orderDetails {
					if detail.IsComment == 0 {
						sql = sql.Values(
							order.AccountID,
							detail.GoodsID,
							detail.ID,
							1,
							"系统默认好评",
							order.GetTime+re,
							order.GetTime+re,
						)
						detail.IsComment = 1
						if err := tx.Save(&detail).Error;err != nil{
							tx.Rollback()
							panic("111")
						}
						tx.Commit()
						flag = true
					}
				}
				if flag{
					if s, args, err := sql.ToSql(); err != nil {
						logUtils.Println(err)
						panic(orderException.OrderDeleteFail())
					} else {
						if err := db.Driver.Exec(s, args...).Error; err != nil {
							logUtils.Println(err)
							panic(orderException.OrderDeleteFail())
						}
					}
				}


				db.Driver.Debug().Save(&order)
			}
		}

		if order.OrderStatus == orderEnums.AfterSales{
			var details []db.TestOrderDetail
			db.Driver.Where("child_order_id = ? and is_after_serve != 0 and is_pass = ?",order.ID,orderEnums.IsPassSUMMIT).Find(&details)
			if len(details) == 0{
				order.OrderStatus = orderEnums.AlreadyAS

				db.Driver.Debug().Save(&order)
			}
		}

		func(data *[]interface{}) {
			*data = append(*data, order.GetInfo())
			defer func() {
				recover()
			}()
		}(&data)
	}
	ctx.JSON(data)
}

func MgetOrderDetail(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "id列表")

	details := db.Driver.GetMany("test_order_deatil",ids,db.TestOrderDetail{})

	data := make([]interface{}, 0, len(ids))

	for _,detail := range details{
		if detail.(db.TestOrderDetail).AccountID != auth.AccountModel().Id && !auth.IsAdmin() {
			continue
		}

		func(data *[]interface{}) {
			*data = append(*data, paramsUtils.ModelToDict(detail,orderDetailField))
			defer func() {
				recover()
			}()
		}(&data)
	}

	ctx.JSON(data)
}

//用户取消订单
func CancelOrder(ctx iris.Context, auth authbase.AuthAuthorization, oid int) {
	auth.CheckLogin()

	var order db.TestOrder
	if err := db.Driver.GetOne("test_order", oid, &order); err != nil {
		panic(orderException.OrderIsNotExsit())
	}

	if order.AccountID != auth.AccountModel().Id {
		panic(accountException.NoPermission())
	}

	//只有未付款可以取消
	if order.PayOrNot {
		panic(orderException.OrderNotCancel())
	}

	Recover(order.ID)

	ctx.JSON(iris.Map{
		"id": order.ID,
	})

}

//订单取消时恢复
func Recover(oid int) {

	tx := db.Driver.Begin()

	//订单状态改变
	var cOrders []db.TestChildOrder
	db.Driver.Where("order_id = ?", oid).Find(&cOrders)

	for _, c := range cOrders {
		c.OrderStatus = orderEnums.CANCEL
		if err := tx.Save(&c).Error; err != nil {
			panic(orderException.OrderCancelFail())
		}

		//库存还原
		var orderDetails []db.TestOrderDetail
		db.Driver.Where("child_order_id = ?", c.ID).Find(&orderDetails)
		for _, d := range orderDetails {
			var goods db.Goods
			db.Driver.GetOne("goods", d.GoodsID, &goods)

			//总库存
			goods.Total += d.PurchaseQty
			if err := tx.Save(&goods).Error; err != nil {
				panic(orderException.OrderCancelFail())
			}

			//对应规格库存
			var specification db.GoodsSpecification
			db.Driver.Where("goods_id = ?", d.GoodsID).First(&specification)

			s := goods.GetInfo()["specification"]

			if goods.Sale {
				for index, v := range s.(specificationEntity.GoodsSpecificationEntity) {
					if v.ID == d.GoodsSpecificationID {
						s.(specificationEntity.GoodsSpecificationEntity)[index].Total += d.PurchaseQty

						if data, err := json.Marshal(s.(specificationEntity.GoodsSpecificationEntity)); err != nil {
							panic(goodsInfoException.SpecificationMarshalFail())
						} else {
							specification.Specification = string(data)
						}
						if err := tx.Save(&specification).Error; err != nil {
							tx.Rollback()
							panic(orderException.OrderCancelFail())
						}
					}
				}
			} else {
				for index, v := range s.(specificationEntity.GoodsSpecificationEntity1) {
					if v.ID == d.GoodsSpecificationID {
						s.(specificationEntity.GoodsSpecificationEntity1)[index].Total += d.PurchaseQty

						if data, err := json.Marshal(s.(specificationEntity.GoodsSpecificationEntity1)); err != nil {
							panic(goodsInfoException.SpecificationMarshalFail())
						} else {
							specification.Specification = string(data)
						}
						if err := tx.Save(&specification).Error; err != nil {
							tx.Rollback()
							panic(orderException.OrderCancelFail())
						}
					}
				}
			}
		}
	}

	tx.Commit()
}

func CheckGet(ctx iris.Context, auth authbase.AuthAuthorization, ooid int, oid int) {
	auth.CheckLogin()

	var o db.TestChildOrder
	if err := db.Driver.Where("order_id = ? and id =?", ooid, oid).First(&o).Error; err != nil {
		panic(orderException.ChildOrderIsNotExsit())
	}

	if o.AccountID != auth.AccountModel().Id {
		panic(accountException.NoPermission())
	}

	if o.OrderStatus != orderEnums.TbShipped && o.OrderStatus != orderEnums.ToPick {
		panic(orderException.OrderPutStatusFail())
	}

	o.OrderStatus = orderEnums.RECEIVED
	o.GetTime = time.Now().Unix()

	db.Driver.Save(&o)

	ctx.JSON(iris.Map{
		"status": "success",
	})
}

//仅供测试，修改子订单状态
func PutOrderStatus(ctx iris.Context, auth authbase.AuthAuthorization, ooid int, oid int) {
	var order db.TestOrder
	if err := db.Driver.GetOne("test_order", ooid, &order); err != nil {
		panic(orderException.OrderIsNotExsit())
	}

	var o db.TestChildOrder
	if err := db.Driver.Where("order_id = ? and id =?", ooid, oid).First(&o).Error; err != nil {
		panic(orderException.OrderIsNotExsit())
	}
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	if params.Has("status") {
		status := params.Int("status", "状态")
		if status != orderEnums.SUBMIT && status != orderEnums.WfShipped && status != orderEnums.TbShipped &&
			status != orderEnums.RECEIVED && status != orderEnums.OVER && status != orderEnums.CANCEL {
			panic(orderException.OrderPutFail())
		}
		o.OrderStatus = int16(status)
	}

	db.Driver.Save(&o)

	ctx.JSON(iris.Map{
		"id": o.ID,
	})
}
