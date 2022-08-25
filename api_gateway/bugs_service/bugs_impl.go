package bugsservice

import (
	"context"

	gen "github.com/cemgurhan/api_gateway/buf/gen/proto_files"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BugServiceServer struct {
	gen.UnimplementedBugServiceServer
}

func getBugServiceClientConn() (gen.BugServiceClient, *grpc.ClientConn, error) {

	bugMicroserviceAddress := ":8080"

	conn, dialError := grpc.Dial(bugMicroserviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if dialError != nil {
		return nil, nil, dialError
	}

	client := gen.NewBugServiceClient(conn)

	return client, conn, nil

}

func (server *BugServiceServer) CreateBug(ctx context.Context, req *gen.CreateBugReq) (*gen.CreateBugRes, error) {

	client, conn, err := getBugServiceClientConn()

	if err != nil {
		return &gen.CreateBugRes{
			Succes:  false,
			Message: "Unable to connect to bug microservice",
		}, err
	}

	defer conn.Close()

	resp, createBugErr := client.CreateBug(ctx, req)

	if createBugErr != nil {
		return &gen.CreateBugRes{
			Succes:  false,
			Message: "Unable to create bug",
		}, createBugErr
	}

	return resp, nil

}
