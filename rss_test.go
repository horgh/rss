package rss

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAsRSS(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		output  *Feed
		success bool
	}{
		{
			name: "well formed XML feed",
			file: "test-data/rss-good.xml",
			output: &Feed{
				Title:       "A Nice Site",
				Link:        "https://example.com",
				Description: "A Nice Website",
				PubDate:     time.Time{},
				Items: []Item{
					{
						Title:       "Nice Title 1",
						Link:        "https://example.com/2020/03/nice-title-1/",
						Description: "<p>should we write something nice?</p>\n",
						PubDate:     time.Date(2020, 3, 6, 18, 15, 47, 0, time.UTC),
						GUID:        "https://example.com/?p=29611",
					},
				},
				Type: "RSS",
			},
			success: true,
		},
		{
			name: "rss feed with no XML declaration",
			file: "test-data/rss-with-no-xml-declaration.xml",
			output: &Feed{
				Title:       "Nice title",
				Link:        "https://blog.example.com/",
				Description: "Recent content on example.com",
				PubDate:     time.Date(2019, 4, 8, 10, 20, 30, 0, time.UTC),
				Items: []Item{
					{
						Title:       "My Nice Post",
						Link:        "https://blog.example.com/post/nice/",
						Description: "hi",
						PubDate:     time.Date(2019, 4, 8, 10, 20, 33, 0, time.UTC),
						GUID:        "https://blog.example.com/post/nice/",
					},
				},
				Type: "RSS",
			},
			success: true,
		},
		{
			name:    "root tag is not rss",
			file:    "test-data/rss-with-different-root-tag.xml",
			success: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf, err := ioutil.ReadFile(test.file)
			require.NoError(t, err, "read file")

			feed, err := ParseFeedXML(buf)
			if !test.success {
				assert.Error(t, err, "error parsing")
				return
			}
			assert.NoError(t, err, "parse feed")
			assert.Equal(t, test.output, feed, "correct feed")
		})
	}
}

func TestParseAsRDF(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		output  *Feed
		success bool
	}{
		{
			"An edited/subset version of a feed from Slashdot.",
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
				PubDate:     time.Date(2017, 1, 17, 21, 30, 14, 0, time.UTC),
				Items: []Item{
					{
						Title:       "Uber Sues City of Seattle To Block Landmark Driver Union Ordinance",
						Link:        "https://tech.slashdot.org/story/17/01/17/197230/uber-sues-city-of-seattle-to-block-landmark-driver-union-ordinance?utm_source=rss1.0mainlinkanon&utm_medium=feed",
						Description: "Seattle's landmark law that lets drivers",
						PubDate:     time.Date(2017, 1, 17, 20, 40, 0, 0, time.UTC),
					},
					{
						Title:       "Netflix is 'Killing' DVD Sales, Research Finds",
						Link:        "https://entertainment.slashdot.org/story/17/01/17/1855219/netflix-is-killing-dvd-sales-research-finds?utm_source=rss1.0mainlinkanon&utm_medium=feed",
						Description: "Netflix has become the go-to destination for many movie",
						PubDate:     time.Date(2017, 1, 17, 20, 0, 0, 0, time.UTC),
					},
				},
				Type: "RDF",
			},
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			feed, err := ParseFeedXML([]byte(test.input))
			if err != nil {
				if !test.success {
					return
				}

				t.Errorf("parseAsAtom(%s) = error %s, wanted success", test.input, err)
				return
			}

			if !test.success {
				t.Errorf("parseAsAtom(%s) = success, wanted error", test.input)
				return
			}

			assert.Equal(t, test.output, feed, "correct feed")
		})
	}
}

func TestParseAsAtom(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		output  *Feed
		success bool
	}{
		{
			"valid feed",
			`<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">

 <title>Test one two</title>
 <link href="http://www.example.com/atom.xml" rel="self"/>
 <link href="http://www.example.com"/>
 <updated>2017-01-11T20:30:23-00:00</updated>
 <id>http://www.example.com-id</id>
 <author>
   <name>John Q. Public</name>
   <email>john@example.com</email>
 </author>

 <entry>
   <title>Test title 1</title>
   <link href="http://www.example.com/test-entry-1"/>
   <updated>2017-01-11T00:00:00-00:00</updated>
   <id>http://www.example.com/test-entry-1-id</id>
   <content type="html">&lt;p&gt;Testing content 1&lt;/p&gt;</content>
</entry>

 <entry>
   <title>Test title 2</title>
   <link href="http://www.example.com/test-entry-2"/>
   <updated>2017-01-12T00:00:00-00:00</updated>
   <id>http://www.example.com/test-entry-2-id</id>
   <content type="html">&lt;p&gt;Testing content 2&lt;/p&gt;</content>
</entry>
</feed>
`,
			&Feed{
				Title:       "Test one two",
				Link:        "http://www.example.com/atom.xml",
				Description: "",
				PubDate:     time.Date(2017, 1, 11, 20, 30, 23, 0, time.UTC),
				Items: []Item{
					{
						Title:       "Test title 1",
						Link:        "http://www.example.com/test-entry-1",
						Description: "<p>Testing content 1</p>",
						PubDate:     time.Date(2017, 1, 11, 0, 0, 0, 0, time.UTC),
						GUID:        "http://www.example.com/test-entry-1-id",
					},
					{
						Title:       "Test title 2",
						Link:        "http://www.example.com/test-entry-2",
						Description: "<p>Testing content 2</p>",
						PubDate:     time.Date(2017, 1, 12, 0, 0, 0, 0, time.UTC),
						GUID:        "http://www.example.com/test-entry-2-id",
					},
				},
				Type: "Atom",
			},
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			feed, err := parseAsAtom([]byte(test.input))
			if err != nil {
				if !test.success {
					return
				}

				t.Errorf("parseAsAtom(%s) = error %s, wanted success", test.input, err)
				return
			}

			if !test.success {
				t.Errorf("parseAsAtom(%s) = success, wanted error", test.input)
				return
			}

			assert.Equal(t, test.output, feed, "correct feed")
		})
	}
}

func TestMakeXML(t *testing.T) {
	tests := []struct {
		name    string
		input   Feed
		output  string
		success bool
	}{
		{
			"success",
			Feed{
				Title:       "Test feed",
				Link:        "https://www.example.com/",
				Description: "A nice feed",
				PubDate: time.Date(2016, 12, 25, 11, 0, 0, 0,
					time.FixedZone("TZ", 0)),
				Items: []Item{
					{
						Title:       "Nice item 1",
						Link:        "https://www.example.com/1",
						Description: "Item 1 is very nice",
						PubDate: time.Date(2016, 12, 25, 11, 01, 0, 0,
							time.FixedZone("TZ", 0)),
					},
					{
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
		t.Run(test.name, func(t *testing.T) {
			buf, err := makeXML(test.input)
			if err != nil {
				if !test.success {
					return
				}

				t.Errorf("makeXML(%#v) = error %s", test.input, err)
				return
			}

			if !test.success {
				t.Errorf("makeXML(%#v) = success, wanted error", test.input)
				return
			}

			if !bytes.Equal(buf, []byte(test.output)) {
				t.Errorf("makeXML(%#v) = %s, wanted %s", test.input, buf, test.output)
				return
			}
		})
	}
}

func TestParseTime(t *testing.T) {
	tests := []struct {
		TimeString string
		Time       time.Time
	}{
		{
			"Sun, 09 Apr 2017 05:06 GMT",
			time.Date(2017, time.April, 9, 5, 6, 0, 0, time.UTC),
		},
	}

	config.Verbose = true

	for _, test := range tests {
		gotTime := parseTime(test.TimeString)

		gotTimeUTC := gotTime.UTC()

		if !gotTimeUTC.Equal(test.Time) {
			t.Errorf("parseTime(%s) = %s, wanted %s", test.TimeString, gotTimeUTC,
				test.Time)
		}
	}
}
