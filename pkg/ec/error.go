package ec

import (
	"github.com/pkg/errors"
	liberror "social/lib/error"
	libutil "social/lib/util"
	pkgproto "social/pkg/proto"
)

func Init() error {
	for k, v := range pkgproto.ERROR_CODE_name {
		if err := liberror.CheckForDuplicates(liberror.NewError(uint32(k), v, v)); err != nil {
			if uint32(k) == liberror.Success.Code {
				continue
			}
			return errors.WithMessage(err, libutil.GetCodeLocation(1).String())
		}
	}
	return nil
}
