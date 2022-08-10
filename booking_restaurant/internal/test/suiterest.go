package test

import (
	"context"
	"os"
	"os/signal"

	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/entities"
	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/repo/restaurantrepo"
	"github.com/google/uuid"
	gc "gopkg.in/check.v1"
)

type SuitBase struct {
	r restaurantrepo.RestaurantStore
}

func (s *SuitBase) SetRest(rest restaurantrepo.RestaurantStore) {
	s.r = rest
}

func (s *SuitBase) TestCreateRest(c *gc.C) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	origin := &entities.Restaurant{
		Title: "Каравелла",
		AverageWaitingTime: 30,
		AverageCheck: 2000,
		AvailableSeats: 3,
	}
	_, err := s.r.Create(ctx, *origin)
	c.Assert(err, gc.IsNil)
	c.Assert(origin.ID, gc.Not(gc.Equals), uuid.Nil, gc.Commentf(""))

	existing := &entities.Restaurant{
		ID:     origin.ID,
		AverageWaitingTime: 30,
		AverageCheck: 2000,
		AvailableSeats: 3,
	}
	_, err = s.r.Create(ctx, *existing)
	c.Assert(err, gc.IsNil)
	c.Assert(existing.ID, gc.Equals, origin.ID, gc.Commentf("link ID changed while upserting"))
	cancel()
}

func (s *SuitBase) TestSearchRest(c *gc.C) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	chtest := make(chan entities.Restaurant, 100)

	rest1 := entities.Restaurant{
		Title: "Каравелла",
		AverageWaitingTime: 30,
		AverageCheck: 2000,
		AvailableSeats: 3,
	}
    desiredData := "09.08.2022"
	desiredTime := "15:00" 
	numbp := "3"
	chtest <- rest1

	_, err := s.r.Create(ctx, rest1)
	c.Assert(err, gc.IsNil)
	c.Assert(rest1.ID, gc.Not(gc.Equals), uuid.Nil, gc.Commentf(""))

	other, err := s.r.SearchRelevantRest(ctx, desiredData, desiredTime, numbp)
	c.Assert(err, gc.IsNil)
	c.Assert(other, gc.DeepEquals, chtest, gc.Commentf("lookup by ID returned the wrong link"))
	cancel()
}
