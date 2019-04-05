package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
func HttpHandle(port string,sender *MessageQueue){
	mux :=http.NewServeMux()
	mux.Handle("/api/mail_send", MailSend(sender))
	mux.Handle("/api/test",http.HandlerFunc(Test))
	go func() {
		err:= http.ListenAndServe(":"+port, mux)
		if err!= nil{
			fmt.Println(err)
		}
	}()


}