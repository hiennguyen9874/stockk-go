package bars

import (
	"github.com/hiennguyen9874/stockk-go/internal"
	"github.com/hiennguyen9874/stockk-go/internal/models"
)

type BarRedisRepository interface {
	internal.RedisRepository[models.Bar]
}
