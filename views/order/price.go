package order

import (
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	specificationEntity "grpc-demo/entity/goods_specification"
	goodsInfoEnum "grpc-demo/enums/goods_info"
	goodsInfoException "grpc-demo/exceptions/goods_info"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

func GetPrice(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	goodsList := params.List("goods_list", "所选商品")
	//[{"goods_id":,"goods_specification":,"goods_total":,"delivery":},{}]

	var totalCoupon,totalExpFare,totalGoodsAmount,totalOrderAmount int
	var childExpFare1 int
	var childExpFare2 int
	var childExpFare3 int

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
			for _, s := range specification.(specificationEntity.GoodsSpecificationEntity) {
				if s.ID == int(gd["goods_specification"].(float64)) {
					ok = true

					totalCoupon += (s.Price - s.ReducedPrice) * goodsTotal
					totalGoodsAmount += s.Price * goodsTotal

					//订单分单
					if delivery == goodsInfoEnum.GetWayExpressage {

						expFare1 = append(expFare1, g.Carriage)
					} else if delivery == goodsInfoEnum.GetWaySameCity {

						expFare2 = append(expFare2, g.Carriage)
					} else {

					}
					break
				}
			}
		} else {
			for _, s := range specification.(specificationEntity.GoodsSpecificationEntity1) {
				if s.ID == int(gd["goods_specification"].(float64)) {
					ok = true

					totalGoodsAmount += s.Price * goodsTotal

					//订单分单
					if delivery == goodsInfoEnum.GetWayExpressage {


						expFare1 = append(expFare1, g.Carriage)
					} else if delivery == goodsInfoEnum.GetWaySameCity {


						expFare2 = append(expFare2, g.Carriage)

					} else {

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

	totalExpFare += childExpFare1 + childExpFare2 + childExpFare3
	totalOrderAmount += totalGoodsAmount + totalExpFare - totalCoupon

	ctx.JSON(iris.Map{
		"total_goods_amount":totalGoodsAmount,
		"total_coupon":totalCoupon,
		"total_exp_fare":totalExpFare,
		"total_order_amount":totalOrderAmount,
	})
}
