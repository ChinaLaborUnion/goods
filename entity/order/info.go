package orderEntity
//
//import (
//	"encoding/json"
//	orderException "grpc-demo/exceptions/order"
//)
//
//type OrderEntity []struct {
//	GoodsID int `json:"goods_id"`
//	GoodsSpecification int `json:"goods_specification"`
//	GoodsTotal int `json:"goods_total"`
//	Message string `json:"message"`
//}
//
//func OrderMarshal(Order []interface{}) bool{
//	//序列化
//	if v, err := json.Marshal(Order); err != nil {
//		panic(orderException.OrderMarshalFail())
//	} else {
//		entity := new(OrderEntity)
//		//反序列化保证格式
//		if err1 := json.Unmarshal(v, entity); err1 != nil {
//
//			panic(orderException.OrderUnmarshalFail())
//		} else {
//			//再次序列化
//			if _, err2 := json.Marshal(entity); err2 != nil {
//				panic(orderException.OrderMarshalFail())
//			} else {
//				return true
//			}
//		}
//	}
//	return false
//}
