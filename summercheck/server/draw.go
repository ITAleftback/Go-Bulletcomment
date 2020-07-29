package main

import (
	"math/rand"
	"net/http"
	"strconv"

)

func Draw(w http.ResponseWriter,r *http.Request){
	//拿去保存user的map进行遍历 看有几个user
	//定义切片存放中奖人

	var slice []int
	index:=Num()
	//假设现在抽n个人中奖
	n:=r.FormValue("n")
	N,_:=strconv.Atoi(n)
	//如果想要抽n个人中奖却大于用户数
	if N>index {
		return
	}
	for i:=0;i<N ;i++  {
		userID:=rand.Intn(index)
		slice=append(slice,userID)

		//然后将中奖名单公布
		msg:=Message{
			Content: "恭喜用户"+strconv.Itoa(userID)+"中奖",
		}
		messageChannel<-msg
	}


}
