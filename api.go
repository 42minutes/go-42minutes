package minutes

import (
	"net/http"
	"strconv"

	gin "github.com/gin-gonic/gin"
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
func (api *API) HandleShows(ctx *gin.Context) {
	qr := ctx.Param("title")
	log.Debug("qr", qr)
	if qr == "" {
		ushs, err := api.ulibrary.GetShows()
		if err != nil {
			log.Error("Could not get ulib shows", err)
			ctx.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		for i := 0; i < len(ushs); i++ {
			gsh, err := api.glibrary.GetShow(ushs[i].ID)
			if err != nil {
				log.Error("Could not get glib show %s", ushs[i].ID)
				ctx.String(http.StatusInternalServerError, "Internal Server Error")
				return
			}
			ushs[i].MergeInPlace(gsh)
		}
		ctx.JSON(http.StatusOK, ushs)
		return
	}

	gshs, err := api.glibrary.QueryShowsByTitle(qr)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	ushs := []*UserShow{}
	for i := 0; i < len(gshs); i++ {
		gsh, err := api.glibrary.GetShow(gshs[i].ID)
		log.Info("%v", gsh)
		if err != nil {
			log.Error("Could not get glib show %s", ushs[i].ID)
			ctx.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		ush := &UserShow{}
		ush.MergeInPlace(gsh)
		ushs = append(ushs, ush)
	}
	ctx.JSON(http.StatusOK, ushs)
}

// HandleShow -
func (api *API) HandleShow(ctx *gin.Context) {
	sid := ctx.Param("show_id")
	ush, _ := api.ulibrary.GetShow(sid)
	// TODO(geoah) Handle error
	gsh, err := api.glibrary.GetShow(sid)
	if err != nil {
		log.Error("Could not get glib show %s", sid)
		// TODO(geoah) Handle error
	}
	ush.MergeInPlace(gsh)
	ctx.JSON(http.StatusOK, ush)
}

type reqShowPost struct {
	ID string `json:"id"`
}

// HandleShowPost -
func (api *API) HandleShowPost(ctx *gin.Context) {
	shr := &reqShowPost{}
	if ctx.BindJSON(&shr) == nil {
		ctx.String(http.StatusBadRequest, "Bad Request")
		return
	}

	ush, err := api.ulibrary.GetShow(shr.ID)
	if err != nil {
		if err != ErrNotFound {
			ctx.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}
	} else if ush != nil {
		gsh, err := api.glibrary.GetShow(shr.ID)
		if err != nil {
			log.Error("Could not get glib show %s", shr.ID)
			// TODO(geoah) Handle error
		}
		ush.MergeInPlace(gsh)
		ctx.JSON(http.StatusOK, ush)
		return
	}

	ush = &UserShow{}

	gsh, err := api.glibrary.GetShow(shr.ID)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	ush.MergeInPlace(gsh)

	if err := api.ulibrary.UpsertShow(ush); err != nil {
		ctx.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	ctx.JSON(http.StatusOK, ush)
}

// HandleSeasons -
func (api *API) HandleSeasons(ctx *gin.Context) {
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
	ctx.JSON(http.StatusOK, rses)
}

// HandleSeason -
func (api *API) HandleSeason(ctx *gin.Context) {
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
	ctx.JSON(http.StatusOK, rse)
}

// HandleEpisodes -
func (api *API) HandleEpisodes(ctx *gin.Context) {
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
	ctx.JSON(http.StatusOK, reps)
}

// HandleEpisode -
func (api *API) HandleEpisode(ctx *gin.Context) {
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
	ctx.JSON(http.StatusOK, rep)
}
