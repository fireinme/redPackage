package main

import (
	"github.com/gin-gonic/gin"
	"weiboRedPackage/controller"
)

/*抢红包小练习*/

func main() {
	app := gin.Default()
	InitRoute(app)
	_ = app.Run()
}
func InitRoute(app *gin.Engine) {
	lotteryController := controller.LotteryController{}
	app.GET("/set", lotteryController.Set) //设置红包金额
	app.GET("/get", lotteryController.Get) //获取红包金额

}
