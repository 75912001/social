package error_code

import (
	"github.com/pkg/errors"
	xrerror "social/pkg/lib/error"
	xrutil "social/pkg/lib/util"
	"social/pkg/proto"
)

func Init() error {
	for k, v := range proto.ERROR_CODE_name {
		if err := xrerror.Register(
			&xrerror.Error{
				Code: uint32(k),
				Name: v,
				Desc: v,
			},
		); err != nil {
			if uint32(k) == xrerror.Success.Code {
				continue
			}
			return errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
		}
	}
	return nil
}
