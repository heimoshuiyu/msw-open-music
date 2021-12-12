package database

func (database *Database) InsertReview(review *Review) error {
	_, err := database.stmt.insertReview.Exec(
		review.UserId,
		review.FileId,
		review.CreatedAt,
		review.Content)
	return err
}
