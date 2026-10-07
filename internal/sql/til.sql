
CREATE TABLE til (
	id         SERIAL8 PRIMARY KEY,
	path       TEXT NOT NULL UNIQUE,
	title      TEXT NOT NULL,
	category   TEXT NOT NULL,
	summary    TEXT,
	text       TEXT,
	created_at DATE DEFAULT CURRENT_DATE,
	updated_at DATE DEFAULT CURRENT_DATE
);