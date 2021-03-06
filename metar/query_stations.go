package metar

import (
	"errors"
	"net/url"
	"strings"
)

func QueryStations(stations []string, hoursBeforeNow float64, mostRecentOnly bool) ([]Metar, error) {
	// stations should be 3 or 4 characters
	for _, s := range stations {
		if len(s) != 3 && len(s) != 4 {
			return []Metar{}, errors.New("station ids must be valid station identifiers: " + s)
		}
	}

	parameters := url.Values{}
	parameters.Add("stationString", strings.Join(stations, ","))

	return queryMetars(parameters, hoursBeforeNow, mostRecentOnly)
}
