package database

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type Database struct {
	filepath   string
	data       *Data
	JWT_SECRET string
	mux        *sync.RWMutex
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Data struct {
	Users         map[int]User    `json:"users"`
	Hashes        map[int]string  `json:"hashes"`
	RefreshTokens map[string]bool `json:"refresh_tokens"`

	Usernames  map[string]int `json:"usernames"`
	NextUserID int            `json:"next_user_id"`
}

func InitDatabase() (*Database, error) {
	fp := os.Getenv("DB_PATH")
	if fp == "" {
		return nil, errors.New(fmt.Sprint("'DB_PATH': ", fp))
	}

	database := Database{
		mux:        &sync.RWMutex{},
		filepath:   fp,
		data:       nil,
		JWT_SECRET: os.Getenv("JWT_SECRET")}

	// Check if debug mode, reset json on startup
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	var err error
	if *dbg {
		fmt.Println("Creating fresh database...")
		err = database.createFreshData()
	} else {
		fmt.Println("Loading database from file...")
		err = database.loadData()
	}
	if err != nil {
		return nil, err
	}

	return &database, nil
}

func (database *Database) TryAddUser(username string, password string) (*User, error) {
	database.mux.Lock()
	defer database.mux.Unlock()

	// Check if username is already registered
	if database.isUsernameTaken(username) {
		return nil, errors.New("username is taken")
	}

	// Get next UID
	uid := database.data.NextUserID
	database.data.NextUserID++

	// Save hash
	hash, err := HashPassword(password)
	if err != nil {
		database.data.NextUserID--
		return nil, err
	}
	database.data.Hashes[uid] = string(hash)

	// Create user
	user := User{
		ID:       uid,
		Username: username}

	// Add user to database
	database.data.Users[uid] = user
	database.data.Usernames[username] = uid

	fmt.Println("Added a new user!")
	return &user, nil
}

func (database *Database) IsValidCredentials(username string, password string) bool {
	// Check user in database
	uid, ok := database.data.Usernames[username]
	if !ok {
		// No username found
		return false
	}

	// Check hash matches password
	err := bcrypt.CompareHashAndPassword([]byte(database.data.Hashes[uid]), []byte(password))
	if err != nil {
		return false
	} else {
		return true
	}
}

func (database *Database) isUsernameTaken(username string) bool {
	_, ok := database.data.Usernames[username]
	if ok {
		return true
	} else {
		return false
	}
}

func (database *Database) createFreshData() error {
	database.mux.Lock()
	defer database.mux.Unlock()

	d := &Data{
		Users:         make(map[int]User),
		Hashes:        make(map[int]string),
		RefreshTokens: make(map[string]bool),
		Usernames:     make(map[string]int),
		NextUserID:    0}

	dat, err := json.Marshal(d)
	if err != nil {
		return err
	}

	database.data = d
	os.WriteFile(database.filepath, dat, 0777)

	return nil
}

func (database *Database) loadData() error {
	database.mux.Lock()
	defer database.mux.Unlock()

	dat, err := os.ReadFile(database.filepath)
	if err != nil {
		return err
	}

	d := &Data{}
	err = json.Unmarshal(dat, &d)
	if err != nil {
		return err
	}

	database.data = d
	return nil
}

func (database *Database) SaveData() error {
	database.mux.RLock()
	defer database.mux.RUnlock()

	dat, err := json.Marshal(database.data)
	if err != nil {
		return err
	}
	os.WriteFile(database.filepath, dat, 0777)

	return nil
}

func HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}
