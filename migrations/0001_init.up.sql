CREATE TABLE IF NOT EXISTS report_stat
(
    id INTEGER PRIMARY KEY,
    report_type_name VARCHAR(100) NOT NULL,
    load_time_sec INTEGER
);
