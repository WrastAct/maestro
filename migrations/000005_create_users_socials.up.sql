CREATE TABLE IF NOT EXISTS nicknames (
    users_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    nickname text NOT NULL,
    PRIMARY KEY (users_id, nickname)
);

CREATE TABLE IF NOT EXISTS users_socials (
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    social_name text NOT NULL,
    social_link text NOT NULL,
    PRIMARY KEY (user_id, social_name)
);

CREATE INDEX idx_socials ON users_socials(social_name);