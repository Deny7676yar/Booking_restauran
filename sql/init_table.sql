CREATE TABLE public.tables (
	id uuid NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	deleted_at timestamptz NULL,
	restaurant_id uuid NOT NULL,
	capacity int2 NULL,
	CONSTRAINT fk_tables_restaurants FOREIGN KEY (restaurant_id) REFERENCES restaurants (id) ON DELETE CASCADE
);