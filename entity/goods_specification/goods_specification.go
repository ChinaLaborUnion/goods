package specificationEntity

import (
	"encoding/json"
	"fmt"
	goodsInfoException "grpc-demo/exceptions/goods_info"
	goodsSpecificationException "grpc-demo/exceptions/goods_specification"
)

type GoodsSpecificationEntity []struct {
	//规格
	Specification map[string]interface{} `json:"specification"`

	//id
	ID int `json:"id"`

	//库存
	Total int `json:"total"`

	//价格
	Price int `json:"price" `

	//重量
	Weight float32 `json:"weight" `

	//成本价
	CostPrice int `json:"cost_price" `

	//图片
	Picture string `json:"picture" `

	//优惠价
	ReducedPrice int `json:"reduced_price" `

}

type GoodsSpecificationEntity1 []struct {
	//规格
	Specification map[string]interface{} `json:"specification"`

	//id
	ID int `json:"id"`

	//库存
	Total int `json:"total"`

	//价格
	Price int `json:"price" `

	//重量
	Weight float32 `json:"weight" `

	//成本价
	CostPrice int `json:"cost_price" `

	//图片
	Picture string `json:"picture" `

}

type TemplateEntity []struct{
	Name string `json:"name"`
	Use bool `json:"use"`
}

type T struct{
	Name string `json:"name"`
	Use bool `json:"use"`
}

func TemplateMarshal(Template []interface{}) string{
	var template string
	//序列化
	if v, err := json.Marshal(Template); err != nil {
		panic(goodsSpecificationException.TemplateMarshalFail())
	} else {
		entity := new(TemplateEntity)
		//反序列化保证格式
		if err1 := json.Unmarshal(v, entity); err1 != nil {
			panic(goodsSpecificationException.TemplateUnMarshalFail())
		} else {
			//再次序列化
			if data, err2 := json.Marshal(entity); err2 != nil {
				panic(goodsSpecificationException.TemplateMarshalFail())
			} else {
				template = string(data)
			}
		}
	}
	return template
}


func TemplateUnmarshal(template string) []T{

	entity := new([]T)
	if err := json.Unmarshal([]byte(template), entity); err != nil {
		panic(goodsSpecificationException.TemplateUnMarshalFail())
	}
	return *entity
}

func SpecificationMarshal(Specification []interface{},sale bool) string{
	var s string
	//序列化
	if v, err := json.Marshal(Specification); err != nil {
		panic(goodsInfoException.SpecificationMarshalFail())
	} else {
		if sale{
			entity := new(GoodsSpecificationEntity)
			//反序列化保证格式
			if err1 := json.Unmarshal(v, entity); err1 != nil {
				fmt.Println(err1)
				panic(goodsInfoException.SpecificationUnMarshalFail())
			} else {
				//再次序列化
				if data, err2 := json.Marshal(entity); err2 != nil {
					panic(goodsInfoException.SpecificationMarshalFail())
				} else {
					s = string(data)
				}
			}
		}else{
			entity := new(GoodsSpecificationEntity1)
			//反序列化保证格式
			if err1 := json.Unmarshal(v, entity); err1 != nil {
				fmt.Println(err1)
				panic(goodsInfoException.SpecificationUnMarshalFail())
			} else {
				//再次序列化
				if data, err2 := json.Marshal(entity); err2 != nil {
					panic(goodsInfoException.SpecificationMarshalFail())
				} else {
					s = string(data)
				}
			}
		}
	}
	return s
}
