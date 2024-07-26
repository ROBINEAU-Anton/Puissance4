package main

import (
	"bufio"
	"log"
	"net"
)

func (g *game) connect() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	g.conn = conn
	go g.read()
	go g.write()
}

func (g *game) read() {
	reader := bufio.NewReader(g.conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Error when reading :", err)
			continue
		}
		g.reader <- msg
	}
}

func (g *game) write() {
	writer := bufio.NewWriter(g.conn)
	for {
		select {
		case str := <-g.writer:
			_, err := writer.WriteString(str + "\n")
			if err != nil {
				log.Fatal("Error when writing :", err)
			}
			err = writer.Flush()
			if err != nil {
				log.Fatal("Error when flushing :", err)
			}
		default:
		}
	}
}
