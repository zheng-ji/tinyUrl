/*
 * Author: zhengji
 */

package global

import (
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
	"time"
)

// Config Struct
type Config struct {
	Bind      string              `yaml:"bind"`
	MonitAddr string              `yaml:"monitaddr"`
	LogCfg    LogConf             `yaml:"logcfg"`
	Redis     RedisConf           `yaml:"redis"`
	MySQL     MySQLConf           `yaml:"mysql"`
    HostName  string              `yaml:"hostname"`
}

// GlobalConfig var,the config of this program
var GlobalConfig Config

// LogConf Struct
type LogConf struct {
	Level  string `yaml:"level"`
	RunLog string `yaml:"runlog"`
}

// RedisConf Struct
type RedisConf struct {
	Addr string `yaml:"addr"`
	DB   int    `yaml:"db"`
}

// MySQLConf Struct
type MySQLConf struct {
	Name               string `yaml:"name"`
	Connection         string `yaml:"connection"`
	MaxIdleConnections int    `yaml:"max_idle_connections"`
}


func ParseConfigFile(filename string) error {
	if config, err := ioutil.ReadFile(filename); err == nil {
		if err = goyaml.Unmarshal(config, &GlobalConfig); err != nil {
			return err
		}
		fmt.Println(GlobalConfig)
	} else {
		return err
	}
	return nil
}
