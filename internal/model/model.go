package model

type Role string

const (
	SuperAdminRole Role = "super_admin"
	AdminRole      Role = "admin"
	UserRole       Role = "user"
	ModeratorRole  Role = "moderator"
)

type RolePermissions string

const (
	ReadPermission    RolePermissions = "read"
	WritePermission   RolePermissions = "write"
	ExecutePermission RolePermissions = "execute" // can add and remove members
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

type TodoListBucketUser struct { // store the relationship between todo list buckets and users
	BucketID   string          `db:"bucket_id"`
	UserEmail  string          `db:"user_email"`
	Permission RolePermissions `db:"permission"`
}

type TodoListItem struct { // store the todo list item information
	ID          string `db:"id"`
	BucketID    string `db:"bucket_id"`
	Description string `db:"description"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
	IsCompleted bool   `db:"is_completed"`
}
