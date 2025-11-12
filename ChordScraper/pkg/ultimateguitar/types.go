package ultimateguitar

type TabType int32

// Here are the rough known tab types
const (
	TabTypeVideo    TabType = 100 // ??
	TabTypeTabs     TabType = 200
	TabTypeChords   TabType = 300
	TabTypeBass     TabType = 400
	TabTypePro      TabType = 500
	TabTypePower    TabType = 600 // ??
	TabTypeDrums    TabType = 700
	TabTypeUkulele  TabType = 800
	TabTypeOfficial TabType = 900
	TabTypeAll      TabType = 1000
)

// TabResult struct - this is what we get when we fetch a tab by id.
type TabResult struct {
	ID                 int     `json:"id"`
	SongName           string  `json:"song_name"`
	ArtistName         string  `json:"artist_name"`
	Type               string  `json:"type"`
	Part               string  `json:"part"`
	Version            int     `json:"version"`
	Votes              int     `json:"votes"`
	Rating             float64 `json:"rating"`
	Date               string  `json:"date"`
	Status             string  `json:"status"`
	PresetID           int     `json:"preset_id"`
	TabAccessType      string  `json:"tab_access_type"`
	TpVersion          int     `json:"tp_version"`
	TonalityName       string  `json:"tonality_name"`
	VersionDescription *string `json:"version_description"`
	Verified           int     `json:"verified"`
	Recording          struct {
		IsAcoustic       int               `json:"is_acoustic"`
		TonalityName     string            `json:"tonality_name"`
		Performance      Performance       `json:"performance"`
		RecordingArtists []RecordingArtist `json:"recording_artists"`
	} `json:"recording"`
	Versions []struct {
		ID                 int     `json:"id"`
		SongName           string  `json:"song_name"`
		ArtistName         string  `json:"artist_name"`
		Type               string  `json:"type"`
		Part               string  `json:"part"`
		Version            int     `json:"version"`
		Votes              int     `json:"votes"`
		Rating             float64 `json:"rating"`
		Date               string  `json:"date"`
		Status             string  `json:"status"`
		PresetID           int     `json:"preset_id"`
		TabAccessType      string  `json:"tab_access_type"`
		TpVersion          int     `json:"tp_version"`
		TonalityName       string  `json:"tonality_name"`
		VersionDescription string  `json:"version_description"`
		Verified           int     `json:"verified"`
		Recording          struct {
			IsAcoustic       int               `json:"is_acoustic"`
			TonalityName     string            `json:"tonality_name"`
			Performance      Performance       `json:"performance"`
			RecordingArtists []RecordingArtist `json:"recording_artists"`
		} `json:"recording"`
	} `json:"versions"`
	Recommended []struct {
		ID                 int     `json:"id"`
		SongName           string  `json:"song_name"`
		ArtistName         string  `json:"artist_name"`
		Type               string  `json:"type"`
		Part               string  `json:"part"`
		Version            int     `json:"version"`
		Votes              int     `json:"votes"`
		Rating             float64 `json:"rating"`
		Date               string  `json:"date"`
		Status             string  `json:"status"`
		PresetID           int     `json:"preset_id"`
		TabAccessType      string  `json:"tab_access_type"`
		TpVersion          int     `json:"tp_version"`
		TonalityName       string  `json:"tonality_name"`
		VersionDescription *string `json:"version_description"`
		Verified           int     `json:"verified"`
		Recording          struct {
			IsAcoustic       int               `json:"is_acoustic"`
			TonalityName     string            `json:"tonality_name"`
			Performance      Performance       `json:"performance"`
			RecordingArtists []RecordingArtist `json:"recording_artists"`
		} `json:"recording"`
	} `json:"recommended"`
	UserRating  int    `json:"userRating"`
	Difficulty  string `json:"difficulty"`
	Tuning      string `json:"tuning"`
	Capo        int    `json:"capo"`
	URLWeb      string `json:"urlWeb"`
	VideosCount int    `json:"videosCount"`
	Contributor struct {
		UserID   int    `json:"user_id"`
		Username string `json:"username"`
	} `json:"contributor"`
	Applicature []Applicature `json:"applicature"`
	Content     string        `json:"content"`
}

type ListCapo struct {
	Fret        int64 `json:"fret"`
	StartString int64 `json:"startString"`
	LastString  int64 `json:"lastString"`
	Finger      int64 `json:"finger"`
}

type Applicature struct {
	Chord      string `json:"chord"`
	Variations []struct {
		ID        string     `json:"id"`
		ListCapos []ListCapo `json:"listCapos"`
		NoteIndex int        `json:"noteIndex"`
		Notes     []int      `json:"notes"`
		Frets     []int      `json:"frets"`
		Fingers   []int      `json:"fingers"`
		Fret      int        `json:"fret"`
	} `json:"variations"`
}

type Performance struct {
	Name      string `json:"name"`
	DateStart int64  `json:"date_start"`
	DateEnd   int64  `json:"date_end"`
	Cancelled int64  `json:"cancelled"`
	Type      string `json:"type"`
	Comment   string `json:"comment"`
}

type TonalityName string

type Recording struct {
	IsAcoustic       int64             `json:"is_acoustic"`
	TonalityName     string            `json:"tonality_name"`
	Performance      interface{}       `json:"performance"`
	RecordingArtists []RecordingArtist `json:"recording_artists"`
}

type RecordingArtist struct {
	JoinField string `json:"join_field"`
	Artist    Artist `json:"artist"`
}

type Artist struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Status string

const (
	Approved Status = "approved"
)

type TabAccessType string

const (
	Public TabAccessType = "public"
)

type Type string

const (
	Power Type = "Power"
)
