package person_in_film

import "database/sql"

type PersonInFilmItem struct {
	IdFilm        int
	IdPerson      int
	IdProfession  int
	CharacterName sql.NullString
}
