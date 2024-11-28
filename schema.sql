-- Create 'User' Table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,                          -- Primary Key
    name VARCHAR(100) NOT NULL,                    -- Users name
    email VARCHAR(255) UNIQUE NOT NULL,            -- Unique email
    password TEXT NOT NULL,                        -- Encrypted password
    profile_image TEXT,                            -- Optional profile image URL
    refresh_token TEXT,                            -- For authentication
    verification_token TEXT,                       -- For account verification
    verification_token_expiry TIMESTAMPTZ ,          -- Expiry for verification token
    token_version INT DEFAULT 1 NOT NULL,          -- Keeps track of password resets
    is_verified BOOLEAN DEFAULT FALSE,            -- Verification status
    created_at TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP, -- Creation timestamp
    updated_at TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP, -- Update timestamp
    is_premium BOOLEAN DEFAULT FALSE,             -- Premium status
    max_urls INT DEFAULT 10,                       -- Max URLs allowed for user
    UNIQUE (email)                                -- Ensuring unique email
);

-- Create 'Url' Table
CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,                        -- Primary Key
    alias VARCHAR(20) UNIQUE NOT NULL,           -- Unique URL alias (shortened URL)
    url TEXT NOT NULL,                             -- Original URL
    created_at TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP, -- Creation timestamp
    updated_at TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP, -- Update timestamp
    expires_at TIMESTAMPTZ  NOT NULL,                          -- Expiration timestamp
    owner INT NOT NULL,                            -- Foreign Key reference to 'User'
    FOREIGN KEY (owner) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Indexes for Performance Optimization
CREATE INDEX IF NOT EXISTS idx_user_email ON users(email); -- Index for email lookup
CREATE INDEX IF NOT EXISTS idx_url_alias ON urls(alias);   -- Index for fast alias lookup
CREATE INDEX IF NOT EXISTS idx_url_owner ON urls(owner);   -- Index for filtering by owner


-- Add Triggers for Automatic Timestamps (if needed)
-- Example for 'Url' Table:
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_url_updated_at
BEFORE UPDATE ON urls
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Increment token_version on password change (Example Trigger)
CREATE OR REPLACE FUNCTION increment_token_version()
RETURNS TRIGGER AS $$
BEGIN
   NEW.token_version = NEW.token_version + 1;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_token_version_on_password_change
BEFORE UPDATE OF password ON users
FOR EACH ROW
WHEN (OLD.password IS DISTINCT FROM NEW.password)
EXECUTE FUNCTION increment_token_version();


-- Delete the expired URLs if they are older than 3 days
CREATE OR REPLACE FUNCTION delete_expired_urls()
RETURNS VOID AS $$
BEGIN
    DELETE FROM urls
    WHERE expires_at IS NOT NULL
      AND expires_at < NOW() - INTERVAL '3 days';
END;
$$ LANGUAGE plpgsql;
