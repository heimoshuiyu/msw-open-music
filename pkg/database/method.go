package database

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

func (database *Database) GetRandomFiles(limit int64) ([]File, error) {
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	rows, err := database.stmt.getRandomFiles.Query(limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	files := make([]File, 0)
	for rows.Next() {
		file := File{
			Db: database,
		}
		err = rows.Scan(&file.ID, &file.Folder_id, &file.Filename, &file.Foldername, &file.Filesize)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

func (database *Database) GetRandomFilesWithTag(tagID, limit int64) ([]File, error) {
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	rows, err := database.stmt.getRandomFilesWithTag.Query(tagID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	files := make([]File, 0)
	for rows.Next() {
		file := File{
			Db: database,
		}
		err = rows.Scan(&file.ID, &file.Folder_id, &file.Filename, &file.Foldername, &file.Filesize)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

func (database *Database) GetFilesInFolder(folder_id int64, limit int64, offset int64) ([]File, error) {
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	rows, err := database.stmt.getFilesInFolder.Query(folder_id, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	files := make([]File, 0)
	for rows.Next() {
		file := File{
			Db:        database,
			Folder_id: folder_id,
		}
		err = rows.Scan(&file.ID, &file.Filename, &file.Filesize, &file.Foldername)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

func (database *Database) SearchFolders(foldername string, limit int64, offset int64) ([]Folder, error) {
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	rows, err := database.stmt.searchFolders.Query("%"+foldername+"%", limit, offset)
	if err != nil {
		return nil, errors.New("Error searching folders at query " + err.Error())
	}
	defer rows.Close()
	folders := make([]Folder, 0)
	for rows.Next() {
		folder := Folder{
			Db: database,
		}
		err = rows.Scan(&folder.ID, &folder.Folder, &folder.Foldername)
		if err != nil {
			return nil, errors.New("Error scanning SearchFolders" + err.Error())
		}
		folders = append(folders, folder)
	}
	return folders, nil
}

func (database *Database) GetFile(id int64) (*File, error) {
	file := &File{
		Db: database,
	}

	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	err := database.stmt.getFile.QueryRow(id).Scan(&file.ID, &file.Folder_id, &file.Realname, &file.Filename, &file.Foldername, &file.Filesize)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (database *Database) ResetFiles() error {
	log.Println("[db] Reset files")
	var err error
	_, err = database.stmt.dropFiles.Exec()
	if err != nil {
		return err
	}
	_, err = database.stmt.initFilesTable.Exec()
	if err != nil {
		return err
	}
	return err
}

func (database *Database) ResetFolder() error {
	log.Println("[db] Reset folders")
	var err error
	_, err = database.stmt.dropFolder.Exec()
	if err != nil {
		return err
	}
	_, err = database.stmt.initFoldersTable.Exec()
	if err != nil {
		return err
	}
	return err
}

func (database *Database) Walk(root string, pattern []string, tagIDs []int64, userID int64) error {
	patternDict := make(map[string]bool)
	for _, v := range pattern {
		patternDict[v] = true
	}
	log.Println("[db] Walk", root, patternDict)

	tags := make([]*Tag, 0)
	for _, tagID := range tagIDs {
		tag, err := database.GetTag(tagID)
		if err != nil {
			return err
		}
		tags = append(tags, tag)
	}

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// check pattern
		ext := filepath.Ext(info.Name())
		if _, ok := patternDict[ext]; !ok {
			return nil
		}

		// insert file, folder will aut created
		fileID, err := database.Insert(path, info.Size())
		if err != nil {
			return err
		}

		for _, tag := range tags {
			err = database.PutTagOnFile(tag.ID, fileID, userID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (database *Database) GetFolder(folderId int64) (*Folder, error) {
	folder := &Folder{
		Db: database,
	}

	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	err := database.stmt.getFolder.QueryRow(folderId).Scan(&folder.Folder)
	if err != nil {
		return nil, err
	}
	return folder, nil
}

func (database *Database) SearchFiles(filename string, limit int64, offset int64) ([]File, error) {
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	rows, err := database.stmt.searchFiles.Query("%"+filename+"%", limit, offset)
	if err != nil {
		return nil, errors.New("Error searching files at query " + err.Error())
	}
	defer rows.Close()
	files := make([]File, 0)
	for rows.Next() {
		var file File = File{
			Db: database,
		}
		err = rows.Scan(&file.ID, &file.Folder_id, &file.Filename, &file.Foldername, &file.Filesize)
		if err != nil {
			return nil, errors.New("Error scanning SearchFiles " + err.Error())
		}
		files = append(files, file)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.New("Error scanning SearchFiles exit without full result" + err.Error())
	}
	return files, nil
}

func (database *Database) FindFolder(folder string) (int64, error) {
	var id int64

	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	err := database.stmt.findFolder.QueryRow(folder).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (database *Database) FindFile(folderId int64, filename string) (int64, error) {
	var id int64

	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	err := database.stmt.findFile.QueryRow(folderId, filename).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (database *Database) InsertFolder(folder string) (int64, error) {

	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	result, err := database.stmt.insertFolder.Exec(folder, filepath.Base(folder))
	if err != nil {
		return 0, err
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

func (database *Database) InsertFile(folderId int64, filename string, filesize int64) (int64, error) {
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	result, err := database.stmt.insertFile.Exec(folderId, filename, filename, filesize)
	if err != nil {
		return 0, err
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

func (database *Database) Insert(path string, filesize int64) (int64, error) {
	folder, filename := filepath.Split(path)
	folderId, err := database.FindFolder(folder)
	if err != nil {
		folderId, err = database.InsertFolder(folder)
		if err != nil {
			return 0, err
		}
	}

	// if file exists, skip it
	lastInsertId, err := database.FindFile(folderId, filename)
	if err == nil {
		return lastInsertId, nil
	}

	lastInsertId, err = database.InsertFile(folderId, filename, filesize)
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

func (database *Database) UpdateFoldername(folderId int64, foldername string) error {
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	_, err := database.stmt.updateFoldername.Exec(foldername, folderId)
	if err != nil {
		return err
	}
	return nil
}

func (database *Database) DeleteFile(fileId int64) error {
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	// begin transaction
	tx, err := database.sqlConn.Begin()
	if err != nil {
		return err
	}

	// delete file
	_, err = tx.Stmt(database.stmt.deleteFile).Exec(fileId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// delete tag on file
	_, err = tx.Stmt(database.stmt.deleteFileReferenceInFileHasTag).Exec(fileId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// delete reviews on file
	_, err = tx.Stmt(database.stmt.deleteFileReferenceInReviews).Exec(fileId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (database *Database) UpdateFilename(fileId int64, filename string) error {
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	_, err := database.stmt.updateFilename.Exec(filename, fileId)
	if err != nil {
		return err
	}
	return nil
}

func (database *Database) ResetFilename(fileId int64) error {
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	_, err := database.stmt.resetFilename.Exec(fileId)
	if err != nil {
		return err
	}
	return nil
}

func (database *Database) ResetFoldername(folderId int64) error {
	folder, err := database.GetFolder(folderId)
	if err != nil {
		return err
	}

	foldername := filepath.Base(folder.Folder)
	err = database.UpdateFoldername(folderId, foldername)
	if err != nil {
		return err
	}
	return nil
}
