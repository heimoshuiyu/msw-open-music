package database

import "errors"

func (database *Database) InsertTag(tag string, description string) (*Tag, error) {
	result, err := database.stmt.insertTag.Exec(tag, description)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return database.GetTag(id)
}

func (database *Database) GetTag(id int64) (*Tag, error) {
	tag := &Tag{}
	err := database.stmt.getTag.QueryRow(id).Scan(&tag.ID, &tag.Name, &tag.Description)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (database *Database) GetTags() ([]Tag, error) {
	tags := []Tag{}
	rows, err := database.stmt.getTags.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tag := Tag{}
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Description)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (database *Database) UpdateTag(tag *Tag) (*Tag, error) {
	result, err := database.stmt.updateTag.Exec(tag.Name, tag.Description, tag.ID)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, errors.New("No rows affected")
	}
	return database.GetTag(tag.ID)
}
