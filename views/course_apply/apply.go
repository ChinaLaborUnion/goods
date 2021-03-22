package courseApply

import (
	"github.com/Masterminds/squirrel"
	"github.com/kataras/iris"
	"github.com/shoogoome/mutils/hash"
	authbase "grpc-demo/core/auth"
	courseEntity "grpc-demo/entity/course"
	courseApplyEnums "grpc-demo/enums/course_apply"
	accountException "grpc-demo/exceptions/account"
	courseException "grpc-demo/exceptions/course"
	courseApplyException "grpc-demo/exceptions/course_apply"
	"grpc-demo/models/db"
	logUtils "grpc-demo/utils/log"
	paramsUtils "grpc-demo/utils/params"
	"strconv"
	"strings"
	"time"
)

func ApplyDirectCreate(ctx iris.Context, auth authbase.AuthAuthorization, cid int) {
	auth.CheckLogin()

	var course db.Course
	if err := db.Driver.GetOne("course", cid, &course); err != nil {
		panic(courseException.CourseIsNotExsit())
	}
	if !course.IsPut {
		panic(courseException.CourseNotPut())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	sessionID := params.Int("session_id", "场次")
	parters := params.List("parters", "参与人")
	people := params.Int("people", "人数")
	number := params.Str("number", "身份证号码")
	if len(parters) != people {
		panic(courseApplyException.PeopleError())
	}
	//todo 格式校验

	courseInfo := course.GetInfo()["session"]
	var ok bool
	ok = false
	var price int
	var limit int

	for _, c := range courseInfo.(courseEntity.Session) {
		if c.ID == sessionID {
			ok = true
			price = c.Money
			limit = c.PeopleLimit
			break
		}
	}
	if !ok {
		panic(courseApplyException.SessionIsNotExsit())
	}

	var applys []db.Apply
	db.Driver.Where("course_id = ? and session_id = ? and status not in (?)", cid, sessionID, []int16{courseApplyEnums.Cancel,courseApplyEnums.AlreadyAS}).Find(&applys)
	var total int
	for _, a := range applys {
		total += a.People
	}
	if people+total > limit {
		panic(courseApplyException.PeopleOverLimit())
	}

	tx := db.Driver.Begin()

	apply := db.Apply{
		AccountID:  auth.AccountModel().Id,
		CourseID:   cid,
		People:     people,
		TotalMoney: price * people,
		SessionID:  sessionID,
		Status:     courseApplyEnums.NoPay,
		IsPreApply: false,
		Number:number,
	}

	if err := tx.Create(&apply).Error; err != nil {
		tx.Rollback()
		panic(courseApplyException.ApplyCreateFail())
	}
	apply.OutTradeNo = strconv.FormatInt(apply.CreateTime, 10) + "-" + hash.GetRandomString(8)
	if err := tx.Save(&apply).Error; err != nil {
		tx.Rollback()
		panic(courseApplyException.ApplyCreateFail())
	}

	//todo 报名本人是否计算在内
	applyAndParter := db.ApplyAndParter{
		ApplyID:   apply.ID,
		AccountID: auth.AccountModel().Id,
		Number:    number,
	}
	if err := tx.Create(&applyAndParter).Error; err != nil {
		tx.Rollback()
		panic(courseApplyException.RelationCreateFail())
	}

	sql := squirrel.Insert("parter_info").Columns(
		"apply_id", "name", "phone", "number", "sex", "birth",
	)
	sql1 := squirrel.Insert("apply_and_parter").Columns(
		"apply_id", "account_id", "number",
	)

	var a db.ApplyAndParter
	for _, parter := range parters {
		p := parter.(map[string]interface{})
		//如果有报名记录，则直接加入
		if err := db.Driver.Where("number = ?", p["number"]).First(&a).Error; err == nil {
			sql1 = sql1.Values(
				apply.ID,
				a.AccountID,
				p["number"],
			)
		}
		//加入信息表
		sql = sql.Values(
			apply.ID,
			p["name"],
			p["phone"],
			p["number"],
			p["sex"],
			p["birth"],
		)
	}

	if s, args, err := sql.ToSql(); err != nil {
		tx.Rollback()
		logUtils.Println(err)
	} else {
		if err := db.Driver.Exec(s, args...).Error; err != nil {
			tx.Rollback()
			logUtils.Println(err)
			return
		}
	}

	if s, args, err := sql1.ToSql(); err != nil {
		tx.Rollback()
		logUtils.Println(err)
	} else {
		if err := db.Driver.Exec(s, args...).Error; err != nil {
			tx.Rollback()
			logUtils.Println(err)
			return
		}
	}

	tx.Commit()

	ctx.JSON(iris.Map{
		"id": apply.ID,
	})

}

func ApplyCreate(ctx iris.Context, auth authbase.AuthAuthorization, pid int) {
	auth.CheckLogin()

	var preApply db.PreApply
	if err := db.Driver.GetOne("pre_apply", pid, &preApply); err != nil {
		panic(courseApplyException.PreApplyIsNotExsit())
	}

	if auth.AccountModel().Id != preApply.AccountID {
		panic(accountException.NoPermission())
	}

	if preApply.Status != courseApplyEnums.NoApply {
		panic(courseApplyException.StatusApplyFail())
	}

	var course db.Course
	if err := db.Driver.GetOne("course", preApply.CourseID, &course); err != nil {
		panic(courseException.CourseIsNotExsit())
	}
	if !course.IsPut {
		panic(courseException.CourseNotPut())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	sessionID := params.Int("session_id", "场次")
	people := params.Int("people", "人数")
	parters := params.List("parters", "参与人")
	number := params.Str("number", "身份证号码")
	//todo 格式校验

	if len(parters) != people {
		panic(courseApplyException.PeopleError())
	}

	courseInfo := course.GetInfo()["session"]
	var ok bool
	ok = false
	var price int
	var limit int

	for _, c := range courseInfo.(courseEntity.Session) {
		if c.ID == sessionID {
			ok = true
			price = c.Money
			limit = c.PeopleLimit
			break
		}
	}
	if !ok {
		panic(courseApplyException.SessionIsNotExsit())
	}

	var applys []db.Apply
	db.Driver.Where("course_id = ? and session_id = ? and status not in (?)", preApply.CourseID, sessionID, []int16{courseApplyEnums.Cancel,courseApplyEnums.AlreadyAS}).Find(&applys)
	var total int
	for _, a := range applys {
		total += a.People
	}
	if people+total > limit {
		panic(courseApplyException.PeopleOverLimit())
	}

	tx := db.Driver.Begin()

	apply := db.Apply{
		AccountID:  auth.AccountModel().Id,
		CourseID:   preApply.CourseID,
		PreApplyID: preApply.ID,
		People:     people,
		TotalMoney: price * people,
		SessionID:  sessionID,
		Status:     courseApplyEnums.NoPay,
		IsPreApply: true,
		Number:number,
	}

	if err := tx.Create(&apply).Error; err != nil {
		tx.Rollback()
		panic(courseApplyException.ApplyCreateFail())
	}
	apply.OutTradeNo = strconv.FormatInt(apply.CreateTime, 10) + "-" + hash.GetRandomString(8)
	if err := tx.Save(&apply).Error; err != nil {
		tx.Rollback()
		panic(courseApplyException.ApplyCreateFail())
	}

	applyAndParter := db.ApplyAndParter{
		ApplyID:   apply.ID,
		AccountID: auth.AccountModel().Id,
		Number:    number,
	}
	if err := tx.Create(&applyAndParter).Error; err != nil {
		tx.Rollback()
		panic(courseApplyException.RelationCreateFail())
	}

	preApply.Status = courseApplyEnums.Apply
	if err := tx.Save(&preApply).Error; err != nil {
		tx.Rollback()
		panic(courseApplyException.ApplyCreateFail())
	}

	sql := squirrel.Insert("parter_info").Columns(
		"apply_id", "name", "phone", "number", "sex", "birth",
	)
	sql1 := squirrel.Insert("apply_and_parter").Columns(
		"apply_id", "account_id", "number",
	)

	var a db.ApplyAndParter
	for _, parter := range parters {
		p := parter.(map[string]interface{})
		if err := db.Driver.Where("number = ?", p["number"]).First(&a).Error; err == nil {
			sql1 = sql1.Values(
				apply.ID,
				a.AccountID,
				p["number"],
			)
		}
		sql = sql.Values(
			apply.ID,
			p["name"],
			p["phone"],
			p["number"],
			p["sex"],
			p["birth"],
		)
	}

	if s, args, err := sql.ToSql(); err != nil {
		tx.Rollback()
		logUtils.Println(err)
	} else {
		if err := db.Driver.Exec(s, args...).Error; err != nil {
			tx.Rollback()
			logUtils.Println(err)
			return
		}
	}

	tx.Commit()

	ctx.JSON(iris.Map{
		"id": apply.ID,
	})

}

func ApplyList(ctx iris.Context, auth authbase.AuthAuthorization) {
	auth.CheckLogin()

	//没有报名记录时，去信息表找所有的报名记录
	//下次登陆时，就会有报名记录
	var applyAndParter db.ApplyAndParter
	if err := db.Driver.Where("account_id = ?", auth.AccountModel().Id).First(&applyAndParter).Error; err != nil && !auth.IsAdmin() {
		params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
		number := params.Str("number", "身份证号码", "")

		// 第一次点进来，此时还没有给身份证号码， 后端检测需要时会告诉前端
		// 前端渲染身份证号码输入框
		if strings.Trim(number, " ") == "" {
			panic(courseApplyException.NeedNumber())
		}

		var parterInfo []db.ParterInfo
		if err := db.Driver.Where("number = ?", number).Find(&parterInfo).Error; err != nil {
			v := db.ApplyAndParter{
				ApplyID:   0,
				AccountID: auth.AccountModel().Id,
				Number:    number,
			}
			db.Driver.Create(&v)
		} else {
			applyIDs := make([]int, 0)
			for _, v := range parterInfo {
				applyIDs = append(applyIDs, v.ApplyID)
			}

			sql := squirrel.Insert("apply_and_parter").Columns(
				"apply_id", "account_id", "number",
			)

			for index, v := range parterInfo {
				sql = sql.Values(
					applyIDs[index],
					auth.AccountModel().Id,
					v.Number,
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

	var lists []struct {
		Id         int   `json:"id"`
		UpdateTime int64 `json:"update_time"`
	}
	var count int

	table := db.Driver.Table("apply")

	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)

	if !auth.IsAdmin() {
		var applyAndParter []db.ApplyAndParter
		db.Driver.Where("account_id = ? and apply_id != ?", auth.AccountModel().Id, 0).Find(&applyAndParter)
		applyIDs := make([]int, 0)
		for _, v := range applyAndParter {
			applyIDs = append(applyIDs, v.ApplyID)
		}
		table = table.Where("id in (?)", applyIDs)
	}

	//所属人过滤
	if author := ctx.URLParamIntDefault("author_id", 0); author != 0 && auth.IsAdmin() {
		var applyAndParter []db.ApplyAndParter
		db.Driver.Where("account_id = ? and apply_id != ?", author, 0).Find(&applyAndParter)
		applyIDs := make([]int, 0)
		for _, v := range applyAndParter {
			applyIDs = append(applyIDs, v.ApplyID)
		}
		table = table.Where("id in (?)", applyIDs)
		//table = table.Where("account_id = ?", author)
	}

	//课程过滤
	if courseID := ctx.URLParamIntDefault("course_id", 0); courseID != 0 {
		table = table.Where("course_id = ?", courseID)
	}

	//状态过滤
	if status := ctx.URLParamIntDefault("status", 0); status != 0 {
		if status == 7{
			table = table.Where("status in (?)", []int16{courseApplyEnums.AfterSales,courseApplyEnums.AlreadyAS})
		}else {
			table = table.Where("status = ?", status)
		}
	}

	table.Count(&count).Order("create_time desc").Offset((page - 1) * limit).Limit(limit).Select("id, update_time").Find(&lists)
	ctx.JSON(iris.Map{
		"applys": lists,
		"total":  count,
		"limit":  limit,
		"page":   page,
	})

}

func ApplyMget(ctx iris.Context, auth authbase.AuthAuthorization) {
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "id列表")

	data := make([]interface{}, 0, len(ids))
	applys := db.Driver.GetMany("apply", ids, db.Apply{})

	for _, apply := range applys {
		if apply.(db.Apply).AccountID != auth.AccountModel().Id && !auth.IsAdmin() {
			continue
		}

		func(data *[]interface{}) {
			*data = append(*data, apply.(db.Apply).GetInfo())
			defer func() {
				recover()
			}()
		}(&data)
	}

	ctx.JSON(data)
}

//只能修改场次
func ApplyPut(ctx iris.Context, auth authbase.AuthAuthorization, aid int) {
	auth.CheckLogin()

	var apply db.Apply
	if err := db.Driver.GetOne("apply", aid, &apply); err != nil {
		panic(courseApplyException.ApplyIsNotExsit())
	}

	if apply.Status != courseApplyEnums.NoPay {
		panic(courseApplyException.ApplyPutFail())
	}

	if apply.AccountID != auth.AccountModel().Id {
		panic(accountException.NoPermission())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	params.Diff(apply)
	if params.Has("session_id") {
		sessionID := params.Int("session_id", "场次")
		var course db.Course
		db.Driver.GetOne("course", apply.CourseID, &course)
		courseInfo := course.GetInfo()["session"]
		var ok bool
		ok = false
		var price int
		var limit int

		for _, c := range courseInfo.(courseEntity.Session) {
			if c.ID == sessionID {
				ok = true
				price = c.Money
				limit = c.PeopleLimit
				break
			}
		}
		if !ok {
			panic(courseApplyException.SessionIsNotExsit())
		}

		var applys []db.Apply
		db.Driver.Where("course_id = ? and session_id = ? and status != ? and id != ?", apply.CourseID, sessionID, courseApplyEnums.Cancel, apply.ID).Find(&applys)
		var total int
		for _, a := range applys {
			total += a.People
		}
		if apply.People+total > limit {
			panic(courseApplyException.PeopleOverLimit())
		}

		apply.SessionID = sessionID
		apply.TotalMoney = price * apply.People
	}

	db.Driver.Save(&apply)

	ctx.JSON(iris.Map{
		"id": apply.ID,
	})
}

func ApplyCancel(ctx iris.Context, auth authbase.AuthAuthorization, aid int) {
	auth.CheckLogin()

	var apply db.Apply
	if err := db.Driver.GetOne("apply", aid, &apply); err != nil {
		panic(courseApplyException.ApplyIsNotExsit())
	}

	if apply.Status != courseApplyEnums.NoPay {
		panic(courseApplyException.ApplyPutFail())
	}

	if apply.AccountID != auth.AccountModel().Id {
		panic(accountException.NoPermission())
	}

	apply.Status = courseApplyEnums.Cancel
	db.Driver.Save(&apply)

	ctx.JSON(iris.Map{
		"id": apply.ID,
	})
}

func BackMoney(ctx iris.Context, auth authbase.AuthAuthorization,aid int){
	auth.CheckLogin()

	var apply db.Apply
	if err := db.Driver.GetOne("apply", aid, &apply); err != nil {
		panic(courseApplyException.ApplyIsNotExsit())
	}

	if apply.AccountID != auth.AccountModel().Id{
		panic(accountException.NoPermission())
	}

	if apply.Status != courseApplyEnums.Paied{
		panic(courseApplyException.StatusNotBackMoney())
	}

	var course db.Course
	db.Driver.GetOne("course",apply.CourseID,&course)
	courseInfo := course.GetInfo()["session"]
	for _, c := range courseInfo.(courseEntity.Session) {
		if c.ID == apply.SessionID{
			if time.Now().Unix() > c.BeginTime + 86400{
				panic(courseApplyException.TimeNotBackMoney())
			}
		}
	}

	apply.Status = courseApplyEnums.AfterSales
	db.Driver.Save(&apply)

	ctx.JSON(iris.Map{
		"id":apply.ID,
	})
}

