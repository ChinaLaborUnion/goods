package orderEnums

import enumsbase "grpc-demo/enums"

const (
	IsPassSUMMIT  = 1 //已提交 / 等待审核
	IsPassPERMIT  = 2 //商家审核 通过
	IsPassREJECT  = 4 //商家审核 拒绝
	IsPassSUCCESS = 8 //成功 / 已退款/货

)

func NewIsPassEnums() enumsbase.EnumBaseInterface {
	return enumsbase.EnumBase{
		Enums: []int{1, 2, 4, 8},
	}
}
