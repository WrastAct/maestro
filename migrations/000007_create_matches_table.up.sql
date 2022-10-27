CREATE TABLE IF NOT EXISTS matches (
    matches_id bigserial NOT NULL,
    tournaments_id bigint NOT NULL REFERENCES tournaments ON DELETE CASCADE,
    match_data text NOT NULL,
    PRIMARY KEY (matches_id, tournaments_id)
);

CREATE TABLE IF NOT EXISTS users_matches (
    users_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    matches_id bigint NOT NULL,
    tournaments_id bigint NOT NULL,
    result text NOT NULL,
    average_stress decimal NOT NULL,
    humidity decimal NOT NULL,
    temperature decimal NOT NULL,
    pressure decimal NOT NULL,
    FOREIGN KEY (matches_id, tournaments_id) REFERENCES matches(matches_id, tournaments_id) ON DELETE CASCADE,
    PRIMARY KEY (users_id, matches_id, tournaments_id)
);