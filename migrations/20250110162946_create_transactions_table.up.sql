CREATE TABLE transactions (
      id VARCHAR(36) NOT NULL PRIMARY KEY,
      campaign_id VARCHAR(36) NULL,
      user_id VARCHAR(36) NULL,
      amount INT NULL,
      status VARCHAR(255) NULL,
      code VARCHAR(255) NULL,
      payment_url VARCHAR(255) NOT NULL,
      created_at TIMESTAMP NULL,
      updated_at TIMESTAMP NULL
);