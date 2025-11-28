// main.go
package main

import (
	"context"
	"math/rand"
	"rhythm-game-picker/database"
	"rhythm-game-picker/models"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
	db  *database.DB
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	db, err := database.NewDB(ctx)
	if err != nil {
		runtime.LogErrorf(ctx, "Database initialization error: %v", err)
	}
	a.db = db

}

func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}

func (a *App) GetAllSongs() ([]models.Song, error) {
	return a.db.GetAllSongs()
}

func (a *App) AddSong(song models.Song) error {
	if song.Color == "" {
		song.Color = a.generateMacaronColor()
	}
	return a.db.AddSong(song)
}

func (a *App) UpdateSong(song models.Song) error {
	return a.db.UpdateSong(song)
}

func (a *App) DeleteSong(id int) error {
	return a.db.DeleteSong(id)
}

func (a *App) ToggleFavorite(id int) error {
	return a.db.ToggleFavorite(id)
}

// 新增：暴露给前端的 ToggleBlacklist 方法
func (a *App) ToggleBlacklist(id int) error {
	return a.db.ToggleBlacklist(id)
}

func (a *App) RandomPick(options models.RandomOptions) ([]models.Song, error) {
	songs, err := a.db.GetAllSongs()
	if err != nil {
		return nil, err
	}

	var filtered []models.Song
	for _, song := range songs {
		// 修改：首先过滤掉黑名单中的歌曲
		if song.IsBlacklisted {
			continue
		}

		if options.OnlyFavorites && !song.IsFavorite {
			continue
		}

		levelInRange := false
		if options.MinLevel == 0 && options.MaxLevel == 0 {
			levelInRange = true
		} else {
			minOK := options.MinLevel == 0 || song.Level >= options.MinLevel
			maxOK := options.MaxLevel == 0 || song.Level <= options.MaxLevel
			if minOK && maxOK {
				levelInRange = true
			}
		}

		if !levelInRange {
			continue
		}

		filtered = append(filtered, song)
	}

	if len(filtered) == 0 {
		return []models.Song{}, nil
	}

	count := options.Count
	if count > len(filtered) {
		count = len(filtered)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(filtered), func(i, j int) {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	})

	return filtered[:count], nil
}

func (a *App) generateMacaronColor() string {
	colors := []string{
		"#FFB3BA", "#FFDFBA", "#FFFFBA", "#BAFFC9", "#BAE1FF",
		"#E0BBE4", "#FFB7D5", "#C7CEEA", "#B4E7F5", "#FED9B7",
		"#A0E7E5", "#F7C6C7", "#C9E4DE", "#FFD1DC", "#E4C1F9",
	}
	return colors[rand.Intn(len(colors))]
}
