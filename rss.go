// Package rss provides helper function for interacting with RSS, RDF, and Atom
// feeds. Primarily this surrounds building and reading/parsing.
package rss

import "time"

// Feed contains information about a feed.
type Feed struct {
	Title       string
	Link        string
	Description string
	PubDate     time.Time
	Items       []Item
	Type        string
}

// Item contains information about an item/entry in a feed.
type Item struct {
	Title       string
	Link        string
	Description string
	PubDate     time.Time
	GUID        string
}

// Config controls package wide settings.
type Config struct {
	// Control whether we have verbose output (or not).
	Verbose bool
}

// Use a global default set of settings.
//
// See package log for a similar approach (global default settings).
var config = Config{
	Verbose: false,
}

// SetVerbose controls the package setting 'Verbose'.
func SetVerbose(verbose bool) {
	config.Verbose = verbose
}
