package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"sync"
	"time"
)

//红包数据结构:  map[用户ID]对应红包（单位:分）
var PackageList = new(sync.Map)

type LotteryController struct {
}
type SetReq struct {
	Money int64 `form:"money" json:"money" binding:"required"`
	Num   uint  `form:"num" json:"num" binding:"required"`
	Uid   int32 `form:"uid" json:"uid" binding:"required"`
}

func (l LotteryController) Set(ctx *gin.Context) {
	var (
		req  SetReq
		err  error
		list []int64
	)
	if err = ctx.ShouldBindQuery(&req); err != nil {
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
	PackageList.Store(req.Uid, list)
	/*PackageList[req.Uid] = list*/
	log.Printf("PackageList:%v", PackageList)
	url := fmt.Sprintf("/get?uid=%d", req.Uid)
	data := gin.H{
		"data": gin.H{"url": url},
		"msg":  "Success",
	}
	ctx.JSON(200, data)
}

type GetReq struct {
	Uid int32 `form:"uid" json:"uid" binding:"required"`
}

func (l LotteryController) Get(ctx *gin.Context) {
	var (
		err error
		req GetReq
	)

	if err = ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(200, gin.H{"msg": "参数错误", "data": ""})
		return
	}
	list, has := PackageList.Load(req.Uid)
	if !has {
		ctx.JSON(200, gin.H{"msg": "红包已经发完", "data": ""})
		return
	}
	iList := list.([]int64)
	money := iList[0]
	if len(iList) == 1 {
		PackageList.Delete(req.Uid)
	} else {
		PackageList.Store(req.Uid, iList[1:])
	}
	log.Printf("PackageList:%v", PackageList)
	ctx.JSON(200, gin.H{"msg": "success", "data": money})
	return
}
