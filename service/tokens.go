package service

//TokenService share a token between the creds handler and hipchat handler
type TokenService struct {
	Token string
}

//NewTokenService instantiates a token service
func NewTokenService() *TokenService {
	return &TokenService{}
}
