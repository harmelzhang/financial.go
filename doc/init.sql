CREATE DATABASE IF NOT EXISTS `financial` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `financial`;

SET NAMES utf8mb4;

DROP TABLE IF EXISTS `index`;
CREATE TABLE `index_sample` (
    `type_code` VARCHAR(20) NOT NULL COMMENT '类型代码（中证指数，www.csindex.com.cn）',
    `type_name` VARCHAR(20) NOT NULL COMMENT '类型名称（沪深300、中证500、上证50....）',
    `stock_code` CHAR(6) NOT NULL COMMENT '股票代码'
) ENGINE=InnoDB DEFAULT CHARSET=UTF8 COMMENT='指数样本信息（来源：中证指数）';

DROP TABLE IF EXISTS `category`;
CREATE TABLE `category` (
    `id` CHAR(8) NOT NULL COMMENT '网易分类ID',
    `name` VARCHAR(20) NOT NULL COMMENT '证券会分类名称',
    `display_order` TINYINT(3) UNSIGNED DEFAULT NULL COMMENT '显示顺序',
    `parent_id` CHAR(8) DEFAULT NULL COMMENT '父分类ID',
    PRIMARY KEY (`id`),
    KEY `i_parent_id` (`parent_id`) USING BTREE,
    CONSTRAINT `fk_parent_id` FOREIGN KEY (`parent_id`) REFERENCES `category` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=UTF8 COMMENT='证监会行业分类';
