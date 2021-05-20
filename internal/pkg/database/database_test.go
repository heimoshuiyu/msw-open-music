package database

import (
	"testing"
)

func TestDatabase(t *testing.T) {
	db, err := NewDatabase("/tmp/test.sqlite3")
	if err != nil {
		t.Fatal("Error creating database" + err.Error())
	}
	t.Log("database open successfully")

	_, err = db.InsertFolder("testfolder")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("insertFolders successfully")

	id, err := db.FindFolder("testfolder")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("folder found", id)

	err = db.Insert("/home/hmsy/go/bin/typora-image-ffmpeg")
	if err != nil {
		t.Fatal(err.Error())
	}

	files, err := db.SearchFiles("ffmpeg", 100, 0)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(files)

	file := files[0]
	t.Log(file.Path())

	err = db.Walk("/home/hmsy/dsa/")
	if err != nil {
		t.Fatal(err.Error())
	}

	//err = db.ResetFiles()
	//if err != nil {
	//	t.Fatal(err.Error())
	//}

	//err = db.ResetFolder()
	//if err != nil {
	//	t.Fatal(err.Error())
	//}

}
