package main

type User struct {
	ID           int    `json:"id"`
	MembershipID string `json:"membership_id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

type SignInCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
