package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/urfave/cli"
)

const serverip string = "127.0.0.1:8080"

func main() {
	cc := make(chan bool)
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "name,n",
			Value: "default",
			Usage: "name set",
		},
	}

	app.Action = func(c *cli.Context) error {

		name := c.String("name")
		if name == "" {
			fmt.Println("name is invalid")
			return nil
		}

		conn, err := net.Dial("tcp", serverip)

		if err != nil {
			panic(err)
		}

		if _, err = conn.Write([]byte(name)); err != nil {
			panic(err)
		}
		go func() {
			defer conn.Close()
			for {
				var readstr string
				Scanf(&readstr)
				if _, err = conn.Write([]byte(readstr)); err != nil {
					panic(err)
				}
			}
		}()
		go func() {
			bytes := make([]byte, 1024)
			var n int
			var err error
			for {
				if n, err = conn.Read(bytes); err != nil {
					panic(err)
				}
				if n > 0 {
					fmt.Println(string(bytes[0:n]))
				}
			}
		}()
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	defer close(cc)
	<-cc
}

func Scanf(a *string) {
	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()
	*a = string(data)
}
