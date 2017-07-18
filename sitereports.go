package sitereports

// App struct
type App struct {
	Properties []Property `json:"properties"`
	Sites      []Site     `json:"sites"`
	Reports    []Report   `json:"reports"`
}
