package service

import (
	"context"
	"errors"

	"wb_labs_l0_backend/internal/domain/order"
	"wb_labs_l0_backend/internal/logger"
	"wb_labs_l0_backend/internal/repository/postgres"
	cachepkg "wb_labs_l0_backend/internal/service/cache"
)

type OrderService struct {
	repo  *postgres.PostgresOrderRepo
	cache cachepkg.Cache
	log   logger.Logger
}

func NewOrderService(repo *postgres.PostgresOrderRepo, cache cachepkg.Cache, log logger.Logger) *OrderService {
	return &OrderService{repo: repo, cache: cache, log: log}
}

func (s *OrderService) SaveOrder(ctx context.Context, o *order.Order) error {
	if err := o.Validate(); err != nil {
		return err
	}
	if err := s.repo.Save(ctx, o); err != nil {
		s.log.Errorf("save to db failed: %v", err)
		return err
	}
	s.cache.Set(o.OrderUID, o)
	return nil
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (*order.Order, error) {
	if id == "" {
		return nil, errors.New("empty id")
	}
	if o, ok := s.cache.Get(id); ok {
		return o, nil
	}
	// fallback to DB if cache miss
	o, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// populate cache
	s.cache.Set(id, o)
	return o, nil
}

// RestoreCache loads all orders from DB into cache at startup.
func (s *OrderService) RestoreCache(ctx context.Context) error {
	list, err := s.repo.FindAll(ctx)
	if err != nil {
		return err
	}
	for _, o := range list {
		if o == nil || o.OrderUID == "" {
			continue
		}
		s.cache.Set(o.OrderUID, o)
	}
	s.log.Infof("cache restored: %d items", s.cache.Len())
	return nil
}
