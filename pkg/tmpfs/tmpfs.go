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
	recordLocks map[string]*sync.Mutex
}

func (tmpfs *Tmpfs) GetObjFilePath(id int64, configName string) (string) {
	return filepath.Join(tmpfs.Config.Root, strconv.FormatInt(id, 10) + "." + configName + ".webm")
}

func (tmpfs *Tmpfs) GetLock(filename string) *sync.Mutex {
	if _, ok := tmpfs.recordLocks[filename]; !ok {
		tmpfs.recordLocks[filename] = &sync.Mutex{}
	}
	return tmpfs.recordLocks[filename]
}

func (tmpfs *Tmpfs) Lock(filename string) {
	tmpfs.GetLock(filename).Lock()
}

func (tmpfs *Tmpfs) Unlock(filename string) {
	tmpfs.GetLock(filename).Unlock()
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
		recordLocks: make(map[string]*sync.Mutex),
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
		for path, lock := range tmpfs.recordLocks {
			lock.Lock()
			recordTime, ok := tmpfs.record[path]
			if !ok {
				lock.Unlock()
				continue
			}
			if now - recordTime > tmpfs.Config.FileLifeTime {
				err = os.Remove(path)
				if err != nil {
					log.Println("[tmpfs] Failed to remove file", err)
					lock.Unlock()
					continue
				}
				log.Println("[tmpfs] Deleted file", path)
				delete(tmpfs.record, path)
				delete(tmpfs.recordLocks, path)
			}
			lock.Unlock()
		}

		time.Sleep(time.Second)
	}
}
