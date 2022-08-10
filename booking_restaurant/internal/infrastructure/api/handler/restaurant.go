package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/entities"
	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/repo/restaurantrepo"
	"github.com/google/uuid"
)

type HandlerRests struct {
	rests *restaurantrepo.Restaurants
}

func NewHandlerRests(rests *restaurantrepo.Restaurants) *HandlerRests {
	h := &HandlerRests{rests: rests}
	return h
}

type Restaurant struct {
	ID                 uuid.UUID `json:"id"`
	Title              string    `json:"title`
	AverageWaitingTime int       `json:"averageWaitingTime"`
	AverageCheck       float64   `json.averageCheck`
	AvailableSeats     int       `json.availableSeats"`
}

func (h *HandlerRests) CreateRest(ctx context.Context, rest Restaurant) (Restaurant, error) {
	r := entities.Restaurant{
		ID:                 rest.ID,
		Title:              rest.Title,
		AverageWaitingTime: rest.AverageWaitingTime,
		AverageCheck:       rest.AverageCheck,
		AvailableSeats:     rest.AvailableSeats,
	}

	rb, err := h.rests.Create(ctx, r)
	if err != nil {
		return Restaurant{}, fmt.Errorf("error when creating: %w", err)
	}

	return Restaurant{
		ID:                 rb.ID,
		Title:              rb.Title,
		AverageWaitingTime: rb.AverageWaitingTime,
		AverageCheck:       rb.AverageCheck,
		AvailableSeats:     rb.AvailableSeats,
	}, nil
}

func (h *HandlerRests) UpdateRest(ctx context.Context, rest Restaurant) (Restaurant, error) {
	restu := entities.Restaurant{
		ID:                 rest.ID,
		Title:              rest.Title,
		AverageWaitingTime: rest.AverageWaitingTime,
		AverageCheck:       rest.AverageCheck,
		AvailableSeats:     rest.AvailableSeats,
	}

	rb, err := h.rests.Update(ctx, restu)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Restaurant{}, ErrRestNotFound
		}
		return Restaurant{}, fmt.Errorf("error when reading restaurant: %w", err)
	}

	return Restaurant{
		ID:                 rb.ID,
		Title:              rb.Title,
		AverageWaitingTime: rb.AverageWaitingTime,
		AverageCheck:       rb.AverageCheck,
		AvailableSeats:     rb.AvailableSeats,
	}, nil
}

var ErrRestNotFound = errors.New("restaurant not found")

func (h *HandlerRests) ReadRest(ctx context.Context, uid uuid.UUID) (Restaurant, error) {
	if (uid == uuid.UUID{}) {
		return Restaurant{}, fmt.Errorf("bad request: uid is empty")
	}

	rb, err := h.rests.Read(ctx, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Restaurant{}, ErrRestNotFound
		}
		return Restaurant{}, fmt.Errorf("error when reading restaurant: %w", err)
	}
	return Restaurant{
		ID:                 rb.ID,
		Title:              rb.Title,
		AverageWaitingTime: rb.AverageWaitingTime,
		AverageCheck:       rb.AverageCheck,
		AvailableSeats:     rb.AvailableSeats,
	}, nil
}

func (h *HandlerRests) DeleteRest(ctx context.Context, uid uuid.UUID) (Restaurant, error) {
	if (uid == uuid.UUID{}) {
		return Restaurant{}, fmt.Errorf("bad request: uid is empty")
	}
	rb, err := h.rests.Delete(ctx, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Restaurant{}, ErrRestNotFound
		}
		return Restaurant{}, fmt.Errorf("error when delete restaurant: %w", err)
	}

	return Restaurant{
		ID:                 rb.ID,
		Title:              rb.Title,
		AverageWaitingTime: rb.AverageWaitingTime,
		AverageCheck:       rb.AverageCheck,
		AvailableSeats:     rb.AvailableSeats,
	}, nil
}

func (h *HandlerRests) SearchRelevantRest(ctx context.Context, desiredData, desiredTime, numbp string, f func(Restaurant) error) error {
	chr, err := h.rests.SearchRelevantRest(ctx, desiredData, desiredTime, numbp)
	if err != nil {
		return fmt.Errorf("error when reading: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case rest, ok := <-chr:
			if !ok {
				return nil
			}
			if err := f(Restaurant{
				ID:                 rest.ID,
				Title:              rest.Title,
				AverageWaitingTime: rest.AverageWaitingTime,
				AverageCheck:       rest.AverageCheck,
				AvailableSeats:     rest.AvailableSeats,
			}); err != nil {
				return err
			}
		}
	}
}
