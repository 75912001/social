package ec

import (
	"github.com/pkg/errors"
	liberror "social/pkg/lib/error"
	libutil "social/pkg/lib/util"
	pkgproto "social/pkg/proto"
)

func Init() error {
	for k, v := range pkgproto.ERROR_CODE_name {
		if err := liberror.Register(
			&liberror.Error{
				Code: uint32(k),
				Name: v,
				Desc: v,
			},
		); err != nil {
			if uint32(k) == liberror.Success.Code {
				continue
			}
			return errors.WithMessage(err, libutil.GetCodeLocation(1).String())
		}
	}
	return nil
}
