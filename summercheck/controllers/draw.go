package controllers

import (
	"math/rand"
	"net/http"
	"strconv"
	"summercheck/server"
)

//思路是 拿去用户的ID  然后生成随机数，对应的ID则中奖



func Draw(w http.ResponseWriter,r *http.Request){
	//拿去保存user的map进行遍历 看有几个user
	//定义切片存放中奖人

	var slice []int
	index:=server.Num()
	//假设现在抽n个人中奖
	n:=r.FormValue("n")
	N,_:=strconv.Atoi(n)
	//如果想要抽n个人中奖却大于用户数
	if N>index {
		msg:="输入的抽奖人数大于已有人数"
		w.Write([]byte(msg))
	}
	for i:=0;i<N ;i++  {
		userID:=rand.Intn(index)
		slice=append(slice,userID)

		//然后将中奖名单公布
		msg:=server.Message{
			Content: "恭喜用户"+strconv.Itoa(userID)+"中奖",
		}
		messageChannel<-msg
	}


}
