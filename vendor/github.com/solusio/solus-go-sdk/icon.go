package solus

type IconsService service

type IconType string

const (
	IconTypeOS          IconType = "os"
	IconTypeApplication IconType = "application"
	IconTypeFlags       IconType = "flags"
)

type Icon struct {
	ID   int      `json:"id"`
	Name string   `json:"name"`
	URL  string   `json:"url"`
	Type IconType `json:"type"`
}
