package domain

type Venue struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	ZoneIDs []string `json:"zone_ids"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Hall struct {
	ID      string           `json:"id"`
	VenueID string           `json:"venue_id"`
	Name    string           `json:"name"`
	ZoneIDs []string         `json:"zone_ids"`
	Layout  map[string]Point `json:"layout"`
}

type Zone struct {
	ID        string `json:"id"`
	HallID    string `json:"hall_id"`
	VenueID   string `json:"venue_id"`
	Name      string `json:"name"`
	Partition string `json:"partition"`
}

type Booth struct {
	ID        string `json:"id"`
	HallID    string `json:"hall_id"`
	ZoneID    string `json:"zone_id"`
	Name      string `json:"name"`
	Powered   bool   `json:"powered"`
	LoadWatts int    `json:"load_watts"`
}

type Gate struct {
	ID       string `json:"id"`
	ZoneID   string `json:"zone_id"`
	Name     string `json:"name"`
	Mode     string `json:"mode"`
	Released bool   `json:"released"`
}

type Scene struct {
	ZoneID      string `json:"zone_id"`
	Name        string `json:"name"`
	Brightness  int    `json:"brightness"`
	Color       string `json:"color"`
	IsEmergency bool   `json:"is_emergency"`
}

type Evacuation struct {
	HallID string `json:"hall_id"`
	Active bool   `json:"active"`
	At     int64  `json:"at"`
}

type Clip struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Sponsor     string `json:"sponsor"`
	DurationSec int    `json:"duration_sec"`
}

type Screen struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Handover struct {
	ScreenID  string `json:"screen_id"`
	Emergency bool   `json:"emergency"`
	At        int64  `json:"at"`
}

type Camera struct {
	ID       string            `json:"id"`
	HallID   string            `json:"hall_id"`
	Name     string            `json:"name"`
	Presets  map[string]Preset `json:"presets"`
	DayNight string            `json:"day_night"`
}

type Preset struct {
	BoothID string `json:"booth_id"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	Pan     int    `json:"pan"`
	Tilt    int    `json:"tilt"`
}

type Checkpoint struct {
	BoothID string `json:"booth_id"`
	ZoneID  string `json:"zone_id"`
	Order   int    `json:"order"`
	Visited bool   `json:"visited"`
}

type Route struct {
	ID          string       `json:"id"`
	ZoneID      string       `json:"zone_id"`
	Name        string       `json:"name"`
	Checkpoints []Checkpoint `json:"checkpoints"`
}

type Quota struct {
	ID         string `json:"id"`
	HallID     string `json:"hall_id"`
	LimitWatts int    `json:"limit_watts"`
	UsedWatts  int    `json:"used_watts"`
}

type AuditRecord struct {
	ID     string `json:"id"`
	Action string `json:"action"`
	Target string `json:"target"`
	Detail string `json:"detail"`
	At     int64  `json:"at"`
}
