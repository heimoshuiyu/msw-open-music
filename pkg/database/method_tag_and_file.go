package database

func (database *Database) PutTagOnFile(tagID, fileID, userID int64) error {
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	_, err := database.stmt.putTagOnFile.Exec(tagID, fileID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (database *Database) GetTagsOnFile(fileID int64) ([]*Tag, error) {
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	rows, err := database.stmt.getTagsOnFile.Query(fileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]*Tag, 0)
	for rows.Next() {
		tag := &Tag{}
		err = rows.Scan(&tag.ID, &tag.Name, &tag.Description, &tag.CreatedByUserId)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (database *Database) DeleteTagOnFile(tagID, fileID int64) error {
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	result, err := database.stmt.deleteTagOnFile.Exec(tagID, fileID)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return ErrTagNotFound
	}
	return nil
}
