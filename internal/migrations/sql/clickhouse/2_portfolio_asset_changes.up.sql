CREATE TABLE IF NOT EXISTS portfolio_asset_changes (
    portfolio_id Int32,
    symbol String,
    amount Float64,
    updated_at DateTime
) ENGINE = MergeTree()
ORDER BY (portfolio_id, symbol, updated_at);