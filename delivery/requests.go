package delivery

type SignupRequest struct {
	Login     string `json:"login"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	BirthDate string `json:"birth_date"`
	Name      string `json:"name"`
}

type SigninRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CommentRequest struct {
	FilmId uint64 `json:"film_id"`
	Rating uint16 `json:"rating"`
	Text   string `json:"text"`
}
