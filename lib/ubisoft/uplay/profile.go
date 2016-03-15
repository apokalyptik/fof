package uplay

type Profile struct {
	UplayID        string `json:"userId"`
	ProfileID      string `json:"profileID"`
	PlatformUserID string `json:"idOnPlatform"`
	PlatformName   string `json:"nameOnPlatform"`
	Platform       string `json:"platformType"`
}
