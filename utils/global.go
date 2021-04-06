package utils

import (
	"fmt"
	"gopkg.in/ini.v1"
	"grpc-demo/models"
)

const filename = "conf.ini"
var GlobalConfig models.SystemConfiguration
func InitGlobal() {
	_ = ini.MapTo(&GlobalConfig, filename)
	fmt.Printf("%#v\n",GlobalConfig)
}
