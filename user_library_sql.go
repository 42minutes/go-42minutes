package minutes

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// SqlUserLibrary is a read-write user-specific library
type SqlUserLibrary struct {
	db *gorm.DB
}

// NewSqlUserLibrary accepts the database directory name and
// returns a new SqlUserLibrary instance
func NewSqlUserLibrary(db *gorm.DB) (*SqlUserLibrary, error) {
	if db.HasTable(UserShow{}) == false {
		if err := db.CreateTable(UserShow{}).Error; err != nil {
			return nil, err
		}
	}
	if db.HasTable(UserSeason{}) == false {
		if err := db.CreateTable(UserSeason{}).Error; err != nil {
			return nil, err
		}
	}
	if db.HasTable(UserEpisode{}) == false {
		if err := db.CreateTable(UserEpisode{}).Error; err != nil {
			return nil, err
		}
	}
	if db.HasTable(UserFile{}) == false {
		if err := db.CreateTable(UserFile{}).Error; err != nil {
			return nil, err
		}
	}
	return &SqlUserLibrary{
		db: db,
	}, nil
}

// UpsertShow adds a new show
// or errors with ErrNotImplemented, or ErrInternalServer
func (squl *SqlUserLibrary) UpsertShow(show *UserShow) error {
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

	if err := squl.db.Model(&ush).Update(show).Error; err != nil {
		return ErrInternalServer
	}
	return nil
}

// UpsertSeason adds or updates a season
// or errors with ErrNotImplemented, or ErrInternalServer, or ErrMissingShow
func (squl *SqlUserLibrary) UpsertSeason(season *UserSeason) error {
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

	if err := squl.db.Model(&useas).Update(season).Error; err != nil {
		return ErrInternalServer
	}
	return nil
}

// UpsertEpisode adds or updates a episode
// or errors with ErrNotImplemented, or ErrInternalServer, ErrMissingShow
// or ErrMissingSeason
func (squl *SqlUserLibrary) UpsertEpisode(episode *UserEpisode) error {
	usep := UserEpisode{}
	err := squl.db.Where(
		"show_id = ? AND season = ? AND number = ?",
		episode.ShowID, episode.Season, episode.Number,
	).Find(&usep).Error
	if err == gorm.ErrRecordNotFound {
		if err := squl.db.Create(episode).Error; err != nil {
			return ErrInternalServer
		}
		if err := squl.upsertFiles(episode); err != nil {
			return ErrInternalServer
		}

		return nil
	}

	if err != nil {
		return ErrInternalServer
	}

	if err = squl.db.Model(&usep).Updates(episode).Error; err != nil {
		return ErrInternalServer
	}
	if err := squl.upsertFiles(episode); err != nil {
		return ErrInternalServer
	}
	return nil
}

// GetShow returns a UserShow
// or errors with ErrNotFound, or ErrInternalServer
func (squl *SqlUserLibrary) GetShow(showID string) (*UserShow, error) {
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
func (squl *SqlUserLibrary) GetShows() ([]*UserShow, error) {
	ushs := []*UserShow{}
	err := squl.db.Find(&ushs).Error
	if err != nil {
		return nil, ErrInternalServer
	}
	return ushs, nil
}

// GetSeasons returns all Seasons for a show
// or errors with ErrNotFound, or ErrInternalServer
func (squl *SqlUserLibrary) GetSeasons(showID string) ([]*UserSeason, error) {
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
func (squl *SqlUserLibrary) GetSeason(showID string, seasonNumber int) (*UserSeason, error) {
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
func (squl *SqlUserLibrary) GetEpisodes(showID string, seasonNumber int) ([]*UserEpisode, error) {
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
func (squl *SqlUserLibrary) GetEpisode(showID string, seasonNumber, episodeNumber int) (*UserEpisode, error) {
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
func (squl *SqlUserLibrary) QueryShowsByTitle(title string) ([]*UserShow, error) {
	ushs := []*UserShow{}
	err := squl.db.Where(
		"title LIKE ?",
		fmt.Sprintf("%%%s%%", title),
	).Find(&ushs).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, ErrInternalServer
	}
	return ushs, nil
}

// QueryEpisodesForFinder
func (squl *SqlUserLibrary) QueryEpisodesForFinder() ([]*UserEpisode, error) {
	usfs := []*UserFile{}
	useps := []*UserEpisode{}
	err := squl.db.Where(
		"status = ? AND retry_time < ? and infohash = ''",
		"pending",
		time.Now().UTC().Unix(),
	).Find(&usfs).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, ErrInternalServer
	}

	for _, file := range usfs {
		usep := UserEpisode{}

		err := squl.db.Where(
			"show_id = ? AND season = ? AND number = ?",
			file.ShowID, file.Season, file.Episode,
		).Find(&usep).Error

		if err != nil {
			return nil, ErrInternalServer
		}
		usep.Files = append(usep.Files, file)
		useps = append(useps, &usep)
	}
	return useps, nil
}

// QueryEpisodesForDownloader
func (squl *SqlUserLibrary) QueryEpisodesForDownloader() ([]*UserEpisode, error) {
	usfs := []*UserFile{}
	useps := []*UserEpisode{}
	err := squl.db.Where(
		"status = ? AND retry_time < ? and infohash <> ''",
		"found",
		time.Now().UTC().Unix(),
	).Find(&usfs).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, ErrInternalServer
	}

	for _, file := range usfs {
		usep := UserEpisode{}

		err := squl.db.Where(
			"show_id = ? AND season = ? AND number = ?",
			file.ShowID, file.Season, file.Episode,
		).Find(&usep).Error

		if err != nil {
			return nil, ErrInternalServer
		}
		usep.Files = append(usep.Files, file)
		useps = append(useps, &usep)
	}
	return useps, nil
}

func (squl *SqlUserLibrary) upsertFiles(episode *UserEpisode) error {
	squl.db.Where(
		"show_id = ? AND season = ? AND episode = ?",
		episode.ShowID, episode.Season, episode.Number,
	).Delete(UserFile{})

	for _, file := range episode.Files {
		file.ShowID = episode.ShowID
		file.Season = episode.Season
		file.Episode = episode.Number
		if err := squl.db.Create(file).Error; err != nil {
			return ErrInternalServer
		}

	}
	return nil
}

func (squl *SqlUserLibrary) Close() error {
	err := squl.db.Close().Error
	if err != nil {
		return ErrInternalServer
	}
	return nil
}
