package main

import (
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main(){
	for i:=0;i<50 ;i++  {
		Gopost()
		time.Sleep(5 *time.Second)
	}
}
func Gopost(){
	go func() {
		conn,err:=net.Dial("tcp",":8080")
		if err != nil {
			panic(err)
		}


		done:=make(chan struct{})
		go func() {
			_, _ = io.Copy(os.Stdout, conn)
			log.Println("done")
			done<- struct{}{}
		}()

		mustCopy(conn,os.Stdin)
		conn.Close()
		//关闭
		<-done
	}()
}

func mustCopy(dst io.Writer,src io.Reader)  {
	if _,err:=io.Copy(dst,src);err!=nil {
		log.Fatal(err)
	}
}