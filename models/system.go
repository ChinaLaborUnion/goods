package models

type Redis struct {
	Username string `ini:"username"`
	Password string `ini:"password"`
	Host     string `ini:"host"`
	Port     string `ini:"port"`
}

type Mysql struct {
	DB       string `ini:"db"`             // db名
	Host     string `ini:"host"`         // 主机名
	Port     string `ini:"port"`         // 端口
	Username string `ini:"username"` // 用户名
	Password string `ini:"password"` // 密码
}


type Qiniu struct {
	Bucket    string `ini:"bucket"`         // 空间名
	Expires   uint64 `ini:"expires"`       // 过期时间
	AccessKey string `ini:"access_key"` // 密钥
	SecretKey string `ini:"secret_key"` // 密钥
}

type SystemConfiguration struct {
	Redis
	Mysql
	Qiniu
}