package rclone

import (
	"fmt"
	"unicode"
)

func buildConfig(config map[string]interface{}) []string {
	rcloneConf := []string{}
	for _, category := range config["categories"].([]interface{}) {
		provider := category.(map[string]interface{})["provider"]
		if provider == "gdrive" {
			provider = "drive"
			clientID := config["gdrive"].(map[string]interface{})["client_id"]
			clientSecret := config["gdrive"].(map[string]interface{})["client_secret"]
			token := map[string]interface{}{
				"access_token":  config["gdrive"].(map[string]interface{})["access_token"],
				"token_type":    "Bearer",
				"refresh_token": config["gdrive"].(map[string]interface{})["refresh_token"],
				"expiry":        "2022-03-27T00:00:00.000+00:00",
			}
			id := category.(map[string]interface{})["id"]
			safeFs := ""
			for _, c := range id.(string) {
				if unicode.IsLetter(c) || unicode.IsNumber(c) {
					safeFs += string(c)
				}
			}
			driveID := category.(map[string]interface{})["drive_id"]
			rcloneConf = append(rcloneConf, fmt.Sprintf("[%s]\ntype = drive\nclient_id = %s\nclient_secret = %s\nscope = drive\nroot_folder_id = %s\ntoken = %s\nteam_drive = %s\n", safeFs, clientID, clientSecret, id, token, driveID))
		} else if provider == "onedrive" {
			token := map[string]interface{}{
				"access_token":  config["onedrive"].(map[string]interface{})["access_token"],
				"token_type":    "Bearer",
				"refresh_token": config["onedrive"].(map[string]interface{})["refresh_token"],
				"expiry":        "2022-03-27T00:00:00.000+00:00",
			}
			id := category.(map[string]interface{})["id"]
			safeFs := ""
			for _, c := range id.(string) {
				if unicode.IsLetter(c) || unicode.IsNumber(c) {
					safeFs += string(c)
				}
			}
			driveID := category.(map[string]interface{})["drive_id"]
			rcloneConf = append(rcloneConf, fmt.Sprintf("[%s]\ntype = onedrive\nscope = drive\nroot_folder_id = %s\ntoken = %s\ndrive_id = %s\ndrive_type = personal", safeFs, id, token, driveID))
		} else if provider == "sharepoint" {
			token := map[string]interface{}{
				"access_token":  config["sharepoint"].(map[string]interface{})["access_token"],
				"token_type":    "Bearer",
				"refresh_token": config["sharepoint"].(map[string]interface{})["refresh_token"],
				"expiry":        "2022-03-27T00:00:00.000+00:00",
			}
			id := category.(map[string]interface{})["id"]
			driveID := category.(map[string]interface{})["drive_id"]
			if id != nil && driveID != nil {
				safeFs := ""
				for _, c := range id.(string) {
					if unicode.IsLetter(c) || unicode.IsNumber(c) {
						safeFs += string(c)
					}
				}
				rcloneConf = append(rcloneConf, fmt.Sprintf("[%s]\ntype = onedrive\nroot_folder_id = %s\ntoken = %s\ndrive_id = %s\ndrive_type = documentLibrary", safeFs, id, token, driveID))
			} else if driveID != nil {
				safeFs := ""
				for _, c := range driveID.(string) {
					if unicode.IsLetter(c) || unicode.IsNumber(c) {
						safeFs += string(c)
					}
				}
				rcloneConf = append(rcloneConf, fmt.Sprintf("[%s]\ntype = onedrive\ntoken = %s\ndrive_id = %s\ndrive_type = documentLibrary", safeFs, token, driveID))
			}
		} else if provider == "local" {
			fsPath := category.(map[string]interface{})["id"]
			if fsPath != nil {
				safeFs := ""
				for _, c := range fsPath.(string) {
					if unicode.IsLetter(c) || unicode.IsNumber(c) {
						safeFs += string(c)
					}
				}
				rcloneConf = append(rcloneConf, fmt.Sprintf("[%s]\ntype = alias\nremote = %s", safeFs, fsPath))
			}
		}
	}
	return rcloneConf
}
