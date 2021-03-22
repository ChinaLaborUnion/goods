package db

//种类标签
type KindTag struct {
	ID int `gorm:"primary_key" json:"id"`

	//父种类
	ParentID int `json:"parent_id" gorm:"default: 0"`

	// 标题
	Title string `json:"title" gorm:"not null"`

	//图片
	Picture string `json:"picture"`
}

//属地标签
type PlaceTag struct {
	ID int `gorm:"primary_key" json:"id"`

	// 标题
	Place string `json:"place" gorm:"not null"`

	//图片
	Picture string `json:"picture"`
}

//销售标签
type SaleTag struct {
	ID int `gorm:"primary_key" json:"id"`

	// 标题
	Title string `json:"title" gorm:"not null"`

}

//商品挂载标签（多对多）
type GoodsAndTag struct {
	ID int `gorm:"primary_key" json:"id"`

	//商品id
	GoodsID int `json:"goods_id" gorm:"not null;index"`

	//标签id
	GoodsTagID int `json:"goods_tag_id" gorm:"not null"`

	//标签种类
	TagType int16 `json:"tag_type" gorm:"not null"`
}
