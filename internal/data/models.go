package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Game        GameModel
	Tokens      TokenModel
	Users       UserModel
	Permissions PermissionModel
	Team        TeamModel
	TeamUsers   TeamUsersModel
	Tournament  TournamentModel
	Match       MatchModel
	UserMatch   UserMatchModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Game:        GameModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Users:       UserModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Team:        TeamModel{DB: db},
		TeamUsers:   TeamUsersModel{DB: db},
		Tournament:  TournamentModel{DB: db},
		Match:       MatchModel{DB: db},
		UserMatch:   UserMatchModel{DB: db},
	}
}
