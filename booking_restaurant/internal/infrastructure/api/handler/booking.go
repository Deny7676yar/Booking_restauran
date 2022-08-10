package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/entities"
	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/repo/bookingrepo"
)

type HandlerBookings struct {
	booking *bookingrepo.Bookings
}

func NewHandlerBookings(booking *bookingrepo.Bookings) *HandlerBookings {
	h := &HandlerBookings{booking: booking}
	return h
}

type Booking struct {
	PeopleNumber int `json:"people_number"`
	// DesiredDatetime представляет дату и время посещения ресторана в рамках брони
	DesiredDatetime time.Time `json:"desired_datetime"`
	// ClientName имя клиента, оформляющего бронь.
	ClientName string `json:"client_name"`
	// ClientPhone телефон клиента, оформляющего бронь.
	ClientPhone string `json:"client_phone"`
}

func (h *HandlerBookings) Create(ctx context.Context, bking Booking) (Booking, error) {
	r := entities.Booking{
		NumberOfPeople: bking.PeopleNumber,
		DataBooking:    bking.DesiredDatetime,
		ClientName:     bking.ClientName,
		ClientPhone:    bking.ClientPhone,
	}
	rb, err := h.booking.Create(ctx, r)
	if err != nil {
		return Booking{}, fmt.Errorf("error when creating: %w", err)
	}

	return Booking{
		PeopleNumber:    rb.NumberOfPeople,
		DesiredDatetime: rb.DataBooking,
		ClientName:      rb.ClientName,
		ClientPhone:     rb.ClientPhone,
	}, nil
}
