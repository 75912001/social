package common

import (
	"fmt"
)

// GenDBName 生成数据库名称
func GenDBName(ZoneID uint32, dbName string) string {
	return fmt.Sprintf("db_%v_%v", dbName, ZoneID)
}
