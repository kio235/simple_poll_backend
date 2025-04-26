-- 创建数据库
DROP DATABASE IF EXISTS poll_db;
CREATE DATABASE poll_db 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;


-- 使用数据库
USE poll_db;

-- 创建questions表
CREATE TABLE IF NOT EXISTS questions (
    id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    deleted_at DATETIME(3) NULL,
    content TEXT NOT NULL
) ENGINE=InnoDB 
  CHARACTER SET utf8mb4
  COLLATE=utf8mb4_unicode_ci;

-- 创建options表
CREATE TABLE IF NOT EXISTS options (
    id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    deleted_at DATETIME(3) NULL,
    question_id INT UNSIGNED NOT NULL,
    text TEXT NOT NULL,
    votes INT NOT NULL DEFAULT 0,
    FOREIGN KEY (question_id) 
        REFERENCES questions(id) 
        ON DELETE CASCADE
) ENGINE=InnoDB 
  CHARACTER SET utf8mb4
  COLLATE=utf8mb4_unicode_ci;

-- 创建索引优化查询性能
CREATE INDEX idx_questions_deleted_at ON questions(deleted_at);
CREATE INDEX idx_options_deleted_at ON options(deleted_at);
CREATE INDEX idx_options_question ON options(question_id);

-- 创建数据库用户并授权（根据实际需要调整host）
CREATE USER IF NOT EXISTS 'poll_db'@'%' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON poll_db.* TO 'poll_db'@'%';
FLUSH PRIVILEGES;