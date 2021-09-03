package main

import (
	"github.com/ProjectAthenaa/proxy-rater/service"
	protos "github.com/ProjectAthenaa/sonic-core/protos/proxy-rater"
	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	var listener net.Listener
	var err error

	if os.Getenv("DEBUG") == "1"{
		listener, err = net.Listen("tcp", ":5000")
		if err != nil{
			log.Fatalln("start listener: ", err)
		}
	}else{
		listener, err = net.Listen("tcp", ":3000")
		if err != nil{
			log.Fatalln("start listener: ", err)
		}
	}

	server := grpc.NewServer()

	protos.RegisterProxyRaterServer(server, service.Server{} )

	if err = server.Serve(listener); err != nil{
		log.Fatalln(err)
	}

}
