package objects

type PasswordParams struct {
	ServiceName        string `validate:"required"`
	MasterPassword     string `validate:"required_without=MasterPasswordFile"`
	MasterPasswordFile string `validate:"required_without=MasterPassword,omitempty,file"`
	Length             uint8  `validate:"omitempty,max=40"`
	Version            int    `validate:"omitempty,min=1"`
}
