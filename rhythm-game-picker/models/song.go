// models/song.go

package models

type Song struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	Artist        string  `json:"artist"`
	Level         float64 `json:"level"` // 改为 float64 类型
	Color         string  `json:"color"`
	IsFavorite    bool    `json:"isFavorite"`
	IsBlacklisted bool    `json:"isBlacklisted"` // <-- 新增：黑名单字段
}

type RandomOptions struct {
	Count         int     `json:"count"`
	MinLevel      float64 `json:"minLevel"` // 改为 float64
	MaxLevel      float64 `json:"maxLevel"` // 改为 float64
	OnlyFavorites bool    `json:"onlyFavorites"`
}
