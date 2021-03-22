package goods_tag

import (
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	goodsTagException "grpc-demo/exceptions/goods_tag"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

//优化方案
//加多一个level属性
func CreateKindTag(ctx iris.Context, auth authbase.AuthAuthorization) {
	auth.CheckAdmin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	title := params.Str("title", "标题")
	picture := params.Str("picture","图片")

	tag := db.KindTag{
		Title: title,
		Picture:picture,
	}

	if params.Has("parent_id") {
		parentId := params.Int("parent_id", "父种类id")
		var kindTag db.KindTag
		if err := db.Driver.GetOne("kind_tag", parentId, &kindTag); err != nil {
			panic(goodsTagException.KindTagIsNotExists())
		} else {
			if kindTag.ParentID != 0 {
				panic(goodsTagException.KindTagOverLevel())
			}
			tag.ParentID = parentId
		}
	}

	db.Driver.Create(&tag)

	ctx.JSON(iris.Map{
		"id": tag.ID,
	})
}

func PutKindTag(ctx iris.Context, auth authbase.AuthAuthorization, kid int) {
	auth.CheckAdmin()

	var tag db.KindTag
	if err := db.Driver.GetOne("kind_tag", kid, &tag); err != nil {
		panic(goodsTagException.KindTagIsNotExists())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	params.Diff(tag)
	tag.Title = params.Str("title", "标题")
	tag.Picture = params.Str("picture","图片")


	if params.Has("parent_id") {
		parentId := params.Int("parent_id", "父种类id")

		//标签下有子标签，层数已经达到2层
		var t db.KindTag
		if err := db.Driver.Where("parent_id = ? and id != ?", kid,kid).First(&t).Error; err == nil{

			panic(goodsTagException.KindTagOverLevel())
		}

		var kindTag db.KindTag
		if err := db.Driver.GetOne("kind_tag", parentId, &kindTag); err != nil {
			panic(goodsTagException.KindTagIsNotExists())
		} else {
			if kindTag.ParentID != 0 {
				panic(goodsTagException.KindTagOverLevel())
			}
			tag.ParentID = parentId
		}
	}

	db.Driver.Save(&tag)
	ctx.JSON(iris.Map{
		"id": tag.ID,
	})

}

func DeleteKindTag(ctx iris.Context, auth authbase.AuthAuthorization, kid int) {
	auth.CheckAdmin()

	var tag db.KindTag
	if err := db.Driver.GetOne("kind_tag", kid, &tag); err == nil {
		tx := db.Driver.Begin()
		if err := tx.Delete(&tag).Error; err != nil {
			tx.Rollback()
			panic(goodsTagException.DeleteKindTagFail())
		}

		if err := tx.Where("parent_id = ?", kid).Delete(db.KindTag{}).Error; err != nil {
			tx.Rollback()
			panic(goodsTagException.DeleteKindTagFail())
		}

		//TODO 删除相应挂载记录

		tx.Commit()
	}

	ctx.JSON(iris.Map{
		"id": kid,
	})
}

func ListKindTag(ctx iris.Context, auth authbase.AuthAuthorization) {
	auth.CheckLogin()

	var lists []struct {
		Id    int    `json:"id"`
		Title string `json:"title"`
	}
	var count int

	table := db.Driver.Table("kind_tag")

	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)

	if parentId := ctx.URLParamIntDefault("parent_id", -1); parentId != -1 {
		table = table.Where("parent_id = ?", parentId)
	}

	table.Count(&count).Offset((page - 1) * limit).Limit(limit).Select("id, title").Find(&lists)
	ctx.JSON(iris.Map{
		"tags":  lists,
		"total": count,
		"limit": limit,
		"page":  page,
	})
}

func MgetKindTag(ctx iris.Context, auth authbase.AuthAuthorization) {
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "id列表")

	data := make([]interface{}, 0, len(ids))
	tags := db.Driver.GetMany("kind_tag", ids, db.KindTag{})
	for _, tag := range tags {
		func(data *[]interface{}) {
			*data = append(*data, paramsUtils.ModelToDict(tag, []string{"ID", "ParentID", "Title","Picture"}))
			defer func() {
				recover()
			}()
		}(&data)
	}
	ctx.JSON(data)
}
