package mongodb

import (
	libfuncmgr "social/lib/func_mgr"
	liblog "social/lib/log"
	"time"
)

const FuncIDUpdateOne uint32 = 0x301
const FuncIDUpdateMany uint32 = 0x302

// FunctionArg 函数,参数
type FunctionArg struct {
	function       libfuncmgr.Function
	funcID         uint32
	collectionName string

	arg []interface{} //0:context.Context //1:*mongo.Collection //2:time.Duration
}

func (p *FunctionArg) GetFuncID() uint32 {
	return p.funcID
}

func (p *FunctionArg) GetCollectionName() string {
	return p.collectionName
}

func (p *FunctionArg) GetTimeOut() time.Duration {
	return p.arg[2].(time.Duration)
}

func (p *FunctionArg) GetArg() (arg []interface{}) {
	return p.arg
}

func (p *FunctionArg) AppendArg(i interface{}) {
	p.arg = append(p.arg, i)
}

// NewFunctionArg 构造新的FunctionArg
//
//	参数:
//		arg:
//			0:context.Context
//			1:*mongo.Collection
//			2:time.Duration
//			...
//			参考 InsertOne, UpdateOne, UpdateMany
func NewFunctionArg(fun libfuncmgr.Function, funcID uint32, collectionName string, arg ...interface{}) *FunctionArg {
	f := &FunctionArg{
		function:       fun,
		funcID:         funcID,
		collectionName: collectionName,
	}
	f.arg = append(f.arg, arg...)
	return f
}

// Invoke 调用函数
func (p *FunctionArg) Invoke() (interface{}, error) {
	i, err := p.function(p.arg...)
	if err != nil {
		liblog.PrintfErr("%v %v %v %v %v %v", ErrorKeyOperateFailure, err, p.funcID, p.function, p.collectionName, p.arg)
	}
	return i, err
}
