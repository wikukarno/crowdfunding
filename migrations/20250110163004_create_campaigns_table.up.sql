CREATE TABLE campaigns (
   id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
   user_id INT NULL,
   name VARCHAR(255) NULL,
   short_description VARCHAR(255) NULL,
   description TEXT NULL,
   perks TEXT NULL,
   backer_count INT NULL,
   goal_amount INT NULL,
   current_amount INT NULL,
   slug VARCHAR(255) NULL,
   created_at DATETIME NULL,
   updated_at DATETIME NULL
);