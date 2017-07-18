package sitereports

// Team struct
type Team struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Properties  []string `json:"properties"`
	Users       []string `json:"users"`
}

// User struct
type User struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Password string   `json:"password"`
	Image    string   `json:"image"`
	Email    string   `json:"email"`
	Teams    []string `json:"teams"`
}
