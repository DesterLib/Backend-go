package rclone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/desterlib/backend-go/config"
)

type RCloneAPI struct {
	Data     map[string]interface{} `json:"data"`
	Index    int                    `json:"index"`
	ID       string                 `json:"id"`
	FS       string                 `json:"fs"`
	Provider string                 `json:"provider"`
	RCLONE   map[string]string      `json:"RCLONE"`
	FSConf   map[string]interface{} `json:"fs_conf"`
	Port     int32                  `json:"port"`
}

func post(url string, body any) any {
	postBody, _ := json.Marshal(body)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return nil
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result
}

func getPaths(data any, paths []any) any {
	for _, path := range paths {
		switch data.(type) {
		case []any:
			data = data.([]any)[path.(int)]
		case map[string]any:
			data = data.(map[string]any)[path.(string)]
		}
	}
	return data
}
func NewRCloneAPI(data map[string]interface{}, index int) *RCloneAPI {
	var id = data["id"]
	if id == nil {
		id = data["drive_id"]
	}
	fs := ""
	for _, c := range id.(string) {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' {
			fs += string(c)
		}
	}
	fs += ":"
	provider := data["provider"]
	if provider == nil {
		provider = "gdrive"
	}
	RCLONE := map[string]string{
		"mkdir":              "operations/mkdir",
		"purge":              "operations/purge",
		"deleteFile":         "operations/deletefile",
		"createPublicLink":   "operations/publiclink",
		"stats":              "core/stats",
		"bwlimit":            "core/bwlimit",
		"moveDir":            "sync/move",
		"moveFile":           "operations/movefile",
		"copyDir":            "sync/copy",
		"copyFile":           "operations/copyfile",
		"cleanUpRemote":      "operations/cleanup",
		"noopAuth":           "rc/noopauth",
		"getRcloneVersion":   "core/version",
		"getRcloneMemStats":  "core/memstats",
		"getOptions":         "options/get",
		"getProviders":       "config/providers",
		"getConfigDump":      "config/dump",
		"getRunningJobs":     "job/list",
		"getStatusForJob":    "job/status",
		"getConfigForRemote": "config/get",
		"createConfig":       "config/create",
		"updateConfig":       "config/update",
		"getFsInfo":          "operations/fsinfo",
		"listRemotes":        "config/listremotes",
		"getFilesList":       "operations/list",
		"getAbout":           "operations/about",
		"deleteConfig":       "config/delete",
		"stopJob":            "job/stop",
		"backendCommand":     "backend/command",
		"coreCommand":        "core/command",
		"transferred":        "core/transferred",
		"getSize":            "operations/size",
		"getFileInfo":        "operations/stat",
		"statsDelete":        "core/stats-delete",
		"statsReset":         "core/stats-reset",
	}
	rclonePort, _ := strconv.ParseInt(config.ValueOf.DatabaseURI, 10, 64)
	fsConf := rcConf(fs[:len(fs)-1], RCLONE, rclonePort)
	return &RCloneAPI{
		Data:     data,
		Index:    index,
		ID:       id.(string),
		FS:       fs,
		Provider: provider.(string),
		RCLONE:   RCLONE,
		FSConf:   fsConf,
	}
}

func (r *RCloneAPI) rcLs(fs string, options map[string]interface{}) []map[string]interface{} {
	rcData := map[string]interface{}{
		"fs":     fs,
		"remote": "",
		"opt":    options,
	}
	result := post("http://localhost:"+strconv.FormatInt(int64(r.Port), 10)+"/"+r.RCLONE["getConfigForRemote"], rcData)
	return result.(map[string]interface{})["list"].([]map[string]interface{})
}

func rcConf(name string, mapPaths map[string]string, port int64) map[string]interface{} {
	rcData := map[string]string{
		"name": name,
	}

	result := post("http://localhost:"+strconv.FormatInt(port, 10)+"/"+mapPaths["getConfigForRemote"], rcData)

	if result.(map[string]interface{})["token"] != nil {
		result.(map[string]interface{})["token"] = json.Unmarshal([]byte(result.(map[string]interface{})["token"].(string)), &map[string]interface{}{})
	}
	return result.(map[string]interface{})
}

func (r *RCloneAPI) FetchMovies() []map[string]interface{} {
	rcLsResult := r.rcLs(r.FS, map[string]interface{}{
		"recurse":   true,
		"filesOnly": false,
	})
	metadata := []map[string]interface{}{}
	dirs := map[string]map[string]interface{}{}
	fileNames := map[string]map[string]interface{}{}
	subIndex := 0
	for _, item := range rcLsResult {
		if !(item["IsDir"].(bool)) && (strings.Contains(item["MimeType"].(string), "video") || strings.HasSuffix(strings.ToLower(item["Name"].(string)), ".mp4") || strings.HasSuffix(strings.ToLower(item["Name"].(string)), ".mkv") || strings.HasSuffix(strings.ToLower(item["Name"].(string)), ".avi") || strings.HasSuffix(strings.ToLower(item["Name"].(string)), ".mov") || strings.HasSuffix(strings.ToLower(item["Name"].(string)), ".webm") || strings.HasSuffix(strings.ToLower(item["Name"].(string)), ".flv")) {
			parentPath := strings.Replace(item["Path"].(string), "/"+item["Name"].(string), "", -1)
			parent := dirs[parentPath]
			var itemID string
			if item["ID"] != nil {
				itemID = item["ID"].(string)
			} else {
				itemID = item["Path"].(string)
			}
			currMetadata := map[string]interface{}{
				"id":            itemID,
				"name":          item["Name"],
				"path":          item["Path"],
				"parent":        parent,
				"mime_type":     item["MimeType"],
				"size":          item["Size"],
				"subtitles":     []map[string]interface{}{},
				"modified_time": item["ModTime"],
			}
			pathWithoutExtension := path.Dir(item["Path"].(string))
			fileName := fileNames[pathWithoutExtension]
			if fileName != nil {
				currMetadata["subtitles"] = fileName["subtitles"]
				fileNames[pathWithoutExtension]["found"] = true
				fileNames[pathWithoutExtension]["index"] = subIndex
			} else {
				fileNames[pathWithoutExtension] = map[string]interface{}{
					"found":     true,
					"index":     subIndex,
					"subtitles": []map[string]interface{}{},
				}
			}
			metadata = append(metadata, currMetadata)
			subIndex += 1
		} else if item["IsDir"].(bool) {
			dirs[item["Path"].(string)] = map[string]interface{}{
				"id":   item["ID"],
				"name": item["Name"],
				"path": item["Path"],
			}
		} else if (!item["IsDir"].(bool)) && strings.HasSuffix(item["Name"].(string), ".vtt") || strings.HasSuffix(item["Name"].(string), ".srt") || strings.HasSuffix(item["Name"].(string), ".ass") || strings.HasSuffix(item["Name"].(string), ".ssa") {
			pathWithoutExtension := path.Dir(item["Path"].(string))
			if pathWithoutExtension[len(pathWithoutExtension)-3] == '.' {
				pathWithoutExtension = pathWithoutExtension[:len(pathWithoutExtension)-3]
			} else if pathWithoutExtension[len(pathWithoutExtension)-4] == '.' {
				pathWithoutExtension = pathWithoutExtension[:len(pathWithoutExtension)-4]
			}
			subMetadata := map[string]interface{}{
				"id":   item["ID"],
				"name": item["Name"],
				"path": item["Path"],
			}
			fileName := fileNames[pathWithoutExtension]
			if fileName != nil {
				if fileName["found"].(bool) {
					metadata[fileName["index"].(int)]["subtitles"] = append(metadata[fileName["index"].(int)]["subtitles"].([]map[string]interface{}), subMetadata)
				} else {
					fileNames[pathWithoutExtension]["subtitles"] = append(fileNames[pathWithoutExtension]["subtitles"].([]map[string]interface{}), subMetadata)
				}
			} else {
				fileNames[pathWithoutExtension] = map[string]interface{}{
					"found":     false,
					"index":     nil,
					"subtitles": []map[string]interface{}{subMetadata},
				}
			}
		}
	}
	return metadata
}

func (r *RCloneAPI) FetchSeries() []map[string]interface{} {
	rcLsResult := r.rcLs(r.FS, map[string]interface{}{
		"recurse":  true,
		"maxDepth": 2,
	})
	metadata := []map[string]interface{}{}
	parentDirs := map[string]map[string]interface{}{
		"": map[string]interface{}{
			"path":      "",
			"depth":     0,
			"json_path": []any{},
		},
	}
	for _, item := range rcLsResult {
		var parentPath string = ""
		if len(strings.Split(item["Path"].(string), "/")) != 1 {
			parentPath = strings.Replace(item["Path"].(string), "/"+item["Name"].(string), "", -1)
		}
		parent := parentDirs[parentPath]
		if !item["IsDir"].(bool) {
			if parent["depth"].(int) == 2 {
				seasonMetadata := getPaths(metadata, parent["json_path"].([]any))
				seasonMetadata.(map[string]interface{})["episodes"] = append(seasonMetadata.(map[string]interface{})["episodes"].([]map[string]interface{}), map[string]interface{}{
					"id":            item["ID"],
					"name":          item["Name"],
					"path":          item["Path"],
					"parent":        parent,
					"mime_type":     item["MimeType"],
					"size":          item["Size"],
					"modified_time": item["ModTime"],
				})
			}
		} else {
			parentDirs[item["Path"].(string)] = map[string]interface{}{
				"id":    item["ID"],
				"name":  item["Name"],
				"path":  item["Path"],
				"depth": parent["depth"].(int) + 1,
			}
			if parent["depth"].(int) == 0 {
				metadata = append(metadata, map[string]interface{}{
					"id":            item["ID"],
					"name":          item["Name"],
					"path":          item["Path"],
					"parent":        parent,
					"mime_type":     item["MimeType"],
					"modified_time": item["ModTime"],
					"seasons":       map[string]interface{}{},
					"json_path":     []any{len(metadata)},
				})
				parentDirs[item["Path"].(string)]["json_path"] = []any{len(metadata) - 1}
			} else if parent["depth"].(int) == 1 {
				seriesMetadata := getPaths(metadata, parent["json_path"].([]any))
				season := regexp.MustCompile(`(?<=Season.|season.|S|s)\d{1,3}|^\d{1,3}$`).FindString(item["Name"].(string))
				if season == "" {
					season = "1"
				}
				if season != "0" {
					season = strings.TrimLeft(season, "0")
				}
				seasonInt, _ := strconv.ParseInt(season, 10, 64)
				seriesMetadata.(map[string]interface{})["seasons"].(map[string]interface{})[season] = map[string]interface{}{
					"id":            item["ID"],
					"name":          item["Name"],
					"path":          item["Path"],
					"parent":        parent,
					"mime_type":     item["MimeType"],
					"modified_time": item["ModTime"],
					"episodes":      []map[string]interface{}{},
					"json_path":     append(parent["json_path"].([]any), seasonInt),
				}
				parentDirs[item["Path"].(string)]["json_path"] = append(parent["json_path"].([]any), "seasons", seasonInt)
			}
		}
	}
	return metadata
}

func (r *RCloneAPI) Size(path string) int {
	options := map[string]interface{}{
		"no-modtime":  true,
		"no-mimetype": true,
	}
	rcData := map[string]interface{}{
		"fs":     r.FS,
		"remote": path,
		"opt":    options,
	}
	result := post("http://localhost:"+strconv.FormatInt(int64(r.Port), 10)+"/"+r.RCLONE["getFileInfo"], rcData)
	return result.(map[string]interface{})["item"].(map[string]interface{})["Size"].(int)
}

func (r *RCloneAPI) Stream(path string) string {
	streamURL := fmt.Sprintf("http://localhost:%d/[%s]/%s", r.Port, r.FS, path)
	return streamURL
}

func (r *RCloneAPI) Thumbnail(id string) string {
	return ""
}
