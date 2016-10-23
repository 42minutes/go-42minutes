package minutes

import (
	"time"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// SqliteUserLibrary is a read-write user-specific library
type SqliteUserLibrary struct {
	databaseDir string
	db          *gorm.DB
}

// NewSqliteUserLibrary accepts the database directory name and
// returns a new SqliteUserLibrary instance
func NewSqliteUserLibrary(DBDir string) *SqliteUserLibrary {
	db, err := gorm.Open("sqlite3", DBDir)
	if err != nil {
		log.Error(err)
	}

	if db.HasTable(UserShow{}) == false {
		if err := db.CreateTable(UserShow{}).Error; err != nil {
			log.Error("Could not create table.", err)
		}
	}
	if db.HasTable(UserSeason{}) == false {
		if err := db.CreateTable(UserSeason{}).Error; err != nil {
			log.Error("Could not create table.", err)
		}
	}
	if db.HasTable(UserEpisode{}) == false {
		if err := db.CreateTable(UserEpisode{}).Error; err != nil {
			log.Error("Could not create table.", err)
		}
	}

	return &SqliteUserLibrary{
		databaseDir: DBDir,
		db:          db,
	}
}

// UpsertShow adds a new show
// or errors with ErrNotImplemented, or ErrInternalServer
func (squl *SqliteUserLibrary) UpsertShow(show *UserShow) error {
	ush := UserShow{}
	err := squl.db.Where("ID = ?", show.ID).Find(&ush).Error

	if err == gorm.ErrRecordNotFound {
		if err := squl.db.Create(show).Error; err != nil {
			return ErrInternalServer
		}
		return nil
	}

	if err != nil {
		return ErrInternalServer
	}

	if err := squl.db.Model(ush).Update(show).Error; err != nil {
		return ErrInternalServer
	}
	return nil
}

// UpsertSeason adds or updates a season
// or errors with ErrNotImplemented, or ErrInternalServer, or ErrMissingShow
func (squl *SqliteUserLibrary) UpsertSeason(season *UserSeason) error {
	useas := UserSeason{}
	err := squl.db.Where(
		"show_id = ? AND number = ?",
		season.ShowID,
		season.Number,
	).Find(&useas).Error

	if err == gorm.ErrRecordNotFound {
		if err := squl.db.Create(season).Error; err != nil {
			return ErrInternalServer
		}
		return nil
	}

	if err != nil {
		return ErrInternalServer
	}

	if err := squl.db.Model(useas).Update(season).Error; err != nil {
		return ErrInternalServer
	}
	return nil
}

// UpsertEpisode adds or updates a episode
// or errors with ErrNotImplemented, or ErrInternalServer, ErrMissingShow
// or ErrMissingSeason
func (squl *SqliteUserLibrary) UpsertEpisode(episode *UserEpisode) error {
	usep := UserEpisode{}
	err := squl.db.Where(
		"show_id = ? AND season = ? AND number = ?",
		episode.ShowID, episode.Season, episode.Number,
	).Find(&usep).Error

	if err == gorm.ErrRecordNotFound {
		if err := squl.db.Create(episode).Error; err != nil {
			return ErrInternalServer
		}
		return nil
	}

	if err != nil {
		return ErrInternalServer
	}
	if err = squl.db.Model(usep).Update(episode).Error; err != nil {
		return ErrInternalServer
	}
	return nil
}

// GetShow returns a UserShow
// or errors with ErrNotFound, or ErrInternalServer
func (squl *SqliteUserLibrary) GetShow(showID string) (*UserShow, error) {
	ush := UserShow{}
	err := squl.db.Where("ID = ?", showID).Find(&ush).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, ErrInternalServer
	}
	return &ush, nil
}

// GetShows returns all Shows
// or errors with ErrNotImplemented, or ErrInternalServer
func (squl *SqliteUserLibrary) GetShows() ([]*UserShow, error) {
	ushs := []*UserShow{}
	err := squl.db.Find(&ushs).Error
	if err != nil {
		return nil, ErrInternalServer
	}
	return ushs, nil
}

// GetSeasons returns all Seasons for a show
// or errors with ErrNotFound, or ErrInternalServer
func (squl *SqliteUserLibrary) GetSeasons(showID string) ([]*UserSeason, error) {
	useas := []*UserSeason{}
	err := squl.db.Where("show_id = ?", showID).Find(&useas).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, ErrInternalServer
	}
	return useas, nil
}

// GetSeason returns a UserSeason given a UserShow's ID and a UserSeason number
// or errors with ErrNotFound, ErrMissingShow, or ErrInternalServer
func (squl *SqliteUserLibrary) GetSeason(showID string, seasonNumber int) (*UserSeason, error) {
	useas := UserSeason{}
	err := squl.db.Where(
		"show_id = ? AND number = ?",
		showID, seasonNumber,
	).Find(&useas).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, ErrInternalServer
	}
	return &useas, nil
}

// GetEpisodes returns all Shows for a show and season number
// or errors with ErrNotFound, or ErrInternalServer
func (squl *SqliteUserLibrary) GetEpisodes(showID string, seasonNumber int) ([]*UserEpisode, error) {
	useps := []*UserEpisode{}
	err := squl.db.Where(
		"show_id = ? AND season = ?",
		showID, seasonNumber,
	).Find(&useps).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, ErrInternalServer
	}
	return useps, nil
}

// GetEpisode returns a UserEpisode  given a UserShow's ID a UserSeason number
// and UserEpisode's number
// or errors with ErrNotFound, ErrMissingShow, or ErrInternalServer
func (squl *SqliteUserLibrary) GetEpisode(showID string, seasonNumber, episodeNumber int) (*UserEpisode, error) {
	usep := UserEpisode{}
	err := squl.db.Where("show_id = ? AND season = ? AND number = ?",
		showID, seasonNumber, episodeNumber,
	).Find(&usep).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, ErrInternalServer
	}
	return &usep, nil
}

// QueryShowsByTitle returns all Shows that match a partial title ordered
// by their probability
// or errors with ErrInternalServer
func (squl *SqliteUserLibrary) QueryShowsByTitle(title string) ([]*UserShow, error) {
	ushs := []*UserShow{}
	err := squl.db.Where("title = ?", title).Find(&ushs).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, ErrInternalServer
	}
	return ushs, nil
}

// QueryEpisodesForFinder
func (squl *SqliteUserLibrary) QueryEpisodesForFinder() ([]*UserEpisode, error) {
	useps := []*UserEpisode{}
	err := squl.db.Where(
		"downloaded = 0 AND retry_time < ? and infohash = ''",
		time.Now().UTC().Unix(),
	).Find(&useps).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, ErrInternalServer
	}
	return useps, nil
}

// QueryEpisodesForDownloader
func (squl *SqliteUserLibrary) QueryEpisodesForDownloader() ([]*UserEpisode, error) {
	useps := []*UserEpisode{}
	err := squl.db.Where(
		"downloaded = 0 AND retry_time < ? and infohash <> ''",
		time.Now().UTC().Unix(),
	).Find(&useps).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, ErrInternalServer
	}
	return useps, nil
}
