package auth

type VModelFormRegister struct {
	Email string `form:"email"`

	FormErrors  RegisterFormErrors
	ServerError string
}

type RegisterFormErrors struct {
	Email    string
	Password string
}

func (f *VModelFormRegister) Validate() bool {
	f.FormErrors = RegisterFormErrors{}
	isValid := true

	if f.Email == "" {
		f.FormErrors.Email = "Email is required"
		isValid = false
	}

	// TODO: implement email format validation

	return isValid
}

func (f *VModelFormRegister) GetName() string {
	return ""
}

func (f *VModelFormRegister) GetEmail() string {
	return f.Email
}

func (f *VModelFormRegister) GetPassword() string {
	return ""
}
