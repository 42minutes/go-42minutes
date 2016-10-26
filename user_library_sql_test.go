package minutes

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	suite "github.com/stretchr/testify/suite"
)

var (
	// shows
	s1 = &UserShow{
		ID:    "1",
		Title: "first-show",
	}
	s2 = &UserShow{
		ID:    "2",
		Title: "second-show",
	}

	// seasons
	s1s1 = &UserSeason{
		ShowID: "1",
		Number: 1,
	}
	s1s2 = &UserSeason{
		ShowID: "1",
		Number: 2,
	}
	s2s1 = &UserSeason{
		ShowID: "2",
		Number: 1,
	}
	s2s2 = &UserSeason{
		ShowID: "2",
		Number: 2,
	}

	// episodes
	s1s1e1 = &UserEpisode{
		ShowID: "1",
		Season: 1,
		Number: 1,
		Title:  "first-episode",
	}
	s1s1e2 = &UserEpisode{
		ShowID: "1",
		Season: 1,
		Number: 2,
		Title:  "second-episode",
	}
	s1s1e3 = &UserEpisode{
		ShowID: "1",
		Season: 1,
		Number: 3,
		Title:  "third-episode",
	}
	s2s2e2 = &UserEpisode{
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
	database *gorm.DB
	library  UserLibrary
}

func (s *UserLibraryPersistenceSuite) SetupSuite() {
	db, err := gorm.Open("sqlite3", "test_data.db")
	if err != nil {
		s.Fail(err.Error())
	}

	// user rw library for single hardcoded user id
	ulib, err := NewSqlUserLibrary(db)
	if err != nil {
		s.Fail(err.Error())
	}

	s.database = db
	s.library = ulib
}

func (s *UserLibraryPersistenceSuite) SetupTest() {

}

func (s *UserLibraryPersistenceSuite) count(table interface{}) int {
	tableCount := 0
	err := s.database.Model(reflect.New(reflect.TypeOf(table)).Interface()).Count(&tableCount).Error
	if err != nil {
		log.Error(err)
		s.Fail(err.Error())
	}
	fmt.Println(tableCount)
	return tableCount
}

func (s *UserLibraryPersistenceSuite) addShows() {
	err := s.library.UpsertShow(s1)
	s.Nil(err)

	err = s.library.UpsertShow(s2)
	s.Nil(err)
	s.Equal(2, s.count(UserShow{}))
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

	s.Equal(4, s.count(UserSeason{}))
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

	s.Equal(4, s.count(UserEpisode{}))
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_UpsertShow_Success() {
	s.addShows()
	s1.Title = "show-first-updated"
	err := s.library.UpsertShow(s1)
	s.Nil(err)

	s.Equal(2, s.count(UserShow{}))

	sh, err := s.library.GetShow(s1.ID)
	s.ueq(s1, sh)
	s.Nil(err)
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_UpsertSeason_Success() {
	s.addSeasons()

	err := s.library.UpsertSeason(s1s2)
	s.Nil(err)

	s.Equal(4, s.count(UserSeason{}))

	se, err := s.library.GetSeason(s1s2.ShowID, s1s2.Number)
	s.ueq(s1s2, se)
	s.Nil(err)

	err = s.library.UpsertSeason(s2s2)
	s.Nil(err)

	s.Equal(4, s.count(UserSeason{}))

	se, err = s.library.GetSeason(s2s2.ShowID, s2s2.Number)
	s.ueq(s2s2, se)
	s.Nil(err)
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_UpsertEpisode_Success() {
	s.addEpisodes()

	s1s1e3.Title = "third-episode-updated"
	err := s.library.UpsertEpisode(s1s1e3)
	s.Nil(err)

	s.Equal(4, s.count(UserEpisode{}))

	ep, err := s.library.GetEpisode(s1s1e3.ShowID, s1s1e3.Season, s1s1e3.Number)
	s.ueq(s1s1e3, ep)
	s.Nil(err)

	s2s2e2.Title = "second-episode-updated"
	err = s.library.UpsertEpisode(s2s2e2)
	s.Nil(err)

	s.Equal(4, s.count(UserEpisode{}))

	ep, err = s.library.GetEpisode(s2s2e2.ShowID, s2s2e2.Season, s2s2e2.Number)
	s.ueq(s2s2e2, ep)
	s.Nil(err)
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_GetShow_Success() {
	s.addShows()

	sh, err := s.library.GetShow(s1.ID)
	s.ueq(s1, sh)
	s.Nil(err)

	sh, err = s.library.GetShow(s2.ID)
	s.ueq(s2, sh)
	s.Nil(err)
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_GetSeason_Success() {
	s.addSeasons()

	se, err := s.library.GetSeason(s1s1.ShowID, s1s1.Number)
	s.ueq(s1s1, se)
	s.Nil(err)

	se, err = s.library.GetSeason(s1s2.ShowID, s1s2.Number)
	s.ueq(s1s2, se)
	s.Nil(err)
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_GetEpisode_Success() {
	s.addEpisodes()

	ep, err := s.library.GetEpisode(s1s1e1.ShowID, s1s1e1.Season, s1s1e1.Number)
	s.ueq(s1s1e1, ep)
	s.Nil(err)

	ep, err = s.library.GetEpisode(s1s1e2.ShowID, s1s1e2.Season, s1s1e2.Number)
	s.ueq(s1s1e2, ep)
	s.Nil(err)

	ep, err = s.library.GetEpisode(s1s1e3.ShowID, s1s1e3.Season, s1s1e3.Number)
	s.ueq(s1s1e3, ep)
	s.Nil(err)
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_GetShows_Success() {
	s.addShows()

	shs, err := s.library.GetShows()
	s.Equal(2, len(shs))
	s.Nil(err)
	s.uconsh(shs, s1)
	s.uconsh(shs, s2)
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_GetSeasonsByShow_Success() {
	s.addShows()
	s.addSeasons()

	ses, err := s.library.GetSeasons(s1.ID)
	s.Equal(2, len(ses))
	s.Nil(err)
	s.uconse(ses, s1s1)
	s.uconse(ses, s1s2)

	ses, err = s.library.GetSeasons(s2.ID)
	s.Equal(2, len(ses))
	s.Nil(err)
	s.uconse(ses, s2s1)
	s.uconse(ses, s2s2)
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_GetEpisodes_Success() {
	s.addEpisodes()

	eps, err := s.library.GetEpisodes(s1s1.ShowID, s1s1.Number)
	s.Equal(3, len(eps))
	s.Nil(err)
	s.ucone(eps, s1s1e1)
	s.ucone(eps, s1s1e2)
	s.ucone(eps, s1s1e3)

	eps, err = s.library.GetEpisodes(s2s2.ShowID, s2s2.Number)
	s.Equal(1, len(eps))
	s.Nil(err)
	s.ucone(eps, s2s2e2)
}

func (s *UserLibraryPersistenceSuite) TestUserLibrary_QueryShowsByTitle_Success() {
	s.addShows()

	shs, err := s.library.QueryShowsByTitle("first")
	s.Equal(1, len(shs))
	s.Nil(err)
	s.uconsh(shs, s1)

	shs, err = s.library.QueryShowsByTitle("second")
	s.Equal(1, len(shs))
	s.Nil(err)
	s.uconsh(shs, s2)

	shs, err = s.library.QueryShowsByTitle("show")
	s.Equal(2, len(shs))
	s.Nil(err)
	s.uconsh(shs, s1)
	s.uconsh(shs, s2)
}

func (s *UserLibraryPersistenceSuite) uconsh(xx []*UserShow, y interface{}) bool {
	for _, x := range xx {
		if s.ueq(x, y) {
			return true
		}
	}
	return false
}

func (s *UserLibraryPersistenceSuite) uconse(xx []*UserSeason, y interface{}) bool {
	for _, x := range xx {
		if s.ueq(x, y) {
			return true
		}
	}
	return false
}

func (s *UserLibraryPersistenceSuite) ucone(xx []*UserEpisode, y interface{}) bool {
	for _, x := range xx {
		if s.ueq(x, y) {
			return true
		}
	}
	return false
}

func (s *UserLibraryPersistenceSuite) ueq(x, y interface{}) bool {
	switch xv := x.(type) {
	case *UserShow:
		yv := y.(*UserShow)
		if xv.ID == yv.ID && xv.Title == yv.Title {
			return true
		}
	case *UserSeason:
		yv := y.(*UserSeason)
		if xv.ShowID == yv.ShowID &&
			xv.Number == yv.Number {
			return true
		}
	case *UserEpisode:
		yv := y.(*UserEpisode)
		if xv.ShowID == yv.ShowID &&
			xv.Season == yv.Season &&
			xv.Number == yv.Number &&
			xv.Title == yv.Title {
			return true
		}
	}
	return false
}
