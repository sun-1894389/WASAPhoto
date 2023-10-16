/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.
To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.
For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
Questo codice serve come middleware tra l'applicazione WASAPhoto e il suo database,fornendo
funzionalità per eseguire operazioni CRUD (Create, Read, Update, Delete) sul database in modo strutturato e organizzato.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// Errors section
var ErrPhotoDoesntExist = errors.New("photo doesn't exist")
var ErrUserBanned = errors.New("user is banned")

/*
var ErrUserAutoLike = errors.New("users can't like their own photos")
var ErrUserAutoFollow = errors.New("users can't follow themselfes")
*/

// Constants che indica foto per home
const PhotosPerUserHome = 3

// AppDatabase è l'interfaccia per i DB
// L'interfaccia AppDatabase definisce un set di metodi che qualsiasi implementazione del database dovrebbe fornire.
type AppDatabase interface {

	// Creates a new user in the database. It returns an error
	CreateUser(User) error

	// Modifies the nickname of a user in the database. It returns an error
	ModifyNickname(User, Nickname) error

	// Searches all the users that match the given name (both identifier and nickname). Returns the list of matching users and an error
	SearchUser(searcher User, userToSearch User) ([]CompleteUser, error)

	// Creates a new photo in the database. It returns the photo identifier and an error
	CreatePhoto(Photo) (int64, error)

	// Inserts a like of a user for a specified photo in the database. It returns an error
	LikePhoto(PhotoId, User) error

	// Removes a like of a user for a specified photo from the database. It returns an error
	UnlikePhoto(PhotoId, User) error

	// Adds a comment from a user to a specified photo in the database. It returns the unique comment id and an error
	CommentPhoto(PhotoId, User, Comment) (int64, error)

	// Deletes a comment from a user from a specified photo in the database. It returns an error
	UncommentPhoto(PhotoId, User, CommentId) error

	// Adds a follower (a) to the user that is being followed (b). It returns an error
	FollowUser(a User, b User) error

	// Removes a follower (a) from the user that is being unfollowed (b). It returns an error
	UnfollowUser(a User, b User) error

	// Adds a user (b) to the banned list of another (a). It returns an error
	BanUser(a User, b User) error

	// Removes a user (b) from the banned list of another (a). It returns an error
	UnbanUser(a User, b User) error

	// Get the a user's stream (photos of people who are followed by the user in reversed chronological order). It returns the photos and an error
	GetStream(User) ([]Photo, error)

	// Removes a photo from the database. The removal includes likes and comments.  It returns an error
	RemovePhoto(User, PhotoId) error

	// ____________________________________  Util Methods ____________________________________

	// Gets the followers list for the specified user. Returns the followers list and an error
	GetFollowers(User) ([]User, error)

	// Gets the following list for the specified user. Returns the following list and an error
	GetFollowing(User) ([]User, error)

	// Gets the photos list of user b for the user a. Returns the photo list and an error
	GetPhotosList(a User, b User) ([]Photo, error)

	// Allows the author of a photo to remove a comment from another user on his/her photo. Returns an error
	UncommentPhotoAuthor(PhotoId, CommentId) error

	// Gets the nickname of a user. Returns the nickname and an error
	GetNickname(User) (string, error)

	// Checks if a user (a) is banned by another (b). Returns a boolean
	BannedUserCheck(a User, b User) (bool, error)

	// Checks if a user (a) exists
	CheckUser(a User) (bool, error)

	// Checks if a photo (via its id) exists. Returns an error
	CheckPhotoExistence(p PhotoId) (bool, error)

	// Ping checks whether the database is available or not (in that case, an error will be returned)
	Ping() error
}

// Questa struttura contiene un campo c che rappresenta la connessione al database.
type appdbimpl struct {
	c *sql.DB
}

// Questa funzione crea e restituisce una nuova istanza di AppDatabase. Se il database non esiste, viene creato.
func New(db *sql.DB) (AppDatabase, error) {
	// non è stata fornita una connessione valida al database.
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Esegue una query SQL per attivare il supporto alle chiavi esterne in SQLite.
	// Questo è importante per garantire l'integrità referenziale tra le tabelle.
	_, errPramga := db.Exec(`PRAGMA foreign_keys= ON`)
	if errPramga != nil {
		return nil, fmt.Errorf("error setting pragmas: %w", errPramga)
	}

	// Esegue una query SQL per verificare se esiste una tabella chiamata 'users' nel database.
	// Il risultato (il nome della tabella) viene memorizzato nella variabile tableName.
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		// Se la tabella 'users' non esiste, chiama la funzione createDatabase per creare tutte le tabelle necessarie nel database.
		err = createDatabase(db)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}
	// restituisce un'istanza dell'implementazione appdbimpl dell'interfaccia AppDatabase
	return &appdbimpl{
		c: db,
	}, nil
}

// Questa funzione verifica se il database è disponibile o meno
func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

// funzione crea tutte le tabelle necessarie per l'applicazione WASAPhoto nel database. 
// Utilizza una serie di stringhe SQL per definire le tabelle e le relazioni tra di esse.
func createDatabase(db *sql.DB) error {
	// Ogni stringa rappresenta una query SQL per creare una tabella specifica nel database.
	// Le 6 stringhe comandi SQL per creare le tabelle users, photos, likes, comments, banned_users e followers se non esistono già.
	tables := [6]string{
		`CREATE TABLE IF NOT EXISTS users (
			id_user VARCHAR(16) NOT NULL PRIMARY KEY,
			nickname VARCHAR(16) NOT NULL
			);`,
		`CREATE TABLE IF NOT EXISTS photos (
			id_photo INTEGER PRIMARY KEY AUTOINCREMENT,
			id_user VARCHAR(16) NOT NULL,
			date DATETIME NOT NULL,
			FOREIGN KEY(id_user) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS  likes (
			id_photo INTEGER NOT NULL,
			id_user VARCHAR(16) NOT NULL,
			PRIMARY KEY (id_photo,id_user),
			FOREIGN KEY(id_photo) REFERENCES photos (id_photo) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS comments (
			id_comment INTEGER PRIMARY KEY AUTOINCREMENT,
			id_photo INTEGER NOT NULL,
			id_user VARCHAR(16) NOT NULL,
			comment VARCHAR(30) NOT NULL,
			FOREIGN KEY(id_photo) REFERENCES photos (id_photo) ON DELETE CASCADE,
			FOREIGN KEY(id_user) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS banned_users (
			banner VARCHAR(16) NOT NULL,
			banned VARCHAR(16) NOT NULL,
			PRIMARY KEY (banner,banned),
			FOREIGN KEY(banner) REFERENCES users (id_user) ON DELETE CASCADE,
			FOREIGN KEY(banned) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS followers(
			follower VARCHAR(16) NOT NULL,
			followed VARCHAR(16) NOT NULL,
			PRIMARY KEY (follower,followed),
			FOREIGN KEY(follower) REFERENCES users (id_user) ON DELETE CASCADE,
			FOREIGN KEY(followed) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
	}

	// Inizia un ciclo for che itera su ogni stringa nell'array tables,per eseguire ogni query SQL nell'array per creare le tabelle.
	for i := 0; i < len(tables); i++ {
		// Assegno la query SQL corrente dall'array tables alla variabile sqlStmt.
		// utilizzo il metodo Exec,se la tabella specificata dalla query esiste già,la query non avrà effetto grazie al IF NOT EXISTS nelle stringhe SQL.
		sqlStmt := tables[i]
		_, err := db.Exec(sqlStmt)

		if err != nil {
			return err
		}
	}
	return nil
}
