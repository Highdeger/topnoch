package credentials

type DeviceAuthType int

var DeviceAuthTypeTitles = []string{"DeviceAuthSsh1", "DeviceAuthSsh2", "DeviceAuthTelnet"}

const (
	DeviceAuthSsh1 DeviceAuthType = iota
	DeviceAuthSsh2
	DeviceAuthTelnet
)

type DeviceCredential struct {
	Username string         `json:","` // username for login
	Password string         `json:","` // password for login
	AuthType DeviceAuthType `json:","` // type of authentication for login
}
