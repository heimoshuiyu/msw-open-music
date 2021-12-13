package database

func (database *Database) Login(username string, password string) (*User, error) {
	user := &User{}

	// get user from database
	err := database.stmt.getUser.QueryRow(username, password).Scan(&user.ID, &user.Username, &user.Role, &user.Active, &user.AvatarId)
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

func (database *Database) Register(username string, password string, usertype int64) error {
	countAdmin, err := database.CountAdmin()
	if err != nil {
		return err
	}

	active := false
	if countAdmin == 0 {
		active = true
	}

	// active normal user by default
	if usertype == 2 {
		active = true
	}

	_, err = database.stmt.insertUser.Exec(username, password, usertype, active, 0)
	if err != nil {
		return err
	}
	return nil
}

func (database *Database) GetUserById(id int64) (*User, error) {
	user := &User{}

	// get user from database
	err := database.stmt.getUserById.QueryRow(id).Scan(&user.ID, &user.Username, &user.Role, &user.Active, &user.AvatarId)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (database *Database) CountAdmin() (int64, error) {
	var count int64
	err := database.stmt.countAdmin.QueryRow().Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (database *Database) GetUsers() ([]*User, error) {
	users := make([]*User, 0)

	rows, err := database.stmt.getUsers.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &User{}
		err = rows.Scan(&user.ID, &user.Username, &user.Role, &user.Active, &user.AvatarId)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (database *Database) UpdateUserActive(id int64, active bool) error {
	_, err := database.stmt.updateUserActive.Exec(active, id)
	if err != nil {
		return err
	}
	return nil
}
