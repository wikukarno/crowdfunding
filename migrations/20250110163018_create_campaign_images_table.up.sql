CREATE TABLE campaign_images (
     id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
     campaign_id INT NULL,
     file_name VARCHAR(255) NULL,
     is_primary TINYINT NULL,
     created_at DATETIME NULL,
     updated_at DATETIME NULL
);