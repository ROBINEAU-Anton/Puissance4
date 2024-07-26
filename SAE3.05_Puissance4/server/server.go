package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	conn1, err := listener.Accept()
	if err != nil {
		log.Fatal("accept error:", err)
	}
	log.Println("Un client s'est connecté")

	conn2, err := listener.Accept()
	if err != nil {
		log.Fatal("accept error:", err)
	}
	log.Println("Un client s'est connecté")

	reader1 := bufio.NewReader(conn1)
	msg1, err := reader1.ReadString('\n')
	log.Println(msg1)

	writer1 := bufio.NewWriter(conn1)
	_, err = writer1.WriteString("GO\n")
	if err != nil {
		log.Println("Error when writing :", err)
	}
	err = writer1.Flush()
	if err != nil {
		log.Println("Error when flushing :", err)
	}

	reader2 := bufio.NewReader(conn2)
	msg2, err := reader2.ReadString('\n')

	log.Println(msg2)

	writer2 := bufio.NewWriter(conn2)
	writer2.WriteString("GO\n")
	writer2.Flush()

	log.Println("connexion établie")

	chanReader1 := make(chan string, 4)
	chanReader2 := make(chan string, 4)
	go func() {
		for {
			str1, err1 := reader1.ReadString('\n')
			if err1 != nil {
				log.Fatal("Error when reading : ", err1)
			}
			log.Println("msg 1 : " + str1)
			chanReader1 <- str1
		}
	}()

	go func() {
		for {
			str2, err2 := reader2.ReadString('\n')
			if err2 != nil {
				log.Fatal("Error when reading : ", err2)
			}
			log.Println("msg 2 : " + str2)
			chanReader2 <- str2
		}
	}()

	for {
		player := rand.Intn(2)
		select {
		case msg1 = <-chanReader1:
			msg1 = strings.ReplaceAll(msg1, "\n", "")
			if msg1 == "play" {
				writer2.WriteString(fmt.Sprint(player) + "\n")
				writer2.Flush()
				if player == 1 {
					writer1.WriteString("0\n")
					writer1.Flush()
				} else {
					writer1.WriteString("1\n")
					writer1.Flush()
				}
			} else {
				msg1 += "\n"
				_, err = writer2.WriteString(msg1)
				if err != nil {
					log.Println("Error : ", err)
				}
				err = writer2.Flush()
				if err != nil {
					log.Println("Error : ", err)
				}
			}
		case msg2 = <-chanReader2:
			msg2 = strings.ReplaceAll(msg2, "\n", "")
			if msg2 == "play" {
				continue
			} else {
				msg2 += "\n"
				writer1.WriteString(msg2)
				writer1.Flush()
			}
		default:
		}
	}
}
