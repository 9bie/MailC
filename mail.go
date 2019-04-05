package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"strings"
	"time"
)


type SendMail struct {
	User     string		`json:"user"`
	Password string		`json:"password"`
	Host     string		`json:"host"`
	Port     string		`json:"port"`
	Auth     smtp.Auth	`json:"auth"`
}

type Attachment struct {
	Name        []string `json:"name"`
	ContentType string	`json:"content_type"`
	WithFile    bool	`json:"with_file"`
}

type Message struct {
	From        string `json:"from"`
	To          []string `json:"to"`
	Cc          []string `json:"cc"`
	Bcc         []string `json:"bcc"`
	Subject     string 	`json:"subject"`
	Body        string	`json:"body"`
	ContentType string	`json:"content_type"`
	Attachment  Attachment `json:"attachment"`
}
//mail := &SendMail{user: "blackguwc@163.com", password: "", host: "smtp.163.com", port: "25"}
//message := Message{
//from:        "blackguwc@163.com",
//to:          []string{"blackguwc@163.com"},
//cc:          nil,
//bcc:         nil,
//subject:     "test",
//body:        "msg body!",
//contentType: "text/plain;charset=utf-8",
//attachment: Attachment{
//name:        []string{"mail.go"},
//contentType: "application/octet-stream",
//withFile:    true,
//},
//}
//err:= mail.Send(message)
//if err != nil{
//fmt.Println(err)
//}

func (mail *SendMail) Authed() {
	mail.Auth = smtp.PlainAuth("", mail.User, mail.Password, mail.Host)
}

func (mail SendMail) Send(message Message) error {
	mail.Authed()
	buffer := bytes.NewBuffer(nil)
	boundary := "GoBoundary"
	Header := make(map[string]string)
	Header["From"] = message.From
	Header["To"] = strings.Join(message.To, ";")
	Header["Cc"] = strings.Join(message.Cc, ";")
	Header["Bcc"] = strings.Join(message.Bcc, ";")
	Header["Subject"] = message.Subject
	Header["Content-Type"] = "multipart/related;boundary=\"" + boundary+"\""
	Header["Date"] = time.Now().String()
	mail.writeHeader(buffer, Header)


	if message.Attachment.WithFile {

		for _, graphname := range message.Attachment.Name {
			attachment := "\r\n--" + boundary + "\r\n"
			attachment += "Content-Transfer-Encoding:base64\r\n"
			attachment += "Content-Type:" + message.Attachment.ContentType  +"\r\n"
			attachment += "Content-Disposition:attachment; filename=\""+ graphname + "\""
			attachment += "Content-ID: <" + graphname + "> \r\n\r\n"
			buffer.WriteString(attachment)

			//拼接成html
			//imgsrc += "<p><img src=\"cid:" + graphname + "\" height=200 width=300></p><br>\r\n\t\t\t"

			defer func() {
				if err := recover(); err != nil {
					fmt.Printf(err.(string))
				}
			}()
			mail.writeFile(buffer, graphname)
		}
	}

	//需要在正文中显示的html格式
	//var template = `
    //<html>
    //    <body>
    //        <p>%s</p><br>
    //        %s
    //    </body>
    //</html>
    //`
	//var content = fmt.Sprintf(template, message.body, imgsrc)
	body := "\r\n--" + boundary + "\r\n"
	body += "Content-Type: text/html; charset=UTF-8 \r\n"
	body +=  message.Body
	buffer.WriteString(body)

	buffer.WriteString("\r\n--" + boundary + "--")
	fmt.Println(buffer.String())
	err := smtp.SendMail(mail.Host+":"+mail.Port, mail.Auth, message.From, message.To, buffer.Bytes())
	if (err!=nil){
		fmt.Println(err)
	}
	return nil
}

func (mail SendMail) writeHeader(buffer *bytes.Buffer, Header map[string]string) string {
	header := ""
	for key, value := range Header {
		header += key + ":" + value + "\r\n"
	}
	header += "\r\n"
	buffer.WriteString(header)
	return header
}

func (mail SendMail) writeFile(buffer *bytes.Buffer, fileName string) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	payload := make([]byte, base64.StdEncoding.EncodedLen(len(file)))
	base64.StdEncoding.Encode(payload, file)
	buffer.WriteString("\r\n")
	for index, line := 0, len(payload); index < line; index++ {
		buffer.WriteByte(payload[index])
		if (index+1)%76 == 0 {
			buffer.WriteString("\r\n")
		}
	}
}
func HandleMail(sender *MessageQueue){
	go func(){
		for{
			data,_ := sender.Dequeue(0)
			//mail := &SendMail{user: data.Auth.json, password: data.Auth.password, host: "smtp.163.com", port: "25"} data = &data.(JsonMessage)
			d := struct{
				JsonMessage

			}{
				data.(JsonMessage),
			}
			mail:= d.Auth
			message:= d.Data
			err:= mail.Send(message)
			if err != nil{
				fmt.Println(err)
			}
		}

	}()
}
