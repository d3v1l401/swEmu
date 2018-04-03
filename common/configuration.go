package common

import (
	"encoding/json"
	"io/ioutil"
)

type Servers struct {
	ID          int    `json:"ID"`
	Flag        int    `json:"Flag"`
	Name        string `json:"Name"`
	IP          string `json:"IP"`
	Port        int    `json:"Port"`
	CanJoin     bool   `json:"CanJoin"`
	PlayerLimit int    `json:"PlayerLimit"`
}

type Configuration struct {
	ServerList []Servers `json:"Servers"`
	AuthServer struct {
		Enabled        bool     `json:"Enabled"`
		IP             string   `json:"IP"`
		Port           int      `json:"Port"`
		AllowAuthProto []string `json:"AllowAuthProto"`
		AdminOnly      bool     `json:"AdminOnly"`
	} `json:"AuthServer"`
	AdminServer struct {
		Enabled        bool     `json:"Enabled"`
		IP             string   `json:"IP"`
		Port           int      `json:"Port"`
		AllowAuthProto []string `json:"AllowAuthProto"`
	} `json:"AdminServer"`
	Database struct {
		UseSignature bool   `json:"UseSignature"`
		IP           string `json:"IP"`
		Port         int    `json:"Port"`
		Username     string `json:"Username"`
		Password     string `json:"Password"`
		Signature    string `json:"Signature"`
	} `json:"Database"`
}

func ImportConfiguration(path string) *Configuration {
	var conf *Configuration

	if len(path) > 0 {
		if data, err := ioutil.ReadFile(path); err == nil {
			if err = json.Unmarshal(data, &conf); err == nil {
				return conf
			}
		}
	}

	return nil
}
