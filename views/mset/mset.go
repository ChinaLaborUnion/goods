package mset

import (
	"github.com/kataras/iris"
)

func Mset(ctx iris.Context){
	//var count int
	//
	//table := db.Driver.Table("country")
	//table.Count(&count)
	//var country db.Country
	//for i:=1; i <= count;i ++{
	//	db.Driver.Where("id = ?",i).First(&country)
	//	key := paramsUtils.CacheBuildKey(constants.DbModel, "country", country.Name)
	//	_, _ = cache.Redis.Do(constants.DbNumberModel, "set", key, country.ID)
	//}
	//
	//ctx.JSON(iris.Map{
	//	"status":"success",
	//})

	//var count int
	//
	//table := db.Driver.Table("province")
	//table.Count(&count)
	//fmt.Println(count)
	//for i:=1; i <= count;i ++{
	//	var province db.Province
	//	db.Driver.GetOne("province",i,&province)
	//	s := fmt.Sprintf("%s-%d",province.Name,province.CountryID)
	//	key := paramsUtils.CacheBuildKey(constants.DbModel, "province", s)
	//	_, _ = cache.Redis.Do(constants.DbNumberModel, "set", key, province.ID)
	//}
	//
	//ctx.JSON(iris.Map{
	//	"status":"success",
	//})

	//var count int
	//
	//table := db.Driver.Table("city")
	//table.Count(&count)
	//
	//for i:=1; i <= count;i ++{
	//	var city db.City
	//	db.Driver.GetOne("city",i,&city)
	//	s := fmt.Sprintf("%s-%d-%d",city.Name,city.CountryID,city.ProvinceID)
	//	key := paramsUtils.CacheBuildKey(constants.DbModel, "city", s)
	//	_, _ = cache.Redis.Do(constants.DbNumberModel, "set", key, city.ID)
	//}
	//
	//ctx.JSON(iris.Map{
	//	"status":"success",
	//})

	//var count int
	//
	//table := db.Driver.Table("district")
	//table.Count(&count)
	//
	//for i:=1; i <= count;i ++{
	//	var district db.District
	//	db.Driver.GetOne("district",i,&district)
	//	s := fmt.Sprintf("%s-%d-%d-%d",district.Name,district.CountryID,district.ProvinceID,district.CityID)
	//	key := paramsUtils.CacheBuildKey(constants.DbModel, "district", s)
	//	_, _ = cache.Redis.Do(constants.DbNumberModel, "set", key, district.ID)
	//}
	//
	//ctx.JSON(iris.Map{
	//	"status":"success",
	//})
}
