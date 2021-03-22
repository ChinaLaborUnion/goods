package courseEntity

import (
	"encoding/json"
	"fmt"
	courseException "grpc-demo/exceptions/course"
)

type CoursePictureEntity []struct{
	Picture string `json:"picture"`
	Order int `json:"order"`
}

func CoursePictureMarshal(CoursePicture []interface{}) string{
	var s string
	//序列化
	if v, err := json.Marshal(CoursePicture); err != nil {
		panic(courseException.PictureMarshalFail())
	} else {
		entity := new(CoursePictureEntity)
		//反序列化保证格式
		if err1 := json.Unmarshal(v, entity); err1 != nil {
			//fmt.Println(err1)
			panic(courseException.PictureUnmarshalFail())
		} else {
			//再次序列化
			if data, err2 := json.Marshal(entity); err2 != nil {
				panic(courseException.PictureMarshalFail())
			} else {
				s = string(data)
			}
		}
	}
	return s
}

type P []struct{
	Picture string `json:"picture"`
	Order int `json:"order"`
}

func CoursePictureUnmarshal(pictures string) P{

	entity := new(P)
	if err := json.Unmarshal([]byte(pictures), entity); err != nil {
		panic(courseException.PictureUnmarshalFail())
	}
	fmt.Println(*entity)
	return *entity
}

