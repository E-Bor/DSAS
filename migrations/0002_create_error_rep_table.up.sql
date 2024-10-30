CREATE TABLE IF NOT EXISTS report_errors
(
    id INTEGER PRIMARY KEY,
    report_type_name VARCHAR(100) NOT NULL,
    trace_id VARCHAR(18) NOT NULL,
    load_error VARCHAR(100) NOT NULL
);
