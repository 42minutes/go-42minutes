package minutes

var (
	fileResolutions = []string{
		"1080p",
		"1080i",
		"720p",
		"720i",
		"hr",
		"576p",
		"480p",
		"368p",
		"360p",
	}

	fileSources = []string{
		"bluray",
		"remux",
		"dvdrip",
		"webdl",
		"hdtv",
		"webrip",
		"bdscr",
		"dvdscr",
		"sdtv",
		"dsr",
		"tvrip",
		"preair",
		"ppvrip",
		"hdrip",
		"r5",
		"tc",
		"ts",
		"cam",
		"workprint",
	}

	fileVideoCodecs = []string{
		"10bit",
		"h265",
		"h264",
		"xvid",
		"divx",
	}

	fileAudioCodecs = []string{
		"truehd",
		"dts",
		"dtshd",
		"flac",
		"ac3",
		"dd5.1",
		"aac",
		"mp3",
	}
)

// Matcher tries to match a filename with an episode
type Matcher interface {
	// Match returns all episodes that match a filename or full path,
	// ordered by their probability
	// or errors with ErrInternalServer
	Match(filename string) ([]*UserEpisode, error)
}
