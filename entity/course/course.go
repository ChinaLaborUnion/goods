package courseEntity

import (
	"encoding/json"
	courseException "grpc-demo/exceptions/course"
)

type CourseSessionEntity []struct{
	BeginTime int64 `json:"begin_time"`
	EndTime int64 	`json:"end_time"`
	Money int `json:"money"`
	PeopleLimit int `json:"people_limit"`
	ID int `json:"id"`
}

func CourseSessionMarshal(CourseSession []interface{}) string{
	var s string
	//序列化
	if v, err := json.Marshal(CourseSession); err != nil {
		panic(courseException.SessionMarshalFail())
	} else {
		entity := new(CourseSessionEntity)
		//反序列化保证格式
		if err1 := json.Unmarshal(v, entity); err1 != nil {
			//fmt.Println(err1)
			panic(courseException.SessionUnmarshalFail())
		} else {
			//再次序列化
			if data, err2 := json.Marshal(entity); err2 != nil {
				panic(courseException.SessionMarshalFail())
			} else {
				s = string(data)
			}
		}
	}
	return s
}

type Session []struct{
	BeginTime int64 `json:"begin_time"`
	EndTime int64 	`json:"end_time"`
	Money int `json:"money"`
	PeopleLimit int `json:"people_limit"`
	ID int `json:"id"`
}

func CourseSessionUnmarshal(session string) Session{

	entity := new(Session)
	if err := json.Unmarshal([]byte(session), entity); err != nil {
		panic(courseException.SessionUnmarshalFail())
	}
	return *entity
}

