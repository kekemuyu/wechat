package main

import (
	"container/list"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/urfave/cli"
)

const serverPort string = "8080"

type Client struct {
	Name string   `json:name` //昵称
	IP   string   `json:ip`   //ip
	Conn net.Conn `json:conn`
}

var (
	Clients = list.New()
)

func main() {
	cc := make(chan bool)
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "port,p",
			Value: "8080",
			Usage: "server port set",
		},
	}

	app.Action = func(c *cli.Context) error {
		port := c.String("port")
		tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+port)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		listener, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		fmt.Println("listen:" + port)
		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					continue
				}

				bytes := make([]byte, 1024)
				var n int
				if n, err = conn.Read(bytes); (n == 0) || (err != nil) { //读取昵称
					conn.Close()
					continue
				}
				fmt.Println("connected client:" + string(bytes[0:n]))
				cli := Client{
					Name: string(bytes[0:n]),
					IP:   conn.RemoteAddr().String(),
					Conn: conn,
				}
				fmt.Println("client:", cli)
				Clients.PushBack(cli)
				go handleClient(cli)
			}
		}()
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	<-cc
}

func handleClient(cli Client) {
	defer cli.Conn.Close()
	for {
		bytes := make([]byte, 1024)

		n, err := cli.Conn.Read(bytes)
		if err != nil {
			cli.Conn.Close()
			listDelet(cli)
			continue
		} else if n <= 0 {
			continue
		}
		tempstr := "\n" + time.Now().Format("2006-01-02 15:04:05") + "\n" + cli.Name + ": " + string(bytes[0:n]) + "\n"

		broadcast([]byte(tempstr), cli)
	}
}

func broadcast(msg []byte, client Client) {
	for cli := Clients.Front(); cli != nil; cli = cli.Next() {
		if cli.Value.(Client) == client {
			continue
		}
		fmt.Println(string(msg))
		cli.Value.(Client).Conn.Write(msg)
	}
}

//查找cli，并删除
func listDelet(client Client) {
	for cli := Clients.Front(); cli != nil; cli = cli.Next() {
		if client == cli.Value.(Client) {
			Clients.Remove(cli)
		}
	}
}
