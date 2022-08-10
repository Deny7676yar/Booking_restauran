package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/entities"
	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/repo/tablerepo"
	"github.com/google/uuid"
)

type HandlerTables struct {
	ts *tablerepo.TableRests
}

func NewHandlerTables(ts *tablerepo.TableRests) *HandlerTables {
	r := &HandlerTables{
		ts: ts,
	}
	return r
}

type TableRest struct {
	ID           uuid.UUID `json:"id"`
	RestaurantID uuid.UUID `json:"restaurant_id`
	CapacityT    int       `json.capacityT"`
}

func (h *HandlerTables) CreateTable(ctx context.Context, table TableRest, rest_id uuid.UUID) (TableRest, error) {
	t := entities.TableRest{
		ID:           table.ID,
		RestaurantID: rest_id,
		CapacityT:    table.CapacityT,
	}

	tb, err := h.ts.Create(ctx, t, rest_id)
	if err != nil {
		return TableRest{}, fmt.Errorf("error when creating: %w", err)
	}

	return TableRest{
		ID:           tb.ID,
		RestaurantID: tb.RestaurantID,
		CapacityT:    tb.CapacityT,
	}, nil
}

func (h *HandlerTables) UpdateTable(ctx context.Context, table TableRest) (TableRest, error) {
	tab := entities.TableRest{
		ID:           table.ID,
		RestaurantID: table.RestaurantID,
		CapacityT:    table.CapacityT,
	}

	rt, err := h.ts.UpdateTable(ctx, tab)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TableRest{}, ErrRestNotFound
		}
		return TableRest{}, fmt.Errorf("error when reading table: %w", err)
	}

	return TableRest{
		ID:           rt.ID,
		RestaurantID: rt.RestaurantID,
		CapacityT:    rt.CapacityT,
	}, nil
}

var ErrTableNotFound = errors.New("restaurant not found")

func (h *HandlerTables) ReadTable(ctx context.Context, uid uuid.UUID) (TableRest, error) {
	if (uid == uuid.UUID{}) {
		return TableRest{}, fmt.Errorf("bad request table: uid is empty")
	}

	rtab, err := h.ts.Read(ctx, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TableRest{}, ErrTableNotFound
		}
		return TableRest{}, fmt.Errorf("error when reading table: %w", err)
	}
	return TableRest{
		ID:           rtab.ID,
		RestaurantID: rtab.RestaurantID,
		CapacityT:    rtab.CapacityT,
	}, nil
}

func (h *HandlerTables) DeleteTable(ctx context.Context, uid uuid.UUID) (TableRest, error) {
	if (uid == uuid.UUID{}) {
		return TableRest{}, fmt.Errorf("bad request: uid is empty")
	}
	rtab, err := h.ts.Delete(ctx, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TableRest{}, ErrTableNotFound
		}
		return TableRest{}, fmt.Errorf("error when delete table: %w", err)
	}

	return TableRest{
		ID:           rtab.ID,
		RestaurantID: rtab.RestaurantID,
		CapacityT:    rtab.CapacityT,
	}, nil
}

func (h *HandlerTables) GetAavailableTable(ctx context.Context, rid uuid.UUID, desiredData, desiredTime string, f func(TableRest) error) error {
	chrt, err := h.ts.GetAavailableTable(ctx, rid, desiredData, desiredTime)
	if err != nil {
		return fmt.Errorf("error when reading: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case table, ok := <-chrt:
			if !ok {
				return nil
			}
			if err := f(TableRest{
				ID:           table.ID,
				RestaurantID: table.RestaurantID,
				CapacityT:    table.CapacityT,
			}); err != nil {
				return err
			}
		}
	}
}
