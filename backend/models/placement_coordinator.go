package models

type PlacementCoordinator struct {
	USN    string `json:"usn"`
	Branch string `json:"branch"`
	Batch  string `json:"batch"`
	User   Admin  `json:"admin"`
}
