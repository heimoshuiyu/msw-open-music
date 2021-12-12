package database

func (database *Database) PutTagOnFile(tagID, fileID, userID int64) error {
	_, err := database.stmt.putTagOnFile.Exec(tagID, fileID, userID)
	return err
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
