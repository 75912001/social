package mongodb

import (
	"context"
	xrerror "dawn-server/impl/xr/lib/error"
	xrutil "dawn-server/impl/xr/lib/util"
	"time"

	"github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertOne 插入文档
//	参数:
//		0:context.Context
//		1:*mongo.Collection
//		2:time.Duration
//		3:document interface{}
//		[!]不支持[4:]:opts ...*options.InsertOneOptions
func InsertOne(arg ...interface{}) (interface{}, error) {
	if 4 != len(arg) {
		return nil, errors.WithMessagef(xrerror.Param, "%v %v %v", ErrorKeyWriteFailure, arg, xrutil.GetCodeLocation(1).String())
	}
	ctx := arg[0].(context.Context)
	collection := arg[1].(*mongo.Collection)
	timeout := arg[2].(time.Duration)
	document := arg[3]

	var opts []*options.InsertOneOptions

	c, cancel := context.WithTimeout(ctx, timeout)
	defer func() {
		cancel()
	}()

	if _, err := collection.InsertOne(c, document, opts...); err != nil {
		return nil, errors.WithMessagef(err, "%v %v %v %v",
			ErrorKeyWriteFailure, document, opts, xrutil.GetCodeLocation(1).String())
	}
	return nil, nil
}
