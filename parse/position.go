package parse

import (
	"errors"
	"fmt"
	"github.com/cragcraig/flight/data"
	"github.com/cragcraig/flight/geo"
	"strconv"
	"strings"
	"unicode"
)

// Parse a position string.
// Several formats are accepted, e.g.,
// 45.42,-105.03
// KBDU
// BJC
// KBDU+5N+3W
// KBDU+8@340
// KBDU+10W+8@340
// 45.42,-105.03+5N+3W
func ParsePos(natfix data.Natfix, pos string) (geo.Coord, error) {
	if len(pos) == 0 {
		return geo.ErrCoord(), errors.New("empty position string")
	}

	// Parse station and any additional modifiers
	var position string
	var modifiers []string
	if split := strings.Split(pos, "+"); len(split) == 0 {
		// Unreachable
		panic("expected non-empty position string")
	} else {
		position = split[0]
		modifiers = split[1:len(split)]
	}

	if c, err := parseStart(natfix, position); err != nil {
		return geo.ErrCoord(), err
	} else {
		return parseAndApplyModifiers(c, modifiers)
	}
}

// A string containing a Lat,Lon or station id
func parseStart(natfix data.Natfix, pos string) (geo.Coord, error) {
	if strings.ContainsRune(pos, ',') {
		// Lon,Lat coordinate
		if c, err := geo.ParseLatLon(pos); err != nil {
			return c, err
		} else {
			return c, nil
		}
	} else if c, err := natfix.GetFix(pos); err != nil {
		return geo.ErrCoord(), err
	} else {
		// Station position
		return c, nil
	}
}

// e.g., 5N 3W 23@340
func ParseGeoVect(v string) (geo.Vect, error) {
	if strings.ContainsAny(v, "NSEWnsew") {
		// Relative positional vector
		// e.g., KBDU+50W
		return parsePosOffset(v)
	} else if strings.ContainsRune(v, '@') {
		// Relative directional vector
		// e.g., KBDU+4@340
		// Makes flat-earth assumption, not valid over large distances
		return parseDirOffset(v)
	} else {
		return geo.Vect{}, errors.New("Invalid vector: " + v)
	}
}

// Apply position relative modifiers
// e.g., +5N +3W +23@340
func parseAndApplyModifiers(c geo.Coord, modifiers []string) (geo.Coord, error) {
	for _, m := range modifiers {
		if v, err := ParseGeoVect(m); err != nil {
			return geo.ErrCoord(), err
		} else {
			c = c.AddToLon(v.X).AddToLat(v.Y)
		}
	}
	return c, nil
}

// e.g., 23N 3W etc
func parseDirOffset(s string) (geo.Vect, error) {
	var v, dir float64
	_, err := fmt.Sscanf(s, "%f@%f", &v, &dir)
	if err != nil {
		return geo.Vect{}, errors.New("Invalid directional vector: " + s)
	}
	// Real angles are in radians, have north at 90 degrees, and go counter-clockwise
	theta := geo.Compass2Rad(dir)
	return geo.HeadingFromAngle(theta).Mult(v), nil
}

// e.g., 26@340
func parsePosOffset(s string) (geo.Vect, error) {
	if len(s) < 2 {
		// Unreachable
		panic("vector too short: " + s)
	}
	v, err := strconv.ParseFloat(s[:len(s)-1], 64)
	if err != nil {
		return geo.Vect{}, errors.New("Invalid cardinal vector: " + s)
	}
	dir := unicode.ToUpper(rune(s[len(s)-1]))
	if dir == 'N' {
		return geo.Vect{0, v}, nil
	} else if dir == 'S' {
		return geo.Vect{0, -1 * v}, nil
	} else if dir == 'E' {
		return geo.Vect{v, 0}, nil
	} else if dir == 'W' {
		return geo.Vect{-1 * v, 0}, nil
	}
	return geo.Vect{}, errors.New("Invalid cardinal direction: " + s)
}
