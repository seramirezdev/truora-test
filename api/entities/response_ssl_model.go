package entities

type Endpoint struct {
	IpAddress     string `json:"ipAddress,omitempty"`
	ServerName    string `json:"serverName,omitempty"`
	StatusMessage string `json:"statusMessage,omitempty"`
	Grade         string `json:"grade,omitempty"`
}

type DataDomain struct {
	Status        string     `json:"status,omitempty"`
	EndPoints     []Endpoint `json:"endpoints,omitempty"`
}
