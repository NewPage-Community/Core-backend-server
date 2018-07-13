package core

import (
	"log"
	"net"
	"os"
	"unicode/utf8"
)

//MaxConnects ...
var MaxConnects = 100

var sersChan = make(map[int]chan string)
var serverNum = 0

//StartTCPServer ...
func StartTCPServer() {
	listener, err := net.Listen("tcp", GetConfig("listenbind"))
	CheckError(err)

	if err != nil {
		log.Println("error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	log.Printf("running ...\n")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		if sersChan[serverNum] == nil {
			go handleClient(conn, serverNum)
		} else {
			log.Println("连接数已达最大限度")
		}

		if serverNum >= 200 {
			serverNum = 0
		} else {
			serverNum++
		}
	}
}

func handleClient(conn net.Conn, num int) {
	var closed = make(chan bool)

	sersChan[num] = make(chan string)

	defer func() {
		conn.Close()
		log.Printf("客户端 %d 关闭连接", num)
		delete(sersChan, num)
	}()

	go func() {
		for {
			data := make([]byte, 4096)
			c, err := conn.Read(data)
			if err != nil {
				log.Println("无法读取数据：", err)
				closed <- true
				return
			}

			//log.Println(string(data[0:c]))

			res, size := GetRightMsg(string(data[0:c]))
			if size == -1 {
				log.Println(string(data[0:c]))
				return
			}

			for i := 0; i < size; i++ {
				go EventHandle(res[i], num)
			}
		}
	}()

	go func() {
		for {
			sendString := <-sersChan[num]
			_, err := conn.Write([]byte(sendString))
			if err != nil {
				closed <- true
				return
			}
		}
	}()

	for {
		if <-closed {
			return
		}
	}
}

//GetRightMsg ...
func GetRightMsg(msg string) ([8]string, int) {
	var res [8]string
	i := -1

	defer func() {
		if err := recover(); err != nil {
			log.Println("Error:", err)
		}
	}()

	index := UnicodeIndex(msg, "}{")
	c := utf8.RuneCountInString(msg)

	if index == -1 {
		res[0] = msg
		return res, 1
	}

	for i = 0; i < 8 && index != -1; i++ {
		res[i] = SubString(msg, 0, index+1)
		if c <= 0 {
			break
		}
		msg = SubString(msg, index+1, c)
		c = utf8.RuneCountInString(msg)
		index = UnicodeIndex(msg, "}{")
	}

	res[i] = msg

	return res, i + 1
}
