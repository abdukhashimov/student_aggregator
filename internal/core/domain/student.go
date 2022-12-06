package domain

const (
	RSS = "RSS"
	WAC = "WAC"
)

type StudentRecord struct {
	Source     string `json:"source" bson:"source"` // RSS, WAC
	Email      string `json:"email" mapstructure:"email" bson:"email,omitempty"`
	Status     string `json:"status" bson:"status,omitempty"`
	StudentRSS `mapstructure:",squash" bson:"student_rss,omitempty"`
	StudentWAC `mapstructure:",squash" bson:"student_wac,omitempty"`
}

type StudentWAC struct {
	JoinDate                 string   `json:"join_date" mapstructure:"join_date" bson:"join_date,omitempty"`
	FullName                 string   `json:"full_name" mapstructure:"full_name" bson:"full_name,omitempty"`
	Location                 string   `json:"location" mapstructure:"location" bson:"location,omitempty"`
	Position                 string   `json:"position" mapstructure:"position" bson:"position,omitempty"`
	Company                  string   `json:"company" mapstructure:"company" bson:"company,omitempty"`
	PrefferedLanguage        []string `json:"preffered_language" mapstructure:"preffered_language" bson:"preffered_language,omitempty"`
	ReceivesCommunityUpdates bool     `json:"receives_community_updates" mapstructure:"receives_community_updates" bson:"receives_community_updates,omitempty"`
	MembershipType           string   `json:"membership_type" mapstructure:"membership_type" bson:"membership_type,omitempty"`
	AttendedEvents           int      `json:"attended_events" mapstructure:"attended_events" bson:"attended_events,omitempty"`
	RegisteredNotVisited     int      `json:"registered_not_visited" mapstructure:"registered_not_visited" bson:"registered_not_visited,omitempty"`
	Registered               int      `json:"registered" mapstructure:"registered" bson:"registered,omitempty"`
}

type Project struct {
	Name       string `json:"name" mapstructure:"name" bson:"name,omitempty"`
	Score      int    `json:"score" mapstructure:"score" bson:"score,omitempty"`
	FinishedAt string `json:"finished_at" mapstructure:"finished_at" bson:"finished_at,omitempty"`
	Deadline   string `json:"deadline" mapstructure:"deadline" bson:"deadline,omitempty"`
}

type StudentRSS struct {
	FirstName       string    `json:"first_name" mapstructure:"first_name" bson:"first_name,omitempty"`
	LastName        string    `json:"last_name" mapstructure:"last_name" bson:"last_name,omitempty"`
	StatusItems     []string  `json:"status_items" mapstructure:"status_items" bson:"status_items,omitempty"`
	ApplicationDate string    `json:"application_date" mapstructure:"application_date" bson:"application_date,omitempty"`
	Projects        []Project `json:"projects" mapstructure:"projects" bson:"projects,omitempty"`
}

type ListStudentsOptions struct {
	Email  string
	Source string
	Sort   map[string]int
	Limit  int
	Skip   int
}
