package model

import "time"

/*
  - CREATE TABLE `emoji_owns` (
    `id` bigint unsigned NOT NULL,
    `public_id` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` bigint unsigned DEFAULT NULL,
    `user_id` bigint unsigned NOT NULL,
    `item_id` bigint unsigned NOT NULL,
    `pack_id` bigint unsigned NOT NULL DEFAULT '0',
    `source_type` tinyint NOT NULL,
    `expire_at` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_emoji_owns_public_id` (`public_id`),
    KEY `idx_emoji_owns_deleted_at` (`deleted_at`),
    KEY `idx_user_item` (`user_id`,`item_id`),
    KEY `idx_user_pack` (`user_id`,`pack_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
*
*
*/
type EmojiOwn struct {
	Model
	UserID     uint64    `gorm:"index:idx_user_item;index:idx_user_pack;not null"`
	ItemID     uint64    `gorm:"index:idx_user_item;not null"`
	PackID     uint64    `gorm:"index:idx_user_pack;default:0;not null"`
	SourceType int8      `gorm:"not null"` // 1系统默认 2自己上传 3购买
	ExpireAt   time.Time `gorm:"default:null"`
}

func (e *EmojiOwn) SetSelfSource() {
	e.SourceType = 2
}
