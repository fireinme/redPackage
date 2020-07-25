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
var (
	PackageList = new(sync.Map)
)

//发红包任务
type Task struct {
	uid       int64
	moneyChan chan int64
}

type LotteryController struct {
	Ch chan Task
}

//参数校验
type SetReq struct {
	Money int64 `form:"money" json:"money" binding:"required"`
	Num   int64 `form:"num" json:"num" binding:"required"`
	Uid   int64 `form:"uid" json:"uid" binding:"required"`
}
type GetReq struct {
	Uid int64 `form:"uid" json:"uid" binding:"required"`
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
	list = l.GenerateRandMoney(leftMoney, num)
	PackageList.Store(req.Uid, list)
	log.Printf("PackageList:%v", PackageList)
	url := fmt.Sprintf("/get?uid=%d", req.Uid)
	data := gin.H{
		"data": gin.H{"url": url},
		"msg":  "Success",
	}
	ctx.JSON(200, data)
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
	mChan := make(chan int64)
	task := Task{
		uid:       req.Uid,
		moneyChan: mChan,
	}
	//派发任务
	l.Ch <- task
	//得到红包
	money := <-task.moneyChan
	log.Printf("PackageList:%v", PackageList)
	ctx.JSON(200, gin.H{"msg": "success", "data": money})
	return
}

func (l LotteryController) GetPackageServer() {
	for {
		task := <-l.Ch
		id := task.uid
		//找出对应红包
		load, ok := PackageList.Load(id)
		if !ok {
			task.moneyChan <- 0
			continue
		}
		int64s := load.([]int64)
		if len(int64s) == 0 {
			PackageList.Delete(id)
			task.moneyChan <- 0
			continue
		}
		task.moneyChan <- int64s[0]
		if len(int64s[1:]) == 0 {
			PackageList.Delete(id)
			continue
		}
		PackageList.Store(id, int64s[1:])
	}

}

//随机生成红包金额
func (l LotteryController) GenerateRandMoney(leftMoney int64, num int64) []int64 {
	var list []int64
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
	return list
}
