package goodsInfoEnum

import enumsbase "grpc-demo/enums"

const (
	PutawayNow = 1  // 立即上架售卖
	PutawayDefine   = 2  // 自定义上架时间
	PutawayNo = 4 //暂不售卖


)

func NewPutawayEnums() enumsbase.EnumBaseInterface {
	return enumsbase.EnumBase{
		Enums: []int{1, 2, 4},
	}
}
