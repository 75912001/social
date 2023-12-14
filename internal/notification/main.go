package notification

import (
	"context"
	pkgserver "social/pkg/server"
)

type Server struct {
	*pkgserver.Normal
}

func (p *Server) OnStop(ctx context.Context) (err error) {
	return nil
}

func (p *Server) OnStart(ctx context.Context) (err error) {
	return nil
}
