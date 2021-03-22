package courseTag

import (
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	courseEnums "grpc-demo/enums/course"
	courseTagException "grpc-demo/exceptions/course_tag"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

func CreateTag(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckAdmin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	name := params.Str("name","名称")

	tag := db.CourseTag{
		Name: name,
	}

	db.Driver.Create(&tag)

	ctx.JSON(iris.Map{
		"id":tag.ID,
	})
}

func PutTag(ctx iris.Context, auth authbase.AuthAuthorization,tid int){
	auth.CheckAdmin()

	var tag db.CourseTag
	if err := db.Driver.GetOne("course_tag",tid,&tag);err != nil{
		panic(courseTagException.CourseTagIsNotExsit())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	params.Diff(tag)
	tag.Name = params.Str("name","名称")
	db.Driver.Save(&tag)

	ctx.JSON(iris.Map{
		"id":tid,
	})
}

func DeleteTag(ctx iris.Context, auth authbase.AuthAuthorization,tid int){
	auth.CheckAdmin()

	var tag db.CourseTag
	if err := db.Driver.GetOne("course_tag",tid,&tag);err == nil{
		db.Driver.Delete(tag)

		db.Driver.Exec("delete from course_and_tag where tag_id = ? and tag_type = ?",tid,courseEnums.CourseType)
	}

	ctx.JSON(iris.Map{
		"id":tid,
	})
}

func ListTag(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckLogin()

	var lists []struct {
		Id         int   `json:"id"`
		Name string `json:"name"`
	}
	var count int

	table := db.Driver.Table("course_tag")

	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)

	table.Count(&count).Offset((page - 1) * limit).Limit(limit).Select("id, name").Find(&lists)
	ctx.JSON(iris.Map{
		"courseTags": lists,
		"total": count,
		"limit": limit,
		"page":  page,
	})

}

func MgetTag(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "id列表")

	data := make([]interface{}, 0, len(ids))
	tags := db.Driver.GetMany("course_tag", ids, db.CourseTag{})
	for _,t  := range tags {
		func(data *[]interface{}) {
			*data = append(*data, paramsUtils.ModelToDict(t,[]string{"ID","Name"}))
			defer func() {
				recover()
			}()
		}(&data)
	}
	ctx.JSON(data)
}