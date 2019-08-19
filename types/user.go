package types

// UserInfo is used for register
type UserInfo struct {
	PID                 string `json:"pid"`
	Username            string `json:"username"`
	Password            string `json:"password"`
	Phone               string `json:"phone"`
	AuthKey             string `json:"auth_key"`
	BeginTimestamp      int    `json:"begin_timestamp"`
	ExpiredTimestamp    int    `json:"expired_timestamp"`
	BeginTimestamp3rd   int    `json:"begin_timestamp3rd"`
	ExpiredTimestamp3rd int    `json:"expired_timestamp3rd"`
}

//Credential is used for login
type Credential struct {
	PID      string `json:"pid"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// Identity used for identity user
type Identity struct {
	PID   string `json:"pid" form:"pid"`
	Phone string `json:"phone" form:"phone"`
}

// ValidPeriod TODO
type ValidPeriod struct {
}

// RegisterObj TODO
type RegisterObj struct {
	UserInfo
	Token string `json:"token"`
}

// PidPhone TODO
type PidPhone struct {
	PID   string `json:"pid"`
	Phone string `json:"phone"`
}

// CheckPopInput TODO
type CheckPopInput struct {
	Token string `json:"token"`
	PidPhone
}

// CheckAuthInput TODO
type CheckAuthInput struct {
	PidPhone
	Period  int    `json:"period"`
	AuthKey string `json:"auth_key"`
}

// AuthClaimInput TODO
type AuthClaimInput struct {
	PidPhone
	DeviceList []string `json:"device_list"`
	Start      int      `json:"start"`
	End        int      `json:"end"`
	Period     int      `json:"period"`
	AuthKey    string   `json:"auth_key"`
	Recevier   string   `json:"receiver"`
}

// AuthList TODO
type AuthList struct {
	PidPhone
	// PID        string   `json:"pid"`
	// Phone      string   `json:"phone"`
	DeviceList []string `json:"device_list"`
	Start      int      `json:"start"`
	End        int      `json:"end"`
	Period     int      `json:"period"`
	Receiver   string   `json:"receiver"`
}
