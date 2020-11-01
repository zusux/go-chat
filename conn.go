package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type connection struct{
	// websocket connection
	ws *websocket.Conn
	send chan []byte
	uid string
	username string
	img string
}

func (c *connection)reader()  {
	for {
		msgType,message,err := c.ws.ReadMessage()
		if err != nil{
			fmt.Println(fmt.Sprintf("readMessage err %s",c.uid),err)
			break
		}
		switch msgType {
			case 1:
				act,err := AnalysisData(message)
				if err != nil{
					fmt.Println("解析message 出错:",err)
					continue
				}else{
					
				}

				if act.Action == "message"{
					toUser := act.Message.ToUser
					if toUser != "*"{
						if c,ok := h.connections[toUser];ok {
							data := SendAction{
								Action: "message",
								Message: Message{
									FromUser: act.Message.FromUser,
									FromUserImg: act.Message.FromUserImg,
									Username:  act.Message.Username,
									ToUser: act.Message.ToUser,
									MsgCode: NormalCode,
									Content:  act.Message.Content,
								},
							}
							databyte,_:= json.Marshal(data)
							c.send <- databyte
						}
					}else{
						BroadMsg("broadcast",act.Message)
					}
				}else if act.Action == "registor"{
					BroadMsg("registor",act.Message)
				}else if act.Action == "unregistor"{
					BroadMsg("unregistor",act.Message)
				}else if act.Action == "set"{

					username,ok := act.Data["username"]
					if ok {
						c.username = username.(string)
					}
					img,ok := act.Data["img"]
					if ok {
						c.img = img.(string)
					}
					fmt.Println("设置信息",c.uid,username,img)
					d := make(map[string]interface{})
					d["user"] = User{
						UserId: c.uid,
						Name: c.username,
						Img: c.img,
					}
					da := SendAction{
						Action: "registor",
						Data:d,
					}
					msgbyte ,_ := json.Marshal(da)
					h.broadcast <- msgbyte
					fmt.Println("广播注册信息...")
				}

		}
	}
	c.ws.Close()
}

func AnalysisData(message []byte) (SendAction,error) {
	action := SendAction{}
	err := json.Unmarshal(message,&action)
	if err != nil{
		return action,err
	}
	return action,nil
}

func BroadMsg(action string,msgStruct Message)  {
	msg := Message{
		FromUser: msgStruct.FromUser,
		FromUserImg: msgStruct.FromUserImg,
		Username:  msgStruct.Username,
		ToUser: "*",
		Content: msgStruct.Content,
		MsgCode: msgStruct.MsgCode,
	}

	act := SendAction{
		Action: action,
		Message: msg,
	}

	msgbyte ,_ := json.Marshal(act)
	h.broadcast <- msgbyte
}

func (c *connection)writer()  {
	for message := range c.send{
		err := c.ws.WriteMessage(websocket.TextMessage,message)
		if err != nil{
			break
		}
	}
	c.ws.Close()

}

var upGrader = &websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}


func wsPage(c *gin.Context){
	ws,err := upGrader.Upgrade(c.Writer,c.Request,nil)
	if err != nil{
		return
	}
	defer ws.Close()
	uid := uuid.NewV4().String()
	conn := &connection{
		ws:ws,
		send: make(chan[]byte,256),
		uid:uid,
	}
	fmt.Println("建立连接",uid)
	h.register <- conn
	defer func() {
		fmt.Println("关闭连接",uid)
		h.unregister <- conn
	}()
	fmt.Println("发送数据1...")
	go conn.writer()
	d := make(map[string]interface{})
	d["uid"] = uid
	d["users"] = h.GetUsers()
	a := SendAction{
		Action: "registor",
		Data: d,
		Message: Message{
			FromUser: uid,
		},
	}
	sendByte,_ := json.Marshal(a)
	fmt.Println("发送数据2...")
	conn.send <- sendByte
	fmt.Println("发送数据3...")
	conn.reader()
}
