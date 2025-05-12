-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_account(
  id UUID NOT NULL,
  email TEXT NOT NULL,
  full_name TEXT NOT NULL,
  github_username TEXT NULL,
  bounty INT NOT NULL DEFAULT 0,
  refresh_token TEXT DEFAULT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT "user_account_pkey" PRIMARY KEY (id)
);
CREATE UNIQUE INDEX idx_user_account__email ON user_account(email);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS otps (
  id UUID NOT NULL,
  email TEXT NOT NULL,
  full_name TEXT NOT NULL,
  otp TEXT NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  
  CONSTRAINT "otps_pkey" PRIMARY KEY (id)
);
-- +goose StatementEnd
