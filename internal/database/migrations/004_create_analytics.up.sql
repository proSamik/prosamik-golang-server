-- Create the 'analytics' table
CREATE TABLE IF NOT EXISTS analytics (
                                         date DATE PRIMARY KEY,
                                         home_views INTEGER DEFAULT 0,
                                         about_views INTEGER DEFAULT 0,
                                         blogs_views INTEGER DEFAULT 0,
                                         projects_views INTEGER DEFAULT 0,
                                         feedback_views INTEGER DEFAULT 0,
                                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert initial row with today's date
INSERT INTO analytics (date, home_views, about_views, blogs_views, projects_views, feedback_views)
VALUES (CURRENT_DATE, 0, 0, 0, 0, 0)
ON CONFLICT DO NOTHING;