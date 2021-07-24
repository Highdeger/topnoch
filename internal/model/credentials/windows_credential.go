package credentials

type WindowsAuthType int

var WindowsAuthTypeTitles = []string{"WindowsAuthLocal", "WindowsDomain"}

const (
	WindowsAuthLocal WindowsAuthType = iota
	WindowsDomain
)

type WindowsCredential struct {
	Username string          `json:","` // username for login
	Password string          `json:","` // password for login
	AuthType WindowsAuthType `json:","` // type of authentication for login
}
