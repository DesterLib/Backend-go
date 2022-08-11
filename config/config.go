package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/anonyindian/logger"
	"github.com/joho/godotenv"
)

const (
	DEFAULT_CONFIG_LOCATION = "./config.json"
	DEFAULT_API_V1_ENDPOINT = "/api/v1"
)

var ValueOf = &config{}

type config struct {
	DatabaseURI      string
	Port             int
	RcloneListenPort int
	DesterDev        bool
	OnHeroku         bool `json:"-"`
}

func (c *config) setupEnvVars() {
	_ = godotenv.Load()
	val := reflect.ValueOf(c).Elem()
	for i := 0; i < val.NumField(); i++ {
		envVal := os.Getenv(toUpperSnakeCase(val.Type().Field(i).Name))
		switch val.Type().Field(i).Type.Kind() {
		case reflect.Int:
			ev, _ := strconv.ParseInt(envVal, 10, 64)
			val.Field(i).SetInt(ev)
		case reflect.String:
			val.Field(i).SetString(envVal)
		case reflect.Bool:
			ev, _ := strconv.ParseBool(envVal)
			val.Field(i).SetBool(ev)
		}
	}
}

func Load(log *logger.Logger) {
	f, err := ioutil.ReadFile(DEFAULT_CONFIG_LOCATION)
	if err == nil {
		err := json.Unmarshal(f, &ValueOf)
		if err != nil {
			ValueOf.setupEnvVars()
		}
	} else {
		ValueOf.setupEnvVars()
	}
	log.ChangeLevel(logger.LevelMain)
	log.Println("LOADED CONFIG VARS")
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toUpperSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToUpper(snake)
}
