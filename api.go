package minutes

import (
	"strconv"

	"github.com/kataras/iris"
)

// API -
type API struct {
	glibrary ShowLibrary
	ulibrary UserLibrary
}

// NewAPI -
func NewAPI(glib ShowLibrary, ulib UserLibrary) *API {
	return &API{
		glibrary: glib,
		ulibrary: ulib,
	}
}

// HandleShows -
func (api *API) HandleShows(ctx *iris.Context) {
	ushs, _ := api.ulibrary.GetShows()
	// TODO(geoah) Handle error
	for i := 0; i < len(ushs); i++ {
		gsh, err := api.glibrary.GetShow(ushs[i].ID)
		if err != nil {
			log.Error("Could not get glib show %s", ushs[i].ID)
			// TODO(geoah) Handle error
		}
		ushs[i].MergeInPlace(gsh)
	}
	ctx.JSON(iris.StatusOK, ushs)
}

// HandleShow -
func (api *API) HandleShow(ctx *iris.Context) {
	sid := ctx.Param("show_id")
	ush, _ := api.ulibrary.GetShow(sid)
	// TODO(geoah) Handle error
	gsh, err := api.glibrary.GetShow(sid)
	if err != nil {
		log.Error("Could not get glib show %s", sid)
		// TODO(geoah) Handle error
	}
	ush.MergeInPlace(gsh)
	ctx.JSON(iris.StatusOK, ush)
}

// HandleSeasons -
func (api *API) HandleSeasons(ctx *iris.Context) {
	sid := ctx.Param("show_id")
	gses, err := api.glibrary.GetSeasons(sid)
	if err != nil {
		log.Error("Could not get gses", err)
		// TODO(geoah) Handle error
	}
	rses := []*UserSeason{}
	for _, gse := range gses {
		rse := &UserSeason{}
		use, err := api.ulibrary.GetSeason(sid, gse.Number)
		if err != nil && err != ErrNotFound {
			log.Error("Could not get ulib season %s %d", sid, gse.Number, err)
			// TODO(geoah) Handle error
		} else if err != ErrNotFound {
			rse = use
		}
		rse.MergeInPlace(gse)
		rses = append(rses, rse)
	}
	ctx.JSON(iris.StatusOK, rses)
}

// HandleSeason -
func (api *API) HandleSeason(ctx *iris.Context) {
	sid := ctx.Param("show_id")
	sen, _ := strconv.Atoi(ctx.Param("season"))
	// TODO(geoah) Handle error
	rse := &UserSeason{}
	gse, _ := api.glibrary.GetSeason(sid, sen)
	// TODO(geoah) Handle error
	rse.MergeInPlace(gse)
	use, err := api.glibrary.GetSeason(sid, sen)
	if err != nil && err != ErrNotFound {
		log.Error("Could not get ulib season %s %d", sid, gse.Number)
		// TODO(geoah) Handle error
	} else if err != ErrNotFound {
		rse.MergeInPlace(use)
	}
	ctx.JSON(iris.StatusOK, rse)
}

// HandleEpisodes -
func (api *API) HandleEpisodes(ctx *iris.Context) {
	sid := ctx.Param("show_id")
	sen, _ := strconv.Atoi(ctx.Param("season"))
	// TODO(geoah) Handle error
	geps, err := api.glibrary.GetEpisodes(sid, sen)
	if err != nil {
		log.Error("Could not get geps", err)
		// TODO(geoah) Handle error
	}
	reps := []*UserEpisode{}
	for _, gep := range geps {
		rep := &UserEpisode{}
		uep, err := api.ulibrary.GetEpisode(sid, sen, gep.Number)
		if err != nil && err != ErrNotFound {
			log.Error("Could not get ulib season %s %d", sid, sen, err)
			// TODO(geoah) Handle error
		} else if err != ErrNotFound {
			rep = uep
		}
		rep.MergeInPlace(gep)
		reps = append(reps, rep)
	}
	ctx.JSON(iris.StatusOK, reps)
}

// HandleEpisode -
func (api *API) HandleEpisode(ctx *iris.Context) {
	sid := ctx.Param("show_id")
	sen, _ := strconv.Atoi(ctx.Param("season"))
	// TODO(geoah) Handle error
	epn, _ := strconv.Atoi(ctx.Param("episode"))
	// TODO(geoah) Handle error
	rep := &UserEpisode{}
	gep, _ := api.glibrary.GetEpisode(sid, sen, epn)
	// TODO(geoah) Handle error
	rep.MergeInPlace(gep)
	uep, err := api.glibrary.GetEpisode(sid, sen, epn)
	if err != nil && err != ErrNotFound {
		log.Error("Could not get ulib episode %s %d %d", sid, sen, epn)
		// TODO(geoah) Handle error
	} else if err != ErrNotFound {
		rep.MergeInPlace(uep)
	}
	ctx.JSON(iris.StatusOK, rep)
}
