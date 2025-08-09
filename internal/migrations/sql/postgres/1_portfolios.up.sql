
CREATE TABLE IF NOT EXISTS portfolios(
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL,
    name VARCHAR(50),
    is_public BOOLEAN,
    created_at TIMESTAMP
);
