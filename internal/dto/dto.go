package dto

type TodoListBucketUserDTO struct {
	UserEmail string `json:"userEmail"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Password  string `db:"password"`
}
