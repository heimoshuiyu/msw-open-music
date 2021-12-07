package database

import (
	"database/sql"
)

var initFilesTableQuery = `CREATE TABLE IF NOT EXISTS files (
	id INTEGER PRIMARY KEY,
	folder_id INTEGER NOT NULL,
	filename TEXT NOT NULL,
	filesize INTEGER NOT NULL,
	FOREIGN KEY(folder_id) REFERENCES folders(id)
);`

var initFoldersTableQuery = `CREATE TABLE IF NOT EXISTS folders (
	id INTEGER PRIMARY KEY,
	folder TEXT NOT NULL,
	foldername TEXT NOT NULL
);`

var initFeedbacksTableQuery = `CREATE TABLE IF NOT EXISTS feedbacks (
	id INTEGER PRIMARY KEY,
	time INTEGER NOT NULL,
	feedback TEXT NOT NULL,
	header TEXT NOT NULL
);`

var initUsersTableQuery = `CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY,
	username TEXT NOT NULL,
	password TEXT NOT NULL
);`

var initAvatarsTableQuery = `CREATE TABLE IF NOT EXISTS avatars (
	id INTEGER PRIMARY KEY,
	avatarname TEXT NOT NULL,
	avatar BLOB NOT NULL
);`

var initTagsTableQuery = `CREATE TABLE IF NOT EXISTS tags (
	id INTEGER PRIMARY KEY,
	tag TEXT NOT NULL
);`

var initFileHasTagTableQuery = `CREATE TABLE IF NOT EXISTS file_has_tag (
	file_id INTEGER NOT NULL,
	tag_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	PRIMARY KEY (file_id, tag_id),
	FOREIGN KEY(user_id) REFERENCES users(id)
	FOREIGN KEY (file_id) REFERENCES files(id),
	FOREIGN KEY (tag_id) REFERENCES tags(id)
);`

var initLikesTableQuery = `CREATE TABLE IF NOT EXISTS likes (
	id INTEGER PRIMARY KEY,
	user_id INTEGER NOT NULL,
	file_id INTEGER NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (file_id) REFERENCES files(id)
);`

var initReviewsTableQuery = `CREATE TABLE IF NOT EXISTS reviews (
	id INTEGER PRIMARY KEY,
	user_id INTEGER NOT NULL,
	time INTEGER NOT NULL,
	modified_time INTEGER DEFAULT 0,
	review TEXT NOT NULL,
	header TEXT NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id)
);`

var initPlaybacksTableQuery = `CREATE TABLE IF NOT EXISTS playbacks (
	id INTEGER PRIMARY KEY,
	user_id INTEGER NOT NULL,
	file_id INTEGER NOT NULL,
	time INTEGER NOT NULL,
	mothod INTEGER NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (file_id) REFERENCES files(id)
);`

var initLogsTableQuery = `CREATE TABLE IF NOT EXISTS logs (
	id INTEGER PRIMARY KEY,
	time INTEGER NOT NULL,
	message TEXT NOT NULL,
	user_id INTEGER NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id)
);`

var insertFolderQuery = `INSERT INTO folders (folder, foldername)
VALUES (?, ?);`

var findFolderQuery = `SELECT id FROM folders WHERE folder = ? LIMIT 1;`

var insertFileQuery = `INSERT INTO files (folder_id, filename, filesize)
VALUES (?, ?, ?);`

var searchFilesQuery = `SELECT
files.id, files.folder_id, files.filename, folders.foldername, files.filesize
FROM files
JOIN folders ON files.folder_id = folders.id
WHERE filename LIKE ?
LIMIT ? OFFSET ?;`

var getFolderQuery = `SELECT folder FROM folders WHERE id = ? LIMIT 1;`

var dropFilesQuery = `DROP TABLE files;`

var dropFolderQuery = `DROP TABLE folders;`

var getFileQuery = `SELECT
files.id, files.folder_id, files.filename, folders.foldername, files.filesize
FROM files
JOIN folders ON files.folder_id = folders.id
WHERE files.id = ?
LIMIT 1;`

var searchFoldersQuery = `SELECT
id, folder, foldername
FROM folders
WHERE foldername LIKE ?
LIMIT ? OFFSET ?;`

var getFilesInFolderQuery = `SELECT
files.id, files.filename, files.filesize, folders.foldername
FROM files
JOIN folders ON files.folder_id = folders.id
WHERE folder_id = ?
LIMIT ? OFFSET ?;`

var getRandomFilesQuery = `SELECT
files.id, files.folder_id, files.filename, folders.foldername, files.filesize
FROM files
JOIN folders ON files.folder_id = folders.id
ORDER BY RANDOM()
LIMIT ?;`

var insertFeedbackQuery = `INSERT INTO feedbacks (time, feedback, header)
VALUES (?, ?, ?);`

type Stmt struct {
	initFilesTable     *sql.Stmt
	initFoldersTable   *sql.Stmt
	initFeedbacksTable *sql.Stmt
	initUsersTable     *sql.Stmt
	initAvatarsTable   *sql.Stmt
	initTagsTable      *sql.Stmt
	initFileHasTag     *sql.Stmt
	initLikesTable     *sql.Stmt
	initReviewsTable   *sql.Stmt
	initPlaybacksTable *sql.Stmt
	initLogsTable      *sql.Stmt
	insertFolder       *sql.Stmt
	insertFile         *sql.Stmt
	findFolder         *sql.Stmt
	searchFiles        *sql.Stmt
	getFolder          *sql.Stmt
	dropFiles          *sql.Stmt
	dropFolder         *sql.Stmt
	getFile            *sql.Stmt
	searchFolders      *sql.Stmt
	getFilesInFolder   *sql.Stmt
	getRandomFiles     *sql.Stmt
	insertFeedback     *sql.Stmt
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

	// init feedbacks tables
	stmt.initFeedbacksTable, err = sqlConn.Prepare(initFeedbacksTableQuery)
	if err != nil {
		return nil, err
	}

	// init users table
	stmt.initUsersTable, err = sqlConn.Prepare(initUsersTableQuery)
	if err != nil {
		return nil, err
	}

	// init avatars table
	stmt.initAvatarsTable, err = sqlConn.Prepare(initAvatarsTableQuery)
	if err != nil {
		return nil, err
	}

	// init tags table
	stmt.initTagsTable, err = sqlConn.Prepare(initTagsTableQuery)
	if err != nil {
		return nil, err
	}

	// init file_has_tag table
	stmt.initFileHasTag, err = sqlConn.Prepare(initFileHasTagTableQuery)
	if err != nil {
		return nil, err
	}

	// init likes table
	stmt.initLikesTable, err = sqlConn.Prepare(initLikesTableQuery)
	if err != nil {
		return nil, err
	}

	// init reviews table
	stmt.initReviewsTable, err = sqlConn.Prepare(initReviewsTableQuery)
	if err != nil {
		return nil, err
	}

	// init playbacks table
	stmt.initPlaybacksTable, err = sqlConn.Prepare(initPlaybacksTableQuery)
	if err != nil {
		return nil, err
	}

	// init logs table
	stmt.initLogsTable, err = sqlConn.Prepare(initLogsTableQuery)
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
	_, err = stmt.initFeedbacksTable.Exec()
	if err != nil {
		return nil, err
	}
	_, err = stmt.initUsersTable.Exec()
	if err != nil {
		return nil, err
	}
	_, err = stmt.initAvatarsTable.Exec()
	if err != nil {
		return nil, err
	}
	_, err = stmt.initTagsTable.Exec()
	if err != nil {
		return nil, err
	}
	_, err = stmt.initFileHasTag.Exec()
	if err != nil {
		return nil, err
	}
	_, err = stmt.initLikesTable.Exec()
	if err != nil {
		return nil, err
	}
	_, err = stmt.initReviewsTable.Exec()
	if err != nil {
		return nil, err
	}
	_, err = stmt.initPlaybacksTable.Exec()
	if err != nil {
		return nil, err
	}
	_, err = stmt.initLogsTable.Exec()
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

	// init getRandomFiles
	stmt.getRandomFiles, err = sqlConn.Prepare(getRandomFilesQuery)
	if err != nil {
		return nil, err
	}

	// init insertFeedback
	stmt.insertFeedback, err = sqlConn.Prepare(insertFeedbackQuery)
	if err != nil {
		return nil, err
	}

	return stmt, err
}
