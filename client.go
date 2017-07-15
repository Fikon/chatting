package main

import(
    "fmt"
    "net"
    "os"
    )

type Message struct{
    from int
    to int
    status int
    userName string
    content string
}

type Client struct{
    conn *net.UDPConn
    msgChan chan Message
    rmsgChan chan string
    userID int
    userName string
}

func (c *Client)FSM(){
    for{
	fmt.Println("send msg(s) or quit(q)")
	var str string
	if fmt.Scanln(&str); str == "q"{
            msg := fmt.Sprintf("%d#%d", 2, c.userID)
            c.conn.Write([]byte(msg))
	    c.conn.Close()
	    return

        }else{
	    fmt.Println("input the userID to receive your msg")
	    var userID int
	    var userName string
	    fmt.Scanln(&userID)
	    fmt.Println("input the message you want to send")
	    fmt.Scanln(&userName)
            msg := fmt.Sprintf("%d#%d#%s", 1, userID, userName)
            c.conn.Write([]byte(msg))
	}
    }
}

func (c *Client)readMSG(){
    var buf [512]byte
    for{
        n,_,err := c.conn.ReadFromUDP(buf[0:])
	if err != nil{
	    fmt.Println("error read from udp", err.Error())
	    return
	}
        msg := string(buf[0:n])
        fmt.Println("receive a msg: ", msg)
    }
}



func main(){
    if len(os.Args) != 2{
	fmt.Println("usage: %s, host:port",os.Args[0])
	os.Exit(1)
    }
    server := os.Args[1]
    addr, err := net.ResolveUDPAddr("udp4", server)
    if err != nil{
	fmt.Println("error resolving udp addr: ", err.Error())
	os.Exit(2)
    }
    var c Client
    fmt.Println("input your userid: ")
    _, err = fmt.Scanln(&c.userID)
    if err != nil{
	fmt.Println("error scaning userid:",err.Error())
	os.Exit(3)
    }
    fmt.Println("input your userName:")
    _, err = fmt.Scanln(&c.userName)
    if err != nil{
	fmt.Println("error scaning userName:",err.Error())
	os.Exit(4)
    }
    c.conn, err = net.DialUDP("udp", nil, addr)
    if err != nil{
	fmt.Println("error dialing udp:",err.Error())
	os.Exit(5)
    }
    defer c.conn.Close()
    msg := fmt.Sprintf("%d#%d#%s", 0, c.userID, c.userName)
    c.conn.Write([]byte(msg))
    go c.readMSG()
    c.FSM()
}


