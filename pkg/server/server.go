package server

type IServer interface {
	Start() (err error)
	Stop() (err error)
}
