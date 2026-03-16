CREATE TABLE jobs (
    id SERIAL PRIMARY KEY,
    external_id TEXT UNIQUE NOT NULL,       -- Unique ID from source or URL hash
    title TEXT NOT NULL,
    company TEXT,
    description TEXT,
    location TEXT,
    url TEXT,
    date_posted DATE,
    source TEXT,
    embedding VECTOR(1536),                 -- Embedding for semantic search
    active BOOLEAN DEFAULT TRUE,
    last_checked TIMESTAMP DEFAULT now()
);