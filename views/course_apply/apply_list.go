package courseApply

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	courseApplyException "grpc-demo/exceptions/course_apply"
	"grpc-demo/models/db"
	logUtils "grpc-demo/utils/log"
	paramsUtils "grpc-demo/utils/params"
	"sort"

	"strings"

)

func ApplyAllList(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckLogin()

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




	//if !auth.IsAdmin() {
		var applyAndParter1 []db.ApplyAndParter
		db.Driver.Where("account_id = ? and apply_id != ?", auth.AccountModel().Id, 0).Find(&applyAndParter1)
		applyIDs := make([]int64, 0)
		for _, v := range applyAndParter1 {
			applyIDs = append(applyIDs, int64(v.ApplyID))
		}

		var v2 v1

		db.Driver.Table("pre_apply").Select("pre_apply.*, apply.*").Joins("left join apply on pre_apply.id = apply.pre_apply_id").Where("pre_apply.account_id = ? or apply.id in (?)",auth.AccountModel().Id,applyIDs).Find(&v2)


		for _,i := range v2{
			fmt.Println(i.PreApply.ID," ",i.Apply.ID," ",i.Apply.CreateTime," ",i.PreApply.CreateTime)
			if i.Apply.ID == 0{
				i.Apply.CreateTime = i.PreApply.CreateTime
				db.Driver.Save(&i.Apply)
			}
			fmt.Println(i.Apply.CreateTime)
		}
		fmt.Println()


		sort.Sort(v2)
		for _,i := range v2{
			fmt.Println(i.PreApply.ID," ",i.Apply.CreateTime)
		}

		fmt.Println()

}

type v1 []struct{
	db.PreApply
	db.Apply
}

func (a v1) Len() int {    // 重写 Len() 方法
	return len(a)
}
func (a v1) Swap(i, j int){     // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a v1) Less(i, j int) bool {    // 重写 Less() 方法， 从大到小排序
	return a[j].Apply.CreateTime < a[i].Apply.CreateTime
}
