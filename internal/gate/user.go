package gate

import (
	"google.golang.org/grpc"
	libactor "social/lib/actor"
)

type User struct {
	*libactor.Normal
	Stream grpc.ServerStream
}
