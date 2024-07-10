package user

/*
Example response

    "properties": {
        "tags": {
            "test": "test"
        },
        "language": "en",
        "timezone_id": "America/New_York",
        "country": "US",
        "first_active": 1719328695,
        "last_active": 1719328929,
        "ip": "73.33.101.85"
    },
*/

// Properties represents the properties of a user
// Properties represents the properties of a user
type Properties struct {
	Tags        map[string]string `json:"tags"`
	Language    string            `json:"language"`
	TimezoneID  string            `json:"timezone_id"`
	Country     string            `json:"country"`
	FirstActive int64             `json:"first_active"`
	LastActive  int64             `json:"last_active"`
	IP          string            `json:"ip"`
}

// UserResponse represents the API response for a user
type User struct {
	Properties Properties `json:"properties"`
}
