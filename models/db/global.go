package db

type Icon struct {
	ID int `gorm:"primary_key" json:"id"`

	// 标题
	Title string `json:"title" gorm:"not null"`

	//名称
	Name string `json:"name" gorm:"not null"`

	//图片
	Picture string `json:"picture" gorm:"not null"`
}
