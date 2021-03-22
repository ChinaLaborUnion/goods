package goodsInfoEntity

import (
	"encoding/json"
	"fmt"
	goodsInfoException "grpc-demo/exceptions/goods_info"
)

type GoodsPictureEntity []struct{
	Picture string `json:"picture"`
	Order int `json:"order"`
}

func GoodsPictureMarshal(GoodsPicture []interface{}) string{
	var s string
	//序列化
	if v, err := json.Marshal(GoodsPicture); err != nil {
		panic(goodsInfoException.PictureMarshalFail())
	} else {
		entity := new(GoodsPictureEntity)
		//反序列化保证格式
		if err1 := json.Unmarshal(v, entity); err1 != nil {
			//fmt.Println(err1)
			panic(goodsInfoException.PictureUnmarshalFail())
		} else {
			//再次序列化
			if data, err2 := json.Marshal(entity); err2 != nil {
				panic(goodsInfoException.PictureMarshalFail())
			} else {
				s = string(data)
			}
		}
	}
	return s
}

type P []struct{
	Picture string `json:"picture"`
	Order int `json:"order"`
}

func PictureUnmarshal(pictures string) P{

	entity := new(P)
	if err := json.Unmarshal([]byte(pictures), entity); err != nil {
		fmt.Println("111")
		fmt.Println("err:",err)
		panic(goodsInfoException.PictureUnmarshalFail())
	}
	fmt.Println(*entity)
	return *entity
}
