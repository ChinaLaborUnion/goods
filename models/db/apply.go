package db

import (
	paramsUtils "grpc-demo/utils/params"
)

//课程报名
type PreApply struct{
	ID int `gorm:"primary_key" json:"id"`

	//用户id
	AccountID int `json:"account_id" gorm:"not null;index"`

	//课程id
	CourseID int `json:"course_id" gorm:"not null;index"`

	//电话
	Phone string `json:"phone" gorm:"not null"`

	//姓名
	Name string `json:"name" gorm:"not null"`

	//人数
	People int `json:"people" gorm:"default: 1"`
	
	//场次
	SessionID int `json:"session_id"`

	//是否报名
	Status int16 `json:"status"`

	// 创建时间
	CreateTime int64 `json:"create_time"`

	// 更新时间
	UpdateTime int64 `json:"update_time"`
}

type Apply struct{
	ID int `gorm:"primary_key" json:"id"`

	//用户id
	AccountID int `json:"account_id" gorm:"not null;index"`

	//课程id
	CourseID int `json:"course_id" gorm:"not null;index"`

	//预报名id
	PreApplyID int `json:"pre_apply_id" gorm:"not null;index"`

	//报名人数
	People int `json:"people" gorm:"default: 1"`

	//总费用
	TotalMoney int `json:"total_money"`

	//场次
	SessionID int `json:"session_id"`

	//身份证号码
	Number string 	`json:"number" gorm:"not null"`

	//状态（已支付、已取消、未支付）
	Status int16 `json:"status"`

	//微信支付id
	WxPayOrderId int `json:"wx_pay_order_id" gorm:"not null;index"`
		
	//支付时间
	PayTime int64 `json:"pay_time"`
	
	//编号
	OutTradeNo string `json:"out_trade_no"`

	//是否预约
	IsPreApply bool `json:"is_pre_apply"`

	// 创建时间
	CreateTime int64 `json:"create_time"`

	// 更新时间
	UpdateTime int64 `json:"update_time"`

}

type ParterInfo struct{
	ID int `gorm:"primary_key" json:"id"`

	//报名id
	ApplyID int `json:"apply_id" gorm:"not null;index"`

	//姓名
	Name string `json:"name" gorm:"not null"`

	//号码
	Phone string `json:"phone" gorm:"not null"`

	//性别
	Sex int `json:"sex"`

	//出生年月
	Birth int64 `json:"birth"`

	//身份证号码
	Number string 	`json:"number" gorm:"not null"`
}

type ApplyAndParter struct {
	ID int `gorm:"primary_key" json:"id"`

	//报名id
	ApplyID int `json:"apply_id" gorm:"not null;index"`

	//账户id
	AccountID int `json:"account_id" gorm:"not null;index"`

	//身份证号码
	Number string 	`json:"number" gorm:"not null"`
}

var preApplyfield = []string{
	"ID","AccountID","CourseID","Phone","Name","Status","CreateTime","UpdateTime","People","SessionID",
}

var applyField = []string{
	"ID", "AccountID", "CourseID", "People", "TotalMoney", "SessionID", "Status", "WxPayOrderId", "IsPreApply","CreateTime", "UpdateTime","Number",
	"PayTime","OutTradeNo",
}

func (a Apply) GetInfo() map[string]interface{}{
	v := paramsUtils.ModelToDict(a, applyField)

	if a.IsPreApply{
		var preApply PreApply
		Driver.GetOne("pre_apply",a.PreApplyID,&preApply)

		v["pre_apply"] = paramsUtils.ModelToDict(preApply, preApplyfield)
	}

	var parters []ParterInfo
	Driver.Where("apply_id = ?",a.ID).Find(&parters)
	data := make([]interface{},len(parters))
	for index,parter := range parters{
		data[index] = paramsUtils.ModelToDict(parter, []string{"Name","Phone","Number","Sex","Birth"})
	}

	v["parters"] = data

	return v
}
