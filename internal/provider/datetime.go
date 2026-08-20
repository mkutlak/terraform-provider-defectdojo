package provider

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// dateTimeLayouts lists the layouts accepted for Terraform string attributes
// that map onto a Go time.Time in the generated client, in priority order.
//
// time.RFC3339 is the canonical form the provider writes back to state; it
// covers "Z", numeric offsets, and (on parse only) fractional seconds. The
// remaining layouts exist so that ordinary Terraform idioms do not fail -
// notably formatdate("YYYY-MM-DD", ...), which produced GitHub issue #23.
// Layouts carrying no zone information are interpreted as UTC, matching
// DefectDojo's default TIME_ZONE.
//
// "2006-01-02 15:04:05" (space separator) is deliberately absent: no Terraform
// builtin emits it, and every accepted layout widens the set of literals
// preserveDateTimeLiteral has to round-trip.
var dateTimeLayouts = []string{
	time.RFC3339,          // 2006-01-02T15:04:05Z07:00
	"2006-01-02T15:04:05", // naive datetime, no offset
	"2006-01-02",          // date only -> midnight UTC
}

// dateLayout is the only layout accepted for attributes backed by an
// openapi_types.Date. See parseDate for why this stays strict.
const dateLayout = "2006-01-02"

// parseDateTime parses a Terraform string attribute into the time.Time the
// generated ddclient struct expects, accepting any layout in dateTimeLayouts.
func parseDateTime(s string) (time.Time, error) {
	for _, layout := range dateTimeLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf(
		"%q is not a recognised datetime: expected RFC3339 (e.g. 2006-01-02T15:04:05Z), "+
			"a datetime without a zone (2006-01-02T15:04:05, read as UTC), "+
			"or a date only (2006-01-02, read as midnight UTC)", s)
}

// parseDate parses a Terraform string attribute into the date the generated
// ddclient struct expects.
//
// Unlike parseDateTime this is deliberately strict. Date-only -> datetime is a
// widening conversion and is lossless, but datetime -> date is narrowing and
// genuinely ambiguous: "2026-07-28T23:30:00-05:00" is 2026-07-29 in UTC, so
// there is no single correct calendar day to pick. Silently choosing one would
// be a data-integrity hazard on required fields like engagement.target_start.
func parseDate(s string) (time.Time, error) {
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return time.Time{}, fmt.Errorf(
			"%q is not a recognised date: expected %s (date only, no time of day)", s, dateLayout)
	}

	return t, nil
}

// preserveDateTimeLiteral keeps a configured date literal that denotes the same
// instant as the one the server sent.
//
// This is what lets target_start = "2026-07-28" survive a round trip through
// DefectDojo, which echoes it back as "2026-07-28T00:00:00Z" (issue #23). It
// also covers a configured non-UTC offset, e.g. "2025-01-01T11:00:00+01:00"
// against a server that answers in UTC. See preserveLiteral.
func preserveDateTimeLiteral(current types.String, server time.Time) types.String {
	return preserveLiteral(current, server.Format(time.RFC3339), func(current, _ string) bool {
		prior, err := parseDateTime(current)
		return err == nil && prior.Equal(server)
	})
}
