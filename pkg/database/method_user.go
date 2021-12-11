package database

func (database *Database) Login(username string, password string) (*User, error) {
	user := &User{}

	// get user from database
	err := database.stmt.getUser.QueryRow(username, password).Scan(&user.ID, &user.Username, &user.Role, &user.AvatarId)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (database *Database) LoginAsAnonymous() (*User, error) {
	user := &User{}

	// get user from database
	err := database.stmt.getAnonymousUser.QueryRow().Scan(&user.ID, &user.Username, &user.Role, &user.AvatarId)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (database *Database) Register(username string, password string, usertype int64) (*User, error) {
	_, err := database.stmt.insertUser.Exec(username, password, usertype, 0)
	if err != nil {
		return nil, err
	}
	return database.Login(username, password)
}

func (database *Database) GetUserById(id int64) (*User, error) {
	user := &User{}

	// get user from database
	err := database.stmt.getUserById.QueryRow(id).Scan(&user.ID, &user.Username, &user.Role, &user.AvatarId)
	if err != nil {
		return user, err
	}
	return user, nil
}
