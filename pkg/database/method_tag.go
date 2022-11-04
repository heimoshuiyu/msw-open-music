package database

import "errors"

func (database *Database) InsertTag(tag *Tag) (int64, error) {
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	result, err := database.stmt.insertTag.Query(tag.Name, tag.Description, tag.CreatedByUserId)
	if err != nil {
		return 0, err
	}
	var id int64
	for result.Next() {
		err = result.Scan(&id)
		if err != nil {
			return 0, err
		}
	}
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (database *Database) GetTag(id int64) (*Tag, error) {
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	tag := &Tag{CreatedByUser: &User{}}
	err := database.stmt.getTag.QueryRow(id).Scan(
		&tag.ID, &tag.Name, &tag.Description,
		&tag.CreatedByUser.ID, &tag.CreatedByUser.Username, &tag.CreatedByUser.Role, &tag.CreatedByUser.AvatarId)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (database *Database) GetTags() ([]*Tag, error) {
	tags := []*Tag{}

	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	rows, err := database.stmt.getTags.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tag := &Tag{CreatedByUser: &User{}}
		err := rows.Scan(
			&tag.ID, &tag.Name, &tag.Description,
			&tag.CreatedByUser.ID, &tag.CreatedByUser.Username, &tag.CreatedByUser.Role, &tag.CreatedByUser.AvatarId)

		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (database *Database) UpdateTag(tag *Tag) error {
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	result, err := database.stmt.updateTag.Exec(tag.Name, tag.Description, tag.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("No rows affected")
	}

	return nil
}

// delete tag and all its references in file_has_tag
func (database *Database) DeleteTag(id int64) error {
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	// begin transaction
	tx, err := database.sqlConn.Begin()
	if err != nil {
		return err
	}

	// delete tag
	_, err = tx.Stmt(database.stmt.deleteTag).Exec(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// delete file_has_tag
	_, err = tx.Stmt(database.stmt.deleteTagReferenceInFileHasTag).Exec(id)
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
