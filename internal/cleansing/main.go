package cleansing

import "context"

type Server struct {
}

func (p *Server) Stop(ctx context.Context) (err error) {
	return nil
}

func (p *Server) Start(ctx context.Context) (err error) {
	return nil
}
