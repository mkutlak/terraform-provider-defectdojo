package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"gotest.tools/assert"
)

func TestParseDateTime(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantRFC    string
		wantUTCLoc bool
	}{
		{
			name:    "RFC3339 UTC",
			input:   "2026-07-28T12:34:56Z",
			wantRFC: "2026-07-28T12:34:56Z",
		},
		{
			name:    "RFC3339 with offset",
			input:   "2026-07-28T12:34:56+02:00",
			wantRFC: "2026-07-28T10:34:56Z",
		},
		{
			// Go's RFC3339 parsing accepts a fractional-second field even
			// though the layout itself omits one.
			name:    "RFC3339 with fractional seconds",
			input:   "2026-07-28T12:34:56.123456Z",
			wantRFC: "2026-07-28T12:34:56Z",
		},
		{
			name:       "zoneless datetime read as UTC",
			input:      "2026-07-28T12:34:56",
			wantRFC:    "2026-07-28T12:34:56Z",
			wantUTCLoc: true,
		},
		{
			// Date only -> midnight UTC. This is the GitHub issue #23 case:
			// formatdate("YYYY-MM-DD", ...) must parse successfully.
			name:       "date only is midnight UTC (issue #23)",
			input:      "2026-07-28",
			wantRFC:    "2026-07-28T00:00:00Z",
			wantUTCLoc: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDateTime(tt.input)
			assert.NilError(t, err)
			assert.Equal(t, got.UTC().Format(time.RFC3339), tt.wantRFC)
			if tt.wantUTCLoc {
				assert.Equal(t, got.Location(), time.UTC)
			}
		})
	}
}

func TestParseDateTimeRejectsUnknownLayouts(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "empty string", input: ""},
		{name: "not a date", input: "yesterday"},
		{name: "wrong field order", input: "28/07/2026"},
		{
			// Deliberately unsupported - see the doc comment on
			// dateTimeLayouts.
			name:  "space separator",
			input: "2026-07-28 12:34:56",
		},
		{name: "invalid month", input: "2026-13-01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseDateTime(tt.input)
			assert.ErrorContains(t, err, fmt.Sprintf("%q", tt.input))
		})
	}
}

func TestParseDate(t *testing.T) {
	got, err := parseDate("2026-07-28")
	assert.NilError(t, err)
	assert.Assert(t, got.Equal(time.Date(2026, 7, 28, 0, 0, 0, 0, time.UTC)))
}

func TestParseDateRejectsDatetime(t *testing.T) {
	// Datetime -> date is a narrowing, ambiguous conversion:
	// "2026-07-28T23:30:00-05:00" is 2026-07-29 in UTC, so a required field
	// like engagement.target_start must not silently guess a calendar day.
	// This test is an executable record of that design decision.
	_, err := parseDate("2026-07-28T00:00:00Z")
	assert.ErrorContains(t, err, `"2026-07-28T00:00:00Z"`)
	assert.ErrorContains(t, err, "2006-01-02")
}

func TestPreserveDateTimeLiteral(t *testing.T) {
	t.Parallel()

	// The server column is a time.Time; the table spells it as the RFC3339
	// literal the API answers with so the rows stay readable.
	runPreserveCases(t, func(current types.String, server string) types.String {
		at, err := time.Parse(time.RFC3339, server)
		assert.NilError(t, err)
		return preserveDateTimeLiteral(current, at)
	}, []preserveCase{
		// Issue #23 round trip: config wrote a date-only literal, the server
		// echoes it back as midnight UTC - keep the practitioner's literal.
		{"date-only literal round-trips (issue #23)", types.StringValue("2026-07-28"), "2026-07-28T00:00:00Z", "2026-07-28"},
		// Same instant, different offset - the latent bug this also fixes.
		{"same instant with non-UTC offset is preserved", types.StringValue("2025-01-01T11:00:00+01:00"), "2025-01-01T10:00:00Z", "2025-01-01T11:00:00+01:00"},
		// Genuinely different instant - real drift is still reported.
		{"different instant is not preserved", types.StringValue("2026-07-28"), "2026-07-29T00:00:00Z", "2026-07-29T00:00:00Z"},
		{"unparsable literal returns canonical", types.StringValue("garbage"), "2026-07-28T00:00:00Z", "2026-07-28T00:00:00Z"},
	})
}
