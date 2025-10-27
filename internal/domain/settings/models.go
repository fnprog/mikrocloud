package settings

type GeneralSettings struct {
	Domain             string `json:"domain"`
	Timezone           string `json:"timezone"`
	IPv4               string `json:"ipv4"`
	IPv6               string `json:"ipv6"`
	AllowRegistrations bool   `json:"allow_registrations"`
	DoNotTrack         bool   `json:"do_not_track"`
}

type AdvancedSettings struct {
	DNSValidation bool   `json:"dns_validation"`
	DNSServers    string `json:"dns_servers"`
	APIAccess     bool   `json:"api_access"`
	AllowedIPs    string `json:"allowed_ips"`
}

type UpdateSettings struct {
	UpdateCheckFrequency string `json:"update_check_frequency"`
	AutoUpdate           bool   `json:"auto_update"`
	AutoUpdateFrequency  string `json:"auto_update_frequency"`
	AutoUpdateTime       string `json:"auto_update_time"`
}

type SMTPSettings struct {
	Enabled   bool   `json:"enabled"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FromEmail string `json:"from_email"`
	FromName  string `json:"from_name"`
}

type UpdateSMTPSettings struct {
	Enabled   bool   `json:"enabled"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FromEmail string `json:"from_email"`
	FromName  string `json:"from_name"`
}

type InstanceInfo struct {
	FQDN string `json:"fqdn"`
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`
}

type DetectedIPs struct {
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`
}
