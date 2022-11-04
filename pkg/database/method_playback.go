package database

func (database *Database) RecordPlayback(playback Playback) error {
	_, err := database.stmt.recordPlaybackStmt.Exec(
		playback.UserID, playback.FileID, playback.Time, playback.Method, playback.Duration)
	return err
}
