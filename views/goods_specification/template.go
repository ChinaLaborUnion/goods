package goods_specification

import (
	"fmt"
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	goods "grpc-demo/entity/goods_specification"
	goodsSpecificationException "grpc-demo/exceptions/goods_specification"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

func CreateTemplate(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckAdmin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	title := params.Str("title","名称")
	template := params.List("template","规格模版")

	s := goods.TemplateMarshal(template)

	tt := db.SpecificationTemplate{
		Title:    title,
		Template: s,
	}

	db.Driver.Create(&tt)

	ctx.JSON(iris.Map{
		"id": tt.ID,
	})
}

func PutTemplate(ctx iris.Context, auth authbase.AuthAuthorization,tid int){
	auth.CheckAdmin()

	var t db.SpecificationTemplate
	if err := db.Driver.GetOne("specification_template",tid,&t);err != nil{
		panic(goodsSpecificationException.TemplateIsNotExsit())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	params.Diff(t)
	t.Title = params.Str("title","名称")

	if params.Has("template"){
		template := params.List("template","规格模版")
		s := goods.TemplateMarshal(template)
		t.Template = s
	}

	db.Driver.Save(&t)
	ctx.JSON(iris.Map{
		"id": t.ID,
	})
}

func DeleteTemplate(ctx iris.Context, auth authbase.AuthAuthorization,tid int){
	auth.CheckAdmin()

	var t db.SpecificationTemplate
	if err := db.Driver.GetOne("specification_template",tid,&t);err == nil{
		var goods db.Goods
		if err := db.Driver.Where("template_id = ?", tid).First(&goods).Error; err == nil {
			panic(goodsSpecificationException.TemplateUsed())
		}
		db.Driver.Delete(t)
	}

	ctx.JSON(iris.Map{
		"id": tid,
	})
}

func ListTemplate(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckAdmin()

	var lists []struct {
		Id         int   `json:"id"`
		Title string `json:"title"`
	}
	var count int

	table := db.Driver.Table("specification_template")

	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)

	if title := ctx.URLParam("title"); len(title) > 0 {
		keyString := fmt.Sprintf("%%%s%%", title)
		table = table.Where("title like ?", keyString)
	}

	table.Count(&count).Offset((page - 1) * limit).Limit(limit).Select("id, title").Find(&lists)
	ctx.JSON(iris.Map{
		"templates": lists,
		"total":    count,
		"limit":    limit,
		"page":     page,
	})
}

func MgetTemplate(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckAdmin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "id列表")

	data := make([]interface{}, 0, len(ids))
	templates := db.Driver.GetMany("specification_template", ids, db.SpecificationTemplate{})
	for _, t := range templates {
		func(data *[]interface{}) {
			*data = append(*data, t.(db.SpecificationTemplate).GetInfo())
			defer func() {
				recover()
			}()
		}(&data)
	}
	ctx.JSON(data)
}








