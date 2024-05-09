/**
 * helper functions related to maps.
 *
**/

package common

import (
	"errors"
	"lp_customer_portal/models"

	"github.com/go-chassis/openlog"
)

func GetZone(lat float64, long float64, zones []models.Zone) (models.Zone, error) {
	openlog.Debug("Getting zone for the lat and long")
	for _, zone := range zones {
		openlog.Debug("checking in zone: " + zone.Name)
		if checkIfPointInPolygon(lat, long, zone.ZoneBoundaries) {
			return zone, nil
		}
	}
	return models.Zone{}, errors.New("zone not found")
}

func checkIfPointInPolygon(lat float64, long float64, polygon []models.ZoneBoundary) bool {
	openlog.Debug("Checking if point is in polygon")
	var isInside bool
	for i, j := 0, len(polygon)-1; i < len(polygon); i++ {
		if polygon[i].Long < long && polygon[j].Long >= long || polygon[j].Long < long && polygon[i].Long >= long {
			if polygon[i].Lat+(long-polygon[i].Long)/(polygon[j].Long-polygon[i].Long)*(polygon[j].Lat-polygon[i].Lat) < lat {
				isInside = !isInside
			}
		}
		j = i
	}
	return isInside
}

func Search_ele(slice []int64, key int64) bool {
	flag := false
	for _, element := range slice {
		if element == key { // check the condition if its true return index
			flag = true
		}
	}
	return flag
}
