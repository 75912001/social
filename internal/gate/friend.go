package gate

import (
	"fmt"
	pkgserver "social/pkg/server"
	pkgservice "social/pkg/service"
)

type Friend struct {
	*pkgservice.Service
}

func (p *Friend) String() string {
	return fmt.Sprintf("%v", pkgserver.NameFriend)
}

func (p *Friend) OnRecv() {

}
