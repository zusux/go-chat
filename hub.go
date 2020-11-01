package main

import (
	"encoding/json"
	"fmt"
)

type hub struct{
	connections map[string]*connection
	broadcast chan []byte
	register chan *connection
	unregister chan *connection
}

var h = hub{
	broadcast: make(chan []byte,1000),
	register : make(chan *connection,1000),
	unregister: make(chan *connection,1000),
	connections: make(map[string]*connection),
}

type MsgCode int
const (
	RegistorCode MsgCode = iota
	UnRegistorCode
	NormalCode
	BroadCaseCode
)

type SendAction struct{
	Action string `json:"action"`
	Message Message `json:"message"`
	Data map[string]interface{} `json:"data"`
}

type User struct{
	UserId string `json:"userId"`
	Img string `json:"img"`
	Name string `json:"name"`
	Content string `json:"content"`
}

type Message struct{
	FromUser string `json:"from_user"`
	FromUserImg string `json:"from_user_img"`
	Username string `json:"username"`
	ToUser string	`json:"to_user"`
	MsgCode MsgCode	`json:"msg_code"`
	Content string `json:"content"`
}

func(h *hub) run(){
	for{
		select{
			case c:= <- h.register:
				h.connections[c.uid] = c
			case c := <- h.unregister:
				if _,ok := h.connections[c.uid];ok{
					delete(h.connections,c.uid)
					close(c.send)


					d := make(map[string]interface{})
					d["user"] = User{
						UserId: c.uid,
						Name: c.username,
						Img: c.img,
					}
					da := SendAction{
						Action: "unregistor",
						Data:d,
					}
					msgbyte ,_ := json.Marshal(da)
					h.broadcast <- msgbyte
					fmt.Println("广播注销信息...")
				}



			case m := <- h.broadcast:
				fmt.Println("广播信息")
				for uid,c := range h.connections{
					select{
						case c.send <- m:
					default:
						delete(h.connections,uid)
						close(c.send)
						d := make(map[string]interface{})
						d["user"] = User{
							UserId: c.uid,
							Name: c.username,
							Img: c.img,
						}
						da := SendAction{
							Action: "unregistor",
							Data:d,
						}
						msgbyte ,_ := json.Marshal(da)
						h.broadcast <- msgbyte
						fmt.Println("广播注销信息...")
					}
			}
		}
	}
}

func (h *hub)GetUsers() []User {
	res := make([]User,0)
	for _,conn := range h.connections{
		if conn.username != ""{
			user := User{
				UserId: conn.uid,
				Img: conn.img,
				Name: conn.username,
			}
			res = append(res,user)
			fmt.Println("发现用户数据...")
		}
	}
	return res
}