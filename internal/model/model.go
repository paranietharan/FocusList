package model

type Role string

const (
	AdminRole     Role = "admin"
	UserRole      Role = "user"
	ModeratorRole Role = "moderator"
)

type User struct { // store the user information
	Email     string `db:"email"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Password  string `db:"password"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
	IsActive  bool   `db:"is_active"`
	Role      Role   `db:"role"`
}

type TodoListBucket struct { // store the todo list category information
	ID        string `db:"id"`
	Name      string `db:"name"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}
