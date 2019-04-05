package main

import "time"

func Helper(){


}
func main(){
	var mq MessageQueue
	mailc:=MailConf{sstg:0}
	mq.Init()

	HttpHandle("8081",&mq,&mailc)
	HandleMail(&mq,&mailc)
	for{
		time.Sleep(2000000)
	}


}