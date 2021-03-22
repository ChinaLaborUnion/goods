package goods_home

import (
	"fmt"
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	goodsHomeException "grpc-demo/exceptions/goods_home"
	goodsInfoException "grpc-demo/exceptions/goods_info"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

func CreateGoodsSlideshow(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckAdmin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))

	goodsId := params.Int("goods_id","商品id")
	var goods db.Goods
	if err := db.Driver.GetOne("goods",goodsId,&goods);err != nil{
		panic(goodsInfoException.GoodsIsNotExsit())
	}

	picture := params.Str("picture","轮播图")

	slideshow := db.GoodsSlideshow{
		GoodsID: goodsId,
		Picture: picture,
	}

	number := params.Int("number","序号")

	var s db.GoodsSlideshow
	if err := db.Driver.Where("number = ?", number).First(&s).Error; err == nil{
		panic(goodsHomeException.OrderExsit())
	}
	slideshow.Number = number


	db.Driver.Create(&slideshow)

	ctx.JSON(iris.Map{
		"id": slideshow.ID,
	})
}

func PutGoodsSlideshow(ctx iris.Context, auth authbase.AuthAuthorization,sid int){
	auth.CheckAdmin()

	var slideshow db.GoodsSlideshow
	if err := db.Driver.GetOne("goods_slideshow",sid,&slideshow);err != nil{
		panic(goodsHomeException.SlideshowIsNotExsit())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	params.Diff(slideshow)
	slideshow.Picture = params.Str("picture","轮播图")

	if params.Has("goods_id"){
		goodsId := params.Int("goods_id","商品id")
		var goods db.Goods
		if err := db.Driver.GetOne("goods",goodsId,&goods);err != nil {
			panic(goodsInfoException.GoodsIsNotExsit())
		}
		slideshow.GoodsID = goodsId
	}

	if params.Has("number"){
		number := params.Int("number","序号")
		var s db.GoodsSlideshow
		if err := db.Driver.Where("number = ?", number).First(&s).Error; err == nil{
			slideshow.Number,s.Number = s.Number,slideshow.Number
		}else{
			slideshow.Number = number
		}
		db.Driver.Save(&s)
	}

	db.Driver.Save(&slideshow)

	ctx.JSON(iris.Map{
		"id": slideshow.ID,
	})

}

func DeleteGoodsSlideshow(ctx iris.Context, auth authbase.AuthAuthorization,sid int){
	auth.CheckAdmin()

	var slideshow db.GoodsSlideshow
	if err := db.Driver.GetOne("goods_slideshow",sid,&slideshow);err == nil{
		db.Driver.Delete(slideshow)
	}

	ctx.JSON(iris.Map{
		"id": sid,
	})
}

func ListGoodsSlideshow(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckLogin()

	var lists []struct {
		Id         int   `json:"id"`
		GoodsID int64 `json:"goods_id"`
	}
	var count int

	table := db.Driver.Table("goods_slideshow")

	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)

	table.Count(&count).Offset((page - 1) * limit).Limit(limit).Select("id, goods_id").Order("number asc").Find(&lists)
	ctx.JSON(iris.Map{
		"slideshow": lists,
		"total": count,
		"limit": limit,
		"page":  page,
	})
}

func MgetGoodsSlideshow(ctx iris.Context, auth authbase.AuthAuthorization){

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	fmt.Println(params)
	ids := params.List("ids", "id列表")

	data := make([]interface{}, 0, len(ids))
	slideshow := db.Driver.GetMany("goods_slideshow", ids, db.GoodsSlideshow{})
	for _, s := range slideshow {
		func(data *[]interface{}) {
			*data = append(*data, paramsUtils.ModelToDict(s,[]string{"ID","GoodsID","Picture"}))
			defer func() {
				recover()
			}()
		}(&data)
	}
	ctx.JSON(data)
}




