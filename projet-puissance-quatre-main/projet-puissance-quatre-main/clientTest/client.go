package main

import (
	"bufio"
	"log"
	"net"
)

func client(network string, port string) {
	conn, err := net.Dial(network, port)
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	defer conn.Close()

	log.Println("connexion établie")

	writer := bufio.NewWriter(conn)
	_, err = writer.WriteString("Bonjour" + "\n")
	if err != nil {
		log.Fatal("Error when writing :", err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatal("Error when flushing :", err)
	}

	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error when reading :", err)
	}
	log.Println(msg)
}

func main() {
	client("tcp", ":8080")
}
