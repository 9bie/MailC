package main

import "time"

func Helper(){


}
func main(){
	var mq MessageQueue
	mq.Init()

	HttpHandle("8081",&mq)
	HandleMail(&mq)
	for{
		time.Sleep(2000000)
	}


}