CREATE TABLE aphelios_website.blog (
	id         int PRIMARY KEY auto_increment,
	path       TEXT NOT NULL UNIQUE,
	title      TEXT NOT NULL,
	category   TEXT NOT NULL,
	summary    TEXT,
	text       TEXT,
	created_at datetime DEFAULT CURRENT_TIMESTAMP,
	updated_at datetime DEFAULT CURRENT_TIMESTAMP
); 