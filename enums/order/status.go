package orderEnums

import enumsbase "grpc-demo/enums"

const (
	SUBMIT     = 1 //提交订单/未付款
	WfShipped  = 2 //待发货
	ToPick     = 9 //待提货
	TbShipped  = 3 //已发货/待收货
	RECEIVED   = 4 //已收货/待评价
	OVER       = 5 //订单结束
	CANCEL     = 6 //已取消
	AfterSales = 7 //申请售后
	AlreadyAS  = 8 //售后完成
)

func NewStatusEnums() enumsbase.EnumBaseInterface {
	return enumsbase.EnumBase{
		Enums: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
	}
}
