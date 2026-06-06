CREATE TABLE campaigns (
   id VARCHAR(36) NOT NULL PRIMARY KEY,
   user_id VARCHAR(36) NULL,
   name VARCHAR(255) NULL,
   short_description VARCHAR(255) NULL,
   description TEXT NULL,
   perks TEXT NULL,
   backer_count INT NULL,
   goal_amount INT NULL,
   current_amount INT NULL,
   slug VARCHAR(255) NULL,
   created_at TIMESTAMP NULL,
   updated_at TIMESTAMP NULL
);