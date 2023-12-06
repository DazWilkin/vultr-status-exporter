package api

type Alert struct {
	ID        string  `json:"id"`
	Subject   string  `json:"subject"`
	Status    string  `json:"status"`
	StartDate string  `json:"start_date"`
	UpdatedAt string  `json:"updated_at"`
	Entries   []Entry `json:"entries"`
}
type Alerts struct {
	ServiceAlerts []ServiceAlert `json:"service_alerts"`
}
type Entry struct {
	UpdatedAt string `json:"updated_at"`
	Message   string `json:"message"`
}
type Region struct {
	Location    string  `json:"location"`
	Country     string  `json:"country"`
	CountryName string  `json:"country_name"`
	Alerts      []Alert `json:"alerts"`
}
type ServiceAlert struct {
	ID        string  `json:"id"`
	Region    string  `json:"region"`
	Subject   string  `json:"subject"`
	StartDate string  `json:"start_date"`
	UpdatedAt string  `json:"updated_at"`
	Status    string  `json:"status"`
	Entries   []Entry `json:"entries"`
}
type Status struct {
	ServiceAlerts []ServiceAlert    `json:"service_alerts"`
	Regions       map[string]Region `json:"regions"`
}
