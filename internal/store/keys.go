package store

const (
	VenueKey  = "venue-"
	HallKey   = "hall-"
	ZoneKey   = "zone-"
	BoothKey  = "booth-"
	GateKey   = "gate-"
	SceneKey  = "scene-"
	ScreenKey = "screen-"
	CameraKey = "camera-"
	RouteKey  = "route-"
	QuotaKey  = "quota-"
)

const (
	EffectiveSceneKey = "effective-"
	DimKey            = "dim-"
	EvacuationKey     = "evacuation-"
	HandoverKey       = "handover-"
	AuditKey          = "audit"
)

func VenueStoreKey(id string) string {
	return VenueKey + id
}

func HallStoreKey(id string) string {
	return HallKey + id
}

func ZoneStoreKey(id string) string {
	return ZoneKey + id
}

func BoothStoreKey(id string) string {
	return BoothKey + id
}

func GateStoreKey(id string) string {
	return GateKey + id
}

func SceneStoreKey(zoneID string) string {
	return SceneKey + zoneID
}

func ScreenStoreKey(id string) string {
	return ScreenKey + id
}

func CameraStoreKey(id string) string {
	return CameraKey + id
}

func RouteStoreKey(id string) string {
	return RouteKey + id
}

func QuotaStoreKey(hallID string) string {
	return QuotaKey + hallID
}

func EffectiveSceneKeyFor(zoneID string) string {
	return EffectiveSceneKey + zoneID
}

func DimKeyFor(zoneID string) string {
	return DimKey + zoneID
}

func EvacuationKeyFor(hallID string) string {
	return EvacuationKey + hallID
}

func HandoverKeyFor(screenID string) string {
	return HandoverKey + screenID
}
