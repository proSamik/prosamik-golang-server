-- Remove the githubme tracking columns
ALTER TABLE analytics
    DROP COLUMN githubme_home,
    DROP COLUMN githubme_about,
    DROP COLUMN githubme_markdown;