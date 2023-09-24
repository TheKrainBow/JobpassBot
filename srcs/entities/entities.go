package entities

import "sync"

var (
	StorePath           = "./srcs/store/store.json"
	DataMutext          sync.Mutex
	Data                *Store
	GitlabProjectMember []GitlabUser
)

type GitlabAccessLevel int

type GitlabUser struct {
	ID          int               `json:"ID"`
	Name        string            `json:"Name"`
	Username    string            `json:"Username"`
	Mail        string            `json:"Mail"`
	AccessLevel GitlabAccessLevel `json:"AccessLevel"`
}

const (
	NoPermissions            GitlabAccessLevel = 0
	MinimalAccessPermissions GitlabAccessLevel = 5
	GuestPermissions         GitlabAccessLevel = 10
	ReporterPermissions      GitlabAccessLevel = 20
	DeveloperPermissions     GitlabAccessLevel = 30
	MaintainerPermissions    GitlabAccessLevel = 40
	OwnerPermissions         GitlabAccessLevel = 50
	AdminPermissions         GitlabAccessLevel = 60
)

type User struct {
	ID              string `json:"ID"`
	GitlabID        int    `json:"GitlabID"`
	GitlabUsername  string `json:"GitlabUsername"`
	DiscordID       string `json:"DiscordID"`
	DiscordUsername string `json:"DiscordUsername"`
}

type Discord struct {
	APIKey        string `json:"APIKey"`
	ApplicationID string `json:"ApplicationID"`
	GuildID       string `json:"GuildID"`
	LogChannelID  string `json:"LogChannelID"`
	Prefix        string `json:"Prefix"`
}

type Gitlab struct {
	APIKey    string `json:"APIKey"`
	ProjectID int    `json:"ProjectID"`
}

type Command struct {
	ID          string `json:"ID"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

type Store struct {
	Gitlab             Gitlab     `json:"Gitlab"`
	Discord            Discord    `json:"Discord"`
	Users              []User     `json:"Users"`
	CommandInitialized bool       `json:"CommandInitialized"`
	Commands           []*Command `json:"Commands"`
}
