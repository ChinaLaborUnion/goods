package goods_info

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	goodsInfo "grpc-demo/entity/goods_info"
	goodsInfoEnum "grpc-demo/enums/goods_info"
	goodsTagEnum "grpc-demo/enums/goods_tag"
	goodsInfoException "grpc-demo/exceptions/goods_info"
	"grpc-demo/models/db"
	logUtils "grpc-demo/utils/log"
	paramsUtils "grpc-demo/utils/params"
	"time"
)

func CreateGoods(ctx iris.Context, auth authbase.AuthAuthorization) {
	auth.CheckAdmin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	name := params.Str("name", "商品名")


	goods := db.Goods{
		Name:       name,

	}

	if params.Has("paid_and_remove") {
		goods.PaidAndRemove = params.Bool("paid_and_remove", "是否付款减库存")
	}
	if params.Has("show_total") {
		goods.ShowTotal = params.Bool("show_total", "是否展示库存")
	}
	if params.Has("exchange") {
		goods.Exchange = params.Bool("exchange", "是否可以换货")
	}
	if params.Has("sale_return") {
		goods.SaleReturn = params.Bool("sale_return", "是否可以退货")
	}
	if params.Has("min_sale") {
		goods.MinSale = params.Int("min_sale", "起售数")
	}

	//收货方式，运费
	getWay := params.Int("get_way", "收货方式")
	getWayEnums := goodsInfoEnum.NewGetWayEnums()
	if !getWayEnums.Has(getWay){
		panic(goodsInfoException.GetWayIsNotExsit())
	}else{
		if getWay != goodsInfoEnum.GetWaySelf {
			carriage := params.Int("carriage", "运费")
			goods.Carriage = carriage
		}
		goods.GetWay = int16(getWay)
	}
	//getWayEnums := goodsInfoEnum.NewGetWayEnums()
	//if !getWayEnums.Has(getWay) {
	//	panic(goodsInfoException.GetWayIsNotExsit())
	//} else {
	//	if getWay != goodsInfoEnum.GetWaySelf {
	//		carriage := params.Float("carriage", "运费")
	//		goods.Carriage = carriage
	//	}
	//	goods.GetWay = int16(getWay)
	//}

	//上架方式，上架时间
	putaway := params.Int("putaway", "上架方式")
	putawayEnums := goodsInfoEnum.NewPutawayEnums()
	if !putawayEnums.Has(putaway) {
		panic(goodsInfoException.PutawayIsNotExsit())
	} else {
		if putaway == goodsInfoEnum.PutawayDefine {
			putawayTime := params.Int("putaway_time", "上架时间")
			if int64(putawayTime) < time.Now().Unix() {
				panic(goodsInfoException.PutawayTimeError())
			}
			goods.PutawayTime = int64(putawayTime)
		}
		goods.Putaway = int16(putaway)
	}

	if params.Has("sale_point") {
		goods.SalePoint = params.Str("sale_point", "卖点")
	}

	//是否优惠，优惠价
	if params.Has("sale") {
		goods.Sale = params.Bool("sale", "是否优惠")

	}

	//是否预售，预售时间
	if params.Has("advance") && params.Bool("advance", "是否预售") {
		advanceTime := params.Int("advance_time", "预售时间")
		if int64(advanceTime) < time.Now().Unix() {
			panic(goodsInfoException.AdvanceTimeError())
		}
		goods.Advance = params.Bool("advance", "是否预售")
		goods.AdvanceTime = int64(advanceTime)
	}

	//是否限购，限购数量
	if params.Has("limit") && params.Bool("limit", "是否限购") {
		limitTotal := params.Int("limit_total", "限购数")
		goods.Limit = params.Bool("limit", "是否限购")
		goods.LimitTotal = limitTotal
	}

	if params.Has("total") {
		goods.Total = params.Int("total", "库存")
	}

	//是否上架
	if params.Has("on_sale") && params.Bool("on_sale", "是否上架") {
		if goods.Total > 0 && ( goods.Putaway == goodsInfoEnum.PutawayNow || (goods.Putaway == goodsInfoEnum.PutawayDefine && goods.PutawayTime <= time.Now().Unix())) {
			goods.OnSale = params.Bool("on_sale", "是否上架")
		}
	}

	if params.Has("view"){
		goods.View = params.Str("view","视频")
	}

	if params.Has("cover"){
		goods.Cover = params.Str("cover","封面")
	}

	if params.Has("detail"){
		goods.Detail = params.Str("detail","详情")
	}

	if params.Has("pictures"){
		pictures := params.List("pictures","图片列表")
		goods.Pictures = goodsInfo.GoodsPictureMarshal(pictures)

	}

	db.Driver.Create(&goods)

	//标签挂载
	if params.Has("place_tag") {
		placeTag := params.List("place_tag", "属地标签")
		placeTagMnt(placeTag, goods)
	}

	if params.Has("sale_tag") {
		saleTag := params.List("sale_tag", "销售标签")
		saleTagMnt(saleTag, goods)
	}

	if params.Has("kind_tag") {
		kindTag := params.List("kind_tag", "种类标签")
		kindTagMnt(kindTag, goods)
	}

	ctx.JSON(iris.Map{
		"id": goods.ID,
	})
}

func PutGoods(ctx iris.Context, auth authbase.AuthAuthorization, gid int) {
	auth.CheckAdmin()

	var goods db.Goods
	if err := db.Driver.GetOne("goods", gid, &goods); err != nil {
		panic(goodsInfoException.GoodsIsNotExsit())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	params.Diff(goods)
	goods.Name = params.Str("name", "商品名")

	goods.ShowTotal = params.Bool("show_total", "是否展示库存")
	goods.Exchange = params.Bool("exchange", "是否可以换货")
	goods.SaleReturn = params.Bool("sale_return", "是否可以退货")
	goods.SalePoint = params.Str("sale_point", "卖点")
	goods.MinSale = params.Int("min_sale", "起售数")
	goods.PaidAndRemove = params.Bool("paid_and_remove", "是否付款减库存")
	goods.Total = params.Int("total", "总数")
	goods.Sale = params.Bool("sale", "是否优惠")


	goods.LimitTotal = params.Int("limit_total", "限购数")

	goods.Carriage = params.Int("carriage", "运费")

	if params.Has("advance_time") {
		advanceTime := params.Int("advance_time", "预售时间")
		if int64(advanceTime) >= time.Now().Unix() {
			goods.AdvanceTime = int64(advanceTime)
		}
	}

	if params.Has("putaway_time") {
		putawayTime := params.Int("putaway_time", "上架时间")
		if int64(putawayTime) >= time.Now().Unix() {
			goods.PutawayTime = int64(putawayTime)
		}
	}


	if params.Has("limit") {
		limit := params.Bool("limit", "是否限购")
		if limit != goods.Limit && limit {
			limitTotal := params.Int("limit_total", "限购数")
			goods.Limit = limit
			goods.LimitTotal = limitTotal
		} else {
			goods.Limit = limit
		}
	}

	if params.Has("advance") {
		advance := params.Bool("advance", "是否预售")
		if advance != goods.Advance && advance {
			atime := params.Int("advance_time", "预售时间")
			if int64(atime) < time.Now().Unix() {
				panic(goodsInfoException.AdvanceTimeError())
			}
			goods.Advance = advance
			goods.AdvanceTime = int64(atime)
		} else {
			goods.Advance = advance
		}
	}

	if params.Has("get_way") {
		getWay := params.Int("get_way", "收货方式")
		if goods.GetWay == goodsInfoEnum.GetWaySelf && int16(getWay) != goods.GetWay {
			carriage := params.Int("carriage", "运费")
			goods.GetWay = int16(getWay)
			goods.Carriage = carriage
		} else {
			goods.GetWay = int16(getWay)
		}
	}

	if params.Has("putaway") {
		putaway := params.Int("putaway", "上架方式")
		if goods.Putaway != goodsInfoEnum.PutawayDefine && int16(putaway) == goodsInfoEnum.PutawayDefine {
			ptime := params.Int("putaway_time", "上架时间")
			if int64(ptime) < time.Now().Unix() {
				panic(goodsInfoException.PutawayTimeError())
			}
			goods.Putaway = int16(putaway)
			goods.PutawayTime = int64(ptime)
		} else {
			goods.Putaway = int16(putaway)
		}
	}

	if params.Has("on_sale") {
		onSale := params.Bool("on_sale", "是否上架")
		if onSale != goods.OnSale && onSale {
			if goods.Total > 0 && (goods.Putaway == goodsInfoEnum.PutawayNow || (goods.Putaway == goodsInfoEnum.PutawayDefine && goods.PutawayTime <= time.Now().Unix())) {
				goods.OnSale = onSale
			} else {
				panic(goodsInfoException.OnSaleFail())
			}
		} else {
			goods.OnSale = onSale
		}
	}

	if params.Has("view"){
		goods.View = params.Str("view","视频")
	}

	if params.Has("cover"){
		goods.Cover = params.Str("cover","封面")
	}

	if params.Has("detail"){
		goods.Detail = params.Str("detail","详情")
	}

	if params.Has("pictures"){
		pictures := params.List("pictures","图片列表")
		goods.Pictures = goodsInfo.GoodsPictureMarshal(pictures)
	}

	db.Driver.Save(&goods)

	//修改标签挂载
	if params.Has("place_tag") {
		placeTag := params.List("place_tag", "属地标签")
		db.Driver.Exec("delete from goods_and_tag where goods_id = ? and tag_type = ?", gid, goodsTagEnum.TagTypePlaceTag)
		placeTagMnt(placeTag, goods)
	}

	if params.Has("sale_tag") {
		saleTag := params.List("sale_tag", "销售标签")
		db.Driver.Exec("delete from goods_and_tag where goods_id = ? and tag_type = ?", gid, goodsTagEnum.TagTypeSaleTag)
		saleTagMnt(saleTag, goods)
	}

	if params.Has("kind_tag") {
		kindTag := params.List("kind_tag", "种类标签")
		db.Driver.Exec("delete from goods_and_tag where goods_id = ? and tag_type = ?", gid, goodsTagEnum.TagTypeKindTag)
		kindTagMnt(kindTag, goods)
	}

	ctx.JSON(iris.Map{
		"id": goods.ID,
	})

}

func DeleteGoods(ctx iris.Context, auth authbase.AuthAuthorization, gid int) {
	auth.CheckAdmin()

	//var goods db.Goods
	//if err := db.Driver.GetOne("goods", gid, &goods); err == nil {
	//	//同时删除标签挂载
	//	db.Driver.Exec("delete from goods_and_tag where goods_id = ?", gid)
	//
	//	//同时删除商品规格
	//	db.Driver.Exec("delete from goods_specification where goods_id = ?", gid)
	//
	//	//同时删除商品轮播
	//	db.Driver.Exec("delete from goods_slideshow where goods_id = ?", gid)
	//
	//	db.Driver.Delete(goods)
	//}

	var goods db.Goods
	if err := db.Driver.GetOne("goods",gid,&goods);err != nil{
		panic(goodsInfoException.GoodsIsNotExsit())
	}

	goods.OnSale = false
	db.Driver.Save(&goods)

	ctx.JSON(iris.Map{
		"id": gid,
	})
}

func ListGoods(ctx iris.Context, auth authbase.AuthAuthorization) {
	//auth.CheckLogin()

	var lists []struct {
		Id         int   `json:"id"`
		UpdateTime int64 `json:"update_time"`
	}
	var count int

	table := db.Driver.Table("goods")

	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)

	//关键词过滤
	if keyword := ctx.URLParam("keyword"); len(keyword) > 0 {
		keyString := fmt.Sprintf("%%%s%%", keyword)
		table = table.Where("name like ?", keyString)
	}

	//种类标签过滤
	if kindTag := ctx.URLParamIntDefault("kind_tag", 0); kindTag != 0 {
		var tag db.KindTag
		if err := db.Driver.GetOne("kind_tag", kindTag, &tag); err == nil {

			ids := make([]interface{}, 0)

			var goods []db.GoodsAndTag
			if err := db.Driver.Where("goods_tag_id = ? and tag_type = ?", kindTag, goodsTagEnum.TagTypeKindTag).Find(&goods).Error; err == nil {
				for _, v := range goods {
					ids = append(ids, v.GoodsID)
				}
			}

			var tags []db.KindTag
			db.Driver.Where("parent_id = ?",kindTag).Find(&tags)
			for _,t := range tags{
				var g []db.GoodsAndTag
				if err := db.Driver.Where("goods_tag_id = ? and tag_type = ?", t.ID, goodsTagEnum.TagTypeKindTag).Find(&g).Error; err == nil {
					for _, v1 := range g {
						ids = append(ids, v1.GoodsID)
					}
				}
			}

			table = table.Where("id in (?)", ids)
		}

	}

	//属地标签过滤
	if placeTag := ctx.URLParamIntDefault("place_tag", 0); placeTag != 0 {
		var tag db.PlaceTag
		if err := db.Driver.GetOne("place_tag", placeTag, &tag); err == nil {
			var goods []db.GoodsAndTag

			if err := db.Driver.Where("goods_tag_id = ? and tag_type = ?", placeTag, goodsTagEnum.TagTypePlaceTag).Find(&goods).Error; err == nil {
				ids := make([]interface{}, len(goods))
				for index, v := range goods {
					ids[index] = v.GoodsID
				}
				table = table.Where("id in (?)", ids)
			}
		}
	}

	//销售标签过滤
	if saleTag := ctx.URLParamIntDefault("sale_tag", 0); saleTag != 0 {
		var tag db.SaleTag
		if err := db.Driver.GetOne("sale_tag", saleTag, &tag); err == nil {
			var goods []db.GoodsAndTag

			if err := db.Driver.Where("goods_tag_id = ? and tag_type = ?", saleTag, goodsTagEnum.TagTypeSaleTag).Find(&goods).Error; err == nil {
				ids := make([]interface{}, len(goods))
				for index, v := range goods {
					ids[index] = v.GoodsID
				}
				table = table.Where("id in (?)", ids)
			}
		}
	}

	table.Count(&count).Order("create_time desc").Offset((page - 1) * limit).Limit(limit).Select("id, create_time").Find(&lists)

	//销量
	if sortWay := ctx.URLParamIntDefault("sort_way", 0); sortWay != 0{
		if sortWay == goodsInfoEnum.SortWayPriceDesc{
			table.Count(&count).Order("min_price desc").Offset((page - 1) * limit).Limit(limit).Select("id, update_time").Find(&lists)
		}else if sortWay == goodsInfoEnum.SortWayPriceAsc{
			table.Count(&count).Order("min_price").Offset((page - 1) * limit).Limit(limit).Select("id, update_time").Find(&lists)
		}else if sortWay == goodsInfoEnum.SortWayPeople{
			table.Count(&count).Order("people desc").Offset((page - 1) * limit).Limit(limit).Select("id, update_time").Find(&lists)
		}else{
			table.Count(&count).Order("create_time desc").Offset((page - 1) * limit).Limit(limit).Select("id, create_time").Find(&lists)
		}
	}

	ctx.JSON(iris.Map{
		"goods": lists,
		"total": count,
		"limit": limit,
		"page":  page,
	})

}

func MgetGoods(ctx iris.Context, auth authbase.AuthAuthorization) {
	//auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "id列表")


	data := make([]interface{}, 0, len(ids))
	goods := db.Driver.GetMany("goods", ids, db.Goods{})

	for _, g := range goods {
		func(data *[]interface{}) {
			*data = append(*data, g.(db.Goods).GetInfo())
			defer func() {
				recover()
			}()
		}(&data)

	}

	ctx.JSON(data)
}

//func LittleMgetGoods(ctx iris.Context, auth authbase.AuthAuthorization){
//	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
//	ids := params.List("ids", "id列表")
//
//	data := make([]interface{}, 0, len(ids))
//	goods := db.Driver.GetMany("goods", ids, db.Goods{})
//	for _, g := range goods {
//		func(data *[]interface{}) {
//			*data = append(*data, g.(db.Goods).GetInfo())
//			defer func() {
//				recover()
//			}()
//		}(&data)
//	}
//	ctx.JSON(data)
//}

func placeTagMnt(p []interface{}, goods db.Goods) {
	var tags []db.PlaceTag
	if err := db.Driver.Where("id in (?)", p).Find(&tags).Error; err != nil || len(tags) == 0 {
		logUtils.Println(err)
		return
	}

	sql := squirrel.Insert("goods_and_tag").Columns(
		"goods_id", "goods_tag_id", "tag_type",
	)

	for _, tag := range tags {
		sql = sql.Values(
			goods.ID,
			tag.ID,
			goodsTagEnum.TagTypePlaceTag,
		)
	}
	if s, args, err := sql.ToSql(); err != nil {
		logUtils.Println(err)
	} else {
		if err := db.Driver.Exec(s, args...).Error; err != nil {
			logUtils.Println(err)
			return
		}
	}
}

func saleTagMnt(s []interface{}, goods db.Goods) {
	var tags []db.SaleTag
	if err := db.Driver.Where("id in (?)", s).Find(&tags).Error; err != nil || len(tags) == 0 {
		logUtils.Println(err)
		return
	}

	sql := squirrel.Insert("goods_and_tag").Columns(
		"goods_id", "goods_tag_id", "tag_type",
	)

	for _, tag := range tags {
		sql = sql.Values(
			goods.ID,
			tag.ID,
			goodsTagEnum.TagTypeSaleTag,
		)
	}
	if s, args, err := sql.ToSql(); err != nil {
		logUtils.Println(err)
	} else {
		if err := db.Driver.Exec(s, args...).Error; err != nil {
			logUtils.Println(err)
			return
		}
	}
}

func kindTagMnt(k []interface{}, goods db.Goods) {
	var tags []db.KindTag
	if err := db.Driver.Where("id in (?)", k).Find(&tags).Error; err != nil || len(tags) == 0 {
		logUtils.Println(err)
		return
	}

	sql := squirrel.Insert("goods_and_tag").Columns(
		"goods_id", "goods_tag_id", "tag_type",
	)

	for _, tag := range tags {
		sql = sql.Values(
			goods.ID,
			tag.ID,
			goodsTagEnum.TagTypeKindTag,
		)
	}
	if s, args, err := sql.ToSql(); err != nil {
		logUtils.Println(err)
	} else {
		if err := db.Driver.Exec(s, args...).Error; err != nil {
			logUtils.Println(err)
			return
		}
	}
}
