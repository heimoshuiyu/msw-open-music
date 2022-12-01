package api

import (
	"errors"
	"log"
	"net/http"
	"os/exec"
	"strconv"
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
	cmd := exec.Command("ffmpeg", "-i", path, "-c:v", "libwebp_anim", "-update", "1", "-f", "image2pipe", "-")
	cmd.Stdout = w
	w.Header().Set("Content-Type", "image/webp")

	err = cmd.Run()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		api.HandleError(w, r, err)
		return
	}
}
