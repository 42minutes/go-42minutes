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
	qr := ctx.URLParam("title")
	log.Debug("qr", qr)
	if qr == "" {
		ushs, err := api.ulibrary.GetShows()
		if err != nil {
			log.Error("Could not get ulib shows", err)
			ctx.EmitError(iris.StatusInternalServerError)
			return
		}
		for i := 0; i < len(ushs); i++ {
			gsh, err := api.glibrary.GetShow(ushs[i].ID)
			if err != nil {
				log.Error("Could not get glib show %s", ushs[i].ID)
				ctx.EmitError(iris.StatusInternalServerError)
				return
			}
			ushs[i].MergeInPlace(gsh)
		}
		ctx.JSON(iris.StatusOK, ushs)
		return
	}

	gshs, err := api.glibrary.QueryShowsByTitle(qr)
	if err != nil {
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}
	ushs := []*UserShow{}
	for i := 0; i < len(gshs); i++ {
		gsh, err := api.glibrary.GetShow(gshs[i].ID)
		log.Info(gsh)
		if err != nil {
			log.Error("Could not get glib show %s", ushs[i].ID)
			ctx.EmitError(iris.StatusInternalServerError)
			return
		}
		ush := &UserShow{}
		ush.MergeInPlace(gsh)
		ushs = append(ushs, ush)
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

type reqShowPost struct {
	ID string `json:"id"`
}

// HandleShowPost -
func (api *API) HandleShowPost(ctx *iris.Context) {
	shr := &reqShowPost{}
	if err := ctx.ReadJSON(shr); err != nil {
		ctx.EmitError(iris.StatusBadRequest)
		return
	}

	ctx.Params = iris.PathParameters{
		iris.PathParameter{
			Key:   "show_id",
			Value: shr.ID,
		},
	}

	ush, err := api.ulibrary.GetShow(shr.ID)
	if err != nil {
		if err != ErrNotFound {
			ctx.EmitError(iris.StatusInternalServerError)
			return
		}
	} else if ush != nil {
		api.HandleShow(ctx)
		return
	}

	ush = &UserShow{}

	gsh, err := api.glibrary.GetShow(shr.ID)
	if err != nil {
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

	ush.MergeInPlace(gsh)

	if err := api.ulibrary.UpsertShow(ush); err != nil {
		ctx.EmitError(iris.StatusInternalServerError)
		return
	}

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
	uep, err := api.ulibrary.GetEpisode(sid, sen, epn)
	if err != nil && err != ErrNotFound {
		log.Error("Could not get ulib episode %s %d %d", sid, sen, epn)
		// TODO(geoah) Handle error
	} else if err != ErrNotFound {
		rep = uep
	}
	rep.MergeInPlace(gep)
	ctx.JSON(iris.StatusOK, rep)
}
