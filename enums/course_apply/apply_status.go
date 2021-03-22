package courseApplyEnums



import enumsbase "grpc-demo/enums"

const (
	PreApplyCancel = 1  // 撤销
	NoApply   = 2  // 预报名
	Apply = 4 //已报名


)

func NewApplyStatusEnums() enumsbase.EnumBaseInterface {
	return enumsbase.EnumBase{
		Enums: []int{1, 2, 4},
	}
}


