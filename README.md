# 引言

遇到的问题：

1.因为对websocket不是太熟，再加上时间紧凑，所以我是基于TCP连接来写。

2.因为弹幕是有滚动效果，前端知识储备不够的我只能想办法让已经发出的消息保留3-5秒实现（想法是如此，但是我写出来的感觉跟聊天室没区别:joy:）。



# 功能

开了三个客户端，下面上效果

客户端1发送消息

![image-20200729223454731](C:\Users\Mechrevo\AppData\Roaming\Typora\typora-user-images\image-20200729223454731.png)

客户端2显示

![image-20200729223507077](C:\Users\Mechrevo\AppData\Roaming\Typora\typora-user-images\image-20200729223507077.png)

客户端3显示



![image-20200729223519091](C:\Users\Mechrevo\AppData\Roaming\Typora\typora-user-images\image-20200729223519091.png)

## 敏感词处理

```
for comment:=range ch{
   //敏感词处理
   keywords:=[]string{"傻逼","傻子","笨猪"}
   for _,keyword:=range keywords{
      comment=strings.ReplaceAll(comment,keyword,"**")
   }
   _, _ = fmt.Fprintln(conn, comment)
}
```

ch是存放用户弹幕的通道，遍历的时候查看有没有敏感词，有的话用”**“取而代之。

效果如图

用户1发送消息：

![image-20200729223718955](C:\Users\Mechrevo\AppData\Roaming\Typora\typora-user-images\image-20200729223718955.png)

用户2收到的：

![image-20200729223733374](C:\Users\Mechrevo\AppData\Roaming\Typora\typora-user-images\image-20200729223733374.png)

用户3收到的

![image-20200729225312146](C:\Users\Mechrevo\AppData\Roaming\Typora\typora-user-images\image-20200729225312146.png)

## 彩色弹幕

思路是color包，用户在输入的消息后面加上自己想要的颜色，然后读取的时候扫描颜色关键字，然后将彩色的字体装入弹幕通道。

很遗憾，学识浅陋，只能纸上谈兵。

## 弹幕红包&红包抽奖（我不怎么玩直播，觉得这两个是一个东西）

```
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

func Num()int{
	//index 用于计数有多少user
	var index int
	for _,_=range users{
		index++
	}
	return index
}
```

在想办法向所有用户发送消息时，我创建了map用来保存加入的用户。

所以我的抽奖思路是，通过Num函数遍历map得知有几个user，用index计数。然后从客户端接收想要设置的中奖人数n，当然n肯定要小于等于index，不然会报错。用for循环n次，每一次用rand.Init创建从1到index个随机数。为了优化，我用到切片保存，倒是再拿切片出来遍历，把中奖的名单消息塞入消息通道进行广播。

当然这样的思路有两个瑕疵：

1.有可能会有人重复中奖（解决方法：每一次有新的中奖人就遍历一次切片查看是否已中奖）。

2.如果用户2离开了房间，房间只有用户1及用户3，但仍然有可能是用户2中奖（还没想到解决办法）。

## 进阶

虽然我没实现但是我有思路:joy:

### 高并发

```
func main(){
   for i:=0;i<50 ;i++  {
      Gopost()
      time.Sleep(10 *time.Millisecond)
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
```

开了50个客户端，发现服务端向客户端发送欢迎的消息有点忙不过来了

![image-20200729230324130](C:\Users\Mechrevo\AppData\Roaming\Typora\typora-user-images\image-20200729230324130.png)

为什么会出现同一个用户多次加入呢？原本想的是阻塞的缘故，可当我加了锁后

```
func GenUserID()int  {
   Lock.Lock()
   ID++
   Lock.Unlock()
   return ID
}
```

仍然如此。

现在的我还没想出原因。

### 弹幕频率

想了下，如果是滚动的弹幕。想控制频率的话，只需加快滚动频率或者减少滚动频率。也就是说我可以设置弹幕存在较长的时间或者是较慢的时间来控制频率。由于我是不怎么会前端的，所以不用设置滚动效果也就更简单了

### 黑名单

很简单，创建map，把想要加入黑名单的放入这个map即可，同时将存放user的map将他删除，这样他就不能发言了（其实就是把他踢出房间:joy:）

# 最大最大的问题

本项目本应是gin+websocket实现，但犹豫对于websocket的不熟悉，我基于TCP来写。这也导致了最大的问题：无法注册路由。在尝试开路由的时候，我遇到了诸多问题如红包抽奖的接口无法使用server里面的map，看了网上的源码发现还是需要ws实现。（回去一定要好好看websocket）

