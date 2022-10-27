CREATE TABLE IF NOT EXISTS games (
    games_id bigserial PRIMARY KEY,
    games_name text NOT NULL
);

CREATE TABLE IF NOT EXISTS tournaments (
    tournaments_id bigserial PRIMARY KEY,
    tournaments_name text NOT NULL,
    games_id bigint NOT NULL REFERENCES games ON DELETE CASCADE,
    start_date date NOT NULL DEFAULT NOW(),
    end_date date NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tournaments_prize_pool (
    tournaments_id bigint NOT NULL REFERENCES tournaments ON DELETE CASCADE,
    place smallint NOT NULL,
    money integer NOT NULL,
    PRIMARY KEY (tournaments_id, place)
);