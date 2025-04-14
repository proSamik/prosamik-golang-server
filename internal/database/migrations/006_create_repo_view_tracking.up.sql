CREATE TABLE IF NOT EXISTS repo_views (
    id SERIAL PRIMARY KEY,
    repo_path TEXT NOT NULL,
    view_date DATE NOT NULL DEFAULT CURRENT_DATE,
    view_count INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE (repo_path, view_date)
);

-- Index for faster lookups by repo_path
CREATE INDEX idx_repo_views_repo_path ON repo_views(repo_path);

-- Index for date range queries
CREATE INDEX idx_repo_views_view_date ON repo_views(view_date); 