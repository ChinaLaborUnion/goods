package courseApplyEnums


import enumsbase "grpc-demo/enums"

const (
	NoPay = 1  // 未支付
	Paied   = 2  // 已支付
	Cancel = 4 //已取消
	AfterSales = 5 //申请售后
	AlreadyAS  = 6 //售后完成
	Over = 3 //已完成


)

func NewStatusEnums() enumsbase.EnumBaseInterface {
	return enumsbase.EnumBase{
		Enums: []int{1, 2, 4},
	}
}

