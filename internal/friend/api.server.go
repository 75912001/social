package friend

import (
	protofriend "social/pkg/proto/friend"
)

type APIServer struct {
	protofriend.UnimplementedServiceServer
}
