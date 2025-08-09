CREATE TABLE IF NOT EXISTS portfolio_assets(
    portfolio_id INTEGER NOT NULL,
    "symbol" VARCHAR(10) NOT NULL,
    amount numeric(18,8) NOT NULL,
    updated_at TIMESTAMP,

    PRIMARY KEY (portfolio_id, "symbol"),
    FOREIGN KEY (portfolio_id) REFERENCES portfolios(id) ON DELETE CASCADE
);

