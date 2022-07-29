package db

import (
	"encoding/json"

	"github.com/desterlib/backend-go/cache"
)

type Config struct {
	Auth0      Auth0      `json:"auth0"`
	App        App        `json:"app"`
	Gdrive     Hosting    `json:"gdrive"`
	Sharepoint Hosting    `json:"sharepoint"`
	Onedrive   Hosting    `json:"onedrive"`
	Tmdb       ApiService `json:"tmdb"`
	Subtitles  ApiService `json:"subtitles"`
	Rclone     Rclone     `json:"rclone"`
	Build      Build      `json:"build"`
}

// {"auth0":{},"categories":[],"gdrive":{},"onedrive":{},"sharepoint":{},"tmdb":{"api_key":"ygy"},"subtitles":{"api_key":"gyg"},"build":{"cron":"gygy"},"rclone":{},"app":{"name":"feuh","title":"uhuhu","description":"hu","domain":"huh","secret_key":"uhhu"}}
type App struct {
	Name        string `json:"name,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Domain      string `json:"domain,omitempty"`
	SecretKey   string `json:"secret_key,omitempty"`
}

type ApiService struct {
	ApiKey string `json:"api_key,omitempty"`
}

type Build struct {
	Cron string `json:"cron,omitempty"`
}

type Rclone struct {
}

type Auth0 struct {
}

type Hosting struct {
}

type ConfigDB struct {
	Id     int `gorm:"primary_key"`
	Config []byte
}

func (c *Config) Json() (b []byte) {
	b, _ = json.Marshal(c)
	return
}

const DEFAULT_ID = 777000

func SaveConfig(config []byte) {
	w := &ConfigDB{
		Id: DEFAULT_ID,
	}
	tx := SESSION.Begin()
	tx.FirstOrCreate(w)
	w.Config = config
	tx.Save(w)
	tx.Commit()
	cacheConfig()
}

func GetConfig() *Config {
	data, err := cache.Cache.Get("config")
	if err != nil {
		data = cacheConfig()
	}
	var config Config
	_ = json.Unmarshal(data, &config)
	return &config
}

func cacheConfig() []byte {
	configDB := &ConfigDB{Id: DEFAULT_ID}
	SESSION.Where("id = ?", DEFAULT_ID).Find(&configDB)
	cache.Cache.Set("config", configDB.Config)
	return configDB.Config
}
