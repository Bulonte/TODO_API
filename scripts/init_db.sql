-- ============================================
-- Go-Todo-API 数据库初始化脚本
-- 版本: 1.0
-- 描述: 创建待办事项管理系统的数据库结构
-- 数据库: MySQL 8.0+
-- 字符集: utf8mb4
-- ============================================

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS `todo_db`
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

-- 使用新创建的数据库
USE `todo_db`;

-- 设置会话变量（确保操作一致性）
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 1. 删除已存在的表（按依赖关系逆序）
DROP TABLE IF EXISTS `todo_tags`;
DROP TABLE IF EXISTS `todos`;
DROP TABLE IF EXISTS `tags`;
DROP TABLE IF EXISTS `users`;

-- 2. 创建用户表 (users)
CREATE TABLE `users` (
                         `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
                         `username` VARCHAR(50) NOT NULL COMMENT '用户名',
                         `email` VARCHAR(100) NOT NULL COMMENT '邮箱',
                         `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希',
                         `avatar_url` VARCHAR(255) DEFAULT NULL COMMENT '头像URL',
                         `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用, 1-正常',
                         `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                         `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                         `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间',
                         PRIMARY KEY (`id`),
                         UNIQUE KEY `uk_username` (`username`) COMMENT '用户名唯一索引',
                         UNIQUE KEY `uk_email` (`email`) COMMENT '邮箱唯一索引',
                         KEY `idx_deleted_at` (`deleted_at`) COMMENT '软删除查询索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 3. 创建待办事项表 (todos)
CREATE TABLE `todos` (
                         `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '待办事项ID',
                         `user_id` INT UNSIGNED NOT NULL COMMENT '用户ID',
                         `title` VARCHAR(200) NOT NULL COMMENT '标题',
                         `description` TEXT COMMENT '描述',
                         `status` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '状态: 0-待办, 1-进行中, 2-已完成',
                         `priority` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '优先级: 1-低, 2-中, 3-高, 4-紧急',
                         `due_date` DATETIME DEFAULT NULL COMMENT '截止时间',
                         `completed_at` DATETIME DEFAULT NULL COMMENT '完成时间',
                         `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                         `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                         `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间',
                         PRIMARY KEY (`id`),
                         KEY `idx_user_id` (`user_id`) COMMENT '用户ID索引',
                         KEY `idx_status` (`status`) COMMENT '状态索引',
                         KEY `idx_priority` (`priority`) COMMENT '优先级索引',
                         KEY `idx_due_date` (`due_date`) COMMENT '截止时间索引',
                         KEY `idx_deleted_at` (`deleted_at`) COMMENT '软删除查询索引',
                         CONSTRAINT `fk_todos_user_id` FOREIGN KEY (`user_id`)
                             REFERENCES `users` (`id`)
                             ON DELETE CASCADE
                             ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='待办事项表';

-- 4. 创建标签表 (tags)
CREATE TABLE `tags` (
                        `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '标签ID',
                        `name` VARCHAR(50) NOT NULL COMMENT '标签名称',
                        `color` VARCHAR(7) NOT NULL DEFAULT '#1890ff' COMMENT '标签颜色',
                        `user_id` INT UNSIGNED NOT NULL COMMENT '用户ID',
                        `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                        `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `uk_user_name` (`user_id`, `name`) COMMENT '同一用户下标签名唯一',
                        KEY `idx_user_id` (`user_id`) COMMENT '用户ID索引',
                        CONSTRAINT `fk_tags_user_id` FOREIGN KEY (`user_id`)
                            REFERENCES `users` (`id`)
                            ON DELETE CASCADE
                            ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='标签表';

-- 5. 创建待办事项-标签关联表 (todo_tags)
CREATE TABLE `todo_tags` (
                             `todo_id` INT UNSIGNED NOT NULL COMMENT '待办事项ID',
                             `tag_id` INT UNSIGNED NOT NULL COMMENT '标签ID',
                             `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '关联创建时间',
                             PRIMARY KEY (`todo_id`, `tag_id`),
                             KEY `idx_tag_id` (`tag_id`) COMMENT '标签ID索引',
                             CONSTRAINT `fk_todo_tags_todo_id` FOREIGN KEY (`todo_id`)
                                 REFERENCES `todos` (`id`)
                                 ON DELETE CASCADE
                                 ON UPDATE CASCADE,
                             CONSTRAINT `fk_todo_tags_tag_id` FOREIGN KEY (`tag_id`)
                                 REFERENCES `tags` (`id`)
                                 ON DELETE CASCADE
                                 ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='待办事项标签关联表';

-- 6. 重新启用外键约束
SET FOREIGN_KEY_CHECKS = 1;