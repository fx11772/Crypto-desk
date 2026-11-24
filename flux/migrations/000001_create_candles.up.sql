CREATE TABLE candles (
    candle_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    exchange VARCHAR(50) NOT NULL,
    symbol VARCHAR(50) NOT NULL,
    timeframe VARCHAR(20) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    open NUMERIC NOT NULL,
    high NUMERIC NOT NULL,
    low NUMERIC NOT NULL,
    close NUMERIC NOT NULL,
    volume NUMERIC NOT NULL,

    CONSTRAINT candles_unique
        UNIQUE (exchange, symbol, timeframe, timestamp)
);