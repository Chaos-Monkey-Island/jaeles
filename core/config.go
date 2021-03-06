package core

import (
	"bytes"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/jaeles-project/jaeles/libs"
	"github.com/jaeles-project/jaeles/utils"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path"
)

// InitConfig init config
func InitConfig(options *libs.Options) {
	options.RootFolder = utils.NormalizePath(options.RootFolder)
	options.Server.DBPath = path.Join(options.RootFolder, "sqlite3.db")
	// init new root folder
	if !utils.FolderExists(options.RootFolder) {
		utils.InforF("Init new config at %v", options.RootFolder)
		os.MkdirAll(options.RootFolder, 0750)
		// cloning default repo
		UpdatePlugins(*options)
		UpdateSignature(*options, "")
	}

	configPath := path.Join(options.RootFolder, "config.yaml")
	v := viper.New()
	v.AddConfigPath(options.RootFolder)
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	if !utils.FileExists(configPath) {
		utils.InforF("Write new config to: %v", configPath)
		// save default config if not exist
		bind := "http://127.0.0.1:5000"
		v.SetDefault("defaultSign", "*")
		v.SetDefault("cors", "*")
		// default credential
		v.SetDefault("username", "jaeles")
		v.SetDefault("password", utils.GenHash(utils.GetTS())[:10])
		v.SetDefault("secret", utils.GenHash(utils.GetTS()))
		v.SetDefault("bind", bind)
		v.WriteConfigAs(configPath)

	} else {
		if options.Debug {
			utils.InforF("Load config from: %v", configPath)
		}
		b, _ := ioutil.ReadFile(configPath)
		v.ReadConfig(bytes.NewBuffer(b))
	}
	// config.defaultSign = fmt.Sprintf("%v", v.Get("defaultSign"))

	// WARNING: change me if you really want to deploy on remote server
	// allow all origin
	options.Server.Cors = v.GetString("cors")
	options.Server.JWTSecret = v.GetString("secret")
	options.Server.Username = v.GetString("username")
	options.Server.Password = v.GetString("password")

	// store default credentials for Burp plugin
	burpConfigPath := path.Join(options.RootFolder, "burp.json")
	if !utils.FileExists(burpConfigPath) {
		jsonObj := gabs.New()
		jsonObj.Set("", "JWT")
		jsonObj.Set(v.GetString("username"), "username")
		jsonObj.Set(v.GetString("password"), "password")
		bind := v.GetString("bind")
		if bind == "" {
			bind = "http://127.0.0.1:5000"
		}
		jsonObj.Set(fmt.Sprintf("http://%v/api/parse", bind), "endpoint")
		utils.WriteToFile(burpConfigPath, jsonObj.String())
		if options.Verbose {
			utils.InforF("Store default credentials for client at: %v", burpConfigPath)
		}
	}
}
