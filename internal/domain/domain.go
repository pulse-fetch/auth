package domain

type PostgresParams struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
	Sslmode  string
}

type User struct {
	Id           int64
	Username     string
	Email        string
	HashPassword string
}

type ResponseUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
