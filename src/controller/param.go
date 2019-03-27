package controller

type ParamGetID struct {
	Type     string `form:"type" binding:"required"` //TODO: 这个参数验证有问题
	ID       int `form:"id" binding:"required"`
	FromTime int    `form:"from" binding:"omitempty,gte=20160101,lte=20300101"`
	ToTime   int    `form:"to" binding:"omitempty,gtefield=FromTime,lte=20300101"`
}

type ParamGetFormula struct {
	Type string `form:"type" binding:"required,eq=tdoll|eq=equip|ep=fairy"` //TODO: 这个参数验证有问题
	Mp int `form:"mp" binding:"required,gte=30,lte=9999"`
	Ammo int `form:"ammo" binding:"required,gte=30,lte=9999"`
	Mre int `form:"mre" binding:"required,gte=30,lte=9999"`
	Part int `form:"part" binding:"required,gte=30,lte=9999"`
	InputLevel int `form:"part" binding:"required,gte=0,lte=3"`
	FromTime int    `form:"from" binding:"omitempty,gte=20160101,lte=20300101"`
	ToTime   int    `form:"to" binding:"omitempty,gtefield=FromTime,lte=20300101"`
}
