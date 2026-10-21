package executor

import "gorm.io/gorm"

// 事务和普通查询都是 *gorm.DB,统一抽象成一个类型别名
// 未来想换成接口只改这里
type Executor = *gorm.DB
