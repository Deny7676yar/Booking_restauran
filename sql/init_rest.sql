CREATE TABLE public.restaurants (
	id uuid NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	deleted_at timestamptz NULL,
	title varchar NOT NULL,
	average_waiting_time int2 NULL,
	average_check bigint NOT NULL,
	available_seats int2 NULL,
	CONSTRAINT restaurants_pk PRIMARY KEY (id)
);