CREATE TABLE IF NOT EXISTS domains (
    id INT PRIMARY KEY DEFAULT unique_rowid(),
    domain STRING NOT NULL,
    servers JSONB,
    servers_changed BOOL,
    ssl_grade STRING,
    previous_ssl_grade STRING,
    logo STRING,
    title STRING,
    is_down BOOL
);

CREATE INDEX IF NOT EXISTS idx_domains_servers ON domains USING gin (servers);


INSERT INTO domains (domain, servers, servers_changed, ssl_grade, previous_ssl_grade, logo, title, is_down) values ('truora.com', '[
    {
        "address": "server1",
        "ssl_grade": "B",
        "country": "US",
        "owner": "Amazon.com, Inc."
    },
    {
        "address": "server2",
        "ssl_grade": "A+",
        "country": "US",
        "owner": "Amazon.com, Inc."
    },
    {
        "address": "server3",
        "ssl_grade": "A",
        "country": "US",
        "owner": "Amazon.com, Inc."
    }
]', true, 'B', 'A+', 'https://server.com/icon.png​', 'Title of the page', false);