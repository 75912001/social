package mongodb

import (
	"context"
	xrlog "dawn-server/impl/xr/lib/log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateBulkWriteModel FunctionArg转换为BulkWrite的model
func CreateBulkWriteModel(funcArg *FunctionArg) mongo.WriteModel {
	switch funcArg.funcID {
	case FuncIDUpdateOne:
		uom := mongo.NewUpdateOneModel()
		uom.SetFilter(funcArg.arg[3])
		uom.SetUpdate(funcArg.arg[4])
		if len(funcArg.arg) > 5 {
			uom.SetArrayFilters(options.ArrayFilters{
				Filters: []interface{}{
					funcArg.arg[5].(bson.M),
				},
			})
		}

		//for _, v := range funcArg.arg[5:] {
		//	v1 := v.(options.UpdateOptions)
		//	if v1.ArrayFilters != nil {
		//		af := options.ArrayFilters{}
		//		af.Filters = v1.ArrayFilters.Filters
		//		af.Registry = v1.ArrayFilters.Registry
		//		uom.SetArrayFilters(af)
		//	}
		//	if v1.Upsert != nil {
		//		uom.SetUpsert(*v1.Upsert)
		//	}
		//	if v1.Hint != nil {
		//		uom.SetHint(v1.Hint)
		//	}
		//	if v1.Collation != nil {
		//		uom.SetCollation(v1.Collation)
		//	}
		//}
		return uom
	case FuncIDUpdateMany:
		umm := mongo.NewUpdateManyModel()
		umm.SetFilter(funcArg.arg[3])
		umm.SetUpdate(funcArg.arg[4])
		if len(funcArg.arg) > 5 {
			umm.SetArrayFilters(options.ArrayFilters{
				Filters: []interface{}{
					funcArg.arg[5].(bson.M),
				},
			})
		}
		return umm
	//case FuncIDInsertOne:
	//	iom := mongo.NewInsertOneModel()
	//	iom.SetDocument(funcArg.arg[3])
	//	return iom
	//case FuncIDDeleteOne:
	//	dom := mongo.NewDeleteOneModel()
	//	dom.SetFilter(funcArg.arg[3])
	//	return dom
	//case FuncIDDeleteMany:
	//	dmm := mongo.NewDeleteManyModel()
	//	dmm.SetFilter(funcArg.arg[3])
	//	return dmm
	default:
		xrlog.PrintErr(ErrorKeyOperateFailure)
	}
	return nil
}

// CreateFunctionArg BulkWrite的model转换为FunctionArg
func CreateFunctionArg(collectionName string, collection *mongo.Collection, model mongo.WriteModel, expiration time.Duration) (functionArg *FunctionArg) {
	ctx := context.Background()

	switch model.(type) {
	case *mongo.UpdateOneModel:
		m1 := model.(*mongo.UpdateOneModel)
		filter := m1.Filter
		update := m1.Update
		if m1.ArrayFilters != nil && len(m1.ArrayFilters.Filters) > 0 {
			functionArg = NewFunctionArg(UpdateOne, FuncIDUpdateOne, collectionName,
				ctx, collection, expiration, filter, update, m1.ArrayFilters.Filters[0])
		} else {
			functionArg = NewFunctionArg(UpdateOne, FuncIDUpdateOne, collectionName,
				ctx, collection, expiration, filter, update)
		}
	case *mongo.UpdateManyModel:
		m1 := model.(*mongo.UpdateManyModel)
		filter := m1.Filter
		update := m1.Update
		if m1.ArrayFilters != nil && len(m1.ArrayFilters.Filters) > 0 {
			functionArg = NewFunctionArg(UpdateMany, FuncIDUpdateMany, collectionName,
				ctx, collection, expiration, filter, update, m1.ArrayFilters.Filters[0])
		} else {
			functionArg = NewFunctionArg(UpdateMany, FuncIDUpdateMany, collectionName,
				ctx, collection, expiration, filter, update)
		}
	//case *mongo.InsertOneModel:
	//	m1 := model.(*mongo.InsertOneModel)
	//	doc := m1.Document
	//	functionArg = NewFunctionArg(InsertOne, FuncIDInsertOne, collectionName,
	//		ctx, collection, expiration, doc)
	//case *mongo.DeleteOneModel, *mongo.DeleteManyModel:
	//	return nil
	default:
		xrlog.PrintErr(ErrorKeyOperateFailure)
		return nil
	}

	return functionArg
}
