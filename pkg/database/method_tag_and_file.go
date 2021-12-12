package database

func (database *Database) PutTagOnFile(tagID, fileID, userID int64) error {
	result, err := database.stmt.putTagOnFile.Exec(tagID, fileID, userID)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return ErrTagNotFound
	}
	return nil
}

func (database *Database) GetTagsOnFile(fileID int64) ([]*Tag, error) {
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
	result, err := database.stmt.deleteTagOnFile.Exec(tagID, fileID)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return ErrTagNotFound
	}
	return nil
}
