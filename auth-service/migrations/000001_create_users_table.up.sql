CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    email text UNIQUE NOT NULL,
    first_name  text NOT NULL,
    last_name text NOT NULL,
    password character varying(60),
    user_active integer DEFAULT 0,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
)


