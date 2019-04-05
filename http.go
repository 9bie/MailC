package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)
type JsonMessage struct{
	Auth SendMail	 `json:"auth"`
	Data Message	 `json:"data"`
}
type StatucCode struct {
	Status      int32 `json:"status"`
	Message     string `json:"message"`
	Successful  bool  `json:"successful"`
}
func MailSend(sender *MessageQueue)http.Handler{
	// write to log
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err:=r.ParseForm()
		if err != nil {
		}
		rpc,ok:= r.Form["rpc"]

		fmt.Println(rpc)

		if ok==false{
			result := StatucCode{-1,"undefind",false}
			b,_ :=json.Marshal(result)
			w.Write(b)
		}else{
			result := StatucCode{0,"undefind",true}

			var jm JsonMessage
			if err := json.Unmarshal([]byte(rpc[0]), &jm); err == nil {
				sender.Enqueue(jm)
				b,_ :=json.Marshal(result)
				w.Write(b)
			}else{
				result := StatucCode{-2,"undefind",false}
				b,_ :=json.Marshal(result)
				w.Write(b)
			}
		}
	})
}
func Test(w http.ResponseWriter,r *http.Request){
	w.Write([]byte("Service is Running...."))
}
func MailControl(mailc *MailConf,mq *MessageQueue)http.Handler{
	return http.HandlerFunc(func (w http.ResponseWriter,r *http.Request){

		err:=r.ParseForm()
		if err != nil {//有错误关窝p事.jpg
		}
		c,ok1:=r.Form["contorl"]
		s,ok2:=r.Form["sstg"]
		if ok1{
			switch c[0] {
			case "stop"://只进不出
				mailc.status=-1
			case "close"://强行中断，不进也不出
				mq.INT(true)
			case "start"://恢复只进不出状态
				mailc.status=0
			case "open"://从中断状态恢复
				mq.INT(false)
			case "save":
				mq.Save()
			default:
				result := StatucCode{-5,"undefind",false}
				b,_ :=json.Marshal(result)
				w.Write(b)
			}
		}
		if ok2{
			i, err := strconv.Atoi(s[0])
			if err!=nil{
				result := StatucCode{-4,"undefind",false}
				b,_ :=json.Marshal(result)
				w.Write(b)
			}
			mailc.sstg=i
		}

	})
}

func GeiWoYeZhengYiGe(w http.ResponseWriter,r *http.Request){

}

func HttpHandle(port string,sender *MessageQueue,mailconfig *MailConf){
	//可以考虑换个RPC,http报文字节浪费太多了，json也是
	mux :=http.NewServeMux()
	mux.Handle("/api/mail_send", MailSend(sender))
	mux.Handle("/api/test",http.HandlerFunc(Test))
	mux.Handle("/api/Mail_Contorl",MailControl(mailconfig,sender))
	go func() {
		err:= http.ListenAndServe(":"+port, mux)
		if err!= nil{
			fmt.Println(err)
		}
	}()


}