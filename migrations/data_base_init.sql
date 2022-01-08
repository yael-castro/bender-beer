CREATE TABLE beers(
    beer_id  SERIAL PRIMARY KEY,
    name     VARCHAR(20) NOT NULL,
    brewery  VARCHAR(50) NOT NULL,
    country  VARCHAR(50) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    price    DOUBLE PRECISION DEFAULT 0,
    UNIQUE (name, brewery)
);