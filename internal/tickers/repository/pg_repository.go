package repository

import (
	"context"
	"fmt"

	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/repository"
	"github.com/hiennguyen9874/stockk-go/internal/tickers"
	"gorm.io/gorm"
)

type TickerPgRepo struct {
	repository.PgRepo[models.Ticker]
}

func CreateTickerPgRepository(db *gorm.DB) tickers.TickerPgRepository {
	return &TickerPgRepo{
		PgRepo: repository.CreatePgRepo[models.Ticker](db),
	}
}

func (r *TickerPgRepo) GetBySymbol(ctx context.Context, symbol string) (*models.Ticker, error) {
	var obj *models.Ticker
	if result := r.DB.WithContext(ctx).First(&obj, "symbol = ?", symbol); result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *TickerPgRepo) UpdateIsActive(ctx context.Context, exp *models.Ticker, isActive bool) (*models.Ticker, error) {
	if result := r.DB.WithContext(ctx).Model(&exp).Select("is_active").
		Updates(map[string]interface{}{"is_active": isActive}); result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}

func (r *TickerPgRepo) GetAllActive(ctx context.Context, isActive bool) ([]*models.Ticker, error) {
	var objs []*models.Ticker
	r.DB.WithContext(ctx).Where("is_active = ?", isActive).Find(&objs)
	return objs, nil
}

func (r *TickerPgRepo) SearchBySymbol(ctx context.Context, symbol string, limit int, exchange string) ([]*models.Ticker, error) {
	var objs []*models.Ticker

	query := r.DB.WithContext(ctx)

	if symbol != "" {
		query = query.Where("symbol LIKE ?", fmt.Sprintf("%v%v", symbol, "%"))
	}

	if exchange != "" {
		query = query.Where("exchange = ?", exchange)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	query.Order("symbol").Find(&objs)

	return objs, nil
}
