package middlewares

const (
	RoleUser  = 1
	RoleAdmin = 2
)

type Role struct {
	Id    int    `db:"id"`
	Title string `db:"title"`
}
