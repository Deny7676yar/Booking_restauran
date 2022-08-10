CREATE TABLE public.bookings (
	id uuid NOT NULL,
	created_at timestamptz NOT NULL,
	deleted_at timestamptz,
	client_name varchar(200) NOT NULL,
	client_phone varchar(20) NOT NULL,
	booked_date date NOT NULL,
	booked_time_from timestamptz NOT NULL,
	booked_time_to timestamptz NOT NULL,
	CONSTRAINT bookings_pk PRIMARY KEY (id)
);