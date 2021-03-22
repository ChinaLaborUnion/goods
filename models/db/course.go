package db

import (
	courseEntity "grpc-demo/entity/course"
	courseEnums "grpc-demo/enums/course"
	paramsUtils "grpc-demo/utils/params"
)


//课程
type Course struct{
	ID int `gorm:"primary_key" json:"id"`

	//课程名
	Name string `json:"name" gorm:"not null"`

	//课程子标题
	SmallName string `json:"small_name"`

	//课程简介
	Describe string `json:"describe" `

	//特色
	Feature string `json:"feature"`

	//课程详情
	Detail string `json:"detail" gorm:"type: text"`

	//注意事项
	Attention string `json:"attention"`

	//适合人群
	Crowd string `json:"crowd"`
	
	//报名人数
	People int `json:"people"`

	//课程安排
	Plan string `json:"plan" gorm:"not null"`

	//课程耗时(以小时为单位)
	Time int `json:"time" gorm:"not null"`

	//封面
	Cover string `json:"cover" gorm:"not null"`

	//最低价钱
	MinPrice int `json:"min_price"`

	//课程开始时间
	BeginTime int64 `json:"begin_time" gorm:"not null"`

	//课程结束时间
	EndTime int64 `json:"end_time" gorm:"not null"`

	//发布状态
	IsPut bool `json:"is_put"`

	//场次
	Session string `json:"session"`

	//图片
	Pictures string `json:"pictures"`

	//是否删除
	IsDelete bool `json:"is_delete"`

	// 创建时间
	CreateTime int64 `json:"create_time"`

	// 更新时间
	UpdateTime int64 `json:"update_time"`
}

//课程标签挂载
type CourseAndTag struct{
	ID int `gorm:"primary_key" json:"id"`

	//课程id
	CourseID int `json:"course_id" gorm:"not null;index"`

	//营地标签id
	TagID int `json:"tag_id" gorm:"not null;index"`

	//标签类型
	TagType int16 `json:"tag_type"`
}

//种类
type Kind struct{
	ID int `gorm:"primary_key" json:"id"`

	//名称
	Name string `json:"name" gorm:"not null"`
}

//课程挂载种类
type CourseAndKind struct {
	ID int `gorm:"primary_key" json:"id"`

	//课程id
	CourseID int `json:"course_id" gorm:"not null;index"`

	//类型id
	KindID int `json:"kind_id" gorm:"not null;index"`
}

//课程标签
type CourseTag struct{
	ID int `gorm:"primary_key" json:"id"`

	//名称
	Name string `json:"name" gorm:"not null"`
}


var courseField = []string{
	"ID","Name","SmallName","Describe","Feature","Detail","Attention","Crowd","Plan","Time","Cover","MinPrice","BeginTime","EndTime","IsPut","CreateTime","UpdateTime",
}

func (c Course) GetInfo() map[string]interface{}{
	v := paramsUtils.ModelToDict(c, courseField)

	session := courseEntity.CourseSessionUnmarshal(c.Session)

	v["session"] = session

	tt := courseEntity.CoursePictureUnmarshal(c.Pictures)

	v["pictures"] = tt

	var placeTag []CourseAndTag
	if err := Driver.Where("course_id = ? and tag_type = ?", c.ID, courseEnums.PlaceType).Find(&placeTag).Error; err == nil {
		place := make([]int,len(placeTag))
		for index,i := range placeTag{
			place[index] = i.TagID
		}
		v["place_tag"] = place
	}

	var courseTag []CourseAndTag
	if err := Driver.Where("course_id = ? and tag_type = ?", c.ID, courseEnums.CourseType).Find(&courseTag).Error; err == nil {
		tags := make([]int,len(courseTag))
		for index,i := range courseTag{
			tags[index] = i.TagID
		}
		v["course_tag"] = tags
	}

	var kind []CourseAndKind
	if err := Driver.Where("course_id = ? ", c.ID).Find(&kind).Error; err == nil {
		k := make([]int,len(kind))
		for index,i := range kind{
			k[index] = i.KindID
		}
		v["kind"] = k
	}

	return v
}

