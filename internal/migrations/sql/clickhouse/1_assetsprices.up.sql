CREATE TABLE IF NOT EXISTS asset_prices (
    symbol String,
    price Float64,
    ts DateTime
) ENGINE = MergeTree
ORDER BY (symbol, ts);
