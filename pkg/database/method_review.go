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

func (database *Database) GetReview(reviewId int64) (*Review, error) {
	row := database.stmt.getReview.QueryRow(reviewId)

	review := &Review{}
	err := row.Scan(
		&review.ID,
		&review.FileId,
		&review.UserId,
		&review.CreatedAt,
		&review.UpdatedAt,
		&review.Content)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (database *Database) UpdateReview(review *Review) error {
	_, err := database.stmt.updateReview.Exec(
		review.Content,
		review.UpdatedAt,
		review.ID)
	return err
}

func (database *Database) DeleteReview(reviewId int64) error {
	_, err := database.stmt.deleteReview.Exec(reviewId)
	return err
}

func (database *Database) GetReviewsByUser(userId int64) ([]*Review, error) {
	rows, err := database.stmt.getReviewsByUser.Query(userId)
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
