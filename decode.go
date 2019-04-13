package rss

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

// rssXML is used for parsing/encoding RSS.
type rssXML struct {
	// If xml.Name is specified and has a tag name, we must have this element as
	// the root. I don't do this though because it is case sensitive. Instead,
	// inspect XMLName manually afterwards.
	XMLName xml.Name
	Channel rssChannelXML `xml:"channel"`
	Version string        `xml:"version,attr"`
}

// rssChannelXML is used for parsing/encoding RSS.
type rssChannelXML struct {
	XMLName     xml.Name     `xml:"channel"`
	Title       string       `xml:"title"`
	Link        string       `xml:"link"`
	Description string       `xml:"description"`
	PubDate     string       `xml:"pubDate"`
	Items       []rssItemXML `xml:"item"`
}

// rssItemXML is used for parsing/encoding RSS.
type rssItemXML struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
	// GUID is optional. Unique identifier.
	GUID string `xml:"guid"`
}

// rdfXML is used for parsing RDF.
type rdfXML struct {
	// Element name. Don't specify here so we can check case insensitively.
	XMLName xml.Name

	Channel rdfChannelXML `xml:"channel"`

	RDFItems []rdfItemXML `xml:"item"`
}

// rdfChannelXML is part of parsing RDF.
type rdfChannelXML struct {
	XMLName     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`
	Links       []string `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"date"`
}

// rdfItemXML is used for parsing <rdf> item XML.
type rdfItemXML struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"date"`
	// RDF doesn't have a unique identifier like guid/id? Or maybe it does, but
	// the only feed I have using RDF doesn't use it, so I'm not looking too hard!
}

// atomXML describes an Atom feed. We use it for parsing. See
// https://tools.ietf.org/html/rfc4287
type atomXML struct {
	// The element name. Enforce it is atom:feed
	XMLName xml.Name `xml:"http://www.w3.org/2005/Atom feed"`

	// Title is human readable. It must be present.
	Title string `xml:"title"`

	// Web resource. Zero or more. Feeds should contain with with rel=self.
	Links []atomLink `xml:"link"`

	// Last time feed was updated.
	Updated string `xml:"updated"`

	Items []atomItemXML `xml:"entry"`
}

// atomLink describes a <link> element.
type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
}

// atomItemXML describes an item/entry in the feed. Atom calls these entries,
// but for consistency with other formats I support, I call them items.
type atomItemXML struct {
	// Human readable title. Must be present.
	Title string `xml:"title"`

	// Web resource. Zero or more.
	Links []atomLink `xml:"link"`

	// Last time entry updated. Must be present.
	Updated string `xml:"updated"`

	// Content is optional.
	Content string `xml:"content"`

	// ID is required. Unique identifier.
	ID string `xml:"id"`
}

// ParseFeedXML takes a feed's raw XML and returns a struct describing the feed.
//
// We support various formats: RSS, RDF, Atom. We try our best to decode the
// feed in one of them.
func ParseFeedXML(data []byte) (*Feed, error) {
	channelRSS, errRSS := parseAsRSS(data)
	if errRSS == nil {
		return channelRSS, nil
	}

	channelRDF, errRDF := parseAsRDF(data)
	if errRDF == nil {
		return channelRDF, nil
	}

	channelAtom, errAtom := parseAsAtom(data)
	if errAtom == nil {
		return channelAtom, nil
	}

	return nil, fmt.Errorf("unable to parse as RSS (%s), RDF (%s), or Atom (%s)",
		errRSS, errRDF, errAtom)
}

// parseAsRSS attempts to parse the buffer as if it contains an RSS feed.
func parseAsRSS(data []byte) (*Feed, error) {
	rssXML := rssXML{}
	if err := newDecoder(data).Decode(&rssXML); err != nil {
		return nil, fmt.Errorf("RSS XML decode error: %v", err)
	}

	if strings.ToLower(rssXML.XMLName.Local) != "rss" {
		return nil, errors.New("base tag is not RSS")
	}

	// Build a channel struct now. It's common to the base formats we support.

	feed := &Feed{
		Title:       rssXML.Channel.Title,
		Link:        rssXML.Channel.Link,
		Description: rssXML.Channel.Description,
		PubDate:     parseTime(rssXML.Channel.PubDate),
		Type:        "RSS",
	}

	if config.Verbose {
		log.Printf("Parsed channel as RSS [%s]", feed.Title)
	}

	for _, item := range rssXML.Channel.Items {
		feed.Items = append(feed.Items,
			Item{
				Title:       item.Title,
				Link:        item.Link,
				Description: item.Description,
				PubDate:     parseTime(item.PubDate),
				GUID:        item.GUID,
			})
	}

	return feed, nil
}

func newDecoder(data []byte) *xml.Decoder {
	d := xml.NewDecoder(bytes.NewBuffer(data))
	d.CharsetReader = charset.NewReaderLabel
	return d
}

// parseAsRDF attempts to parse the buffer as if it contains an RDF feed.
//
// See parseAsRSS() for a similar function, but for RSS.
func parseAsRDF(data []byte) (*Feed, error) {
	rdfXML := rdfXML{}
	if err := newDecoder(data).Decode(&rdfXML); err != nil {
		return nil, fmt.Errorf("RDF XML decode error: %v", err)
	}

	if strings.ToLower(rdfXML.XMLName.Local) != "rdf" {
		return nil, errors.New("base tag is not RDF")
	}

	link := ""
	if len(rdfXML.Channel.Links) > 0 {
		link = rdfXML.Channel.Links[0]
	}

	feed := &Feed{
		Title:       rdfXML.Channel.Title,
		Link:        link,
		Description: rdfXML.Channel.Description,
		PubDate:     parseTime(rdfXML.Channel.PubDate),
		Type:        "RDF",
	}

	if config.Verbose {
		log.Printf("Parsed channel as RDF [%s]", feed.Title)
	}

	for _, item := range rdfXML.RDFItems {
		feed.Items = append(feed.Items,
			Item{
				Title:       item.Title,
				Link:        item.Link,
				Description: item.Description,
				PubDate:     parseTime(item.PubDate),
			})
	}

	return feed, nil
}

// parseAsAtom attempts to parse the buffer as Atom.
//
// See parseAsRSS() and parseAsRDF() for similar parsing. Also I omit comments
// that would be repeated here if they are in those functions.
func parseAsAtom(data []byte) (*Feed, error) {
	atomXML := atomXML{}
	if err := newDecoder(data).Decode(&atomXML); err != nil {
		return nil, fmt.Errorf("Atom XML decode error: %v", err)
	}

	// May have multiple <link> elements. Look for rel=self.
	link := ""
	for _, l := range atomXML.Links {
		if l.Rel != "self" {
			continue
		}
		link = l.Href
		break
	}

	feed := &Feed{
		Title:   atomXML.Title,
		Link:    link,
		PubDate: parseTime(atomXML.Updated),
		Type:    "Atom",
	}

	if config.Verbose {
		log.Printf("Parsed channel as Atom [%s]", feed.Title)
	}

	for _, item := range atomXML.Items {
		link := ""
		// Take the first. Probably we can be more intelligent.
		if len(item.Links) > 0 {
			link = item.Links[0].Href
		}

		feed.Items = append(feed.Items, Item{
			Title:       item.Title,
			Link:        link,
			Description: item.Content,
			PubDate:     parseTime(item.Updated),
			GUID:        item.ID,
		})
	}

	return feed, nil
}

// parseTime tries to retrieve a publication date for the item.
//
// We try parsing using multiple formats, and fall back to a default of the
// current time if none succeed.
func parseTime(pubDate string) time.Time {
	if len(pubDate) == 0 {
		if config.Verbose {
			log.Print("No publication date on channel/item. Defaulting to now.")
		}
		return time.Now()
	}

	pubDate = strings.TrimSpace(pubDate)

	// Use RFC1123 time format for parsing. This appears to be what is present in
	// the Slashdot feed, though I expect this could vary in other feed
	// sources...
	//
	// Slashdot's feed: Sat, 29 Jun 2013 18:20:00 GMT
	pubDateTimeParsed, err := time.Parse(time.RFC1123, pubDate)
	// We use the parsed time only if we had no errors parsing it.
	if err == nil {
		return pubDateTimeParsed
	}

	// Try another format.
	//
	// Torrentfreak RSS feed format:
	//
	// Sun, 30 Jun 2013 21:26:26 +0000
	//
	// Mon, 10 Jun 2013 21:04:57 +0000
	pubDateTimeParsed, err = time.Parse(time.RFC1123Z, pubDate)
	// We use the parsed time only if we had no errors parsing it.
	if err == nil {
		return pubDateTimeParsed
	}

	// Slashdot RDF format: 2015-03-03T21:29:00+00:00
	//
	// NOTE: RFC3339 is not exactly this it seems?
	pubDateTimeParsed, err = time.Parse(time.RFC3339, pubDate)
	if err == nil {
		return pubDateTimeParsed
	}

	// commandlinefanatic.com: Wednesday, October 28, 2015, 15:24 -0700
	layout0 := "Monday, January _2, 2006, 15:04 -0700"
	pubDateTimeParsed, err = time.Parse(layout0, pubDate)
	if err == nil {
		return pubDateTimeParsed
	}

	// commandlinefanatic.com: Sat, May 7 2011 12:56:00 +0000
	layout1 := "Mon, Jan 2 2006 15:04:05 -0700"
	pubDateTimeParsed, err = time.Parse(layout1, pubDate)
	if err == nil {
		return pubDateTimeParsed
	}

	// commandlinefanatic.com: Tuesday, August 14, 2012 20:45 +0000
	layout2 := "Monday, January _2, 2006 15:04 -0700"
	pubDateTimeParsed, err = time.Parse(layout2, pubDate)
	if err == nil {
		return pubDateTimeParsed
	}

	// commandlinefanatic.com: Monday, September 5 2011 13:54:00 +0000
	layout3 := "Monday, January _2 2006 15:04:05 -0700"
	pubDateTimeParsed, err = time.Parse(layout3, pubDate)
	if err == nil {
		return pubDateTimeParsed
	}

	// commandlinefanatic.com: Wed, July 18 2011 09:59:00 +0000
	layout4 := "Mon, January _2 2006 15:04:05 -0700"
	pubDateTimeParsed, err = time.Parse(layout4, pubDate)
	if err == nil {
		return pubDateTimeParsed
	}

	// yarchive.net: Sun, 09 Apr 2017 05:06 GMT
	yarchive := "Mon, _2 Jan 2006 15:04 MST"
	pubDateTimeParsed, err = time.Parse(yarchive, pubDate)
	if err == nil {
		return pubDateTimeParsed
	}

	log.Printf("No format worked for date [%s]. Defaulting to now.", pubDate)

	return time.Now()
}
