package user

type EmitUserCreated struct {
	ID string `json:"id"`
}

type EmitUserCreatedError struct {
	ID string `json:"id"`
}