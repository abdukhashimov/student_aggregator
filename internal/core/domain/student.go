package domain

const (
	RSS = "RSS"
	WAC = "WAC"
)

type StudentRecord struct {
	Source string `json:"source" bson:"source"` // RSS, WAC
	Email  string `json:"email" bson:"email,omitempty"`
	Status string `json:"status" bson:"status,omitempty"`
	StudentRSS
	StudentWAC
}

type StudentWAC struct {
	JoinDate                 string   `json:"join_date" bson:"join_date,omitempty"`
	FullName                 string   `json:"full_name" bson:"full_name,omitempty"`
	Location                 string   `json:"location" bson:"location,omitempty"`
	Position                 string   `json:"position" bson:"position,omitempty"`
	Company                  string   `json:"company" bson:"company,omitempty"`
	PrefferedLanguage        []string `json:"preffered_language" bson:"preffered_language,omitempty"`
	ReceivesCommunityUpdates bool     `json:"receives_community_updates" bson:"receives_community_updates,omitempty"`
	MembershipType           string   `json:"membership_type" bson:"membership_type,omitempty"`
	AttendedEvents           int      `json:"attended_events" bson:"attended_events,omitempty"`
	RegisteredNotVisited     int      `json:"registered_not_visited" bson:"registered_not_visited,omitempty"`
	Registered               int      `json:"registered" bson:"registered,omitempty"`
}

type Project struct {
	Name       string `json:"name" bson:"name,omitempty"`
	Score      int    `json:"score" bson:"score,omitempty"`
	FinishedAt string `json:"finished_at" bson:"finished_at,omitempty"`
	Deadline   string `json:"deadline" bson:"deadline,omitempty"`
}

type StudentRSS struct {
	FirstName       string    `json:"first_name" bson:"first_name,omitempty"`
	LastName        string    `json:"last_name" bson:"last_name,omitempty"`
	StatusItems     []string  `json:"status_items" bson:"status_items,omitempty"`
	ApplicationDate string    `json:"application_date" bson:"application_date,omitempty"`
	Projects        []Project `json:"projects" bson:"projects,omitempty"`
}
