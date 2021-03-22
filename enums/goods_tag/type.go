package goodsTagEnum

import enumsbase "grpc-demo/enums"

const (
	TagTypeKindTag = 1  // 种类标签
	TagTypePlaceTag   = 2  // 属地标签
	TagTypeSaleTag   = 4  // 销售标签

)

func NewTagTypeEnums() enumsbase.EnumBaseInterface {
	return enumsbase.EnumBase{
		Enums: []int{1, 2, 4},
	}
}
