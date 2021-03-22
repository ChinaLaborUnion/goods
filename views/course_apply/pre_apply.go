package courseApply


import (

	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	courseEntity "grpc-demo/entity/course"
	courseApplyEnums "grpc-demo/enums/course_apply"
	accountException "grpc-demo/exceptions/account"
	courseException "grpc-demo/exceptions/course"
	courseApplyException "grpc-demo/exceptions/course_apply"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

var field = []string{
	"ID","AccountID","CourseID","Phone","Name","Status","CreateTime","UpdateTime","People","SessionID",
}

func PreApplyCreate(ctx iris.Context, auth authbase.AuthAuthorization,cid int){
	auth.CheckLogin()

	var course db.Course
	if err := db.Driver.GetOne("course", cid, &course); err != nil {
		panic(courseException.CourseIsNotExsit())
	}
	if !course.IsPut {
		panic(courseException.CourseNotPut())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	phone := params.Str("phone","号码")
	name := params.Str("name","姓名")
	people := params.Int("people","人数")
	sessionID := params.Int("session_id","场次id")

	courseInfo := course.GetInfo()["session"]
	var ok bool
	ok = false
	for _, c := range courseInfo.(courseEntity.Session) {
		if c.ID == sessionID {
			ok = true
			break
		}
	}
	if !ok {
		panic(courseApplyException.SessionIsNotExsit())
	}

	preApply := db.PreApply{
		AccountID:  auth.AccountModel().Id,
		CourseID:   cid,
		Phone:      phone,
		Name:       name,
		Status:    courseApplyEnums.NoApply,
		People:people,
		SessionID:sessionID,
	}
	db.Driver.Create(&preApply)

	ctx.JSON(iris.Map{
		"id":preApply.ID,
	})

}

func PreApplyList(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckLogin()

	var lists []struct {
		Id    int    `json:"id"`
		UpdateTime int64 `json:"update_time"`
	}
	var count int

	table := db.Driver.Table("pre_apply")

	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)

	if !auth.IsAdmin(){
		table = table.Where("account_id = ?", auth.AccountModel().Id)
	}

	//所属人过滤
	if author := ctx.URLParamIntDefault("author_id", 0); author != 0 && auth.IsAdmin() {
		table = table.Where("account_id = ?", author)
	}

	//课程过滤
	if courseID := ctx.URLParamIntDefault("course_id", 0); courseID != 0{
		table = table.Where("course_id = ?", courseID)
	}

	//是否报名过滤
	if status := ctx.URLParamIntDefault("status", -1); status != -1{
		table = table.Where("status = ?", status)
	}

	table.Count(&count).Order("create_time desc").Offset((page - 1) * limit).Limit(limit).Select("id, update_time").Find(&lists)
	ctx.JSON(iris.Map{
		"preApplys": lists,
		"total": count,
		"limit": limit,
		"page":  page,
	})

}

func PreApplyMget(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "id列表")

	data := make([]interface{}, 0, len(ids))
	preApplys := db.Driver.GetMany("pre_apply", ids, db.PreApply{})

	for _,p := range preApplys{
		if p.(db.PreApply).AccountID != auth.AccountModel().Id && !auth.IsAdmin(){
					continue
		}

		func(data *[]interface{}) {
			*data = append(*data, paramsUtils.ModelToDict(p,field))
			defer func() {
				recover()
			}()
		}(&data)
	}

	ctx.JSON(data)
}
//不能修改课程id
func PreApplyPut(ctx iris.Context, auth authbase.AuthAuthorization,pid int){
	auth.CheckLogin()

	var preApply db.PreApply
	if err := db.Driver.GetOne("pre_apply",pid,&preApply);err != nil{
		panic(courseApplyException.PreApplyIsNotExsit())
	}

	if preApply.Status != courseApplyEnums.NoApply{
		panic(courseApplyException.StatusPutFail())
	}

	if preApply.AccountID != auth.AccountModel().Id{
		panic(accountException.NoPermission())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	params.Diff(preApply)
	preApply.Phone = params.Str("phone","号码")
	preApply.Name = params.Str("name","姓名")
	preApply.People = params.Int("people","人数")


	if params.Has("session_id"){
		var course db.Course
		db.Driver.GetOne("course",preApply.CourseID,&course)

		preApply.SessionID = params.Int("session_id","场次id")
		courseInfo := course.GetInfo()["session"]
		var ok bool
		ok = false
		for _, c := range courseInfo.(courseEntity.Session) {
			if c.ID == preApply.SessionID {
				ok = true
				break
			}
		}
		if !ok {
			panic(courseApplyException.SessionIsNotExsit())
		}
	}


	db.Driver.Save(&preApply)

	ctx.JSON(iris.Map{
		"id":preApply.ID,
	})
}

func PreApplyCancel(ctx iris.Context, auth authbase.AuthAuthorization,pid int){
	auth.CheckLogin()

	var preApply db.PreApply
	if err := db.Driver.GetOne("pre_apply",pid,&preApply);err != nil{
		panic(courseApplyException.PreApplyIsNotExsit())
	}

	if preApply.Status != courseApplyEnums.NoApply{
		panic(courseApplyException.StatusPutFail())
	}

	if preApply.AccountID != auth.AccountModel().Id{
		panic(accountException.NoPermission())
	}

	preApply.Status = courseApplyEnums.PreApplyCancel

	db.Driver.Save(&preApply)

	ctx.JSON(iris.Map{
		"id":preApply.ID,
	})

}

