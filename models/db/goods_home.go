package db

type GoodsSlideshow struct {
	ID int `gorm:"primary_key" json:"id"`

	//商品id
	GoodsID int `json:"goods_id" gorm:"not null;index"`

	//图片
	Picture string `json:"picture"`

	//序号
	Number int `json:"number" gorm:"not null"`
}
