package goodsInfoEnum

import enumsbase "grpc-demo/enums"

const (
	SortWayMain      = 1 // 综合
	SortWayPriceDesc = 2 //价格降序
	SortWayPriceAsc  = 3 //价格升序
	SortWayPeople    = 4 // 销量

)

func NewSortWayEnums() enumsbase.EnumBaseInterface {
	return enumsbase.EnumBase{
		Enums: []int{1, 2, 3, 4},
	}
}
