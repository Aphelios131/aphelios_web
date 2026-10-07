 CREATE TABLE pages (
	id         SERIAL primary key,
	name       TEXT NOT NULL,
	version    INT  NOT NULL,
	content    TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



