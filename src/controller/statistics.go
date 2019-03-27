package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"model"
	"util"
	"net/http"
)

func GetStatsInfo(c *gin.Context) {
	info, err := model.GetInfo()
	if err != nil {
		log.Printf("GetInfo error: %s", err.Error())
		util.Error(c, http.StatusInternalServerError, "Server error!")
		return
	}

	util.Success(c, info)
}

func GetRecordsByID(c *gin.Context) {

	// 参数检查
	var param ParamGetID
	if err := c.ShouldBindQuery(&param); err != nil {
		util.Error(c, http.StatusBadRequest, "Bad query!")
		return
	}

	// 判断查询数据类型(人形、装备、精灵)
	var t int

	switch param.Type {
	case "tdoll":
		t = model.RecordTypeTDoll
	case "equip":
		t = model.RecordTypeEquip
	case "fairy":
		t = model.RecordTypeFairy
	}

	if param.ToTime == 0 {
		param.ToTime = util.GetDate() // 如果没有设置结束时间，设为今天
	}

	if param.FromTime == 0 {
		param.FromTime = model.GetAddedDate(t, param.ID) // 起始时间设为新加入当天
	}

	log.Printf("%+v\n", param)

	ret, err := model.GetIDStatsData(t, param.ID, param.FromTime, param.ToTime) // 获得数据

	if err != nil {
		log.Printf("Server error in GetRecordsByID: %s", err.Error())
		util.Error(c, http.StatusInternalServerError, "Server error.")
		return
	}

	util.Success(c, ret)
}

func GetRecordsByFormula(c *gin.Context) {
	var param ParamGetFormula
	if err := c.ShouldBindQuery(&param); err != nil {
		util.Error(c, http.StatusBadRequest, "Bad query!")
		return
	}

	var t int
	switch param.Type {
	case "tdoll":
		t = model.RecordTypeTDoll
	case "equip":
		t = model.RecordTypeEquip
	case "fairy":
		t = model.RecordTypeFairy
	}

	if param.ToTime == 0 {
		param.ToTime = util.GetDate() // 设为今天
	}

	formula := model.Formula{
		Mp: param.Mp,
		Ammo: param.Ammo,
		Mre: param.Mre,
		Part: param.Part,
		InputLevel: param.InputLevel,
	}

	ret, err := model.GetFormulaStatsData(t, formula, param.FromTime, param.ToTime)

	if err != nil {
		log.Printf("Server error in GetRecordsByFormula: %s", err.Error())
		util.Error(c, http.StatusInternalServerError, "Server error.")
		return
	}

	util.Success(c, ret)
}
