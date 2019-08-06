package types

//KeyPair store pk,sk pair
type KeyPair struct{
	PK string `json:"pk"`
	SK string `json:"sk"`
}

// AuthKey for update authkey
type AuthKey struct {
	Identity
	AuthKey string `json:"auth_key"`
}