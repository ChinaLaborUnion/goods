package course

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	courseEntity "grpc-demo/entity/course"
	courseEnums "grpc-demo/enums/course"
	courseException "grpc-demo/exceptions/course"
	"grpc-demo/models/db"
	logUtils "grpc-demo/utils/log"
	paramsUtils "grpc-demo/utils/params"
)

func Create(ctx iris.Context, auth authbase.AuthAuthorization) {
	auth.CheckAdmin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	name := params.Str("name", "名称")
	describe := params.Str("describe", "简介")
	plan := params.Str("plan", "课程安排")
	time := params.Int("time", "课程耗时")
	detail := params.Str("detail", "详情")
	cover := params.Str("cover", "封面图片")
	beginTime := params.Int("begin_time", "课程开始时间")
	endTime := params.Int("end_time", "课程结束时间")
	isPut := params.Bool("is_put", "是否发布")

	tx := db.Driver.Begin()

	course := db.Course{
		Name:      name,
		Describe:  describe,
		Detail:    detail,
		Plan:      plan,
		Time:      time,
		Cover:     cover,
		BeginTime: int64(beginTime),
		EndTime:   int64(endTime),
		IsPut:     isPut,
		IsDelete:  false,
	}

	if err := tx.Debug().Create(&course).Error; err != nil {
		tx.Rollback()
		panic(courseException.CourseCreateFail())
	}

	//[{"session":,"money":,"people_limit":},{}]
	session := params.List("session", "场次")
	minPrice := int(session[0].(map[string]interface{})["money"].(float64))

	for _, v1 := range session {
		if int(v1.(map[string]interface{})["money"].(float64)) < minPrice {
			minPrice = int(v1.(map[string]interface{})["money"].(float64))
		}
	}
	course.MinPrice = minPrice

	for index, v := range session {
		v.(map[string]interface{})["id"] = index + 1
	}

	course.Session = courseEntity.CourseSessionMarshal(session)

	if params.Has("small_name") {
		smallName := params.Str("small_name", "子标题")
		course.SmallName = smallName
	}

	if params.Has("feature") {
		feature := params.Str("feature", "特点")
		course.Feature = feature
	}

	if params.Has("attention") {
		attention := params.Str("attention", "注意事项")
		course.Attention = attention
	}

	if params.Has("crowd") {
		crowd := params.Str("crowd", "适合人群")
		course.Crowd = crowd
	}

	if params.Has("pictures"){
		pictures := params.List("pictures","图片列表")
		course.Pictures = courseEntity.CoursePictureMarshal(pictures)
	}

	if err := tx.Debug().Save(&course).Error; err != nil {
		tx.Rollback()
		panic(courseException.CourseCreateFail())
	}

	tx.Commit()

	if params.Has("place_tag") {
		placeTag := params.List("place_tag", "属地标签")
		TagMnt(course, placeTag, courseEnums.PlaceType)
	}

	if params.Has("course_tag") {
		courseTag := params.List("course_tag", "课程标签")
		TagMnt(course, courseTag, courseEnums.CourseType)
	}

	if params.Has("kind") {
		kind := params.List("kind", "课程类型")
		KindMnt(course, kind)
	}

	ctx.JSON(iris.Map{
		"id": course.ID,
	})
}

func Put(ctx iris.Context, auth authbase.AuthAuthorization, cid int) {
	auth.CheckAdmin()

	var course db.Course
	if err := db.Driver.GetOne("course", cid, &course); err != nil {
		panic(courseException.CourseIsNotExsit())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	params.Diff(course)

	course.Name = params.Str("name", "名称")
	course.Describe = params.Str("describe", "简介")
	course.Feature = params.Str("feature", "特色")
	course.Attention = params.Str("attention", "注意事项")
	course.Crowd = params.Str("crowd", "适合人群")
	course.Plan = params.Str("plan", "课程安排")
	course.Time = params.Int("time", "课程耗时")
	course.SmallName = params.Str("small_name", "子标题")
	course.Detail = params.Str("detail", "课程详情")
	course.Cover = params.Str("cover", "封面图片")
	course.BeginTime = int64(params.Int("begin_time", "课程开始时间"))
	course.EndTime = int64(params.Int("end_time", "课程结束时间"))
	if params.Has("is_put") {
		course.IsPut = params.Bool("is_put", "是否发布")
	}

	if params.Has("pictures"){
		pictures := params.List("pictures","图片列表")
		course.Pictures = courseEntity.CoursePictureMarshal(pictures)
	}

	if params.Has("place_tag") {
		placeTag := params.List("place_tag", "属地标签")
		db.Driver.Exec("delete from course_and_tag where course_id = ? and tag_type = ?", cid, courseEnums.PlaceType)
		TagMnt(course, placeTag, courseEnums.PlaceType)
	}

	if params.Has("course_tag") {
		courseTag := params.List("course_tag", "课程标签")
		db.Driver.Exec("delete from course_and_tag where course_id = ? and tag_type = ?", cid, courseEnums.CourseType)
		TagMnt(course, courseTag, courseEnums.CourseType)
	}

	if params.Has("kind") {
		kind := params.List("kind", "课程类型")
		db.Driver.Exec("delete from course_and_kind where course_id = ? ", cid)
		KindMnt(course, kind)
	}

	if params.Has("session") {
		session := params.List("session", "场次")
		minPrice := int(session[0].(map[string]interface{})["money"].(float64))

		for _, v1 := range session {
			if int(v1.(map[string]interface{})["money"].(float64)) < minPrice {
				minPrice = int(v1.(map[string]interface{})["money"].(float64))
			}
		}
		course.MinPrice = minPrice
		for index, v := range session {
			v.(map[string]interface{})["id"] = index + 1
		}

		course.Session = courseEntity.CourseSessionMarshal(session)

	}

	db.Driver.Save(&course)

	ctx.JSON(iris.Map{
		"id": course.ID,
	})
}

//todo 取消发布
func Delete(ctx iris.Context, auth authbase.AuthAuthorization, cid int) {
	auth.CheckAdmin()

	var course db.Course
	if err := db.Driver.GetOne("course", cid, &course); err == nil {

		//db.Driver.Exec("delete from course_and_tag where course_id = ?", cid)
		//
		//db.Driver.Exec("delete from course_and_kind where course_id = ?", cid)
		//
		//db.Driver.Delete(course)
		course.IsPut = false
		course.IsDelete = true
		db.Driver.Save(&course)

	}

	ctx.JSON(iris.Map{
		"id": cid,
	})
}

func List(ctx iris.Context, auth authbase.AuthAuthorization) {
	auth.CheckLogin()

	var lists []struct {
		Id         int   `json:"id"`
		UpdateTime int64 `json:"update_time"`
	}
	var count int

	table := db.Driver.Table("course")

	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)

	table = table.Where("is_delete = ?", false)

	if isPut := ctx.URLParamIntDefault("is_put", -1); isPut != -1 {
		table = table.Where("is_put = ?", isPut)
	}

	//名称过滤
	if name := ctx.URLParam("name"); len(name) > 0 {
		keyString := fmt.Sprintf("%%%s%%", name)
		table = table.Where("name like ?", keyString)
	}

	//适合人群过滤
	if crowd := ctx.URLParam("crowd"); len(crowd) > 0 {
		keyString := fmt.Sprintf("%%%s%%", crowd)
		table = table.Where("crowd like ?", keyString)
	}

	//课程标签过滤
	if courseTag := ctx.URLParamIntDefault("course_tag", 0); courseTag != 0 {
		var tag db.CourseTag
		if err := db.Driver.GetOne("course_tag", courseTag, &tag); err == nil {
			var courses []db.CourseAndTag

			if err := db.Driver.Where("tag_id = ? and tag_type = ?", courseTag, courseEnums.CourseType).Find(&courses).Error; err == nil {
				ids := make([]interface{}, len(courses))
				for index, v := range courses {
					ids[index] = v.CourseID
				}
				table = table.Where("id in (?)", ids)
			}
		}
	}

	//属地标签过滤
	if placeTag := ctx.URLParamIntDefault("place_tag", 0); placeTag != 0 {
		var tag db.PlaceTag
		if err := db.Driver.GetOne("place_tag", placeTag, &tag); err == nil {
			var courses []db.CourseAndTag

			if err := db.Driver.Where("tag_id = ? and tag_type = ?", placeTag, courseEnums.PlaceType).Find(&courses).Error; err == nil {
				ids := make([]interface{}, len(courses))
				for index, v := range courses {
					ids[index] = v.CourseID
				}
				table = table.Where("id in (?)", ids)
			}
		}
	}

	//课程类型过滤
	if kind := ctx.URLParamIntDefault("kind", 0); kind != 0 {
		var k db.Kind
		if err := db.Driver.GetOne("kind", kind, &k); err == nil {
			var courses []db.CourseAndKind

			if err := db.Driver.Where("kind_id = ?", kind).Find(&courses).Error; err == nil {
				ids := make([]interface{}, len(courses))
				for index, v := range courses {
					ids[index] = v.CourseID
				}
				table = table.Where("id in (?)", ids)
			}
		}
	}

	table.Count(&count).Order("create_time desc").Offset((page - 1) * limit).Limit(limit).Select("id, update_time").Find(&lists)
	ctx.JSON(iris.Map{
		"courses": lists,
		"total":   count,
		"limit":   limit,
		"page":    page,
	})
}

func Mget(ctx iris.Context, auth authbase.AuthAuthorization) {
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "id列表")

	data := make([]interface{}, 0, len(ids))
	courses := db.Driver.GetMany("course", ids, db.Course{})
	for _, c := range courses {
		func(data *[]interface{}) {
			*data = append(*data, c.(db.Course).GetInfo())
			defer func() {
				recover()
			}()
		}(&data)
	}
	ctx.JSON(data)
}

func TagMnt(course db.Course, t []interface{}, tagType int) {
	if tagType == courseEnums.PlaceType {
		var tags []db.PlaceTag
		if err := db.Driver.Where("id in (?)", t).Find(&tags).Error; err != nil || len(tags) == 0 {
			logUtils.Println(err)
			return
		}
		sql := squirrel.Insert("course_and_tag").Columns(
			"course_id", "tag_id", "tag_type",
		)

		for _, tag := range tags {
			sql = sql.Values(
				course.ID,
				tag.ID,
				tagType,
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
	} else {
		var tags []db.CourseTag
		if err := db.Driver.Where("id in (?)", t).Find(&tags).Error; err != nil || len(tags) == 0 {
			logUtils.Println(err)
			return
		}
		sql := squirrel.Insert("course_and_tag").Columns(
			"course_id", "tag_id", "tag_type",
		)

		for _, tag := range tags {
			sql = sql.Values(
				course.ID,
				tag.ID,
				tagType,
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

}

func KindMnt(course db.Course, k []interface{}) {
	var kinds []db.Kind
	if err := db.Driver.Where("id in (?)", k).Find(&kinds).Error; err != nil || len(kinds) == 0 {
		logUtils.Println(err)
		return
	}

	sql := squirrel.Insert("course_and_kind").Columns(
		"course_id", "kind_id",
	)

	for _, kind := range kinds {
		sql = sql.Values(
			course.ID,
			kind.ID,
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
