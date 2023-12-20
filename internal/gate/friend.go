package gate

import (
	"fmt"
	"google.golang.org/grpc"
	libbench "social/lib/bench"
	pkgserver "social/pkg/server"
)

type Friend struct {
	key           string
	Stream        grpc.ServerStream
	EtcdValueJson libbench.EtcdValueJson
}

func (p *Friend) String() string {
	return fmt.Sprintf("%v", pkgserver.NameFriend)
}
