package models

type JobType struct {
	ID                  string               `json:"id" bson:"_id,omitempty"`
	Identifier          string               `json:"identifier" bson:"identifier"`
	Name                string               `json:"name" bson:"name"`
	ConfigurationFields []ConfigurationField `json:"configurationFields" bson:"configurationFields"`
}

type ConfigurationField struct {
	Name        string `json:"name" bson:"name"`
	Type        string `json:"type" bson:"type"`
	Description string `json:"description" bson:"description"`
}
