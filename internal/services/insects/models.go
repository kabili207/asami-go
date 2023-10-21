package insects

type InsectSummary struct {
	Type InsectType `json:"type"`
	Name string     `json:"name"`
	Key  string     `json:"key"`
}

// Insect data
type Insect struct {
	ID                 string             `json:"id"`
	Name               string             `json:"name"`
	ScientificName     string             `json:"scientific_name"`
	Description        string             `json:"description"`
	Family             string             `json:"family,omitempty"`
	Size               string             `json:"size,omitempty"`
	Wingspan           string             `json:"wingspan,omitempty"`
	Habitat            string             `json:"habitat,omitempty"`
	CaterpillarFood    string             `json:"caterpillar_food,omitempty"`
	FlightSeason       string             `json:"flight_season,omitempty"`
	ConservationStatus ConservationStatus `json:"conservation_status,omitempty"`
	Distribution       DistributionStatus `json:"distribution,omitempty"`
	Pictures           []PictureInfo      `json:"pictures"`
}

type PictureInfo struct {
	Url         string
	Description string
	Credit      string
}

type ConservationStatus struct {
	UK_BAP  string `json:"uk_bap,omitempty"`
	General string `json:"general,omitempty"`
}

type DistributionStatus struct {
	Countries       []string `json:"countries"`
	Localities      []string `json:"localities,omitempty"`
	TrendSince1970s string   `json:"trendSince1970s,omitempty"`
	General         string   `json:"general"`
}
