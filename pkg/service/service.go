package service

import (
	"google.golang.org/grpc"
	libbench "social/lib/bench"
)

type Service struct {
	key    string
	Stream grpc.ServerStream
	libbench.EtcdValueJson
}

func (p *Service) OnBidirectionalRecv() {

}
