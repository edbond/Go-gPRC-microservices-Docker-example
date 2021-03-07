package ports

// Description of a Port structure.
type Port struct {
	Key         string    `json:"key,omitempty"`
	Name        string    `json:"name,omitempty"`
	City        string    `json:"city,omitempty"`
	Country     string    `json:"country,omitempty"`
	Alias       []string  `json:"alias,omitempty"`
	Regions     []string  `json:"regions,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty"`
	Province    string    `json:"province,omitempty"`
	Timezone    string    `json:"timezone,omitempty"`
	Unlocks     []string  `json:"unlocks,omitempty"`
	Code        string    `json:"code,omitempty"`
}


func (port *Port) ToTransport() *PortTransport {
	return &PortTransport{
		Key:         port.Key,
		Name:        port.Name,
		City:        port.City,
		Country:     port.Country,
		Alias:       port.Alias,
		Regions:     port.Regions,
		Coordinates: port.Coordinates,
		Province:    port.Province,
		Timezone:    port.Timezone,
		Unlocks:     port.Unlocks,
		Code:        port.Code,
	}
}


func (portTransport *PortTransport) ToValue() *Port {
	return &Port{
		Key:         portTransport.Key,
		Name:        portTransport.Name,
		City:        portTransport.City,
		Country:     portTransport.Country,
		Alias:       portTransport.Alias,
		Regions:     portTransport.Regions,
		Coordinates: portTransport.Coordinates,
		Province:    portTransport.Province,
		Timezone:    portTransport.Timezone,
		Unlocks:     portTransport.Unlocks,
		Code:        portTransport.Code,
	}
}