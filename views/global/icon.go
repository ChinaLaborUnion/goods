package global

import (
	"fmt"
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	globalException "grpc-demo/exceptions/global"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

func Create(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckAdmin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	title := params.Str("title","标题")
	name := params.Str("name","名称")
	picture := params.Str("picture","图片")

	icon := db.Icon{
		Title:    title,
		Name:     name,
		Picture: picture,
	}

	db.Driver.Create(&icon)

	ctx.JSON(iris.Map{
		"id": icon.ID,
	})
}

func Put(ctx iris.Context, auth authbase.AuthAuthorization,id int){
	auth.CheckAdmin()

	var icon db.Icon
	if err := db.Driver.GetOne("icon",id,&icon);err != nil{
		panic(globalException.IconIsNotExsit())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	params.Diff(icon)
	icon.Title = params.Str("title","标题")
	icon.Name = params.Str("name","名称")
	icon.Picture = params.Str("picture","图片")

	db.Driver.Save(&icon)

	ctx.JSON(iris.Map{
		"id": icon.ID,
	})
}

func Delete(ctx iris.Context, auth authbase.AuthAuthorization,id int){
	auth.CheckAdmin()

	var icon db.Icon
	if err := db.Driver.GetOne("icon",id,&icon);err == nil{
		db.Driver.Delete(icon)
	}

	ctx.JSON(iris.Map{
		"id": id,
	})
}

func List(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckLogin()

	var lists []struct {
		Id         int   `json:"id"`
		Name string `json:"name"`
	}
	var count int

	table := db.Driver.Table("icon")

	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)

	//title过滤
	if title := ctx.URLParam("title"); len(title) > 0 {
		keyString := fmt.Sprintf("%%%s%%", title)
		table = table.Where("title like ?", keyString)
	}

	table.Count(&count).Offset((page - 1) * limit).Limit(limit).Select("id, name").Find(&lists)
	ctx.JSON(iris.Map{
		"icons": lists,
		"total": count,
		"limit": limit,
		"page":  page,
	})
}

func Mget(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "id列表")

	data := make([]interface{}, 0, len(ids))
	icons := db.Driver.GetMany("icon", ids, db.Icon{})
	for _, v := range icons {
		func(data *[]interface{}) {
			*data = append(*data, paramsUtils.ModelToDict(v,[]string{"ID","Title","Name","Picture"}))
			defer func() {
				recover()
			}()
		}(&data)
	}
	ctx.JSON(data)
}



