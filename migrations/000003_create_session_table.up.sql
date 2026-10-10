CREATE TABLE sessions (
  token char(43) NOT NULL,
  data bytea NOT NULL,
  expiry timestamp(6) NOT NULL,
  PRIMARY KEY (token)
);

CREATE INDEX sessions_expiry_idx ON sessions(expiry);