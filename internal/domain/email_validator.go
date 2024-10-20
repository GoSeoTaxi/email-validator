package domain

type EmailValidator interface {
	Validate(email string) (bool, string)
}
