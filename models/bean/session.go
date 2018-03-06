package bean

type CreateSession struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type OutPutSession struct {
	Uid   int    `json:"uid"`
	Status  int `json:"status"`
	Token string `json:"token"`
}
