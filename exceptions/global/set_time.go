package globalException

import "grpc-demo/models"

func SetTimeError() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status:  false,
		ErrCode: 5412,
		Message: "设置失败",
	}
}
