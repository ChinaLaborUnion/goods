package goods_info

import (
	"github.com/kataras/iris"
	authbase "grpc-demo/core/auth"
	goods "grpc-demo/entity/goods_specification"
	goodsInfoException "grpc-demo/exceptions/goods_info"
	goodsSpecificationException "grpc-demo/exceptions/goods_specification"
	"grpc-demo/models/db"
	paramsUtils "grpc-demo/utils/params"
)

func CreateGoodsSpecification(ctx iris.Context, auth authbase.AuthAuthorization, gid int) {
	auth.CheckAdmin()

	var g db.Goods
	if err := db.Driver.GetOne("goods", gid, &g); err != nil {
		panic(goodsInfoException.GoodsIsNotExsit())
	}

	//存在商品规格则不能创建
	var gs db.GoodsSpecification
	if err := db.Driver.Where("goods_id = ?", gid).First(&gs).Error; err == nil {
		panic(goodsInfoException.SpecificationExsit())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	t := params.Int("template_id", "模版id")
	var template db.SpecificationTemplate
	if err := db.Driver.GetOne("specification_template", t, &template); err != nil {
		panic(goodsSpecificationException.TemplateIsNotExsit())
	}

	//找出模版中启用的参数
	tt := goods.TemplateUnmarshal(template.Template)
	data := make([]string, 0)
	for _, v := range tt {
		if v.Use == true {
			data = append(data, v.Name)
		}
	}

	specification := params.List("specification", "规格")

	minPrice := int(specification[0].(map[string]interface{})["price"].(float64))
	total := 0

	for index, v1 := range specification {
		p := v1.(map[string]interface{})["specification"]



		if int(v1.(map[string]interface{})["price"].(float64)) < minPrice{
			minPrice = int(v1.(map[string]interface{})["price"].(float64))
		}


		v1.(map[string]interface{})["id"] = index+1
		total += int(v1.(map[string]interface{})["total"].(float64))

		//判断规格是否包含所有启用的规格参数
		if len(data) != len(p.(map[string]interface{})) {
			panic(goodsInfoException.SpecificationParamsError())
		}
		for _, v2 := range data {
			if _, ok := p.(map[string]interface{})[v2]; !ok {
				panic(goodsInfoException.SpecificationParamsError())
			}
		}



	}

	s := goods.SpecificationMarshal(specification,g.Sale)
	goodsSpecification := db.GoodsSpecification{
		GoodsID:       gid,
		TemplateID:    t,
		Specification: s,
	}

	g.Total = total
	g.MinPrice = minPrice
	//g.Price = float64(minPrice)

	//同时改变商品中的总库存
	tx := db.Driver.Begin()
	if err := tx.Save(&g).Error; err != nil {
		tx.Rollback()
		panic(goodsInfoException.SpecificationCreateError())
	}

	if err := tx.Create(&goodsSpecification).Error; err != nil {
		tx.Rollback()
		panic(goodsInfoException.SpecificationCreateError())
	}

	tx.Commit()


	ctx.JSON(iris.Map{
		"id": g.ID,
	})
}

func PutGoodsSpecification(ctx iris.Context, auth authbase.AuthAuthorization, gid int) {
	auth.CheckAdmin()

	var both struct {
		db.Goods
		db.GoodsSpecification
	}

	if err := db.Driver.Select("goods.*, goods_specification.*").Table("goods, goods_specification").Where("goods.id = ? and goods_specification.goods_id = goods.id", gid).Find(&both).Error; err != nil {
		panic(goodsInfoException.SpecificationIsNotExsit())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	if params.Has("template_id"){
		t := params.Int("template_id", "模版id")
		var template db.SpecificationTemplate
		if err := db.Driver.GetOne("specification_template", t, &template); err != nil {
			panic(goodsSpecificationException.TemplateIsNotExsit())
		}
		both.GoodsSpecification.TemplateID = t
	}


		if params.Has("specification") {
		specification := params.List("specification", "规格")


		var t1 db.SpecificationTemplate
		if err := db.Driver.GetOne("specification_template", both.GoodsSpecification.TemplateID, &t1); err != nil {
			panic(goodsSpecificationException.TemplateIsNotExsit())
		}

		//找出模版中启用的参数
		tt := goods.TemplateUnmarshal(t1.Template)
		data := make([]string, 0)
		for _, v := range tt {
			if v.Use == true {
				data = append(data, v.Name)
			}
		}



		total := 0
		minPrice := int(specification[0].(map[string]interface{})["price"].(float64))
		//minPrice := specification[0].(map[string]interface{})["price"].(float64)
		for index, v1 := range specification {
			p := v1.(map[string]interface{})["specification"]

			v1.(map[string]interface{})["id"] = index+1
			//if v1.(map[string]interface{})["price"].(float64) < minPrice{
			//	minPrice = v1.(map[string]interface{})["price"].(float64)
			//}

			if int(v1.(map[string]interface{})["price"].(float64)) < minPrice{
				minPrice = int(v1.(map[string]interface{})["price"].(float64))
			}

			total += int(v1.(map[string]interface{})["total"].(float64))

			//判断规格是否包含所有启用的规格参数
			if len(data) != len(p.(map[string]interface{})) {
				panic(goodsInfoException.SpecificationParamsError())
			}
			for _, v2 := range data {
				if _, ok := p.(map[string]interface{})[v2]; !ok {
					panic(goodsInfoException.SpecificationParamsError())
				}
			}
		}
		str1 := goods.SpecificationMarshal(specification,both.Goods.Sale)
		both.Goods.Total = total
		both.Goods.MinPrice = minPrice
		//both.Goods.Price = float64(minPrice)

		both.GoodsSpecification.Specification = str1

		tx := db.Driver.Begin()
		if err := tx.Save(&both.Goods).Error;err != nil{
			tx.Rollback()
			panic(goodsInfoException.SpecificationPutError())
		}

		if err := tx.Save(&both.GoodsSpecification).Error;err != nil{
			tx.Rollback()
			panic(goodsInfoException.SpecificationPutError())
		}
		tx.Commit()

	}

	ctx.JSON(iris.Map{
		"id": both.Goods.ID,
	})
}

func DeleteGoodsSpecification(ctx iris.Context, auth authbase.AuthAuthorization, gid int) {
	auth.CheckAdmin()

	var s db.GoodsSpecification
	if err := db.Driver.Where("goods_id = ?", gid).First(&s).Error; err == nil {

		var goods db.Goods
		db.Driver.GetOne("goods", gid, &goods)
		goods.Total = 0
		tx := db.Driver.Begin()
		//同时改变商品中的总库存
		if err := tx.Save(&goods).Error; err != nil {
			tx.Rollback()
			panic(goodsInfoException.DeleteSpecificationFail())
		}

		if err := tx.Delete(s).Error; err != nil {
			tx.Rollback()
			panic(goodsInfoException.DeleteSpecificationFail())
		}

		tx.Commit()
	}

	ctx.JSON(iris.Map{
		"id": gid,
	})
}
