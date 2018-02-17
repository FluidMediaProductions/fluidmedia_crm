CREATE TABLE IF NOT EXISTS contacts (
  id SERIAL PRIMARY KEY,
  image TEXT NOT NULL,
  name TEXT NOT NULL,
  state INTEGER NOT NULL,
  email TEXT NOT NULL,
  phone TEXT NOT NULL,
  mobile TEXT NOT NULL,
  website TEXT NOT NULL,
  twitter TEXT NOT NULL,
  address TEXT NOT NULL,
  description TEXT NOT NULL);

CREATE TABLE IF NOT EXISTS organisations (
  id SERIAL PRIMARY KEY,
  image TEXT NOT NULL,
  name TEXT NOT NULL,
  email TEXT NOT NULL,
  phone TEXT NOT NULL,
  website TEXT NOT NULL,
  twitter TEXT NOT NULL,
  address TEXT NOT NULL,
  description TEXT NOT NULL);