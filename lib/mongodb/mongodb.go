package mongodb

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	liblog "social/lib/log"
	libruntime "social/lib/runtime"
	"sync/atomic"
	"time"
)

const (
	StatusNormal   uint32 = 0 // 正常状态
	StatusAbnormal uint32 = 1 // 异常状态
)

//var (
//	instance *Mgr
//	once     sync.Once
//)
//
//// GetInstance 获取
//func GetInstance() *Mgr {
//	once.Do(func() {
//		instance = new(Mgr)
//	})
//	return instance
//}
//
//// IsEnable 是否 开启
//func IsEnable() bool {
//	if instance == nil {
//		return false
//	}
//	return GetInstance().client != nil
//}

type Mgr struct {
	options *Options

	client   *mongo.Client   //Client is a handle representing a pool of connections to a MongoDB deployment
	database *mongo.Database //MongoDB database

	status uint32 //状态 0:正常 1:异常
}

// Connect 连接
func (p *Mgr) Connect(ctx context.Context, opts ...*Options) error {
	p.options = merge(opts...)
	if err := configure(p.options); err != nil {
		return errors.WithMessagef(err, "%v %v", ErrorKeyConnectFailure, libruntime.Location())
	}

	uri := p.options.genURI()

	c, cancel := context.WithTimeout(context.Background(), *p.options.timeoutDuration)
	defer func() {
		cancel()
	}()

	opt := options.Client().ApplyURI(uri).
		SetMaxPoolSize(*p.options.maxPoolSize).
		SetMinPoolSize(*p.options.minPoolSize).
		SetTimeout(*p.options.timeoutDuration).
		SetMaxConnIdleTime(*p.options.maxConnIdleTime).
		SetMaxConnecting(*p.options.maxConnecting)

	var err error
	if p.client, err = mongo.Connect(c, opt); err != nil {
		return errors.WithMessagef(err, "%v %v", ErrorKeyConnectFailure, libruntime.Location())
	}
	if err = p.client.Ping(ctx, nil); err != nil {
		return errors.WithMessagef(err, "%v %v", ErrorKeyConnectFailure, libruntime.Location())
	}
	return nil
}

// EnableSharding 启用分片
func (p *Mgr) EnableSharding(ctx context.Context, dbName string) error {
	c, cancel := context.WithTimeout(ctx, *p.options.timeoutDuration)
	defer func() {
		cancel()
	}()
	if _, err := p.client.Database("admin").RunCommand(
		c,
		bson.D{{"enableSharding", dbName}},
		//bsonx.Doc{{"enableSharding", bsonx.String(dbName)}},
	).DecodeBytes(); err != nil {
		return errors.WithMessagef(err, "%v %v", ErrorKeyOperateFailure, libruntime.Location())
	}
	return nil
}

// ShardCollectionHash 分片数据集-hash
//
//	参数:
//		numInitialChunks:chuck数量 e.g.:1024
func (p *Mgr) ShardCollectionHash(ctx context.Context, dbName string, collectionName string, key string, numInitialChunks int) error {
	c, cancel := context.WithTimeout(ctx, *p.options.timeoutDuration)
	defer func() {
		cancel()
	}()
	//https://docs.mongodb.com/manual/reference/command/shardCollection/?_ga=2.124298636.780010457.1647225755-1073769207.1627904584
	cmd := bson.D{
		{"shardCollection", dbName + "." + collectionName},
		{"key", bson.D{{key, "hashed"}}},
		{"unique", false},
		{"numInitialChunks", numInitialChunks},
	}
	if err := p.client.Database("admin").RunCommand(c, cmd).Err(); err != nil {
		return errors.WithMessagef(err, "%v %v", ErrorKeyOperateFailure, libruntime.Location())
	}
	return nil
}

func (p *Mgr) SwitchedDatabase(name string) {
	p.database = p.client.Database(name)
}

func (p *Mgr) DropDatabase() error {
	err := p.database.Drop(context.TODO())
	if err != nil {
		return errors.WithMessagef(err, "%v %v", ErrorKeyOperateFailure, libruntime.Location())
	}
	return nil
}

func (p *Mgr) SwitchedCollection(name string) *mongo.Collection {
	return p.database.Collection(name)
}

// IndexesCreateOne 创建索引
//
//	参数:
//		field:索引字段
//		order:[1:ascending order 按升序创建索引. -1:descending order 按降序来创建索引]
//		unique: 建立的索引是否唯一.指定为true创建唯一索引.默认值为false
func (p *Mgr) IndexesCreateOne(ctx context.Context, collection *mongo.Collection, field string, order int, unique bool) (newIndexName string, err error) {
	c, cancel := context.WithTimeout(ctx, *p.options.timeoutDuration)
	defer func() {
		cancel()
	}()
	newIndexName, err = collection.Indexes().CreateOne(c,
		mongo.IndexModel{
			Keys: bson.M{
				field: order,
			},
			Options: options.Index().SetUnique(unique),
		},
	)
	if err != nil {
		return "", errors.WithMessagef(err, "%v %v", ErrorKeyOperateFailure, libruntime.Location())
	}
	return newIndexName, nil
}

func (p *Mgr) GetClient() *mongo.Client {
	return p.client
}

// Disconnect 断开链接
func (p *Mgr) Disconnect(ctx context.Context) error {
	c, cancel := context.WithTimeout(ctx, *p.options.timeoutDuration)
	defer func() {
		cancel()
	}()
	if err := p.client.Disconnect(c); err != nil {
		return errors.WithMessagef(err, "%v %v", ErrorKeyDisconnectFailure, libruntime.Location())
	}
	p.client = nil
	return nil
}

// GetTimeOutDuration 超时时间
func (p *Mgr) GetTimeOutDuration() time.Duration {
	return *p.options.timeoutDuration
}

// SetAbnormal 设置异常 false:正常 true:异常
func (p *Mgr) SetAbnormal(abnormal bool) {
	if abnormal {
		atomic.StoreUint32(&p.status, StatusAbnormal)
	} else {
		atomic.StoreUint32(&p.status, StatusNormal)
	}
}

// IsAbnormal 是否异常 false:正常 true:异常
func (p *Mgr) IsAbnormal() bool {
	return atomic.LoadUint32(&p.status) == StatusAbnormal
}

//	IsErrClientDisconnected (只做参考) 错误 使用断开连接的客户端运行操作
//	Deprecated:弃用
//func IsErrClientDisconnected(err error) bool {
//	return err == mongo.ErrClientDisconnected ||
//		strings.Contains(err.Error(), mongo.ErrClientDisconnected.Error())
//}

// StartTransaction 开启MongoDB事务
func StartTransaction(sessionContext mongo.SessionContext) error {
	err := sessionContext.StartTransaction()
	if err != nil {
		return errors.WithMessagef(err, "%v %v", ErrorKeyOperateFailure, libruntime.Location())
	}
	return nil
}

// EndTransaction 结束MongoDB事务
func EndTransaction(sessionContext mongo.SessionContext, err error) {
	if err != nil {
		e := sessionContext.AbortTransaction(context.Background())
		if e != nil {
			liblog.PrintErr(ErrorKeyOperateFailure, e)
		}
	} else {
		e := sessionContext.CommitTransaction(context.Background())
		if e != nil {
			liblog.PrintErr(ErrorKeyOperateFailure, e)
		}
	}
}
