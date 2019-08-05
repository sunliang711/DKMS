package types

// UserInfo is used for register
type UserInfo struct {
	PID                 string `json:"pid"`
	Username            string `json:"username"`
	Password            string `json:"password"`
	Phone               string `json:"phone"`
	AuthKey             string `json:"auth_key"`
	ExpiredTimestamp    int    `json:"expired_timestamp"`
	ExpiredTimestamp3rd int    `json:"expired_timestamp3rd"`
}

//Credential is used for login
type Credential struct {
	PID      string `json:"pid"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
