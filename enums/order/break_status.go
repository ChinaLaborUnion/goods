package orderEnums

import enumsbase "grpc-demo/enums"

const (
	BreakStatusBreak = 1  // 已拆单
	BreakStatusNo   = 2  // 未拆单
)

func NewBreakStatusEnums() enumsbase.EnumBaseInterface {
	return enumsbase.EnumBase{
		Enums: []int{1, 2},
	}
}
