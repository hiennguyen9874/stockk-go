package stockssnapshot

import (
	"github.com/hiennguyen9874/stockk-go/internal"
	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type StockSnapshotRedisRepository interface {
	internal.RedisRepository[models.StockSnapshot]
}
