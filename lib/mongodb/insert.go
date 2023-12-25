package mongodb

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	liberror "social/lib/error"
	libruntime "social/lib/runtime"
	"time"
)

// InsertOne 插入文档
//
//	参数:
//		0:context.Context
//		1:*mongo.Collection
//		2:time.Duration
//		3:document interface{}
//		[!]不支持[4:]:opts ...*options.InsertOneOptions
func InsertOne(arg ...interface{}) (interface{}, error) {
	if 4 != len(arg) {
		return nil, errors.WithMessagef(liberror.Param, "%v %v %v", ErrorKeyWriteFailure, arg, libruntime.Location())
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
			ErrorKeyWriteFailure, document, opts, libruntime.Location())
	}
	return nil, nil
}
