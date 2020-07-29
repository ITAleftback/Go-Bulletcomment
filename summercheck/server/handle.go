package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	sync "sync"
)




type User struct {
	ID int//用户的标识
	BulletcommentChannel chan string//发送弹幕的通道
	MessageChannel chan string//用户收到的系统消息
}
func handleConn(conn net.Conn){
	defer conn.Close()

	//创建一个用户

	user:=&User{
		ID:             GenUserID(),
		BulletcommentChannel: make(chan string,8),
		MessageChannel: make(chan string,8),
	}

	msg:=Message{
		OwnerID: user.ID,
		Content: "用户："+strconv.Itoa(user.ID)+"加入了房间",
	}
	go sendBulletcomment(conn,user.BulletcommentChannel)
	go sendMessage(conn,user.MessageChannel)
	//将用户加入的消息广播
	user.MessageChannel<-"欢迎您的到来"
	messageChannel<-msg
	//把新加入的user放到enteringChannel，这样map就会添加保存了
	enteringChannel<-user

	//读取用户输入的东西
	input:=bufio.NewScanner(conn)
	for input.Scan(){
		msg.Content=input.Text()
		BulletcommentChannel<-msg
	}

	//如果用户离开，如法炮制，就像有用户加入一样
	leavingChannel<-user
	msg.Content="用户："+strconv.Itoa(user.ID)+"离开了房间"
	messageChannel<-msg


}
//用于给用户发送弹幕
func sendBulletcomment(conn net.Conn,ch <-chan string){
	for comment:=range ch{
		//敏感词处理
		keywords:=[]string{"傻逼","傻子","笨猪"}
		for _,keyword:=range keywords{
			comment=strings.ReplaceAll(comment,keyword,"**")
		}
		_, _ = fmt.Fprintln(conn, comment)
	}
}
//这个是发送的系统消息
func sendMessage(conn net.Conn,ch <-chan string){
	for msg:=range ch{
		_, _ = fmt.Fprintln(conn, msg)
	}
}


var (
	ID int
	Lock sync.Mutex
)
//因为可能同时有多个人加入，所以用到锁防止阻塞
func GenUserID()int  {
	Lock.Lock()
	ID++
	Lock.Unlock()
	return ID
}
