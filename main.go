package main

import (
	"github.com/gin-gonic/gin"
	"weiboRedPackage/controller"
)

/*抢红包小练习*/

func main() {
	app := gin.Default()
	lotteryController := controller.LotteryController{}
	//单独处理红包发放
	ch := make(chan controller.Task)
	lotteryController.Ch = ch
	go lotteryController.GetPackageServer()
	//启动路由
	InitRoute(app, lotteryController)
	_ = app.Run()
}
func InitRoute(app *gin.Engine, lotteryController controller.LotteryController) {
	app.GET("/set", lotteryController.Set) //设置红包金额
	app.GET("/get", lotteryController.Get) //获取红包金额

}
