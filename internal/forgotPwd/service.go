package forgotPwd

type IForgotPwdService interface {
	GenerateAndSendMagicLink(email string) error
}

type ForgotPwdService struct{}

func (*ForgotPwdService) GenerateAndSendMagicLink(email string) error {
	return nil
}
