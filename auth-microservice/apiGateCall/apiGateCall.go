package apigatecall

import (
	"context"
	"fmt"

	usercontext "github.com/cemgurhan/auth-microservice/context"
)

func CallApiGateway(ctx context.Context) {

	fmt.Printf("Context user: %v", usercontext.UserFromContext(ctx))

}
