package models

type Token struct {
	Account_id    string
	Token_id      string
}

func NewToken(accountId string, token string) *Token {
	return &Token{accountId, token}
}