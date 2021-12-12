package database

func (database *Database) InsertReview(review *Review) error {
	_, err := database.stmt.insertReview.Exec(
		review.UserId,
		review.FileId,
		review.CreatedAt,
		review.Content)
	return err
}

func (database *Database) GetReviewsOnFile(fileId int64) ([]*Review, error) {
	rows, err := database.stmt.getReviewsOnFile.Query(fileId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := make([]*Review, 0)
	for rows.Next() {
		review := &Review{
			User: &User{},
			File: &File{},
		}
		err := rows.Scan(
			&review.ID,
			&review.CreatedAt,
			&review.UpdatedAt,
			&review.Content,
			&review.User.ID,
			&review.User.Username,
			&review.User.Role,
			&review.User.AvatarId,
			&review.File.ID,
			&review.File.Filename)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}
