package goods_tag

import (
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	goodsTagEnum "grpc-demo/enums/goods_tag"
	goodsTagException "grpc-demo/exceptions/goods_tag"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

func CreateSaleTag(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckAdmin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	title := params.Str("title","标题")

	tag := db.SaleTag{
		Title:title,
	}
	db.Driver.Create(&tag)
	ctx.JSON(iris.Map{
		"id": tag.ID,
	})
}

func PutSaleTag(ctx iris.Context, auth authbase.AuthAuthorization,sid int){
	auth.CheckAdmin()

	var tag db.SaleTag
	if err := db.Driver.GetOne("sale_tag",sid,&tag);err != nil{
		panic(goodsTagException.SaleTagIsNotExists())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	params.Diff(tag)
	tag.Title = params.Str("title","标题")


	db.Driver.Save(&tag)
	ctx.JSON(iris.Map{
		"id": tag.ID,
	})
}

func DeleteSaleTag(ctx iris.Context, auth authbase.AuthAuthorization,sid int){
	auth.CheckAdmin()

	var tag db.SaleTag
	if err := db.Driver.GetOne("sale_tag",sid,&tag);err == nil{

		db.Driver.Exec("delete from goods_and_tag where goods_tag_id = ? and tag_type = ?", sid, goodsTagEnum.TagTypeSaleTag)

		db.Driver.Delete(tag)
	}

	ctx.JSON(iris.Map{
		"id": sid,
	})
}

func ListSaleTag(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckLogin()

	var lists []struct {
		Id         int   `json:"id"`
		Title string `json:"title"`
	}
	var count int

	table := db.Driver.Table("sale_tag")

	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)

	table.Count(&count).Offset((page - 1) * limit).Limit(limit).Select("id, title").Find(&lists)
	ctx.JSON(iris.Map{
		"tags": lists,
		"total":    count,
		"limit":    limit,
		"page":     page,
	})

}

func MgetSaleTag(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "id列表")

	data := make([]interface{}, 0, len(ids))
	tags := db.Driver.GetMany("sale_tag", ids, db.SaleTag{})
	for _, tag := range tags {
		func(data *[]interface{}) {
			*data = append(*data, paramsUtils.ModelToDict(tag,[]string{"ID","Title"}))
			defer func() {
				recover()
			}()
		}(&data)
	}
	ctx.JSON(data)
}
