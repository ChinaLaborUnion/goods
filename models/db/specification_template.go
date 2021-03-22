package db

import (
	"encoding/json"
	goods "grpc-demo/entity/goods_specification"
	goodsSpecificationException "grpc-demo/exceptions/goods_specification"
	paramsUtils "grpc-demo/utils/params"
)

//规格模版
type SpecificationTemplate struct {
	ID int `gorm:"primary_key" json:"id"`

	//模版名称
	Title string `json:"title" gorm:"not null"`

	//规格模版（一句话描述）
	Template string `json:"template" gorm:"type:text"`
}

var field = []string{
	"ID", "Title",
}

func (t SpecificationTemplate) GetInfo() map[string]interface{} {
	v := paramsUtils.ModelToDict(t, field)

	entity := new(goods.TemplateEntity)
	if err := json.Unmarshal([]byte(t.Template), entity); err != nil {
		panic(goodsSpecificationException.TemplateUnMarshalFail())
	}

	v["template"] = *entity

	return v
}
