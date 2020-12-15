package config

/*

import (
	"log"

	"github.com/meteocima/vfs/fs"
	"github.com/spf13/viper"
)

type Config struct {
	Filesystems map[string]*fs.Filesystem
	MountPoints map[string]fs.MountPoint
}

var Cfg Config

func (cfg *Config) Init() {

	cfg.Filesystems = map[string]*fs.Filesystem{}
	fss := viper.GetStringMap("filesystem")

	var getPropS = func(props map[string]interface{}, name string) string {
		if val, ok := props[name]; ok {
			return val.(string)
		}
		return ""
	}

	var getPropI = func(props map[string]interface{}, name string) int64 {
		if val, ok := props[name]; ok {
			return val.(int64)
		}
		return 0
	}

	var getPropStrArr = func(props map[string]interface{}, name string) []string {
		if val, ok := props[name]; ok {
			arr := val.([]interface{})
			res := make([]string, len(arr))
			for idx, item := range arr {
				res[idx] = item.(string)
			}
			return res
		}
		return nil
	}

	for name, fsInst := range fss {
		fsMap := fsInst.(map[string]interface{})
		cfg.Filesystems[name] = &fs.Filesystem{
			Type:        fs.FsType(getPropI(fsMap, "type")),
			Name:        name,
			BackupHosts: getPropStrArr(fsMap, "backup-hosts"),
			Host:        getPropS(fsMap, "host"),
			Port:        getPropI(fsMap, "port"),
			User:        getPropS(fsMap, "user"),
			Password:    getPropS(fsMap, "password"),
			Key:         getPropS(fsMap, "key"),
		}
	}

	mountPoints := viper.GetStringMap("mountpoint")
	cfg.MountPoints = map[string]fs.MountPoint{}
	for name, mountpoint := range mountPoints {
		fsMap := mountpoint.(map[string]interface{})
		fsInst, ok := cfg.Filesystems[getPropS(fsMap, "fs")]
		if !ok {
			log.Fatalf("cannot find file system `%s`", getPropS(fsMap, "fs"))
		}
		cfg.MountPoints[name] = fs.MountPoint{
			Filesystem: fsInst,
			Name:       name,
			Root:       getPropS(fsMap, "root"),
		}
	}

}
*/
