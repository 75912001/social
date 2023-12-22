package mongodb

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	liberror "social/lib/error"
	libruntime "social/lib/runtime"
	"time"
)

// UpdateOne
//
//	参数:
//		0:context.Context
//		1:*mongo.Collection
//		2:time.Duration
//		3:filter interface{}
//		4:update interface{}
//		[5:]:opts ...*options.UpdateOptions(目前只支持一个 opt-ArrayFilters)
//	返回值:
//		updateResult *mongo.UpdateResult
func UpdateOne(arg ...interface{}) (updateResult interface{}, err error) {
	if len(arg) < 5 || 6 < len(arg) {
		return nil, errors.WithMessagef(liberror.Param, "%v %v %v", ErrorKeyWriteFailure, arg, libruntime.Location())
	}
	ctx := arg[0].(context.Context)
	collection := arg[1].(*mongo.Collection)
	timeout := arg[2].(time.Duration)

	filter := arg[3]
	update := arg[4]
	var opts []*options.UpdateOptions
	if 5 < len(arg) {
		opt := &options.UpdateOptions{
			ArrayFilters: &options.ArrayFilters{
				Filters: []interface{}{
					arg[5].(bson.M),
				},
			},
		}
		opts = append(opts, opt)
	}

	c, cancel := context.WithTimeout(ctx, timeout)
	defer func() {
		cancel()
	}()

	if updateResult, err = collection.UpdateOne(c, filter, update, opts...); err != nil {
		return updateResult, errors.WithMessagef(err, "%v %v %v %v %v",
			ErrorKeyWriteFailure, filter, update, opts, libruntime.Location())
	} else {
		return updateResult, nil
	}
}

// UpdateMany
//
//	参数:
//		0:context.Context
//		1:*mongo.Collection
//		2:time.Duration
//		3:filter interface{}
//		4:update interface{}
//		[5:]:opts ...*options.UpdateOptions(目前只支持一个opt-ArrayFilters)
//	返回值:
//		updateResult *mongo.UpdateResult
func UpdateMany(arg ...interface{}) (updateResult interface{}, err error) {
	if len(arg) < 5 || 6 < len(arg) {
		return nil, errors.WithMessagef(liberror.Param, "%v %v %v", ErrorKeyWriteFailure, arg, libruntime.Location())
	}

	ctx := arg[0].(context.Context)
	collection := arg[1].(*mongo.Collection)
	timeout := arg[2].(time.Duration)

	filter := arg[3]
	update := arg[4]
	var opts []*options.UpdateOptions
	if 5 < len(arg) {
		opt := &options.UpdateOptions{
			ArrayFilters: &options.ArrayFilters{
				Filters: []interface{}{
					arg[5].(bson.M),
				},
			},
		}
		opts = append(opts, opt)
	}

	c, cancel := context.WithTimeout(ctx, timeout)
	defer func() {
		cancel()
	}()

	if updateResult, err = collection.UpdateMany(c, filter, update, opts...); err != nil {
		return updateResult, errors.WithMessagef(err, "%v %v %v %v %v",
			ErrorKeyWriteFailure, filter, update, opts, libruntime.Location())
	} else {
		return updateResult, nil
	}
}
