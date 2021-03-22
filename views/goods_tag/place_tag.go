package goods_tag

import (
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	goodsTagEnum "grpc-demo/enums/goods_tag"
	goodsTagException "grpc-demo/exceptions/goods_tag"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

func CreatePlaceTag(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckAdmin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	place := params.Str("place","属地")
	picture := params.Str("picture","图片")

	tag := db.PlaceTag{
		Place: place,
		Picture:picture,
	}
	db.Driver.Create(&tag)
	ctx.JSON(iris.Map{
		"id": tag.ID,
	})
}

func PutPlaceTag(ctx iris.Context, auth authbase.AuthAuthorization,pid int){
	auth.CheckAdmin()

	var tag db.PlaceTag
	if err := db.Driver.GetOne("place_tag",pid,&tag);err != nil{
		panic(goodsTagException.PlaceTagIsNotExists())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	params.Diff(tag)
	tag.Place = params.Str("place","属地")
	tag.Picture = params.Str("picture","图片")

	db.Driver.Save(&tag)
	ctx.JSON(iris.Map{
		"id": tag.ID,
	})
}

func DeletePlaceTag(ctx iris.Context, auth authbase.AuthAuthorization,pid int){
	auth.CheckAdmin()

	var tag db.PlaceTag
	if err := db.Driver.GetOne("place_tag",pid,&tag);err == nil{
		db.Driver.Exec("delete from goods_and_tag where goods_tag_id = ? and tag_type = ?", pid, goodsTagEnum.TagTypePlaceTag)

		db.Driver.Delete(tag)
	}

	ctx.JSON(iris.Map{
		"id": pid,
	})
}

func ListPlaceTag(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckLogin()

	var lists []struct {
		Id         int   `json:"id"`
		Place string `json:"place"`
	}
	var count int

	table := db.Driver.Table("place_tag")

	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)

	table.Count(&count).Offset((page - 1) * limit).Limit(limit).Select("id, place").Find(&lists)
	ctx.JSON(iris.Map{
		"tags": lists,
		"total":    count,
		"limit":    limit,
		"page":     page,
	})

}

func MgetPlaceTag(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "id列表")

	data := make([]interface{}, 0, len(ids))
	tags := db.Driver.GetMany("place_tag", ids, db.PlaceTag{})
	for _, tag := range tags {
		func(data *[]interface{}) {
			*data = append(*data, paramsUtils.ModelToDict(tag,[]string{"ID","Place","Picture"}))
			defer func() {
				recover()
			}()
		}(&data)
	}
	ctx.JSON(data)
}