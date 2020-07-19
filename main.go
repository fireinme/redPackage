package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"time"
)

/*抢红包小练习*/

//红包数据结构:  map[用户ID]对应红包（单位:分）
var PackageList = make(map[int32][]int64)

type LotteryController struct {
}

func main() {
	lotteryController := LotteryController{}
	app := gin.Default()
	app.GET("/set", lotteryController.Set)
	app.Run()
}

type SettingReq struct {
	Money int64 `form:"money" json:"money" binding:"required"`
	Num   uint  `form:"num" json:"num" binding:"required"`
	Uid   int32 `form:"uid" json:"uid" binding:"required"`
}

func (l LotteryController) Set(ctx *gin.Context) {
	var (
		req  SettingReq
		err  error
		list []int64
	)
	if err = ctx.ShouldBindQuery(&req); err != nil {
		log.Print("参数错误")
		ctx.JSON(200, gin.H{"msg": "参数错误", "data": ""})
		return
	}
	log.Printf("money:%d,num:%d", req.Money, req.Num)
	source := rand.NewSource(time.Now().Unix())
	rand.New(source)
	//生成随机红包
	leftMoney := req.Money //剩余金额 分
	num := req.Num
	maxRate := 0.24 //每个红包占余额的是最大比率
	for i := num; i > 0; i-- {
		var pMoney int64
		//最后一个红包 剩余全部给
		if i == 1 {
			pMoney = leftMoney
		} else {
			floorMoney := int64(float64(leftMoney) * maxRate)
			pMoney = rand.Int63n(floorMoney - 1) //随机生成金额
			//不能有空红包
			if pMoney == 0 {
				continue
			}
			leftMoney -= pMoney
		}
		list = append(list, pMoney)
	}
	PackageList[req.Uid] = list
	log.Printf("PackageList:%v", PackageList)
	url := fmt.Sprintf("/get?id=%d", req.Uid)
	data := gin.H{
		"data": gin.H{"url": url},
		"msg":  "Success",
	}
	ctx.JSON(200, data)
}
