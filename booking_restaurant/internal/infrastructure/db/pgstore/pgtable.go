package pgstore

import (
	"context"
	"database/sql"
	"time"

	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/entities"
	"github.com/Deny7676yar/booking_restaurant/booking_restaurant/internal/repo/tablerepo"
	"github.com/google/uuid"
)

var _ tablerepo.TableRestStore = &TableRests{}

type TableRests struct {
	db *sql.DB
}

type DBPgTable struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
	RestaurantID uuid.UUID
	CapacityT    int
}

func NewTableRests(dsn string) (*TableRests, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tables (
		id uuid NOT NULL,
		created_at timestamptz NOT NULL,
		updated_at timestamptz NOT NULL,
		deleted_at timestamptz NULL,
		restaurant_id uuid NOT NULL,
		capacity int2 NULL,
		CONSTRAINT fk_tables_restaurants FOREIGN KEY (restaurant_id) REFERENCES restaurants (id) ON DELETE CASCADE
	)`)
	if err != nil {
		db.Close()
		return nil, err
	}
	ts := &TableRests{
		db: db,
	}
	return ts, nil
}

func (tbl *TableRests) Close() {
	tbl.db.Close()
}

func (tbl *TableRests) Create(ctx context.Context, table entities.TableRest, rest_id uuid.UUID) (*uuid.UUID, error) {
	dbt := &DBPgTable{
		ID:           table.ID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		RestaurantID: rest_id,
		CapacityT:    table.CapacityT,
	}

	_, err := tbl.db.ExecContext(ctx, `INSERT INTO tables 
	(id, created_at, updated_at, deleted_at, restaurant_id, capacity)
	values ($1, $2, $3, $4, $5, $6)`,
		dbt.ID,
		dbt.CreatedAt,
		dbt.UpdatedAt,
		nil,
		dbt.RestaurantID,
		dbt.CapacityT,
	)
	if err != nil {
		return nil, err
	}

	return &table.ID, nil
}

func (tbl *TableRests) Read(ctx context.Context, uid uuid.UUID) (*entities.TableRest, error) {
	dbt := &DBPgTable{}
	rows, err := tbl.db.QueryContext(ctx,
		`SELECT id, created_at, updated_at, deleted_at, restaurant_id, capacity
	FROM tables WHERE id = $1`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&dbt.ID,
			&dbt.CreatedAt,
			&dbt.UpdatedAt,
			&dbt.DeletedAt,
			&dbt.RestaurantID,
			&dbt.CapacityT,
		); err != nil {
			return nil, err
		}
	}

	return &entities.TableRest{
		ID:           dbt.ID,
		RestaurantID: dbt.ID,
		CapacityT:    dbt.CapacityT,
	}, nil
}

func (tbl *TableRests) UpdateTable(ctx context.Context, table entities.TableRest) error {
	dbt := &DBPgTable{
		ID:           table.ID,
		UpdatedAt:    time.Now(),
		RestaurantID: table.RestaurantID,
		CapacityT:    table.CapacityT,
	}
	_, err := tbl.db.ExecContext(ctx, `UPDATE tables SET updated_at = $2, restaurant_id = $3, capacity = $4 WHERE id = $1`,
		dbt.ID,
		dbt.UpdatedAt,
		dbt.RestaurantID,
		dbt.CapacityT,
	)
	return err
}

func (tbl *TableRests) Delete(ctx context.Context, uid uuid.UUID) error {
	_, err := tbl.db.ExecContext(ctx, `UPDATE tables SET deleted_at = $2, WHERE id = $1`,
		uid, time.Now(),
	)
	return err
}

func (tbl *TableRests) GetAavailableTable(ctx context.Context, rid uuid.UUID, desiredData, desiredTime string) (chan entities.TableRest, error) {
	choutbl := make(chan entities.TableRest, 100)
	go func() {
		defer close(choutbl)
		dbt := &DBPgTable{}
		rows, err := tbl.db.QueryContext(ctx, `SELECT * FROM tables(date = $2, time = $3) WHERE restaurant_id = $1`,
			rid, desiredData, desiredTime)
		if err != nil {
			return
		}
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(
				&dbt.ID,
				&dbt.CreatedAt,
				&dbt.UpdatedAt,
				&dbt.DeletedAt,
				&dbt.RestaurantID,
				&dbt.CapacityT,
			); err != nil {
				return
			}
			choutbl <- entities.TableRest{
				ID:           dbt.ID,
				RestaurantID: dbt.RestaurantID,
				CapacityT:    dbt.CapacityT,
			}
		}
	}()
	return choutbl, nil
}
