package rss

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestParseAsRDF(t *testing.T) {
	tests := []struct {
		input   string
		output  *Feed
		success bool
	}{
		{
			// An edited/subset version of a feed from Slashdot.
			`<?xml version="1.0" encoding="ISO-8859-1"?>
<?xml-stylesheet type="text/xsl" media="screen" href="/~d/styles/rss1full.xsl"?><?xml-stylesheet type="text/css" media="screen" href="http://rss.slashdot.org/~d/styles/itemcontent.css"?><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns="http://purl.org/rss/1.0/" xmlns:slash="http://purl.org/rss/1.0/modules/slash/" xmlns:content="http://purl.org/rss/1.0/modules/content/" xmlns:taxo="http://purl.org/rss/1.0/modules/taxonomy/" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:syn="http://purl.org/rss/1.0/modules/syndication/" xmlns:admin="http://webns.net/mvcb/">

<channel rdf:about="https://slashdot.org/">
<title>Slashdot</title>
<link>https://slashdot.org/</link>
<description>News for nerds, stuff that matters</description>
<dc:language>en-us</dc:language>
<dc:rights>Copyright 1997-2016, SlashdotMedia. All Rights Reserved.</dc:rights>
<dc:date>2017-01-17T21:30:14+00:00</dc:date>
<dc:publisher>Dice</dc:publisher>
<dc:creator>help@slashdot.org</dc:creator>
<dc:subject>Technology</dc:subject>
<syn:updateBase>1970-01-01T00:00+00:00</syn:updateBase>
<syn:updateFrequency>1</syn:updateFrequency>
<syn:updatePeriod>hourly</syn:updatePeriod>
<items>
 <rdf:Seq>
  <rdf:li rdf:resource="https://tech.slashdot.org/story/17/01/17/197230/uber-sues-city-of-seattle-to-block-landmark-driver-union-ordinance?utm_source=rss1.0mainlinkanon&amp;utm_medium=feed" />
  <rdf:li rdf:resource="https://entertainment.slashdot.org/story/17/01/17/1855219/netflix-is-killing-dvd-sales-research-finds?utm_source=rss1.0mainlinkanon&amp;utm_medium=feed" />
 </rdf:Seq>
</items>

<image rdf:resource="http://a.fsdn.com/sd/topics/topicslashdot.gif" />
<textinput rdf:resource="https://slashdot.org/search.pl" />
<atom10:link xmlns:atom10="http://www.w3.org/2005/Atom" rel="self" type="application/rdf+xml" href="http://rss.slashdot.org/slashdot/slashdotMain" /><feedburner:info xmlns:feedburner="http://rssnamespace.org/feedburner/ext/1.0" uri="slashdot/slashdotmain" /><atom10:link xmlns:atom10="http://www.w3.org/2005/Atom" rel="hub" href="http://pubsubhubbub.appspot.com/" />
</channel>

<image rdf:about="http://a.fsdn.com/sd/topics/topicslashdot.gif">
<title>Slashdot</title>
<url>http://a.fsdn.com/sd/topics/topicslashdot.gif</url>
<link>https://slashdot.org/</link>
</image>

<item rdf:about="https://tech.slashdot.org/story/17/01/17/197230/uber-sues-city-of-seattle-to-block-landmark-driver-union-ordinance?utm_source=rss1.0mainlinkanon&amp;utm_medium=feed">
<title>Uber Sues City of Seattle To Block Landmark Driver Union Ordinance</title>
<link>https://tech.slashdot.org/story/17/01/17/197230/uber-sues-city-of-seattle-to-block-landmark-driver-union-ordinance?utm_source=rss1.0mainlinkanon&amp;utm_medium=feed</link>
<description>Seattle's landmark law that lets drivers</description>
<dc:creator>msmash</dc:creator>
<dc:date>2017-01-17T20:40:00+00:00</dc:date>
<dc:subject>transportation</dc:subject>
<slash:department>tussle-continues</slash:department>
<slash:section>technology</slash:section>
<slash:comments>42</slash:comments>
<slash:hit_parade>42,42,27,22,3,0,0</slash:hit_parade>
</item>
<item rdf:about="https://entertainment.slashdot.org/story/17/01/17/1855219/netflix-is-killing-dvd-sales-research-finds?utm_source=rss1.0mainlinkanon&amp;utm_medium=feed">
<title>Netflix is 'Killing' DVD Sales, Research Finds</title>
<link>https://entertainment.slashdot.org/story/17/01/17/1855219/netflix-is-killing-dvd-sales-research-finds?utm_source=rss1.0mainlinkanon&amp;utm_medium=feed</link>
<description>Netflix has become the go-to destination for many movie</description>
<dc:creator>msmash</dc:creator>
<dc:date>2017-01-17T20:00:00+00:00</dc:date>
<dc:subject>movies</dc:subject>
<slash:department>how-things-work</slash:department>
<slash:section>entertainment</slash:section>
<slash:comments>101</slash:comments>
<slash:hit_parade>101,100,66,55,17,8,2</slash:hit_parade>
</item>
</rdf:RDF>
`,
			&Feed{
				Title:       "Slashdot",
				Link:        "https://slashdot.org/",
				Description: "News for nerds, stuff that matters",
				PubDate:     time.Date(2017, 1, 17, 21, 30, 14, 0, time.FixedZone("TZ", 0)),
				Items: []Item{
					Item{
						Title:       "Uber Sues City of Seattle To Block Landmark Driver Union Ordinance",
						Link:        "https://tech.slashdot.org/story/17/01/17/197230/uber-sues-city-of-seattle-to-block-landmark-driver-union-ordinance?utm_source=rss1.0mainlinkanon&utm_medium=feed",
						Description: "Seattle's landmark law that lets drivers",
						PubDate: time.Date(2017, 1, 17, 20, 40, 0, 0,
							time.FixedZone("TZ", 0)),
					},
					Item{
						Title:       "Netflix is 'Killing' DVD Sales, Research Finds",
						Link:        "https://entertainment.slashdot.org/story/17/01/17/1855219/netflix-is-killing-dvd-sales-research-finds?utm_source=rss1.0mainlinkanon&utm_medium=feed",
						Description: "Netflix has become the go-to destination for many movie",
						PubDate: time.Date(2017, 1, 17, 20, 0, 0, 0,
							time.FixedZone("TZ", 0)),
					},
				},
				Type: "RDF",
			},
			true,
		},
	}

	for _, test := range tests {
		feed, err := ParseFeedXML([]byte(test.input))
		if err != nil {
			if !test.success {
				continue
			}

			t.Errorf("parseAsAtom(%s) = error %s, wanted success", test.input, err)
			continue
		}

		if !test.success {
			t.Errorf("parseAsAtom(%s) = success, wanted error", test.input)
			continue
		}

		err = feedEqual(feed, test.output)
		if err != nil {
			t.Errorf("parseAsAtom(%s): %s", test.input, err)
			t.Errorf("Got:    %#v", feed)
			t.Errorf("Wanted: %#v", test.output)
			continue
		}
	}
}

func TestParseAsAtom(t *testing.T) {
	tests := []struct {
		input   string
		output  *Feed
		success bool
	}{
		{
			`<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">

 <title>Test one two</title>
 <link href="http://www.example.com/atom.xml" rel="self"/>
 <link href="http://www.example.com"/>
 <updated>2017-01-11T20:30:23-05:00</updated>
 <id>http://www.example.com-id</id>
 <author>
   <name>John Q. Public</name>
   <email>john@example.com</email>
 </author>

 <entry>
   <title>Test title 1</title>
   <link href="http://www.example.com/test-entry-1"/>
   <updated>2017-01-11T00:00:00-05:00</updated>
   <id>http://www.example.com/test-entry-1-id</id>
   <content type="html">&lt;p&gt;Testing content 1&lt;/p&gt;</content>
</entry>

 <entry>
   <title>Test title 2</title>
   <link href="http://www.example.com/test-entry-2"/>
   <updated>2017-01-12T00:00:00-05:00</updated>
   <id>http://www.example.com/test-entry-2-id</id>
   <content type="html">&lt;p&gt;Testing content 2&lt;/p&gt;</content>
</entry>
</feed>
`,
			&Feed{
				Title:       "Test one two",
				Link:        "http://www.example.com/atom.xml",
				Description: "",
				PubDate: time.Date(2017, 1, 11, 20, 30, 23, 0,
					time.FixedZone("TZ", -5*60*60)),
				Items: []Item{
					Item{
						Title:       "Test title 1",
						Link:        "http://www.example.com/test-entry-1",
						Description: "<p>Testing content 1</p>",
						PubDate: time.Date(2017, 1, 11, 0, 0, 0, 0,
							time.FixedZone("TZ", -5*60*60)),
						GUID: "http://www.example.com/test-entry-1-id",
					},
					Item{
						Title:       "Test title 2",
						Link:        "http://www.example.com/test-entry-2",
						Description: "<p>Testing content 2</p>",
						PubDate: time.Date(2017, 1, 12, 0, 0, 0, 0,
							time.FixedZone("TZ", -5*60*60)),
						GUID: "http://www.example.com/test-entry-2-id",
					},
				},
				Type: "Atom",
			},
			true,
		},
	}

	for _, test := range tests {
		feed, err := parseAsAtom([]byte(test.input))
		if err != nil {
			if !test.success {
				continue
			}

			t.Errorf("parseAsAtom(%s) = error %s, wanted success", test.input, err)
			continue
		}

		if !test.success {
			t.Errorf("parseAsAtom(%s) = success, wanted error", test.input)
			continue
		}

		err = feedEqual(feed, test.output)
		if err != nil {
			t.Errorf("parseAsAtom(%s): %s", test.input, err)
			t.Errorf("Got:    %#v", feed)
			t.Errorf("Wanted: %#v", test.output)
			continue
		}
	}
}

func feedEqual(a, b *Feed) error {
	if a.Title != b.Title {
		return fmt.Errorf("feed title mismatch")
	}

	if a.Link != b.Link {
		return fmt.Errorf("feed link mismatch")
	}

	if a.Description != b.Description {
		return fmt.Errorf("feed description mismatch")
	}

	if !a.PubDate.Equal(b.PubDate) {
		return fmt.Errorf("feed pubdate mismatch")
	}

	if len(a.Items) != len(b.Items) {
		return fmt.Errorf("feed items count mismatch")
	}

	for i := range a.Items {
		ai := a.Items[i]
		bi := b.Items[i]

		if ai.Title != bi.Title {
			return fmt.Errorf("item %d title mismatch", i)
		}

		if ai.Link != bi.Link {
			return fmt.Errorf("item %d link mismatch", i)
		}

		if ai.Description != bi.Description {
			return fmt.Errorf("item %d description mismatch", i)
		}

		if !ai.PubDate.Equal(bi.PubDate) {
			return fmt.Errorf("item %d pubdate mismatch", i)
		}
	}

	if a.Type != b.Type {
		return fmt.Errorf("feed type mismatch, %s vs. %s", a.Type, b.Type)
	}

	return nil
}

func TestMakeXML(t *testing.T) {
	tests := []struct {
		input   Feed
		output  string
		success bool
	}{
		{
			Feed{
				Title:       "Test feed",
				Link:        "https://www.example.com/",
				Description: "A nice feed",
				PubDate: time.Date(2016, 12, 25, 11, 0, 0, 0,
					time.FixedZone("TZ", 0)),
				Items: []Item{
					Item{
						Title:       "Nice item 1",
						Link:        "https://www.example.com/1",
						Description: "Item 1 is very nice",
						PubDate: time.Date(2016, 12, 25, 11, 01, 0, 0,
							time.FixedZone("TZ", 0)),
					},
					Item{
						Title:       "Nice item 2",
						Link:        "https://www.example.com/2",
						Description: "Item 2 is very nice",
						PubDate: time.Date(2016, 12, 25, 10, 01, 0, 0,
							time.FixedZone("TZ", 0)),
					},
				},
			},
			`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Test feed</title>
    <link>https://www.example.com/</link>
    <description>A nice feed</description>
    <pubDate>Sun, 25 Dec 2016 11:00:00 +0000</pubDate>
    <lastBuildDate>Sun, 25 Dec 2016 11:00:00 +0000</lastBuildDate>
    <item>
      <title>Nice item 1</title>
      <link>https://www.example.com/1</link>
      <description>Item 1 is very nice</description>
      <pubDate>Sun, 25 Dec 2016 11:01:00 +0000</pubDate>
      <guid>https://www.example.com/1</guid>
    </item>
    <item>
      <title>Nice item 2</title>
      <link>https://www.example.com/2</link>
      <description>Item 2 is very nice</description>
      <pubDate>Sun, 25 Dec 2016 10:01:00 +0000</pubDate>
      <guid>https://www.example.com/2</guid>
    </item>
  </channel>
</rss>`,
			true,
		},
	}

	for _, test := range tests {
		buf, err := makeXML(test.input)
		if err != nil {
			if !test.success {
				continue
			}

			t.Errorf("makeXML(%#v) = error %s", test.input, err)
			continue
		}

		if !test.success {
			t.Errorf("makeXML(%#v) = success, wanted error", test.input)
			continue
		}

		if !bytes.Equal(buf, []byte(test.output)) {
			t.Errorf("makeXML(%#v) = %s, wanted %s", test.input, buf, test.output)
			continue
		}
	}
}

func TestGetEncodingName(t *testing.T) {
	tests := []struct {
		input     []byte
		output    string
		wantError error
	}{
		{
			input:     []byte(`<?xml version="1.0" encoding="UTF-8" ?>`),
			output:    "UTF-8",
			wantError: nil,
		},
		{
			input:     []byte(`<?xml version="1.0" encoding='UTF-8' ?>`),
			output:    "UTF-8",
			wantError: nil,
		},
		{
			input:     []byte(`<?xml version="1.0" encoding="UTF-8"?>`),
			output:    "UTF-8",
			wantError: nil,
		},
		{
			input:     []byte(`<?xml version=`),
			output:    "",
			wantError: fmt.Errorf("buffer does not have XML header, or header is malformed"),
		},
		{
			// Note 'xersion' instead of 'version'
			input:     []byte(`<?xml xersion="1.0" encoding="UTF-8"?>`),
			output:    "",
			wantError: fmt.Errorf("buffer does not have XML header, or header is malformed"),
		},
		{
			input:     []byte(`<?xml version="1.0"`),
			output:    "",
			wantError: fmt.Errorf("buffer does not have XML header, or header is malformed"),
		},
		{
			input:     []byte(`<?xml version="1.0" encoding=UTF-8`),
			output:    "",
			wantError: fmt.Errorf("buffer does not have XML header, or header is malformed"),
		},
		{
			input:     []byte(`<?xml version="1.0" encoding=`),
			output:    "",
			wantError: fmt.Errorf("buffer does not have XML header, or header is malformed"),
		},
		{
			input:     []byte(`<?xml version="1.0" encoding= ?>`),
			output:    "",
			wantError: fmt.Errorf("buffer does not have XML header, or header is malformed"),
		},
		{
			input:     []byte(`<?xml version="1.0" encoding="UTF-8' ?>`),
			output:    "",
			wantError: fmt.Errorf("buffer does not have XML header, or header is malformed"),
		},
		{
			input:     []byte(`<?xml version="1.0" encoding="UTF-8 ?>`),
			output:    "",
			wantError: fmt.Errorf("buffer does not have XML header, or header is malformed"),
		},
		{
			input:     []byte(`<?xml version="1.0" encoding="" ?>`),
			output:    "",
			wantError: fmt.Errorf("buffer does not have XML header, or header is malformed"),
		},
		{
			input:     []byte(`<?xml version="1.0"?>`),
			output:    "UTF-8",
			wantError: nil,
		},
		{
			input:     []byte(`<?xml version="1.0" ?>`),
			output:    "UTF-8",
			wantError: nil,
		},
	}

	for _, test := range tests {
		name, err := getEncodingName(test.input)

		if !errorsEqual(err, test.wantError) {
			t.Errorf("getEncodingName(%s) error = %s, wanted %s", test.input,
				err, test.wantError)
			continue
		}

		if name != test.output {
			t.Errorf("getEncodingName(%s) name = %s, wanted %s", test.input, name,
				test.output)
			continue
		}
	}
}

func TestConvertToUTF8(t *testing.T) {
	tests := []struct {
		input         []byte
		encodingLabel string
		output        []byte
		wantError     error
	}{
		{
			input:         []byte(`<?xml version="1.0" encoding="UTF-8" ?> 😀`),
			encodingLabel: "UTF-8",
			output:        []byte(`<?xml version="1.0" encoding="UTF-8" ?> 😀`),
			wantError:     nil,
		},
		{
			input:         []byte(`<?xml version="1.0" encoding="ISO-8859-1" ?> :)`),
			encodingLabel: "ISO-8859-1",
			output:        []byte(`<?xml version="1.0" encoding="ISO-8859-1" ?> :)`),
			wantError:     nil,
		},
	}

	for _, test := range tests {
		converted, err := convertToUTF8(test.input, test.encodingLabel)

		if !errorsEqual(err, test.wantError) {
			t.Errorf("convertToUTF8(%s, %s) error = %s, wanted %s", test.input,
				test.encodingLabel, err, test.wantError)
			continue
		}

		if !byteSlicesEqual(converted, test.output) {
			t.Errorf("convertToUTF8(%s, %s) = %s, wanted %s", test.input,
				test.encodingLabel, converted, test.output)
			continue
		}
	}
}

func TestCleanXMLv1(t *testing.T) {
	tests := []struct {
		input     []byte
		output    []byte
		wantError error
	}{
		{
			input:     []byte(`abc`),
			output:    []byte(`abc`),
			wantError: nil,
		},
		{
			input:     []byte{'a', 'b', 'c', 0x0b},
			output:    []byte(`abc`),
			wantError: nil,
		},
		{
			input:     []byte{'a', 0x0b, 'b', 0x0a, 'c', 0x0b},
			output:    []byte{'a', 'b', 0x0a, 'c'},
			wantError: nil,
		},
		{
			input:     []byte("😀"),
			output:    []byte("😀"),
			wantError: nil,
		},
	}

	config.Verbose = true

	for _, test := range tests {
		output, err := cleanXMLv1(test.input)

		if !errorsEqual(err, test.wantError) {
			t.Errorf("cleanXMLv1(%s) = error %s, wanted %s", test.input, err,
				test.wantError)
			continue
		}

		if !byteSlicesEqual(output, test.output) {
			t.Errorf("cleanXMLv1(%s) = %s, wanted %s", test.input, output,
				test.output)
			continue
		}
	}

	config.Verbose = false
}

func TestUpdateXMLv1HeaderToUTF8(t *testing.T) {
	tests := []struct {
		input     []byte
		output    []byte
		wantError error
	}{
		{
			input:     []byte(`<?xml version="1.0" encoding="UTF-8"?>h`),
			output:    []byte(`<?xml version="1.0" encoding="UTF-8"?>h`),
			wantError: nil,
		},
		{
			input:     []byte(`<?xml version="1.0" encoding="UTF-8"?>`),
			output:    nil,
			wantError: fmt.Errorf("document ends with XML header"),
		},
		{
			input:     []byte(`<?xml version="1.0" encoding="ISO-8859-1"?>h`),
			output:    []byte(`<?xml version="1.0" encoding="UTF-8"?>h`),
			wantError: nil,
		},
		{
			input:     []byte(`<?xml version="1.0" encoding="ISO-8859-1" ?>h`),
			output:    []byte(`<?xml version="1.0" encoding="UTF-8"?>h`),
			wantError: nil,
		},
		{
			input:     []byte(`<?xml version="1.0" encoding="ISO-8859-1" h`),
			output:    nil,
			wantError: fmt.Errorf("could not find end of XML header"),
		},
	}

	for _, test := range tests {
		output, err := updateXMLv1HeaderToUTF8(test.input)

		if !errorsEqual(err, test.wantError) {
			t.Errorf("updateXMLv1HeaderToUTF8(%s) error = %s, wanted %s", test.input,
				err, test.wantError)
			continue
		}

		if !byteSlicesEqual(output, test.output) {
			t.Errorf("updateXMLv1HeaderToUTF8(%s) = %s, wanted %s", test.input,
				output, test.output)
			continue
		}
	}
}

func errorsEqual(e0, e1 error) bool {
	if e0 == nil && e1 == nil {
		return true
	}

	if e0 == nil && e1 != nil {
		return false
	}

	if e0 != nil && e1 == nil {
		return false
	}

	s0 := fmt.Sprintf("%s", e0)
	s1 := fmt.Sprintf("%s", e1)

	return s0 == s1
}

func byteSlicesEqual(b0, b1 []byte) bool {
	if b0 == nil && b1 == nil {
		return true
	}

	if b0 == nil && b1 != nil {
		return false
	}

	if b0 != nil && b1 == nil {
		return false
	}

	if len(b0) != len(b1) {
		return false
	}

	for i := range b0 {
		if b0[i] != b1[i] {
			return false
		}
	}

	return true
}
