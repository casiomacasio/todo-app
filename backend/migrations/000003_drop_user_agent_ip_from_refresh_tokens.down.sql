ALTER TABLE refresh_tokens
ADD COLUMN user_agent TEXT,
ADD COLUMN ip_address TEXT;
