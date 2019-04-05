// message queue
package main
//数据安全不关我的事.jpg
//说起来我为啥要整这个，直接用channel不好么emmmmm
import (
    "github.com/kataras/iris/core/errors"
    "time"
)

type Node struct{
    nextp *Node
    value interface{}
    lastp *Node
}
type MessageQueue struct{
    init bool
    length int
    now *Node//当前指针
    head *Node
    in chan interface{}// 传入数据通道
    out chan interface{} // 传出数据通道
    flag bool//这个我忘了是干啥用的了
    flag_int bool//中断标识
}
func (this *MessageQueue)Init()error{
    this.in = make(chan interface{})
    this.out = make(chan interface{})
    this.now = nil
    this.flag = false
    this.length = 0
    this.flag_int = false
    go func(){
        for{
            select {
            case value:=<- this.in:
                node := &Node{value:value,nextp:this.now,lastp:nil}
                if this.flag == true{
                    this.out<-value
                }
                if this.length == 0{
                    this.head = node
                    this.now = node
                }else{
                    this.now = node
                    this.now.nextp.lastp = this.now
                }
                this.length++
            }
        }
    }()
    this.init = true
    return nil
}
//查询中断状态
func (this *MessageQueue)Status()bool{
    return this.flag_int
}
//查询队列内容
func (this *MessageQueue)Length()int{
    return this.length
}
// 非堵塞，不管是否成功
func (this *MessageQueue)Enqueue(data interface{}){


    //var ok chan  bool
    go func() {
        select{
        case this.in<-data:
        }

    }()


}
//堵塞，等待结果，这个东西似乎没啥用.jpg
func (this *MessageQueue)Enqueue2(data interface{},timeout int)error{

    var errorz error
    ok:=make(chan bool)



    go func() {
       select{
       case this.in<-data:
          errorz = nil
          ok<-true
       }

    }()
    select{
    case flag:=<-ok:
        if flag==true{
            return nil
        }else{
            errorz = errors.New("Have a error")
         return errorz
        }
    }


}

func (this *MessageQueue)Dequeue(timeoutSecs int)(interface{},error){
    var errorz error
    if this.flag_int{
        errorz = errors.New("INT")
        return nil,errorz
    }
    if this.length == 0 && timeoutSecs!=0{
        time.Sleep(time.Duration(timeoutSecs))
        if this.length == 0{
            errorz = errors.New("queue is null")
            return nil,errorz
        }

    }else if timeoutSecs ==0 {
        for{

            time.Sleep(200)
            if this.length !=0{
                break//手动堵塞
            }else if this.flag_int{
                break
            }
        }
    }

    value:= this.head
    this.head = this.head.lastp
    this.length--
    return value.value,nil
}
//强制中断
func (this *MessageQueue)INT(yes_or_no bool){
    if yes_or_no{
        this.flag_int = true
    }else{
        this.flag_int = false
    }

}
func (this *MessageQueue)Clear(){
    this.head = nil
    this.now = nil
    this.length =0
    }
func (this *MessageQueue)Save(){

}
func (this *MessageQueue)Load() {

}