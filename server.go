package main

import (
	"github.com/ProjectAthenaa/proxy-rater/service"
	protos "github.com/ProjectAthenaa/sonic-core/protos/proxy-rater"
	"github.com/prometheus/common/log"
	"net"
	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", "3000")
	if err != nil{
		log.Fatalln("start listener: ", err)
	}

	server := grpc.NewServer()

	protos.RegisterProxyRaterServer(server, service.Server{} )

	if err = server.Serve(listener); err != nil{
		log.Fatalln(err)
	}

}
