package db

import (
	"encoding/json"
	goodsInfo "grpc-demo/entity/goods_info"
	goods "grpc-demo/entity/goods_specification"
	goodsTagEnum "grpc-demo/enums/goods_tag"
	goodsInfoException "grpc-demo/exceptions/goods_info"
	paramsUtils "grpc-demo/utils/params"
)

//TODO 分享 会员
//商品
type Goods struct {
	ID int `gorm:"primary_key" json:"id"`

	//最低价
	MinPrice int `json:"min_price"`

	//商品名
	Name string `json:"name" gorm:"not null"`

	//商品卖点
	SalePoint string `json:"sale_point" gorm:"default: null"`

	//发货地
	//Address string `json:"address" gorm:"default: null"`

	//月销量
	MonthSale int `json:"month_sale" `

	//付款人数
	People int `json:"people" `

	//是否有优惠
	Sale bool `json:"sale" gorm:"default: false"`

	//是否付款减库存
	PaidAndRemove bool `json:"paid_and_remove" `

	//总库存
	Total int `json:"total" `

	//是否展示库存
	ShowTotal bool `json:"show_total" `

	//运费
	Carriage int `json:"carriage" gorm:"default: null"`

	//收货方式
	GetWay int16 `json:"get_way" gorm:"not null"`

	//上架时间
	PutawayTime int64 `json:"putaway_time" gorm:"default: null"`

	//上架方式
	Putaway int16 `json:"putaway" gorm:"not null"`

	//商品状态
	OnSale bool `json:"on_sale" `

	//是否支持换货
	Exchange bool `json:"exchange" `

	//是否支持7天无理由退货
	SaleReturn bool `json:"sale_return" `

	//起售数量
	MinSale int `json:"min_sale" gorm:"default: 1"`

	//是否预售
	Advance bool `json:"advance" gorm:"default: false"`

	//预售时间
	AdvanceTime int64 `json:"advance_time" gorm:"default: null"`

	//是否限购
	Limit bool `json:"limit" gorm:"default: false"`

	//限购数量
	LimitTotal int `json:"limit_total" gorm:"default: null"`

	//图片
	Pictures string `json:"pictures"`

	//封面图片
	Cover string `json:"cover"`

	//视频
	View string `json:"view"`

	//TODO
	//商品详情
	Detail string `json:"detail" gorm:"type: text"`

	// 创建时间
	CreateTime int64 `json:"create_time"`

	// 更新时间
	UpdateTime int64 `json:"update_time"`
}

//商品规格
type GoodsSpecification struct {
	ID int `gorm:"primary_key" json:"id"`

	//商品id
	GoodsID int `json:"goods_id" gorm:"not null;index"`

	//模版id
	TemplateID int `json:"template_id" gorm:"not null"`

	//规格（一句话描述）
	Specification string `json:"specification"`
}

var f = []string{
	"ID", "Name", "SalePoint", "MonthSale", "People", "Sale", "PaidAndRemove",
	"Total", "ShowTotal", "Carriage", "GetWay", "PutawayTime", "Putaway", "OnSale", "Exchange",
	"SaleReturn", "MinSale", "Advance", "AdvanceTime", "Limit", "LimitTotal", "View", "Cover", "Detail", "CreateTime",
	"UpdateTime",
}

var f1 = []string{
	"ID", "TemplateID",
}

func (g Goods) GetInfo() map[string]interface{} {
	v := paramsUtils.ModelToDict(g, f)

	tt := goodsInfo.PictureUnmarshal(g.Pictures)

	v["pictures"] = tt

	var placeTag []GoodsAndTag
	if err := Driver.Where("goods_id = ? and tag_type = ?", g.ID, goodsTagEnum.TagTypePlaceTag).Find(&placeTag).Error; err == nil {
		place := make([]int,len(placeTag))
		for index,i := range placeTag{
			place[index] = i.GoodsTagID
		}
		v["place_tag"] = place
	}

	var saleTag []GoodsAndTag
	if err := Driver.Where("goods_id = ? and tag_type = ?", g.ID, goodsTagEnum.TagTypeSaleTag).Find(&saleTag).Error; err == nil {
		sale := make([]int,len(saleTag))
		for index,i := range saleTag{
			sale[index] = i.GoodsTagID
		}
		v["sale_tag"] = sale
	}

	var kindTag []GoodsAndTag
	if err := Driver.Where("goods_id = ? and tag_type = ?", g.ID, goodsTagEnum.TagTypeKindTag).Find(&kindTag).Error; err == nil {
		kind := make([]int,len(kindTag))
		for index,i := range kindTag{
			kind[index] = i.GoodsTagID
		}
		v["kind_tag"] = kind
	}
	var s GoodsSpecification
	if err := Driver.Where("goods_id = ?", g.ID).First(&s).Error; err == nil {
		v["template_id"] = s.TemplateID
		var template SpecificationTemplate
		Driver.GetOne("specification_template", s.TemplateID, &template)
		t := goods.TemplateUnmarshal(template.Template)
		data := make([]string, 0)
		for _, v := range t {
			if v.Use == true {
				data = append(data, v.Name)
			}
		}
		v["template"] = data
		if g.Sale {
			entity2 := new(goods.GoodsSpecificationEntity)
			if err := json.Unmarshal([]byte(s.Specification), entity2); err != nil {
				panic(goodsInfoException.SpecificationUnMarshalFail())
			}
			v["specification"] = *entity2
		} else {
			entity2 := new(goods.GoodsSpecificationEntity1)
			if err := json.Unmarshal([]byte(s.Specification), entity2); err != nil {
				panic(goodsInfoException.SpecificationUnMarshalFail())
			}
			v["specification"] = *entity2
		}
	}
	return v
}
