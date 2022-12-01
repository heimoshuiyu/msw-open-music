package api

import (
	"bytes"
	"errors"
	"io"
	"log"
	"msw-open-music/pkg/database"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

func (api *API) HandelGetFileAvatar(w http.ResponseWriter, r *http.Request) {
	var err error
	q := r.URL.Query()
	ids := q["id"]
	if len(ids) == 0 {
		err = errors.New(`parameter "id" can't be empty`)
		api.HandleError(w, r, err)
		return
	}
	id, err := strconv.Atoi(ids[0])
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	file, err := api.Db.GetFile(int64(id))
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	path, err := file.Path()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	log.Println("[api] Get avatar of file", path)
	buff := make([]byte, 0)
	cache := bytes.NewBuffer(buff)
	cmd := exec.Command("ffmpeg", "-i", path, "-c:v", "libwebp_anim", "-update", "1", "-f", "image2pipe", "-")
	cmd.Stdout = cache

	err = cmd.Run()
	if err != nil {
		api.HandleGetAlternativeFileAvatar(w, r, file)
		return
	}

	w.Header().Set("Content-Type", "image/webp")
	io.Copy(w, cache)
}

func (api *API) HandleGetAlternativeFileAvatar(w http.ResponseWriter, r *http.Request, f *database.File) {
	var err error
	dir, err := f.Dir()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	log.Println("[api] Get alternative avatar in dir", dir)
	files, err := os.ReadDir(dir)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	avatar, err := findAvatarFile(files)
	avatarPath := path.Join(dir, avatar)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	cmd := exec.Command("ffmpeg", "-i", avatarPath, "-c:v", "libwebp_anim", "-f", "image2pipe", "-")
	cmd.Stdout = w
	w.Header().Set("Content-Type", "image/webp")

	err = cmd.Run()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
}

func findAvatarFile(files []os.DirEntry) (string, error) {
	for _, file := range files {
		if isAvatarType(file.Name()) {
			return file.Name(), nil
		}
	}
	return "", errors.New("Cannot find avatar file")
}

var avatarFileTypes = []string{
	".jpg",
	".png",
}

func isAvatarType(filename string) bool {
	for _, t := range avatarFileTypes {
		if strings.HasSuffix(strings.ToLower(filename), t) {
			return true
		}
	}
	return false
}
