package database

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var initFilesTableQuery = `CREATE TABLE IF NOT EXISTS files (
	id INTEGER PRIMARY KEY,
	folder_id INTEGER NOT NULL,
	filename TEXT NOT NULL,
	filesize INTEGER NOT NULL
);`
var initFoldersTableQuery = `CREATE TABLE IF NOT EXISTS folders (
	id INTEGER PRIMARY KEY,
	folder TEXT NOT NULL,
	foldername TEXT NOT NULL
);`
var insertFolderQuery = `INSERT INTO folders (folder, foldername) VALUES (?, ?);`
var findFolderQuery = `SELECT id FROM folders WHERE folder = ? LIMIT 1;`
var insertFileQuery = `INSERT INTO files (folder_id, filename, filesize) VALUES (?, ?, ?);`
var searchFilesQuery = `SELECT files.id, files.folder_id, files.filename, folders.foldername, files.filesize FROM files JOIN folders ON files.folder_id = folders.id WHERE filename LIKE ? LIMIT ? OFFSET ?;`
var getFolderQuery = `SELECT folder FROM folders WHERE id = ? LIMIT 1;`
var dropFilesQuery = `DROP TABLE files;`
var dropFolderQuery = `DROP TABLE folders;`
var getFileQuery = `SELECT files.id, files.folder_id, files.filename, folders.foldername, files.filesize FROM files JOIN folders ON files.folder_id = folders.id WHERE files.id = ? LIMIT 1;`
var searchFoldersQuery = `SELECT id, folder, foldername FROM folders WHERE foldername LIKE ? LIMIT ? OFFSET ?;`
var getFilesInFolderQuery = `SELECT id, filename, filesize FROM files WHERE folder_id = ? LIMIT ? OFFSET ?;`

type Database struct {
	sqlConn *sql.DB
	stmt *Stmt
}

type Stmt struct {
	initFilesTable *sql.Stmt
	initFoldersTable *sql.Stmt
	insertFolder *sql.Stmt
	insertFile *sql.Stmt
	findFolder *sql.Stmt
	searchFiles *sql.Stmt
	getFolder *sql.Stmt
	dropFiles *sql.Stmt
	dropFolder *sql.Stmt
	getFile *sql.Stmt
	searchFolders *sql.Stmt
	getFilesInFolder *sql.Stmt
}

type File struct {
	Db *Database `json:"-"`
	ID int64 `json:"id"`
	Folder_id int64 `json:"folder_id"`
	Foldername string `json:"foldername"`
	Filename string `json:"filename"`
	Filesize int64 `json:"filesize"`
}

type Folder struct {
	Db *Database `json:"-"`
	ID int64 `json:"id"`
	Folder string `json:"folder"`
	Foldername string `json:"foldername"`
}

func (database *Database) GetFilesInFolder(folder_id int64, limit int64, offset int64) ([]File, error) {
	rows, err := database.stmt.getFilesInFolder.Query(folder_id, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	files := make([]File, 0)
	for rows.Next() {
		file := File{
			Db: database,
			Folder_id: folder_id,
		}
		err = rows.Scan(&file.ID, &file.Filename, &file.Filesize)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

func (database *Database) SearchFolders(foldername string, limit int64, offset int64) ([]Folder, error) {
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
	err := database.stmt.getFile.QueryRow(id).Scan(&file.ID, &file.Folder_id, &file.Filename, &file.Foldername, &file.Filesize)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (database *Database) ResetFiles() (error) {
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

func (database *Database) ResetFolder() (error) {
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

func (database *Database) Walk(root string, pattern []string) (error) {
	patternDict := make(map[string]bool)
	for _, v := range pattern {
		patternDict[v] = true
	}
	log.Println(patternDict)
	return filepath.Walk(root, func (path string, info os.FileInfo, err error) (error) {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// check pattern
		ext := filepath.Ext(info.Name())
		if _, ok := patternDict[ext]; !ok {
			log.Println("False", ext, info.Name())
			return nil
		}
		log.Println("True", ext, info.Name())

		// insert file, folder will aut created
		err = database.Insert(path, info.Size())
		if err != nil {
			return err
		}
		return nil
	})
}

func (f *File) Path() (string, error) {
	folder, err := f.Db.GetFolder(f.Folder_id)
	if err != nil {
		return "", err
	}
	return filepath.Join(folder.Folder, f.Filename), nil
}

func (database *Database) GetFolder(folderId int64) (*Folder, error) {
	folder := &Folder{
		Db: database,
	}
	err := database.stmt.getFolder.QueryRow(folderId).Scan(&folder.Folder)
	if err != nil {
		return nil, err
	}
	return folder, nil
}

func (database *Database) SearchFiles(filename string, limit int64, offset int64) ([]File, error) {
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
	err := database.stmt.findFolder.QueryRow(folder).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (database *Database) InsertFolder(folder string) (int64, error) {
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

func (database *Database) InsertFile(folderId int64, filename string, filesize int64) (error) {
	_, err := database.stmt.insertFile.Exec(folderId, filename, filesize)
	if err != nil {
		return err
	}
	return nil
}

func (database *Database) Insert(path string, filesize int64) (error) {
	folder, filename := filepath.Split(path)
	folderId, err := database.FindFolder(folder)
	if err != nil {
		folderId, err = database.InsertFolder(folder)
		if err != nil {
			return err
		}
	}
	err = database.InsertFile(folderId, filename, filesize)
	if err != nil {
		return err
	}
	return nil
}

func NewPreparedStatement(sqlConn *sql.DB) (*Stmt, error) {
	var err error

	stmt := &Stmt{}

	// init files table
	stmt.initFilesTable, err = sqlConn.Prepare(initFilesTableQuery)
	if err != nil {
		return nil, err
	}

	// init folders table
	stmt.initFoldersTable, err = sqlConn.Prepare(initFoldersTableQuery)
	if err != nil {
		return nil, err
	}

	// run init statement
	_, err = stmt.initFilesTable.Exec()
	if err != nil {
		return nil, err
	}
	_, err = stmt.initFoldersTable.Exec()
	if err != nil {
		return nil, err
	}

	// init insert folder statement
	stmt.insertFolder, err = sqlConn.Prepare(insertFolderQuery)
	if err != nil {
		return nil, err
	}

	// init findFolder statement
	stmt.findFolder, err = sqlConn.Prepare(findFolderQuery)
	if err != nil {
		return nil, err
	}

	// init insertFile stmt
	stmt.insertFile, err = sqlConn.Prepare(insertFileQuery)
	if err != nil {
		return nil, err
	}

	// init searchFile stmt
	stmt.searchFiles, err = sqlConn.Prepare(searchFilesQuery)
	if err != nil {
		return nil, err
	}

	// init getFolder stmt
	stmt.getFolder, err = sqlConn.Prepare(getFolderQuery)
	if err != nil {
		return nil, err
	}

	// init dropFolder stmt
	stmt.dropFolder, err = sqlConn.Prepare(dropFolderQuery)
	if err != nil {
		return nil, err
	}

	// init dropFiles stmt
	stmt.dropFiles, err = sqlConn.Prepare(dropFilesQuery)
	if err != nil {
		return nil, err
	}

	// init getFile stmt
	stmt.getFile, err = sqlConn.Prepare(getFileQuery)
	if err != nil {
		return nil, err
	}

	// init searchFolder stmt
	stmt.searchFolders, err = sqlConn.Prepare(searchFoldersQuery)
	if err != nil {
		return nil, err
	}

	// init getFilesInFolder stmt
	stmt.getFilesInFolder, err = sqlConn.Prepare(getFilesInFolderQuery)
	if err != nil {
		return nil, err
	}

	return stmt, err
}

func NewDatabase(dbName string) (*Database, error) {
	var err error

	// open database
	sqlConn, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	// prepare statement
	stmt, err := NewPreparedStatement(sqlConn)
	if err != nil {
		return nil, err
	}

	// new database
	database := &Database{
		sqlConn: sqlConn,
		stmt: stmt,
	}

	return database, nil
}
