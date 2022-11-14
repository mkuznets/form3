package models

type Resource struct {
	ID             string `json:"id,omitempty"`
	OrganisationID string `json:"organisation_id,omitempty"`
	Type           string `json:"type,omitempty"`
	Version        *int   `json:"version,omitempty"`
}
