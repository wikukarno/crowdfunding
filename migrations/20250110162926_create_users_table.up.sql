CREATE TABLE users (
   id VARCHAR(36) NOT NULL PRIMARY KEY,
   name VARCHAR(255) NULL,
   occupation VARCHAR(255) NULL,
   email VARCHAR(255) NULL,
   password_hash VARCHAR(255) NULL,
   avatar_file_name VARCHAR(255) NULL,
   role VARCHAR(255) NULL,
   token VARCHAR(255) NULL,
   created_at TIMESTAMP NULL,
   updated_at TIMESTAMP NULL
);