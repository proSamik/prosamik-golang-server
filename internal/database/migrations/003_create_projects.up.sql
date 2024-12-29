-- Create a sequence for projects
CREATE SEQUENCE IF NOT EXISTS projects_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MAXVALUE
    NO CYCLE;
-- This creates an auto-incrementing sequence for unique project IDs, similar to how blogs_id_seq works

-- Create the 'projects' table
CREATE TABLE IF NOT EXISTS projects (
                                     id INTEGER PRIMARY KEY DEFAULT nextval('projects_id_seq'),
                                     title VARCHAR(255) NOT NULL UNIQUE,
                                     path TEXT NOT NULL UNIQUE,
                                     description TEXT,
                                     tags TEXT,
                                     views_count INTEGER DEFAULT 0
);
-- Create indexes for frequently accessed columns
CREATE INDEX IF NOT EXISTS idx_projects_title ON projects(title);
CREATE INDEX IF NOT EXISTS idx_projects_path ON projects(path);
CREATE INDEX IF NOT EXISTS idx_projects_technologies ON projects(tags);