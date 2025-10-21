package maintenance

type SystemStatusResponse struct {
	Body struct {
		Status     string `json:"status" example:"healthy"`
		Components struct {
			Database  string `json:"database" example:"ok"`
			Storage   string `json:"storage" example:"ok"`
			Container string `json:"container" example:"ok"`
		} `json:"components"`
		Timestamp string `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	}
}

type SystemInfoResponse struct {
	Body struct {
		Version     string `json:"version" example:"0.1.0"`
		Platform    string `json:"platform" example:"linux/amd64"`
		GoVersion   string `json:"go_version" example:"1.21"`
		BuildTime   string `json:"build_time,omitempty" example:"2024-01-01T00:00:00Z"`
		Environment string `json:"environment" example:"production"`
	}
}

type ResourcesResponse struct {
	Body struct {
		Projects     int `json:"projects" example:"5"`
		Applications int `json:"applications" example:"10"`
		Databases    int `json:"databases" example:"3"`
		Services     int `json:"services" example:"8"`
	}
}

type HealthCheckResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

type DomainListResponse struct {
	Domains []DomainInfo `json:"domains"`
	Total   int          `json:"total"`
}

type DomainInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Verified    bool   `json:"verified"`
	SSLEnabled  bool   `json:"ssl_enabled"`
	SSLExpiry   string `json:"ssl_expiry,omitempty"`
	ServiceID   string `json:"service_id,omitempty"`
	ServiceName string `json:"service_name,omitempty"`
	CreatedAt   string `json:"created_at"`
}

type AddDomainRequest struct {
	Name      string `json:"name"`
	ServiceID string `json:"service_id,omitempty"`
}

type EnableSSLRequest struct {
	Provider string `json:"provider"`
	Email    string `json:"email,omitempty"`
}
