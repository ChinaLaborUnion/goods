package courseEnums

import enumsbase "grpc-demo/enums"

const (
	PlaceType = 1 //属地
	CourseType = 2 //课程性质

)

func NewTagTypeEnums() enumsbase.EnumBaseInterface {
	return enumsbase.EnumBase{
		Enums: []int{1, 2},
	}
}
