package controllers

type UserIdParam struct {
	UserId string `in:"query=userId"`
}

// TODO: finish using httpin...

type UserIdPath struct {
	UserId string `in:"query=userId"`
}
