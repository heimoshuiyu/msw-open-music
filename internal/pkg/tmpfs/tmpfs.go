package tmpfs

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type Tmpfs struct {
	record map[string]int64
	Config TmpfsConfig
	wg sync.WaitGroup
}

func (tmpfs *Tmpfs) GetObjFilePath(id int64, configName string) (string) {
	return filepath.Join(tmpfs.Config.Root, strconv.FormatInt(id, 10) + "." + configName + ".ogg")
}

type TmpfsConfig struct {
	FileLifeTime int64 `json:"file_life_time"`
	CleanerInternal int64 `json:"cleaner_internal"`
	Root string `json:"root"`
}

func NewTmpfsConfig() (*TmpfsConfig) {
	config := &TmpfsConfig{}
	return config
}

func NewTmpfs(config TmpfsConfig) *Tmpfs {
	tmpfs := &Tmpfs{
		record: make(map[string]int64),
		Config: config,
	}
	tmpfs.wg.Add(1)
	go tmpfs.Cleaner()
	return tmpfs
}

func (tmpfs *Tmpfs) Record(filename string) {
	tmpfs.record[filename] = time.Now().Unix()
}

func (tmpfs *Tmpfs) Exits(filename string) (bool) {
	_, ok := tmpfs.record[filename]
	return ok
}

func (tmpfs *Tmpfs) Cleaner() {
	var err error
	for {
		now := time.Now().Unix()
		for key, value := range tmpfs.record {
			if now - value > tmpfs.Config.FileLifeTime {
				err = os.Remove(key)
				if err != nil {
					log.Println("[tmpfs] Failed to remove file", err)
				}
				log.Println("[tmpfs] Deleted file", key)
				delete(tmpfs.record, key)
			}
		}

		time.Sleep(time.Second)
	}
}
