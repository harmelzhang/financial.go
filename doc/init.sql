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
    CONSTRAINT `fk_parent_id` FOREIGN KEY (`parent_id`) REFERENCES `category` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=UTF8 COMMENT='证监会行业分类';

DROP TABLE IF EXISTS `stock`;
CREATE TABLE `stock` (
    `code` CHAR(6) NOT NULL COMMENT '股票代码',
    `stock_name` VARCHAR(16) DEFAULT NULL COMMENT '股票名称',
    `stock_name_pinyin` VARCHAR(16) DEFAULT NULL COMMENT '股票名称（拼音）',
    `company_name` VARCHAR(50) DEFAULT NULL COMMENT '公司名称',
    `organization` VARCHAR(16) DEFAULT NULL COMMENT '组织形式（民营、国营...）',
    `region` VARCHAR(8) DEFAULT NULL COMMENT '地域（省份）',
    `address` TEXT DEFAULT NULL COMMENT '办公地址',
    `website` VARCHAR(50) DEFAULT NULL COMMENT '公司网站',
    `main_business` TEXT DEFAULT NULL COMMENT '主营业务',
    `business_scope` TEXT DEFAULT NULL COMMENT '经营范围',
    `date_of_incorporation` DATE DEFAULT NULL COMMENT '成立日期',
    `listing_date` DATE DEFAULT NULL COMMENT '上市日期',
    `main_underwriter` VARCHAR(50) DEFAULT NULL COMMENT '主承销商',
    `sponsor` VARCHAR(50) DEFAULT NULL COMMENT '上市保荐人',
    `accounting_firm` VARCHAR(50) DEFAULT NULL COMMENT '会计师事务所',
    `market_place` CHAR(2) DEFAULT NULL COMMENT '交易市场（上海、深圳、北京）',
    `category_id` CHAR(8) DEFAULT NULL COMMENT '行业分类ID',
    PRIMARY KEY (`code`),
    KEY `i_stock_name` (`stock_name`),
    KEY `i_stock_name_pinyin` (`stock_name_pinyin`),
    CONSTRAINT `fk_category_id` FOREIGN KEY (`category_id`) REFERENCES `category` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=UTF8 COMMENT='股票信息';
