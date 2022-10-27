CREATE TABLE IF NOT EXISTS teams (
    teams_id bigserial PRIMARY KEY,
    teams_name text NOT NULL,
    teams_description text NOT NULL,
    teams_region text NOT NULL
);

CREATE TABLE IF NOT EXISTS teams_users (
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    teams_id bigint NOT NULL REFERENCES teams ON DELETE CASCADE,
    join_date date NOT NULL DEFAULT NOW(),
    leave_date date NOT NULL DEFAULT NOW(),
    role text NOT NULL,
    PRIMARY KEY (user_id, teams_id, join_date)
);
