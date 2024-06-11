CREATE TABLE IF NOT EXIST users(
    id bigserial PRIMARY KEY,
    email citext UNIQUE NOT NULL,
    first_name  text NOT NULL
    last_name text NOT NULL
    password character varying(60),
    user_active integer DEFAULT 0,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
)


