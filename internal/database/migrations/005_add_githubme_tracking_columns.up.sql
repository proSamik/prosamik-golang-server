-- Add new columns for githubme tracking
ALTER TABLE analytics
    ADD COLUMN githubme_home INTEGER DEFAULT 0,
    ADD COLUMN githubme_about INTEGER DEFAULT 0,
    ADD COLUMN githubme_markdown INTEGER DEFAULT 0;

-- Update existing rows to have default values
UPDATE analytics
SET
    githubme_home = 0,
    githubme_about = 0,
    githubme_markdown = 0
WHERE
    githubme_home IS NULL OR
    githubme_about IS NULL OR
    githubme_markdown IS NULL;