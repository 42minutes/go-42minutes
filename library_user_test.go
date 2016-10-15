package minutes

import (
	"fmt"
	"os"
	"testing"

	rethink "github.com/dancannon/gorethink"
	suite "github.com/stretchr/testify/suite"
)

const (
	database = "test_library"
)

var (
	// shows
	s1 = &Show{
		ID:    "1",
		Title: "first-show",
	}
	s2 = &Show{
		ID:    "2",
		Title: "second-show",
	}

	// seasons
	s1s1 = &Season{
		ID:       "3",
		ShowID:   "1",
		Number:   1,
		Overview: "first-season",
	}
	s1s2 = &Season{
		ID:       "4",
		ShowID:   "1",
		Number:   2,
		Overview: "second-season",
	}
	s2s1 = &Season{
		ID:       "5",
		ShowID:   "2",
		Number:   1,
		Overview: "first-season",
	}
	s2s2 = &Season{
		ID:       "6",
		ShowID:   "2",
		Number:   2,
		Overview: "second-season",
	}

	// episodes
	s1s1e1 = &Episode{
		ID:     "7",
		ShowID: "1",
		Season: 1,
		Number: 1,
		Title:  "first-episode",
	}
	s1s1e2 = &Episode{
		ID:     "8",
		ShowID: "1",
		Season: 1,
		Number: 2,
		Title:  "second-episode",
	}
	s1s1e3 = &Episode{
		ID:     "9",
		ShowID: "1",
		Season: 1,
		Number: 3,
		Title:  "third-episode",
	}
	s2s2e2 = &Episode{
		ID:     "10",
		ShowID: "2",
		Season: 2,
		Number: 2,
		Title:  "second-episode",
	}
)

func TestUserLibrarySuite(t *testing.T) {
	suite.Run(t, new(UserLibraryPersistenceSuite))
}

type UserLibraryPersistenceSuite struct {
	suite.Suite
	rethink *rethink.Session
	library Library
}

func (s *UserLibraryPersistenceSuite) SetupSuite() {
	host := os.Getenv("RETHINKDB_PORT_28015_TCP_ADDR")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("RETHINKDB_PORT_28015_TCP_PORT")
	if port == "" {
		port = "28015"
	}

	re, errConnecting := rethink.Connect(rethink.ConnectOpts{
		Address:  host + ":" + port,
		Database: database,
	})

	if errConnecting != nil {
		fmt.Println("Could not connect to db", errConnecting)
	}

	s.library = &UserLibrary{
		rethinkdb: re,
	}
	s.rethink = re

	rethink.DBDrop(database).Run(s.rethink)
	rethink.DBCreate(database).Run(s.rethink)

	rethink.DB(database).TableCreate(tableShows).Run(s.rethink)
	rethink.DB(database).TableCreate(tableSeasons).Run(s.rethink)
	rethink.DB(database).TableCreate(tableEpisodes).Run(s.rethink)

	rethink.DB(database).Table(tableSeasons).IndexCreate("show_id").Run(s.rethink)
	rethink.DB(database).Table(tableSeasons).IndexCreate("show_id").Run(s.rethink)
	rethink.DB(database).Table(tableEpisodes).IndexCreate("show_id").Run(s.rethink)

	rethink.DB(database).Table(tableEpisodes).IndexCreate("season_id").Run(s.rethink)

	rethink.DB(database).Table(tableShows).IndexWait().Exec(s.rethink)
	rethink.DB(database).Table(tableSeasons).IndexWait().Exec(s.rethink)
	rethink.DB(database).Table(tableEpisodes).IndexWait().Exec(s.rethink)
}

func (s *UserLibraryPersistenceSuite) SetupTest() {
	rethink.DB(database).Table(tableShows).Delete().RunWrite(s.rethink)
	rethink.DB(database).Table(tableSeasons).Delete().RunWrite(s.rethink)
	rethink.DB(database).Table(tableEpisodes).Delete().RunWrite(s.rethink)
}

func (s *UserLibraryPersistenceSuite) count(tbl string) int {
	cursor, err := rethink.Table(tbl).Count().Run(s.rethink)
	if err != nil {
		s.Fail(err.Error())
	}
	var cnt int
	cursor.One(&cnt)
	cursor.Close()
	return cnt
}

func (s *UserLibraryPersistenceSuite) addShows() {
	err := s.library.UpsertShow(s1)
	s.Nil(err)

	err = s.library.UpsertShow(s2)
	s.Nil(err)

	s.Equal(2, s.count(tableShows))
}

func (s *UserLibraryPersistenceSuite) addSeasons() {
	err := s.library.UpsertSeason(s1s1)
	s.Nil(err)

	err = s.library.UpsertSeason(s1s2)
	s.Nil(err)

	err = s.library.UpsertSeason(s2s1)
	s.Nil(err)

	err = s.library.UpsertSeason(s2s2)
	s.Nil(err)

	s.Equal(4, s.count(tableSeasons))
}

func (s *UserLibraryPersistenceSuite) addEpisodes() {
	err := s.library.UpsertEpisode(s1s1e1)
	s.Nil(err)

	err = s.library.UpsertEpisode(s1s1e2)
	s.Nil(err)

	err = s.library.UpsertEpisode(s1s1e3)
	s.Nil(err)

	err = s.library.UpsertEpisode(s2s2e2)
	s.Nil(err)

	s.Equal(4, s.count(tableEpisodes))
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_UpsertShow_Success() {
	s.addShows()

	s1.Title = "show-first-updated"
	err := s.library.UpsertShow(s1)
	s.Nil(err)

	s.Equal(2, s.count(tableShows))

	sh, err := s.library.GetShow(s1.ID)
	s.Equal(s1, sh)
	s.Nil(err)
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_UpsertSeason_Success() {
	s.addSeasons()

	s1s2.Overview = "second-season-updated"
	err := s.library.UpsertSeason(s1s2)
	s.Nil(err)

	s.Equal(4, s.count(tableSeasons))

	se, err := s.library.GetSeason(s1s2.ID)
	s.Equal(s1s2, se)
	s.Nil(err)

	s2s2.Overview = "second-season-updated"
	err = s.library.UpsertSeason(s2s2)
	s.Nil(err)

	s.Equal(4, s.count(tableSeasons))

	se, err = s.library.GetSeason(s2s2.ID)
	s.Equal(s2s2, se)
	s.Nil(err)
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_UpsertEpisode_Success() {
	s.addEpisodes()

	s1s1e3.Overview = "third-episode-updated"
	err := s.library.UpsertEpisode(s1s1e3)
	s.Nil(err)

	s.Equal(4, s.count(tableEpisodes))

	ep, err := s.library.GetEpisode(s1s1e3.ID)
	s.Equal(s1s1e3, ep)
	s.Nil(err)

	s2s2e2.Overview = "second-episode-updated"
	err = s.library.UpsertEpisode(s2s2e2)
	s.Nil(err)

	s.Equal(4, s.count(tableEpisodes))

	ep, err = s.library.GetEpisode(s2s2e2.ID)
	s.Equal(s2s2e2, ep)
	s.Nil(err)
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_GetShow_Success() {
	s.addShows()

	sh, err := s.library.GetShow(s1.ID)
	s.Equal(s1, sh)
	s.Nil(err)

	sh, err = s.library.GetShow(s2.ID)
	s.Equal(s2, sh)
	s.Nil(err)
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_GetSeason_Success() {
	s.addSeasons()

	se, err := s.library.GetSeason(s1s1.ID)
	s.Equal(s1s1, se)
	s.Nil(err)

	se, err = s.library.GetSeason(s1s2.ID)
	s.Equal(s1s2, se)
	s.Nil(err)
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_GetEpisode_Success() {
	s.addEpisodes()

	ep, err := s.library.GetEpisode(s1s1e1.ID)
	s.Equal(s1s1e1, ep)
	s.Nil(err)

	ep, err = s.library.GetEpisode(s1s1e2.ID)
	s.Equal(s1s1e2, ep)
	s.Nil(err)

	ep, err = s.library.GetEpisode(s1s1e3.ID)
	s.Equal(s1s1e3, ep)
	s.Nil(err)
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_GetShows_Success() {
	s.addShows()

	shs, err := s.library.GetShows()
	s.Equal(2, len(shs))
	s.Nil(err)
	s.Contains(shs, s1)
	s.Contains(shs, s2)
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_GetSeasonsByShow_Success() {
	s.addShows()
	s.addSeasons()

	shs, err := s.library.GetSeasonsByShow(s1.ID)
	s.Equal(2, len(shs))
	s.Nil(err)
	s.Contains(shs, s1s1)
	s.Contains(shs, s1s2)

	shs, err = s.library.GetSeasonsByShow(s2.ID)
	s.Equal(2, len(shs))
	s.Nil(err)
	s.Contains(shs, s2s1)
	s.Contains(shs, s2s2)
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_GetSeasonByNumber_Success() {
	s.addShows()
	s.addSeasons()

	se, err := s.library.GetSeasonByNumber(s1.ID, s1s1.Number)
	s.Equal(s1s1.Overview, se.Overview)
	s.Nil(err)

	se, err = s.library.GetSeasonByNumber(s1.ID, s1s2.Number)
	s.Equal(s1s2.Overview, se.Overview)
	s.Nil(err)
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_GetEpisodesBySeason_Success() {
	// TODO(geoah) Method will most likely be removed, not worth writing tests
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_GetEpisodesBySeasonNumber_Success() {
	s.addEpisodes()

	eps, err := s.library.GetEpisodesBySeasonNumber(s1s1.ShowID, s1s1.Number)
	s.Equal(3, len(eps))
	s.Nil(err)
	s.Contains(eps, s1s1e1)
	s.Contains(eps, s1s1e2)
	s.Contains(eps, s1s1e3)

	eps, err = s.library.GetEpisodesBySeasonNumber(s2s2.ShowID, s2s2.Number)
	s.Equal(1, len(eps))
	s.Nil(err)
	s.Contains(eps, s2s2e2)
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_GetEpisodeByNumber_Success() {
	s.addEpisodes()

	ep, err := s.library.GetEpisodeByNumber(s1s1e1.ShowID, s1s1e1.Season, s1s1e1.Number)
	s.Equal(s1s1e1, ep)
	s.Nil(err)

	ep, err = s.library.GetEpisodeByNumber(s2s2e2.ShowID, s2s2e2.Season, s2s2e2.Number)
	s.Equal(s2s2e2, ep)
	s.Nil(err)
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_QueryShowsByTitle_Success() {
	s.addShows()

	shs, err := s.library.QueryShowsByTitle("first")
	s.Equal(1, len(shs))
	s.Nil(err)
	s.Contains(shs, s1)

	shs, err = s.library.QueryShowsByTitle("second")
	s.Equal(1, len(shs))
	s.Nil(err)
	s.Contains(shs, s2)

	shs, err = s.library.QueryShowsByTitle("show")
	s.Equal(2, len(shs))
	s.Nil(err)
	s.Contains(shs, s1)
	s.Contains(shs, s2)
}
