CREATE TABLE IF NOT EXISTS zoll_data(
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    bpm INT,
    time timestamp
);