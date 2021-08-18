package monitoring

import "github.com/ridebeam/go-common/monitor"

const (
	TagKeyNameCityID = "city_id"
)

var (
	TagKeyCityID = monitor.NewTagKey(TagKeyNameCityID, true, true)
)
