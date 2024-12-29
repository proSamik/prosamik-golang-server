-- Create a sequence for blog posts
CREATE SEQUENCE IF NOT EXISTS blogs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MAXVALUE
    NO CYCLE;

-- Create the blog table
CREATE TABLE IF NOT EXISTS blogs (
                                     id INTEGER PRIMARY KEY DEFAULT nextval('blogs_id_seq'),
                                     title VARCHAR(255) NOT NULL UNIQUE,
                                     path TEXT NOT NULL UNIQUE,
                                     description TEXT,
                                     tags TEXT,
                                     views_count INTEGER DEFAULT 0
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_blogs_title ON blogs(title);
CREATE INDEX IF NOT EXISTS idx_blogs_path ON blogs(path);
CREATE INDEX IF NOT EXISTS idx_blogs_tags ON blogs(tags);