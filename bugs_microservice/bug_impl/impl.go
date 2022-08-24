package bugimpl

import (
	"context"

	gen "github.com/cemgurhan/bugs_microservice/buf/proto/proto-gen-files"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BugServiceServer struct {
	gen.UnimplementedBugServiceServer
}

func (server *BugServiceServer) CreateBug(ctx context.Context, req *gen.CreateBugReq) (*gen.CreateBugRes, error) {

	bugName := req.GetBug().GetName()
	bugType := req.GetBug().GetId()

	client, newClientErr := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	if newClientErr != nil {
		return (&gen.CreateBugRes{
			Succes:  false,
			Message: "Unable to establish new mongodb client",
		}), newClientErr

	}

	dbConnectErr := client.Connect(ctx)

	if dbConnectErr != nil {
		return (&gen.CreateBugRes{
			Succes:  false,
			Message: "Unable to connect to mongoDB",
		}), dbConnectErr
	}

	pingError := client.Ping(ctx, nil)

	if pingError != nil {
		return (&gen.CreateBugRes{
			Succes:  false,
			Message: "MongoDB server unresponsive",
		}), pingError

	}

	bugTrackerDatabase := client.Database("bug-tracker-db")
	bugsCollection := bugTrackerDatabase.Collection("bugs")

	bugsCollection.InsertOne(context.Background(), bson.D{
		{Key: "name", Value: bugName},
		{Key: "type", Value: bugType},
	})

	return &gen.CreateBugRes{
		Succes:  true,
		Message: "Added",
	}, nil

}
