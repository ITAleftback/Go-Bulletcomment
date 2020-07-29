package main


var (
	// 新用户到来，通过该 channel 进行登记
	enteringChannel = make(chan *User)
	// 用户离开，通过该 channel 进行登记
	leavingChannel = make(chan *User)
	// 广播弹幕的 channel
	BulletcommentChannel = make(chan Message, 8)
	//
	messageChannel = make(chan Message,8)

	users=make(map[*User]struct{})
)
// 发送消息的时候，如果是自己的话会看见两次，这是因为服务器遍历了所有用户 给发消息的人又发了一次
//所以可以用ID判断 如果ID相等就不发了
type Message struct {
	OwnerID int
	Content string
}


//记录用户， 并将用户的弹幕发送
func broadcaster()  {
	//用map保存用户清单，加入的时候添加 退出的时候删除

	//map的key 一开始尝试直接写User 但后续添加的时候报了错
	//想了下， 如果不是指针的话 map会为User的成员每个分配内存这是不合理的
	//所以用到指针


	for  {
		select {
		case user:=<-enteringChannel:
			//说明有新用户加入 往map添加
			users[user]= struct{}{}
		case user:=<-leavingChannel:
			//删除map里的user即可
			delete(users,user)
		case Bulletcomment:=<-BulletcommentChannel:
			//如果通道有东西，广播
			for user:=range users{
				if user.ID==Bulletcomment.OwnerID{
					//
					continue
				}
				user.BulletcommentChannel<-Bulletcomment.Content
			}
		case msg:=<-messageChannel:
			//如果通道有消息，广播
			for user:=range users{
				user.MessageChannel<-msg.Content
			}

		}
	}

}
//

func Num()int{
	//index 用于计数有多少user
	var index int
	for _,_=range users{
		index++
	}
	return index
}
