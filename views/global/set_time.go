package global

import (
	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris"
	"grpc-demo/constants"
	authbase "grpc-demo/core/auth"
	"grpc-demo/core/cache"
	accountException "grpc-demo/exceptions/account"
	globalException "grpc-demo/exceptions/global"
	paramsUtils "grpc-demo/utils/params"
)

func SetTime(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckLogin()

	if !auth.IsAdmin(){
		panic(accountException.NoPermission())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	//未付款多久自动取消订单
	//自提货多久自动收货
	//待收货多久自动收货
	//待评价多久自动变成已完成
	//已完成多久不能追评
	autoCancel := params.Int("auto_cancel","未付款自动取消时间",86400)
	pickAutoGet := params.Int("pick_auto_get","自取自动收货时间",604800)
	tranAutoGet := params.Int("tran_auto_get","发货自动收货时间",604800)
	toOver := params.Int("to_over","自动完成时间",604800)
	noSecondComment := params.Int("no_second_comment","限制追评时间",7776000)

	if _, err := redis.Bytes(cache.Redis.Do(
		constants.DbNumberOther,
		"set",
		paramsUtils.CacheBuildKey(constants.GlobalOrderTime, "auto_cancel"),
		autoCancel)); err != nil {
		panic(globalException.SetTimeError())
	}

	if _, err := redis.Bytes(cache.Redis.Do(
		constants.DbNumberOther,
		"set",
		paramsUtils.CacheBuildKey(constants.GlobalOrderTime, "pick_auto_get"),
		pickAutoGet)); err != nil {
		panic(globalException.SetTimeError())
	}

	if _, err := redis.Bytes(cache.Redis.Do(
		constants.DbNumberOther,
		"set",
		paramsUtils.CacheBuildKey(constants.GlobalOrderTime, "tran_auto_get"),
		tranAutoGet)); err != nil {
		panic(globalException.SetTimeError())
	}

	if _, err := redis.Bytes(cache.Redis.Do(
		constants.DbNumberOther,
		"set",
		paramsUtils.CacheBuildKey(constants.GlobalOrderTime, "to_over"),
		toOver)); err != nil {
		panic(globalException.SetTimeError())
	}

	if _, err := redis.Bytes(cache.Redis.Do(
		constants.DbNumberOther,
		"set",
		paramsUtils.CacheBuildKey(constants.GlobalOrderTime, "no_second_comment"),
		noSecondComment)); err != nil {
		panic(globalException.SetTimeError())
	}

	ctx.JSON(iris.Map{
		"status":"success",
	})
}

func GetTime(ctx iris.Context, auth authbase.AuthAuthorization){
	auth.CheckLogin()

	if !auth.IsAdmin(){
		panic(accountException.NoPermission())
	}

	key1 := paramsUtils.CacheBuildKey(constants.GlobalOrderTime, "auto_cancel")
	autoCancel, err1 := redis.Int64(cache.Redis.Do(constants.DbNumberOther, "get", key1))
	if err1 != nil || autoCancel == 0 {
		autoCancel = 86400
	}

	key2 := paramsUtils.CacheBuildKey(constants.GlobalOrderTime, "pick_auto_get")
	pickAutoGet, err2 := redis.Int64(cache.Redis.Do(constants.DbNumberOther, "get", key2))
	if err2 != nil || pickAutoGet == 0 {
		pickAutoGet = 604800
	}

	key3 := paramsUtils.CacheBuildKey(constants.GlobalOrderTime, "tran_auto_get")
	tranAutoGet, err3 := redis.Int64(cache.Redis.Do(constants.DbNumberOther, "get", key3))
	if err3 != nil || tranAutoGet == 0 {
		tranAutoGet = 604800
	}

	key4 := paramsUtils.CacheBuildKey(constants.GlobalOrderTime, "to_over")
	toOver, err4 := redis.Int64(cache.Redis.Do(constants.DbNumberOther, "get", key4))
	if err4 != nil || toOver == 0 {
		toOver = 604800
	}

	key5 := paramsUtils.CacheBuildKey(constants.GlobalOrderTime, "no_second_comment")
	noSecondComment, err5 := redis.Int64(cache.Redis.Do(constants.DbNumberOther, "get", key5))
	if err5 != nil || noSecondComment == 0 {
		noSecondComment = 7776000
	}

	ctx.JSON(iris.Map{
		"auto_cancel":autoCancel,
		"pick_auto_get":pickAutoGet,
		"tran_auto_get":tranAutoGet,
		"to_over":toOver,
		"no_second_comment":noSecondComment,
	})
}
