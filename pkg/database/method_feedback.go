package database

func (database *Database) InsertFeedback(time int64, content string, userID int64, header string) error {
	_, err := database.stmt.insertFeedback.Exec(time, content, userID, header)
	if err != nil {
		return err
	}
	return nil
}

func (database *Database) GetFeedbacks() ([]*Feedback, error) {
	rows, err := database.stmt.getFeedbacks.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feedbacks := make([]*Feedback, 0)
	for rows.Next() {
		feedback := &Feedback{
			User: &User{},
		}
		err := rows.Scan(
			&feedback.ID, &feedback.Time, &feedback.Content, &feedback.Header,
			&feedback.User.ID, &feedback.User.Username, &feedback.User.Role, &feedback.User.Active, &feedback.User.AvatarId)
		if err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, feedback)
	}
	return feedbacks, nil
}

func (database *Database) DeleteFeedback(id int64) error {
	_, err := database.stmt.deleteFeedback.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
