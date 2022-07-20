package database

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func (database *Database) Login(username string, password string) (*User, error) {
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	user := &User{}

	// get user from database
	err := database.stmt.getUser.QueryRow(username).Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.Active, &user.AvatarId)
	if err != nil {
		return user, err
	}

	// validate password
	err = database.ComparePassword(user.Password, password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (database *Database) LoginAsAnonymous() (*User, error) {
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

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

	// encrypt password
	password = database.EncryptPassword(password)

	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	_, err = database.stmt.insertUser.Exec(username, password, usertype, active, 0)
	if err != nil {
		return err
	}
	return nil
}

func (database *Database) GetUserById(id int64) (*User, error) {
	user := &User{}

	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	// get user from database
	err := database.stmt.getUserById.QueryRow(id).Scan(&user.ID, &user.Username, &user.Role, &user.Active, &user.AvatarId)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (database *Database) CountAdmin() (int64, error) {
	var count int64

	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	err := database.stmt.countAdmin.QueryRow().Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (database *Database) GetUsers() ([]*User, error) {
	users := make([]*User, 0)

	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

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
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	_, err := database.stmt.updateUserActive.Exec(active, id)
	if err != nil {
		return err
	}
	return nil
}

func (database *Database) UpdateUsername(id int64, username string) error {
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	_, err := database.stmt.updateUsername.Exec(username, id)
	if err != nil {
		return err
	}
	return nil
}

func (database *Database) UpdateUserPassword(id int64, password string) error {
	database.singleThreadLock.Lock()
	defer database.singleThreadLock.Unlock()

	// encrypt password
	password = database.EncryptPassword(password)

	_, err := database.stmt.updateUserPassword.Exec(password, id)
	if err != nil {
		return err
	}
	return nil
}

func (database *Database) EncryptPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println("[database] Failed to hash password, using plaintext password")
		return password
	}

	return string(hash)
}

func (database *Database) ComparePassword(hashedPassword string, plainTextPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainTextPassword))
	return err
}
