-- Create sequence for auto-incrementing the `id` field
CREATE SEQUENCE IF NOT EXISTS newsletter_subscriptions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MAXVALUE
    NO CYCLE;

-- Create the `newsletter_subscriptions` table
CREATE TABLE IF NOT EXISTS newsletter_subscriptions (
                                                        id INTEGER PRIMARY KEY DEFAULT nextval('newsletter_subscriptions_id_seq'),
                                                        email VARCHAR(255) UNIQUE NOT NULL,
                                                        registration_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                                        verified BOOLEAN DEFAULT FALSE
);

-- Reset the sequence to the correct value, based on existing data
SELECT setval('newsletter_subscriptions_id_seq',
              COALESCE((SELECT MAX(id) FROM newsletter_subscriptions), 1), false);

-- Create an index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_newsletter_email ON newsletter_subscriptions(email);
