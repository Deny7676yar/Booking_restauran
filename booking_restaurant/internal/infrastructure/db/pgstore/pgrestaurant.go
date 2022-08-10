package pgstore

import (
	"context"
	"database/sql"
	"time"

	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/entities"
	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/repo/restaurantrepo"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib" // Postgresql driver
)

var _ restaurantrepo.RestaurantStore = &Restaurants{}

type Restaurants struct {
	db *sql.DB
}

type DBPgRest struct {
	ID                 uuid.UUID
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time
	Title              string
	AverageWaitingTime int
	AverageCheck       float64
	AvailableSeats     int
}

func NewRestaurants(dsn string) (*Restaurants, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS restaurants (
		id uuid NOT NULL,
		created_at timestamptz NOT NULL,
		updated_at timestamptz NOT NULL,
		deleted_at timestamptz NULL,
		title varchar NOT NULL,
		average_waiting_time int2 NULL,
		average_check bigint NOT NULL,
		available_seats int2 NULL,
		CONSTRAINT restaurants_pk PRIMARY KEY (id)
	)`)
	if err != nil {
		db.Close()
		return nil, err
	}
	rs := &Restaurants{
		db: db,
	}
	return rs, nil
}

func (r *Restaurants) Close() {
	r.db.Close()
}

func (r *Restaurants) Create(ctx context.Context, rest entities.Restaurant) (*uuid.UUID, error) {
	dbu := &DBPgRest{
		ID:                 rest.ID,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		Title:              rest.Title,
		AverageWaitingTime: rest.AverageWaitingTime,
		AverageCheck:       rest.AverageCheck,
		AvailableSeats:     rest.AvailableSeats,
	}

	_, err := r.db.ExecContext(ctx, `INSERT INTO restaurants 
	(id, created_at, updated_at, deleted_at, title, average_waiting_time, average_check, available_seats)
	values ($1, $2, $3, $4, $5, $6, $7, $8)`,
		dbu.ID,
		dbu.CreatedAt,
		dbu.UpdatedAt,
		nil,
		dbu.Title,
		dbu.AverageWaitingTime,
		dbu.AverageCheck,
		dbu.AvailableSeats,
	)
	if err != nil {
		return nil, err
	}

	return &rest.ID, nil
}

func (r *Restaurants) Read(ctx context.Context, uid uuid.UUID) (*entities.Restaurant, error) {
	dbu := &DBPgRest{}
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, created_at, updated_at, deleted_at, title, average_waiting_time, average_check, available_seats 
	FROM restaurants WHERE id = $1`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&dbu.ID,
			&dbu.CreatedAt,
			&dbu.UpdatedAt,
			&dbu.DeletedAt,
			&dbu.Title,
			&dbu.AverageWaitingTime,
			&dbu.AverageCheck,
			&dbu.AvailableSeats,
		); err != nil {
			return nil, err
		}
	}

	return &entities.Restaurant{
		ID:                 dbu.ID,
		Title:              dbu.Title,
		AverageWaitingTime: dbu.AverageWaitingTime,
		AverageCheck:       dbu.AverageCheck,
		AvailableSeats:     dbu.AvailableSeats,
	}, nil
}

func (r *Restaurants) Update(ctx context.Context, rest entities.Restaurant) error {
	dbu := &DBPgRest{
		ID:                 rest.ID,
		UpdatedAt:          time.Now(),
		Title:              rest.Title,
		AverageWaitingTime: rest.AverageWaitingTime,
		AverageCheck:       rest.AverageCheck,
		AvailableSeats:     rest.AvailableSeats,
	}
	_, err := r.db.ExecContext(ctx, `UPDATE restaurants SET updated_at = $2, title = $3, average_waiting_time = $4, average_check = $5, available_seats = $6,  WHERE id = $1`,
		dbu.ID,
		dbu.UpdatedAt,
		dbu.Title,
		dbu.AverageWaitingTime,
		dbu.AverageCheck,
		dbu.AvailableSeats,
	)
	return err
}

// Удаление происходит за счет отметки времени в поле deleted_at
// для оптимизации, это экономит ресурсы, БД не нужно перестраиватся
func (r *Restaurants) Delete(ctx context.Context, uid uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE restaurants SET deleted_at = $2, WHERE id = $1`,
		uid, time.Now(),
	)
	return err
}

func (r *Restaurants) SearchRelevantRest(ctx context.Context, desiredData, desiredTime, numbp string) (chan entities.Restaurant, error) {
	chout := make(chan entities.Restaurant, 100)
	go func() {
		defer close(chout)
		dbu := &DBPgRest{}

		rows, err := r.db.QueryContext(
			ctx, `SELECT restaurant.title, restaurant.average_waiting_time, restaurant.average_check, SUM(table.capacity) as available_capacity
			        FROM get_available_tables(data = $1, time = $2) table
					JOIN table ON restaurant.id = table.restaurant_id
					GROUP BY restaurant.id, restaurant.title, restaurant.average_waiting_time, restaurant.average_check
					HAVING SUM(table.capacity) > $1
					ORDER BY restaurant.average_waiting_time, restaurant.average_check and deleted_at is null`,
			desiredData, desiredTime)
		if err != nil {
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(
				&dbu.ID,
				&dbu.Title,
				&dbu.AverageWaitingTime,
				&dbu.AverageCheck,
				&dbu.AvailableSeats,
			); err != nil {
				return
			}
			chout <- entities.Restaurant{
				ID:                 dbu.ID,
				Title:              dbu.Title,
				AverageWaitingTime: dbu.AverageWaitingTime,
				AverageCheck:       dbu.AverageCheck,
				AvailableSeats:     dbu.AvailableSeats,
			}
		}
	}()

	return chout, nil
}
