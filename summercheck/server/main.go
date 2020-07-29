package main

import (
	"log"
	"net"
)

func main(){
	listener,err:=net.Listen("tcp",":8080")
	if err != nil {
		panic(err)
	}

	//用于记录用户，以及将他的弹幕进行广播
	go broadcaster()

	for {
		conn,err:=listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}

}
