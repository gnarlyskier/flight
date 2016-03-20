package metar

import (
	"net/url"
	"strings"
)

func QueryStations(stations []string, hoursBeforeNow float64, mostRecentOnly bool) ([]Metar, error) {
	parameters := url.Values{}
	parameters.Add("stationString", strings.Join(stations, ","))

	queryUrl := buildQueryUrl(parameters, hoursBeforeNow, mostRecentOnly)
	body, err := queryXml(queryUrl)
	if err != nil {
		return []Metar{}, err
	}
	return unmarshalXml(body)
}
