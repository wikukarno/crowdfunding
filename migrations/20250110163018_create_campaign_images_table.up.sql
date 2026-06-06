CREATE TABLE campaign_images (
     id VARCHAR(36) NOT NULL PRIMARY KEY,
     campaign_id VARCHAR(36) NULL,
     file_name VARCHAR(255) NULL,
     is_primary SMALLINT NULL,
     created_at TIMESTAMP NULL,
     updated_at TIMESTAMP NULL
);