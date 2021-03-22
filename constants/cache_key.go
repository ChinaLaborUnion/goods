package constants

// 缓存key
const (
	DbModel                      = "model"                           // 数据库缓存键
	NicknameModel                = "account:nickname"                // 昵称缓存键
	AccountVerificationEmail     = "account:verification:email"      // 用户发送邮件token缓存键
	AccountVerificationEmailTime = "account:verification:email:time" // 用户发送邮件延迟缓存键
	GlobalOrderTime = "global:order:time"
)


const (
	DbNumberModel = 1 // 数据库缓存服务
	DbNumberEmail = 2 // 邮箱缓存服务
	DbNumberOther = 3 // 其他杂七杂八的缓存
)