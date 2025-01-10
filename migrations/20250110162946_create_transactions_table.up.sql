CREATE TABLE transactions (
      id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
      campaign_id INT NULL,
      user_id INT NULL,
      amount INT NULL,
      status VARCHAR(255) NULL,
      code VARCHAR(255) NULL,
      payment_url VARCHAR(255) NOT NULL,
      created_at DATETIME NULL,
      updated_at DATETIME NULL
);