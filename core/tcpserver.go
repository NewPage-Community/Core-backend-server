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
		log.Println("无法绑定监听:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	log.Printf("欢迎使用 NCS \n 后台服务正在运行 ...\n")

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

			res, size := GetRightMsg(string(data[0:c]))
			if size == -1 {
				log.Println(string(data[0:c]))
				return
			}

			for i := 0; i < size; i++ {
				go EventHandle(res[i], num, string(data[0:c]))
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
func GetRightMsg(msg string) ([32]string, int) {
	var json [32]string
	i := -1

	index := UnicodeIndex(msg, "}{")
	c := utf8.RuneCountInString(msg)

	defer func() {
		if err := recover(); err != nil {
			log.Println("Json分割错误:", err)
			log.Println("未处理数据:", msg)
			log.Println("处理数据:", json)
			log.Println("查找:", index)
			log.Println("长度:", c)
		}
	}()

	if index == -1 {
		json[0] = msg
		return json, 1
	}

	for i = 0; i < 32 && index != -1; i++ {
		json[i] = SubString(msg, 0, index+1)
		if c <= 0 {
			break
		}
		msg = SubString(msg, index+1, c)
		c = utf8.RuneCountInString(msg)
		index = UnicodeIndex(msg, "}{")
	}

	json[i] = msg

	return json, i + 1
}
