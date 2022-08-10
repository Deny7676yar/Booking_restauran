package server

import (
	"context"
	"net/http"
	"time"

	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/repo/bookingrepo"
	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/repo/restaurantrepo"
	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/repo/tablerepo"
)

type Server struct {
	srv   http.Server
	rest  *restaurantrepo.Restaurants
	table *tablerepo.TableRests
	booking *bookingrepo.Bookings
}

func NewServer(addr string, h http.Handler) *Server {
	s := &Server{}

	s.srv = http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
	}
	return s
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	s.srv.Shutdown(ctx)
	cancel()
}

func (s *Server) Start(rest *restaurantrepo.Restaurants, table *tablerepo.TableRests, booking *bookingrepo.Bookings) {
	s.rest = rest
	s.table = table
	s.booking = booking
	// TODO: migrations
	go s.srv.ListenAndServe()
}
