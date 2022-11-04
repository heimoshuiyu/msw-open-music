package database

import (
	"database/sql"
	"log"
)

var initFilesTableQuery = `CREATE TABLE IF NOT EXISTS files (
	id SERIAL PRIMARY KEY,
	folder_id INTEGER NOT NULL REFERENCES folders(id),
	realname TEXT NOT NULL,
	filename TEXT NOT NULL,
	filesize INTEGER NOT NULL
);`

var initFoldersTableQuery = `CREATE TABLE IF NOT EXISTS folders (
	id SERIAL PRIMARY KEY,
	folder TEXT NOT NULL,
	foldername TEXT NOT NULL
);`

var initFeedbacksTableQuery = `CREATE TABLE IF NOT EXISTS feedbacks (
	id SERIAL PRIMARY KEY,
	time INTEGER NOT NULL,
	content TEXT NOT NULL,
	user_id INTEGER NOT NULL,
	header TEXT NOT NULL
);`

// User table schema definition
// role: 0 - Anonymous User, 1 - Admin, 2 - User
// postgres avatar references problem
var initUsersTableQuery = `CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	username TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	role INTEGER NOT NULL,
	active BOOLEAN NOT NULL,
	avatar_id INTEGER NOT NULL DEFAULT 0
);`

var initAvatarsTableQuery = `CREATE TABLE IF NOT EXISTS avatars (
	id SERIAL PRIMARY KEY,
	avatarname TEXT NOT NULL,
	avatar BYTEA NOT NULL
);`

var initTagsTableQuery = `CREATE TABLE IF NOT EXISTS tags (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	description TEXT NOT NULL,
	created_by_user_id INTEGER NOT NULL REFERENCES users(id)
);`

var initFileHasTagTableQuery = `CREATE TABLE IF NOT EXISTS file_has_tag (
	file_id INTEGER NOT NULL REFERENCES files(id),
	tag_id INTEGER NOT NULL REFERENCES tags(id),
	user_id INTEGER NOT NULL REFERENCES users(id),
	PRIMARY KEY (file_id, tag_id)
);`

var initLikesTableQuery = `CREATE TABLE IF NOT EXISTS likes (
	user_id INTEGER NOT NULL REFERENCES users(id),
	file_id INTEGER NOT NULL REFERENCES files(id),
	PRIMARY KEY (user_id, file_id)
);`

var initReviewsTableQuery = `CREATE TABLE IF NOT EXISTS reviews (
	id SERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL REFERENCES users(id),
	file_id INTEGER NOT NULL REFERENCES files(id),
	created_at INTEGER NOT NULL,
	updated_at INTEGER NOT NULL DEFAULT 0,
	content TEXT NOT NULL
);`

var initPlaybacksTableQuery = `CREATE TABLE IF NOT EXISTS playbacks (
	id SERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL REFERENCES users(id),
	file_id INTEGER NOT NULL REFERENCES files(id),
	time TIMESTAMP NOT NULL,
	method INTEGER NOT NULL,
	duration INTERVAL NOT NULL
);`

var initLogsTableQuery = `CREATE TABLE IF NOT EXISTS logs (
	id SERIAL PRIMARY KEY,
	time INTEGER NOT NULL,
	message TEXT NOT NULL,
	user_id INTEGER NOT NULL REFERENCES users(id)
);`

var initTmpfsTableQuery = `CREATE TABLE IF NOT EXISTS tmpfs (
	id SERIAL PRIMARY KEY,
	path TEXT NOT NULL,
	size INTEGER NOT NULL,
	file_id INTEGER NOT NULL REFERENCES files(id),
	ffmpeg_config TEXT NOT NULL,
	created_time INTEGER NOT NULL,
	accessed_time INTEGER NOT NULL
);`

var insertFolderQuery = `INSERT INTO folders (folder, foldername)
VALUES ($1, $2)
RETURNING id;
;`

var findFolderQuery = `SELECT id FROM folders WHERE folder = $1 LIMIT 1;`

var findFileQuery = `SELECT id FROM files WHERE folder_id = $1 AND realname = $2 LIMIT 1;`

var insertFileQuery = `INSERT INTO files (folder_id, realname, filename, filesize)
VALUES ($1, $2, $3, $4)
RETURNING id;`

var searchFilesQuery = `SELECT
files.id, files.folder_id, files.filename, folders.foldername, files.filesize
FROM files
JOIN folders ON files.folder_id = folders.id
WHERE filename LIKE $1
ORDER BY folders.foldername, files.filename
LIMIT $2 OFFSET $3;`

var getFolderQuery = `SELECT folder FROM folders WHERE id = $1 LIMIT 1;`

var dropFilesQuery = `DROP TABLE files;`

var dropFolderQuery = `DROP TABLE folders;`

var getFileQuery = `SELECT
files.id, files.folder_id, files.realname, files.filename, folders.foldername, files.filesize
FROM files
JOIN folders ON files.folder_id = folders.id
WHERE files.id = $1
LIMIT 1;`

var searchFoldersQuery = `SELECT
id, folder, foldername
FROM folders
WHERE foldername LIKE $1
ORDER BY foldername
LIMIT $2 OFFSET $3;`

var getFilesInFolderQuery = `SELECT
files.id, files.filename, files.filesize, folders.foldername, folders.folder
FROM files
JOIN folders ON files.folder_id = folders.id
WHERE folder_id = $1
ORDER BY files.filename
LIMIT $2 OFFSET $3;`

var getRandomFilesQuery = `SELECT
files.id, files.folder_id, files.filename, folders.foldername, files.filesize
FROM files
JOIN folders ON files.folder_id = folders.id
ORDER BY RANDOM()
LIMIT $1;`

var getRandomFilesWithTagQuery = `SELECT
files.id, files.folder_id, files.filename, folders.foldername, files.filesize
FROM file_has_tag
JOIN files ON file_has_tag.file_id = files.id
JOIN folders ON files.folder_id = folders.id
WHERE file_has_tag.tag_id = $1
ORDER BY RANDOM()
LIMIT $2;`

var insertFeedbackQuery = `INSERT INTO feedbacks (time, content, user_id, header)
VALUES ($1, $2, $3, $4);`

var getFeedbacksQuery = `SELECT
feedbacks.id, feedbacks.time, feedbacks.content, feedbacks.header,
users.id, users.username, users.role, users.active, users.avatar_id
FROM feedbacks
JOIN users ON feedbacks.user_id = users.id
ORDER BY feedbacks.time
;`

var deleteFeedbackQuery = `DELETE FROM feedbacks WHERE id = $1;`

var insertUserQuery = `INSERT INTO users (username, password, role, active, avatar_id)
VALUES ($1, $2, $3, $4, $5);`

var countUserQuery = `SELECT count(*) FROM users;`

var countAdminQuery = `SELECT count(*) FROM users WHERE role= 1;`

var getUserQuery = `SELECT id, username, password, role, active, avatar_id FROM users WHERE username = $1 LIMIT 1;`

var getUsersQuery = `SELECT id, username, role, active, avatar_id FROM users;`

var getUserByIdQuery = `SELECT id, username, role, active, avatar_id FROM users WHERE id = $1 LIMIT 1;`

var updateUserActiveQuery = `UPDATE users SET active = $1 WHERE id = $2;`

var updateUsernameQuery = `UPDATE users SET username = $1 WHERE id = $2;`

var updateUserPasswordQuery = `UPDATE users SET password = $1 WHERE id = $2;`

var getAnonymousUserQuery = `SELECT id, username, role, avatar_id FROM users WHERE role = 0 LIMIT 1;`

var insertTagQuery = `INSERT INTO tags (name, description, created_by_user_id) VALUES ($1, $2, $3) RETURNING id;`

var deleteTagQuery = `DELETE FROM tags WHERE id = $1;`

var getTagQuery = `SELECT
tags.id, tags.name, tags.description,
users.id, users.username, users.role, users.avatar_id
FROM tags
JOIN users ON tags.created_by_user_id = users.id
WHERE tags.id = $1 LIMIT 1;`

var getTagsQuery = `SELECT
tags.id, tags.name, tags.description,
users.id, users.username, users.role, users.avatar_id
FROM tags
JOIN users ON tags.created_by_user_id = users.id
ORDER BY tags.name
;`

var updateTagQuery = `UPDATE tags SET name = $1, description = $2 WHERE id = $3;`

// postgres INSERT IGNORE
var putTagOnFileQuery = `INSERT INTO file_has_tag (tag_id, file_id, user_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING;`

var getTagsOnFileQuery = `SELECT
tags.id, tags.name, tags.description, tags.created_by_user_id
FROM file_has_tag
JOIN tags ON file_has_tag.tag_id = tags.id
WHERE file_has_tag.file_id = $1
ORDER BY tags.name
;`

var deleteTagOnFileQuery = `DELETE FROM file_has_tag WHERE tag_id = $1 AND file_id = $2;`

var deleteTagReferenceInFileHasTagQuery = `DELETE FROM file_has_tag WHERE tag_id = $1;`

var updateFoldernameQuery = `UPDATE folders SET foldername = $1 WHERE id = $2;`

var insertReviewQuery = `INSERT INTO reviews (user_id, file_id, created_at, content)
VALUES ($1, $2, $3, $4);`

var getReviewsOnFileQuery = `SELECT
reviews.id, reviews.created_at, reviews.updated_at, reviews.content,
users.id, users.username, users.role, users.avatar_id,
files.id, files.filename
FROM reviews
JOIN users ON reviews.user_id = users.id
JOIN files ON reviews.file_id = files.id
WHERE reviews.file_id = $1
ORDER BY reviews.created_at
;`

var getReviewQuery = `SELECT id, file_id, user_id, created_at, updated_at, content FROM reviews WHERE id = $1 LIMIT 1;`

var updateReviewQuery = `UPDATE reviews SET content = $1, updated_at = $2 WHERE id = $3;`

var deleteReviewQuery = `DELETE FROM reviews WHERE id = $1;`

var getReviewsByUserQuery = `SELECT
reviews.id, reviews.created_at, reviews.updated_at, reviews.content,
users.id, users.username, users.role, users.avatar_id,
files.id, files.filename
FROM reviews
JOIN users ON reviews.user_id = users.id
JOIN files ON reviews.file_id = files.id
WHERE reviews.user_id = $1
ORDER BY reviews.created_at
;`

var deleteFileQuery = `DELETE FROM files WHERE id = $1;`

var deleteFileReferenceInFileHasTagQuery = `DELETE FROM file_has_tag WHERE file_id = $1;`

var deleteFileReferenceInReviewsQuery = `DELETE FROM reviews WHERE file_id = $1;`

var updateFilenameQuery = `UPDATE files SET filename = $1 WHERE id = $2;`

var resetFilenameQuery = `UPDATE files SET filename = realname WHERE id = $1;`

var recordPlaybackQuery = `INSERT INTO playbacks (user_id, file_id, time, method, duration) VALUES ($1, $2, $3, $4, $5);`

type Stmt struct {
	initFilesTable                  *sql.Stmt
	initFoldersTable                *sql.Stmt
	initFeedbacksTable              *sql.Stmt
	initUsersTable                  *sql.Stmt
	initAvatarsTable                *sql.Stmt
	initTagsTable                   *sql.Stmt
	initFileHasTag                  *sql.Stmt
	initLikesTable                  *sql.Stmt
	initReviewsTable                *sql.Stmt
	initPlaybacksTable              *sql.Stmt
	initLogsTable                   *sql.Stmt
	initTmpfsTable                  *sql.Stmt
	insertFolder                    *sql.Stmt
	insertFile                      *sql.Stmt
	findFolder                      *sql.Stmt
	findFile                        *sql.Stmt
	searchFiles                     *sql.Stmt
	getFolder                       *sql.Stmt
	dropFiles                       *sql.Stmt
	dropFolder                      *sql.Stmt
	getFile                         *sql.Stmt
	searchFolders                   *sql.Stmt
	getFilesInFolder                *sql.Stmt
	getRandomFiles                  *sql.Stmt
	getRandomFilesWithTag           *sql.Stmt
	insertFeedback                  *sql.Stmt
	getFeedbacks                    *sql.Stmt
	deleteFeedback                  *sql.Stmt
	insertUser                      *sql.Stmt
	countUser                       *sql.Stmt
	countAdmin                      *sql.Stmt
	getUser                         *sql.Stmt
	getUsers                        *sql.Stmt
	getUserById                     *sql.Stmt
	updateUserActive                *sql.Stmt
	updateUsername                  *sql.Stmt
	updateUserPassword              *sql.Stmt
	getAnonymousUser                *sql.Stmt
	insertTag                       *sql.Stmt
	deleteTag                       *sql.Stmt
	getTag                          *sql.Stmt
	getTags                         *sql.Stmt
	updateTag                       *sql.Stmt
	putTagOnFile                    *sql.Stmt
	getTagsOnFile                   *sql.Stmt
	deleteTagOnFile                 *sql.Stmt
	deleteTagReferenceInFileHasTag  *sql.Stmt
	updateFoldername                *sql.Stmt
	insertReview                    *sql.Stmt
	getReviewsOnFile                *sql.Stmt
	getReview                       *sql.Stmt
	updateReview                    *sql.Stmt
	deleteReview                    *sql.Stmt
	getReviewsByUser                *sql.Stmt
	deleteFile                      *sql.Stmt
	deleteFileReferenceInFileHasTag *sql.Stmt
	deleteFileReferenceInReviews    *sql.Stmt
	updateFilename                  *sql.Stmt
	resetFilename                   *sql.Stmt
	recordPlaybackStmt              *sql.Stmt
}

func NewPreparedStatement(sqlConn *sql.DB) (*Stmt, error) {
	var err error

	stmt := &Stmt{}

	// init folders table
	stmt.initFoldersTable, err = sqlConn.Prepare(initFoldersTableQuery)
	if err != nil {
		return nil, err
	}
	_, err = stmt.initFoldersTable.Exec()
	if err != nil {
		return nil, err
	}

	// init files table
	stmt.initFilesTable, err = sqlConn.Prepare(initFilesTableQuery)
	if err != nil {
		return nil, err
	}
	_, err = stmt.initFilesTable.Exec()
	if err != nil {
		return nil, err
	}

	// init avatars table
	stmt.initAvatarsTable, err = sqlConn.Prepare(initAvatarsTableQuery)
	if err != nil {
		return nil, err
	}
	_, err = stmt.initAvatarsTable.Exec()
	if err != nil {
		return nil, err
	}

	// init users table
	stmt.initUsersTable, err = sqlConn.Prepare(initUsersTableQuery)
	if err != nil {
		return nil, err
	}
	_, err = stmt.initUsersTable.Exec()
	if err != nil {
		return nil, err
	}

	// init feedbacks tables
	stmt.initFeedbacksTable, err = sqlConn.Prepare(initFeedbacksTableQuery)
	if err != nil {
		return nil, err
	}
	_, err = stmt.initFeedbacksTable.Exec()
	if err != nil {
		return nil, err
	}

	// init tags table
	stmt.initTagsTable, err = sqlConn.Prepare(initTagsTableQuery)
	if err != nil {
		return nil, err
	}
	_, err = stmt.initTagsTable.Exec()
	if err != nil {
		return nil, err
	}

	// init file_has_tag table
	stmt.initFileHasTag, err = sqlConn.Prepare(initFileHasTagTableQuery)
	if err != nil {
		return nil, err
	}
	_, err = stmt.initFileHasTag.Exec()
	if err != nil {
		return nil, err
	}

	// init likes table
	stmt.initLikesTable, err = sqlConn.Prepare(initLikesTableQuery)
	if err != nil {
		return nil, err
	}
	_, err = stmt.initLikesTable.Exec()
	if err != nil {
		return nil, err
	}

	// init reviews table
	stmt.initReviewsTable, err = sqlConn.Prepare(initReviewsTableQuery)
	if err != nil {
		return nil, err
	}
	_, err = stmt.initReviewsTable.Exec()
	if err != nil {
		return nil, err
	}

	// init playbacks table
	stmt.initPlaybacksTable, err = sqlConn.Prepare(initPlaybacksTableQuery)
	if err != nil {
		return nil, err
	}
	_, err = stmt.initPlaybacksTable.Exec()
	if err != nil {
		return nil, err
	}

	// init logs table
	stmt.initLogsTable, err = sqlConn.Prepare(initLogsTableQuery)
	if err != nil {
		return nil, err
	}
	_, err = stmt.initLogsTable.Exec()
	if err != nil {
		return nil, err
	}

	// init tmpfs table
	stmt.initTmpfsTable, err = sqlConn.Prepare(initTmpfsTableQuery)
	if err != nil {
		return nil, err
	}
	_, err = stmt.initTmpfsTable.Exec()
	if err != nil {
		return nil, err
	}

	log.Println("Init tables finished")

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

	// init findFile statement
	stmt.findFile, err = sqlConn.Prepare(findFileQuery)
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

	// init getRandomFilesWithTag
	stmt.getRandomFilesWithTag, err = sqlConn.Prepare(getRandomFilesWithTagQuery)
	if err != nil {
		return nil, err
	}

	// init insertFeedback
	stmt.insertFeedback, err = sqlConn.Prepare(insertFeedbackQuery)
	if err != nil {
		return nil, err
	}

	// init getFeedbacks
	stmt.getFeedbacks, err = sqlConn.Prepare(getFeedbacksQuery)
	if err != nil {
		return nil, err
	}

	// init deleteFeedback
	stmt.deleteFeedback, err = sqlConn.Prepare(deleteFeedbackQuery)
	if err != nil {
		return nil, err
	}

	// init insertUser
	stmt.insertUser, err = sqlConn.Prepare(insertUserQuery)
	if err != nil {
		return nil, err
	}

	// init countUser
	stmt.countUser, err = sqlConn.Prepare(countUserQuery)
	if err != nil {
		return nil, err
	}

	// init countAdmin
	stmt.countAdmin, err = sqlConn.Prepare(countAdminQuery)
	if err != nil {
		return nil, err
	}

	// init getUser
	stmt.getUser, err = sqlConn.Prepare(getUserQuery)
	if err != nil {
		return nil, err
	}

	// init getUsers
	stmt.getUsers, err = sqlConn.Prepare(getUsersQuery)
	if err != nil {
		return nil, err
	}

	// init getUserById
	stmt.getUserById, err = sqlConn.Prepare(getUserByIdQuery)
	if err != nil {
		return nil, err
	}

	// init updateUserActive
	stmt.updateUserActive, err = sqlConn.Prepare(updateUserActiveQuery)
	if err != nil {
		return nil, err
	}

	// init updateUsername
	stmt.updateUsername, err = sqlConn.Prepare(updateUsernameQuery)
	if err != nil {
		return nil, err
	}

	// init updateUserPassword
	stmt.updateUserPassword, err = sqlConn.Prepare(updateUserPasswordQuery)
	if err != nil {
		return nil, err
	}

	// init getAnonymousUser
	stmt.getAnonymousUser, err = sqlConn.Prepare(getAnonymousUserQuery)
	if err != nil {
		return nil, err
	}

	// insert Anonymous user if users is empty
	userCount := 0
	err = stmt.countUser.QueryRow().Scan(&userCount)
	if err != nil {
		return nil, err
	}
	if userCount == 0 {
		_, err = stmt.insertUser.Exec("Anonymous user", "", 0, 1, 0)
		if err != nil {
			return nil, err
		}
	}

	// init insertTag
	stmt.insertTag, err = sqlConn.Prepare(insertTagQuery)
	if err != nil {
		return nil, err
	}

	// init deleteTag
	stmt.deleteTag, err = sqlConn.Prepare(deleteTagQuery)
	if err != nil {
		return nil, err
	}

	// init getTag
	stmt.getTag, err = sqlConn.Prepare(getTagQuery)
	if err != nil {
		return nil, err
	}

	// init getTags
	stmt.getTags, err = sqlConn.Prepare(getTagsQuery)
	if err != nil {
		return nil, err
	}

	// init updateTag
	stmt.updateTag, err = sqlConn.Prepare(updateTagQuery)
	if err != nil {
		return nil, err
	}

	// init putTagOnFile
	stmt.putTagOnFile, err = sqlConn.Prepare(putTagOnFileQuery)
	if err != nil {
		return nil, err
	}

	// init getTagsOnFile
	stmt.getTagsOnFile, err = sqlConn.Prepare(getTagsOnFileQuery)
	if err != nil {
		return nil, err
	}

	// init deleteTagOnFile
	stmt.deleteTagOnFile, err = sqlConn.Prepare(deleteTagOnFileQuery)
	if err != nil {
		return nil, err
	}

	// init deleteTagReferenceInFileHasTag
	stmt.deleteTagReferenceInFileHasTag, err = sqlConn.Prepare(
		deleteTagReferenceInFileHasTagQuery)
	if err != nil {
		return nil, err
	}

	// init updateFoldername
	stmt.updateFoldername, err = sqlConn.Prepare(updateFoldernameQuery)
	if err != nil {
		return nil, err
	}

	// init insertReview
	stmt.insertReview, err = sqlConn.Prepare(insertReviewQuery)
	if err != nil {
		return nil, err
	}

	// init getReviewsOnFile
	stmt.getReviewsOnFile, err = sqlConn.Prepare(getReviewsOnFileQuery)
	if err != nil {
		return nil, err
	}

	// init getReview
	stmt.getReview, err = sqlConn.Prepare(getReviewQuery)
	if err != nil {
		return nil, err
	}

	// init updateReview
	stmt.updateReview, err = sqlConn.Prepare(updateReviewQuery)
	if err != nil {
		return nil, err
	}

	// init deleteReview
	stmt.deleteReview, err = sqlConn.Prepare(deleteReviewQuery)
	if err != nil {
		return nil, err
	}

	// init getReviewsByUser
	stmt.getReviewsByUser, err = sqlConn.Prepare(getReviewsByUserQuery)
	if err != nil {
		return nil, err
	}

	// init deleteFile
	stmt.deleteFile, err = sqlConn.Prepare(deleteFileQuery)
	if err != nil {
		return nil, err
	}

	// init deleteFileReferenceInFileHasTag
	stmt.deleteFileReferenceInFileHasTag, err = sqlConn.Prepare(
		deleteFileReferenceInFileHasTagQuery)
	if err != nil {
		return nil, err
	}

	// init deleteFileReferenceInReviews
	stmt.deleteFileReferenceInReviews, err = sqlConn.Prepare(
		deleteFileReferenceInReviewsQuery)
	if err != nil {
		return nil, err
	}

	// init updateFilename
	stmt.updateFilename, err = sqlConn.Prepare(updateFilenameQuery)
	if err != nil {
		return nil, err
	}

	// init resetFilename
	stmt.resetFilename, err = sqlConn.Prepare(resetFilenameQuery)
	if err != nil {
		return nil, err
	}

	stmt.recordPlaybackStmt, err = sqlConn.Prepare(recordPlaybackQuery)
	if err != nil {
		return nil, err
	}

	log.Println("Init statements finished")

	return stmt, err
}
