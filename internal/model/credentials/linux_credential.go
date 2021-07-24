package credentials

type LinuxAuthType int

var LinuxAuthTypeTitles = []string{"LinuxAuthLocal", "LinuxAuthDomain"}

const (
	LinuxAuthLocal  = 0
	LinuxAuthDomain = 1
)

type LinuxCredential struct {
	Username string        `json:","` // username for login
	Password string        `json:","` // password for login
	AuthType LinuxAuthType `json:","` // type of authentication for login
}
