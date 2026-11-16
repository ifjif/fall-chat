package model

/*
  - CREATE TABLE `emoji_packs` (
    `id` bigint unsigned NOT NULL,
    `public_id` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` bigint unsigned DEFAULT NULL,
    `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `cover_url` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT ”,
    `owner_id` bigint unsigned NOT NULL,
    `source` tinyint NOT NULL DEFAULT '1',
    `price_type` tinyint NOT NULL DEFAULT '1',
    `price` decimal(2,0) NOT NULL DEFAULT '0',
    `status` tinyint NOT NULL DEFAULT '1',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_emoji_packs_public_id` (`public_id`),
    KEY `idx_emoji_packs_deleted_at` (`deleted_at`),
    KEY `idx_emoji_packs_owner_id` (`owner_id`),
    KEY `idx_source_status` (`source`,`status`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
*
*
*/

var NilEmojiPack *EmojiPack = nil

type EmojiPack struct {
	Model
	Name      string  `gorm:"size:64;not null"`
	CoverURL  string  `gorm:"size:512;default:''"`
	OwnerID   uint64  `gorm:"index;not null"`
	Source    int8    `gorm:"index:idx_source_status;default:1;not null"` // 1:系统 2:用户
	PriceType int8    `gorm:"default:1;not null"`                         // 1免费 2付费
	Price     float64 `gorm:"precision:2;default:0;not null"`
	Status    int8    `gorm:"index:idx_source_status;default:1;not null"` // 0:下架 1:草稿 2:上架 3:审核中
}

func (m *EmojiPack) GetSourceText() string {
	switch m.Source {
	case 1:
		return "system"
	case 2:
		return "customer"
	}

	return "unknown"
}

func (m *EmojiPack) EncodeOwnerID() string {
	return EncodeID(int64(m.OwnerID))
}

func (m *EmojiPack) GetPriceTypeText() string {
	switch m.PriceType {
	case 1:
		return "free"
	case 2:
		return "payment"
	}
	return "unknown"
}

func (m *EmojiPack) GetStatusText() string {
	switch m.Status {
	case 0:
		return "soldout"
	case 1:
		return "draft"
	case 2:
		return "added"
	case 3:
		return "audit"
	}

	return "unknown"
}

func (m *EmojiPack) GetDraftStatusNilSafe() int {
	return 1
}

func (m *EmojiPack) GetAddedStatusNilSafe() int {
	return 2
}
