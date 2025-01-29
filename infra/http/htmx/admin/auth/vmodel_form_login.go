package auth

type VModelFormLogin struct {
	Email              string `form:"email"`
	GRecaptchaResponse string `form:"g-recaptcha-response"`

	FormErrors  LoginFormErrors
	ServerError string
}

type LoginFormErrors struct {
	Email    string
	Password string
}

func (f *VModelFormLogin) Validate() bool {
	f.FormErrors = LoginFormErrors{}
	isValid := true

	if f.Email == "" {
		f.FormErrors.Email = "Email is required"
		isValid = false
	}

	// TODO: implement email format validation

	return isValid
}

func (f *VModelFormLogin) GetName() string {
	return ""
}

func (f *VModelFormLogin) GetEmail() string {
	return f.Email
}

func (f *VModelFormLogin) GetPassword() string {
	return ""
}
