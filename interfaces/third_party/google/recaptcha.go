package google

type IRecaptchaAPI interface {
	Verify(responseToken string) (bool, error)
}
