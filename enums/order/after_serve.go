package orderEnums

import enumsbase "grpc-demo/enums"

const (
	BackMoney         = 1 //退款
	BackMoneyAndThing = 2 //退货退款
)

func NewAfterServeEnums() enumsbase.EnumBaseInterface {
	return enumsbase.EnumBase{
		Enums: []int{1, 2},
	}
}
