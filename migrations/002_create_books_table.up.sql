CREATE TABLE IF NOT EXISTS books (
   id SERIAL PRIMARY KEY,
   isbn VARCHAR(255) UNIQUE NOT NULL,
    title VARCHAR(255),
    author VARCHAR(255),
    published INT,
    pages INT,
    status INT DEFAULT 0
);