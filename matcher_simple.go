package minutes

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ls -R  -p | awk '
// /:$/&&f{s=$0;f=0}
// /:$/&&!f{sub(/:$/,"");s=$0;f=1;next}
// NF&&f{ print s"/"$0 }'

var (
	matcherEpisodeRegexs = []*regexp.Regexp{
		regexp.MustCompile(`(?P<show>.*?)[sS](?P<season>[0-9]+)[\._ ]*[eE](?P<ep>[0-9]+)([- ]?[Ee+](?P<secondEp>[0-9]+))?`),                              // S03E04-E05
		regexp.MustCompile(`(?P<show>.*?)[sS](?P<season>[0-9]{2})[\._\- ]+(?P<ep>[0-9]+)`),                                                               // S03-03
		regexp.MustCompile(`(?P<show>.*?)([^0-9]|^)(?P<season>[0-9]{1,2})[Xx](?P<ep>[0-9]+)(-[0-9]+[Xx](?P<secondEp>[0-9]+))?`),                          // 3x03
		regexp.MustCompile(`(.*?)[^0-9a-z](?P<season>[0-9]{1,2})(?P<ep>[0-9]{2})([\.\-][0-9]+(?P<secondEp>[0-9]{2})([ \-_\.]|$)[\.\-]?)?([^0-9a-z%]|$)`), // .602.
	}

	matcherStandaloneEpisodeRegexs = []*regexp.Regexp{
		regexp.MustCompile(`(.*?)( \(([0-9]+)\))? - ([0-9]+)+x([0-9]+)(-[0-9]+[Xx]([0-9]+))?( - (.*))?`),     // Newzbin style, no _UNPACK_
		regexp.MustCompile(`(.*?)( \(([0-9]+)\))?[Ss]([0-9]+)+[Ee]([0-9]+)(-[0-9]+[Xx]([0-9]+))?( - (.*))?`), // standard s00e00
	}

	matcherSeasonRegexs = []*regexp.Regexp{
		regexp.MustCompile(`.*?(?P<season>[0-9]+)$`), // folder for a season
	}

	matcherJustEpisodeRegexs = []*regexp.Regexp{
		regexp.MustCompile(`(?P<ep>[0-9]{1,3})[\. -_]of[\. -_]+[0-9]{1,3}`),       // 01 of 08
		regexp.MustCompile(`^(?P<ep>[0-9]{1,3})[^0-9]`),                           // 01 - Foo
		regexp.MustCompile(`e[a-z]*[ \.\-_]*(?P<ep>[0-9]{2,3})([^0-9c-uw-z%]|$)`), // Blah Blah ep234
		regexp.MustCompile(`.*?[ \.\-_](?P<ep>[0-9]{2,3})[^0-9c-uw-z%]+`),         // Flah - 04 - Blah
		regexp.MustCompile(`.*?[ \.\-_](?P<ep>[0-9]{2,3})$`),                      // Flah - 04
		regexp.MustCompile(`.*?[^0-9x](?P<ep>[0-9]{2,3})$`),                       // Flah707
	}
)

// SimpleMatch uses regular expressions to match against a show, season,
// and episode
type SimpleMatch struct {
	globalLibrary Library
}

// NewSimpleMatch returns a SimpleMatch
func NewSimpleMatch(glib Library) (*SimpleMatch, error) {
	return &SimpleMatch{
		globalLibrary: glib,
	}, nil
}

type matchedEpisode struct {
	Show          string
	Season        string
	Episode       string
	SecondEpisode string
}

// Match returns all episodes that match a filename or full path,
// ordered by their probability
// or errors with ErrInternalServer
func (m *SimpleMatch) Match(fp string) (eps []*Episode, err error) {
	// TODO(geoah) This needs a complete rewrite at some point but for now should be ok
	fp = strings.ToLower(fp)
	ps := strings.Split(fp, "/") // TODO(geoah) get os seperator
	fn := ps[len(ps)-1]
	ds := ps[:len(ps)-1]
	me := &matchedEpisode{}

	// try to match standalone episodes from the filename
	// from this, at the very least we should get season and episode
	if mer := m.match(matcherStandaloneEpisodeRegexs, fn); mer != nil {
		me.Show = mer.Show
		me.Season = mer.Season
		me.Episode = mer.Episode
		me.SecondEpisode = mer.SecondEpisode
	}

	// we can now try to match it as a full episode just in case we
	// are missing something
	if mer := m.match(matcherEpisodeRegexs, fn); mer != nil {
		if mer.Show != "" && me.Show == "" {
			me.Show = mer.Show
		}
		if mer.Season != "" && me.Season == "" {
			me.Season = mer.Season
		}
		if mer.Episode != "" && me.Episode == "" {
			me.Episode = mer.Episode
		}
	}

	// if they have at least one parent dir
	if len(ds) > 0 {
		pd := ds[len(ds)-1]
		// check if the parent dir is a season
		mer := m.match(matcherSeasonRegexs, pd)
		if mer != nil {
			if mer.Season != "" && me.Season == "" {
				me.Season = mer.Season
			}
		}

		// and check if it has the show name
		mer = m.match(matcherEpisodeRegexs, pd)
		if mer != nil {
			if mer.Show != "" && me.Show == "" {
				me.Show = mer.Show
			}
			if mer.Season != "" && me.Season == "" {
				me.Season = mer.Season
			}
			if mer.Episode != "" && me.Episode == "" {
				me.Episode = mer.Episode
			}
		}
	}

	if len(ds) > 1 {
		pd := ds[len(ds)-2]
		// check if it has the show name
		mer := m.match(matcherEpisodeRegexs, pd)
		if mer != nil {
			if mer.Show != "" && me.Show == "" {
				me.Show = mer.Show
			}
			if mer.Season != "" && me.Season == "" {
				me.Season = mer.Season
			}
			if mer.Episode != "" && me.Episode == "" {
				me.Episode = mer.Episode
			}
		}
	}

	if me == nil || me.Show == "" {
		return
	}

	clsh := strings.Replace(me.Show, ".", " ", -1)
	shs, err := m.globalLibrary.QueryShowsByTitle(clsh)
	if err != nil {
		log.Infof("> Could not find match for file '%s'", fp)
		return
	}

	if len(shs) == 0 {
		err = errors.New("Not enough shows")
		return
	}
	epn, _ := strconv.Atoi(me.Episode)
	sen, _ := strconv.Atoi(me.Season)
	ep := &Episode{
		ShowID: fmt.Sprintf("%d", shs[0].IDs.Trakt),
		Season: sen,
		Number: epn,
	}
	log.Infof("> Got match '%s' S%02dE%02d from file '%s'", me.Show, epn, sen, fp)
	eps = []*Episode{ep}

	return
}

func (m *SimpleMatch) match(rxs []*regexp.Regexp, fn string) *matchedEpisode {
	for _, rx := range matcherEpisodeRegexs {
		ms := rx.FindAllStringSubmatch(fn, -1)
		if len(ms) > 0 {
			if me := m.parseMatches(rx, ms); me != nil {
				return me
			}
		}
	}
	return nil
}

func (m *SimpleMatch) parseMatches(rx *regexp.Regexp, ms [][]string) *matchedEpisode {
	me := &matchedEpisode{}
	ns := rx.SubexpNames()
	for i, n := range ns {
		switch n {
		case "show":
			me.Show = ms[0][i]
		case "ep":
			me.Episode = ms[0][i]
		case "secondEp":
			me.SecondEpisode = ms[0][i]
		case "season":
			me.Season = ms[0][i]
		}
	}

	return me
}
