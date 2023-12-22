package mongodb

import (
	"context"
	xrerror "dawn-server/impl/xr/lib/error"
	xrutil "dawn-server/impl/xr/lib/util"
	"time"

	"github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindOne 查找一个文档/查找一个空文档
//	参数:
//		0:context.Context
//		1:*mongo.Collection
//		2:time.Duration
//		3:filter interface{}
//		4:record interface{}
//		[5:]:opts ...[]*options.FindOneOptions
//	返回值:
//		exist bool 是否有 true:有, false:无
func FindOne(arg ...interface{}) (exist interface{}, err error) {
	if len(arg) < 5 {
		return nil, errors.WithMessagef(xrerror.Param, "%v %v %v", ErrorKeyReadFailure, arg, xrutil.GetCodeLocation(1).String())
	}
	ctx := arg[0].(context.Context)
	collection := arg[1].(*mongo.Collection)
	timeout := arg[2].(time.Duration)
	filter := arg[3]
	record := arg[4]
	var opts []*options.FindOneOptions
	if 5 < len(arg) {
		for _, v := range arg[5:] {
			for _, v1 := range v.([]*options.FindOneOptions) {
				opts = append(opts, v1)
			}
		}
	}

	c, cancel := context.WithTimeout(ctx, timeout)
	defer func() {
		cancel()
	}()

	singleResult := collection.FindOne(c, filter, opts...)
	if err = singleResult.Err(); err != nil {
		if mongo.ErrNoDocuments == err {
			return false, nil
		}
		return false, errors.WithMessagef(err, "%v %v %v %v", ErrorKeyReadFailure, filter, opts, xrutil.GetCodeLocation(1).String())
	}
	if err = singleResult.Decode(record); err != nil {
		return false, errors.WithMessagef(err, "%v %v %v", ErrorKeyReadFailure, singleResult, xrutil.GetCodeLocation(1).String())
	}
	return true, nil
}

// CountDocuments 获取文档总数
//	参数:
//		0:context.Context
//		1:*mongo.Collection
//		2:time.Duration
//		3:filter interface{}
//	返回值:
//		count int64 数量
func CountDocuments(arg ...interface{}) (count interface{}, err error) {
	if 4 != len(arg) {
		return nil, errors.WithMessagef(xrerror.Param, "%v %v %v", ErrorKeyReadFailure, arg, xrutil.GetCodeLocation(1).String())
	}
	ctx := arg[0].(context.Context)
	collection := arg[1].(*mongo.Collection)
	timeout := arg[2].(time.Duration)
	filter := arg[3]

	c, cancel := context.WithTimeout(ctx, timeout)
	defer func() {
		cancel()
	}()

	if count, err = collection.CountDocuments(c, filter); err != nil {
		return count, errors.WithMessagef(err, "%v %v %v", ErrorKeyReadFailure, filter, xrutil.GetCodeLocation(1).String())
	} else {
		return count, nil
	}
}

// NextValue 返回更新后的值
//	参数:
//		0:context.Context
//		1:*mongo.Collection
//		2:time.Duration
//		3:filter interface{}
//		4:update interface{}
//	返回值:
//		singleResult *mongo.SingleResult
func NextValue(arg ...interface{}) (singleResult interface{}, err error) {
	if 5 != len(arg) {
		return nil, errors.WithMessagef(xrerror.Param, "%v %v %v", ErrorKeyWriteFailure, arg, xrutil.GetCodeLocation(1).String())
	}
	ctx := arg[0].(context.Context)
	collection := arg[1].(*mongo.Collection)
	timeout := arg[2].(time.Duration)
	filter := arg[3]
	update := arg[4]

	c, cancel := context.WithTimeout(ctx, timeout)
	defer func() {
		cancel()
	}()

	returnDocument := options.After
	singleResult = collection.FindOneAndUpdate(c, filter, update, &options.FindOneAndUpdateOptions{
		ReturnDocument: &returnDocument,
	})
	if err = singleResult.(*mongo.SingleResult).Err(); err != nil {
		return singleResult, errors.WithMessagef(err, "%v %v %v %v %v",
			ErrorKeyWriteFailure, filter, update, singleResult, xrutil.GetCodeLocation(1).String())
	}
	return singleResult, nil
}

// Find 查找所有文档
//	参数:
//		0:context.Context
//		1:*mongo.Collection
//		2:time.Duration
//		3:filter interface{}
//		[4:]:opts ...[]*options.FindOptions
//	返回值:
//		[]bson.M 多个文档数据
func Find(arg ...interface{}) ([]bson.M, error) {
	if len(arg) < 5 {
		return nil, errors.WithMessagef(xrerror.Param, "%v %v %v", ErrorKeyReadFailure, arg, xrutil.GetCodeLocation(1).String())
	}
	ctx := arg[0].(context.Context)
	collection := arg[1].(*mongo.Collection)
	timeout := arg[2].(time.Duration)
	filter := arg[3]
	var opts []*options.FindOptions
	if 4 < len(arg) {
		for _, v := range arg[4:] {
			opts = append(opts, v.(*options.FindOptions))
		}
	}

	c, cancel := context.WithTimeout(ctx, timeout)
	defer func() {
		cancel()
	}()

	cursor, err := collection.Find(c, filter, opts...)
	if err != nil {
		return nil, errors.WithMessagef(err, "%v %v %v %v",
			ErrorKeyReadFailure, filter, opts, xrutil.GetCodeLocation(1).String())
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, errors.WithMessagef(err, "%v %v %v", ErrorKeyReadFailure, cursor, xrutil.GetCodeLocation(1).String())
	}

	return results, nil
}
