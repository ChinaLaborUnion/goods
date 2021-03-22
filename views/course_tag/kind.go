package courseTag

import (
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	courseTagException "grpc-demo/exceptions/course_tag"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

func CreateKind(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckAdmin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	name := params.Str("name","名称")

	kind := db.Kind{
		Name: name,
	}

	db.Driver.Create(&kind)

	ctx.JSON(iris.Map{
		"id":kind.ID,
	})
}

func PutKind(ctx iris.Context, auth authbase.AuthAuthorization,kid int){
	auth.CheckAdmin()

	var kind db.Kind
	if err := db.Driver.GetOne("kind",kid,&kind);err != nil{
		panic(courseTagException.CourseKindIsNotExsit())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	params.Diff(kind)
	kind.Name = params.Str("name","名称")
	db.Driver.Save(&kind)

	ctx.JSON(iris.Map{
		"id":kid,
	})
}

func DeleteKind(ctx iris.Context, auth authbase.AuthAuthorization,kid int){
	auth.CheckAdmin()

	var kind db.Kind
	if err := db.Driver.GetOne("kind",kid,&kind);err == nil{
		db.Driver.Delete(kind)

		db.Driver.Exec("delete from course_and_kind where kind_id = ?",kid)
	}

	ctx.JSON(iris.Map{
		"id":kid,
	})
}

func ListKind(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckLogin()

	var lists []struct {
		Id         int   `json:"id"`
		Name string `json:"name"`
	}
	var count int

	table := db.Driver.Table("kind")

	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)

	table.Count(&count).Offset((page - 1) * limit).Limit(limit).Select("id, name").Find(&lists)
	ctx.JSON(iris.Map{
		"courseKinds": lists,
		"total": count,
		"limit": limit,
		"page":  page,
	})

}

func MgetKind(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "id列表")

	data := make([]interface{}, 0, len(ids))
	kinds := db.Driver.GetMany("kind", ids, db.Kind{})
	for _,k  := range kinds {
		func(data *[]interface{}) {
			*data = append(*data, paramsUtils.ModelToDict(k,[]string{"ID","Name"}))
			defer func() {
				recover()
			}()
		}(&data)
	}
	ctx.JSON(data)
}