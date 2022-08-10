package pgstore

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/entities"
	"github.com/google/uuid"
)

type Bookings struct {
	db *sql.DB
}

type DBPgBooking struct {
	BookingID        uuid.UUID
	CreatedAt        time.Time
	DeletedAt        *time.Time
	ClientName       string
	ClientPhone      string
	DataBooking      time.Time
	BookingTimeFirst time.Time
	BookingTimeTo    time.Time
}

func NewBookings(dsn string) (*Bookings, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS bookings (
		id uuid NOT NULL,
		created_at timestamptz NOT NULL,
		deleted_at timestamptz,
		client_name varchar(200) NOT NULL,
		client_phone varchar(20) NOT NULL,
		booked_date date NOT NULL,
		booked_time_from timestamptz NOT NULL,
		booked_time_to timestamptz NOT NULL,
		CONSTRAINT bookings_pk PRIMARY KEY (id)
	)`)
	if err != nil {
		db.Close()
		return nil, err
	}
	rs := &Bookings{
		db: db,
	}
	return rs, nil
}

func (b *Bookings) Close() {
	b.db.Close()
}


func (b *Bookings) Create(ctx context.Context, bking entities.Booking, tableIDs ...uuid.UUID) (*uuid.UUID, error) {
	fail := func(err error) (*uuid.UUID, error) {
		return nil, fmt.Errorf("create booking: %w", err)
	}

	begtx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}
	defer begtx.Rollback()

	dbbooking := &DBPgBooking{
		BookingID:        bking.BookingID,
		CreatedAt:        time.Now(),
		ClientName:       bking.ClientName,
		ClientPhone:      bking.ClientPhone,
		DataBooking:      bking.DataBooking,
		BookingTimeFirst: bking.BookingTimeFirst,
		BookingTimeTo:    bking.BookingTimeFirst.Add(2 * time.Hour),
	}
	_, err = b.db.ExecContext(ctx, `INSERT INTO bookings 
	(id, created_at, deleted_at, client_name, client_phone, booked_date, booked_time_from, booked_time_to)
	values ($1, $2, $3, $4, $5, $6, $7, $8)`,
		dbbooking.BookingID,
		dbbooking.CreatedAt,
		dbbooking.ClientName,
		dbbooking.ClientPhone,
		dbbooking.DataBooking,
		dbbooking.BookingTimeFirst,
		dbbooking.BookingTimeTo,
	)
	if err != nil {
		return nil, err
	}
	return &bking.BookingID, nil
}
