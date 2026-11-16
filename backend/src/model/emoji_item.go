package model

/*
  - CREATE TABLE `emoji_items` (
    `id` bigint unsigned NOT NULL,
    `public_id` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` bigint unsigned DEFAULT NULL,
    `pack_id` bigint unsigned NOT NULL,
    `media_type` tinyint NOT NULL,
    `media_url` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT ”,
    `font_name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT ”,
    `thumb_url` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT ”,
    `width` smallint unsigned NOT NULL DEFAULT '0',
    `height` smallint unsigned NOT NULL DEFAULT '0',
    `file_size` bigint unsigned NOT NULL DEFAULT '0',
    `sort_no` bigint NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_emoji_items_public_id` (`public_id`),
    KEY `idx_emoji_items_deleted_at` (`deleted_at`),
    KEY `idx_emoji_items_pack_id` (`pack_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
*
*
*/
type EmojiItem struct {
	Model
	PackID    uint64 `gorm:"index;not null"`
	MediaType int8   `gorm:"not null"` // 1静态图 2动态图 3字体 4unicode
	MediaURL  string `gorm:"size:512;default:'';not null"`
	FontName  string `gorm:"size:128;default:'';not null"`
	ThumbURL  string `gorm:"size:512;default:'';not null"`
	Width     uint16 `gorm:"default:0;not null"`
	Height    uint16 `gorm:"default:0;not null"`
	FileSize  uint   `gorm:"default:0;not null"`
	SortNO    int    `gorm:"default:0;not null"`
}

func (m *EmojiItem) EncodePackID() string {
	return EncodeID(int64(m.PackID))
}

func (m *EmojiItem) GetMediaTypeText() string {
	switch m.MediaType {
	case 1:
		return "static"
	case 2:
		return "dynamic"
	case 3:
		return "font"
	case 4:
		return "unicode"
	}

	return "unknown"
}
