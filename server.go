package main

import(
    "fmt"
    "net"
    "strings"
    "strconv"
    )

const port = ":4399"

type User struct{
    userID int
    userName string
    userAddr *net.UDPAddr
}

type Server struct{
    conn *net.UDPConn
    msgChan chan Message
    Users map [int]User
}

type Message struct{
    from int
    to int
    status int
    userName string
    content string
}


func (s *Server) handleMSG(){
    var buf [512]byte
    n, addr, err := s.conn.ReadFromUDP(buf[0:])
    if err != nil{
	fmt.Println("error reading from udp", err.Error())
	return
    }
    msg := string(buf[0:n])
    fmt.Println("got a msg : ", msg)
    message := s.parseMSG(msg)
    switch message.status{
	case 0:
	    var user User
	    user.userAddr = addr
	    user.userName = message.userName
	    user.userID = message.from
	    s.Users[user.userID] = user
	    fmt.Println("a new user come in: ", user.userName)
	case 1:
	    fmt.Println("case 1")
            if _, ok := s.Users[message.to]; ok{
		s.msgChan <- message
	    }else{
		fmt.Println("invalid userid")
	    }
	case 2:
	    fmt.Println("case 2")
	    delete(s.Users, message.from)
	default:
		fmt.Println("fuck you, man!!!!")
    }
}

func (s *Server) parseMSG(msg string) (message Message){
    strs := strings.Split(msg, "#")
    status,_ := strconv.Atoi(strs[0])
    message.status = status
    switch status{
	case 0:
	    fmt.Println("new user")
	    message.from,_ = strconv.Atoi(strs[1])
	    message.userName = strs[2]
	    return
	case 1:
	    fmt.Println("send msg")
	    message.to,_ = strconv.Atoi(strs[1])
	    message.content = strs[2]
	    return
	case 2:
	    fmt.Println("user quit")
	    message.from,_ = strconv.Atoi(strs[1])
	    return
    }
    return
}

func (s *Server) sendMSG(message Message){
    addr := s.Users[message.to].userAddr
    fmt.Println("the msg to send: ", message.content)
    fmt.Printf("the addr to send: %+v", addr)
    _, err := s.conn.WriteToUDP([]byte(message.content), addr)
    if err != nil{
	fmt.Println("error writing to udp: ", err.Error())
	return
    }
    fmt.Println("send msg done")
}

func (s *Server)send(){
    for{
        msg := <- s.msgChan
        fmt.Println(msg.content)
        s.sendMSG(msg)
    }
}


func main(){
    addr, err := net.ResolveUDPAddr("udp", port)
    if err != nil{
	fmt.Println("error resolve udp addr: ", err.Error())
	return
    }
    var s Server
    s.Users = make(map[int]User, 0)
    s.conn, err = net.ListenUDP("udp", addr)
    if err != nil{
	fmt.Println("error listen udp: ", err.Error())
	return
    }
   s.msgChan = make(chan Message)
   go s.send()
   for{
       s.handleMSG()
   }

}


