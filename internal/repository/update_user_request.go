package repository

type UpdateUserRequest struct {
	Username  *string
	Password  *string
	Role      *string
	Active    *bool
	UpdatedBy *string
}
