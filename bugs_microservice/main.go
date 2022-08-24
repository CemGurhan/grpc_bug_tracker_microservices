package main

import (
	"fmt"
	"log"
	"net"

	gen "github.com/cemgurhan/bugs_microservice/buf/proto/proto-gen-files"
	impl "github.com/cemgurhan/bugs_microservice/bug_impl"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	gRPC_Server := grpc.NewServer()
	bug_Server := &impl.BugServiceServer{}

	gen.RegisterBugServiceServer(gRPC_Server, bug_Server)

	address := ":8080"

	bugMicroservicePort, serveListenErr := net.Listen("tcp", address)

	if serveListenErr != nil {
		log.Fatalf("Could not listen to address %v", address)

	} else {
		fmt.Printf("Successfully listening to address %v", address)
	}

	reflection.Register(gRPC_Server)

	log.Fatalln(gRPC_Server.Serve(bugMicroservicePort))

}
