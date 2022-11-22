package domain

var RSS = "RSS"
var WAC = "WAC"

type StudentRecord struct {
	Source string `json:"source" bson:"source,omitempty"` // RSS, WAC
	StudentRSS
	StudentWAC
}

type StudentRSS struct {
}

type StudentWAC struct {
}
