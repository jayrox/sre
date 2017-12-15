package main

import "bytes"
import "encoding/json"
import "flag"
import "fmt"
import "io/ioutil"
import "log"
import "math/rand"
import "net/http"
import "os"
import "strconv"
import "strings"
import "time"

//import "github.com/davecgh/go-spew/spew"
import "github.com/kardianos/osext"

type sonarrMissingEpisodes struct {
	Page          int    `json:"page"`
	PageSize      int    `json:"pageSize"`
	SortKey       string `json:"sortKey"`
	SortDirection string `json:"sortDirection"`
	TotalRecords  int    `json:"totalRecords"`
	Records       []struct {
		SeriesID                 int       `json:"seriesId"`
		EpisodeFileID            int       `json:"episodeFileId"`
		SeasonNumber             int       `json:"seasonNumber"`
		EpisodeNumber            int       `json:"episodeNumber"`
		Title                    string    `json:"title"`
		AirDate                  string    `json:"airDate"`
		AirDateUtc               time.Time `json:"airDateUtc"`
		Overview                 string    `json:"overview"`
		HasFile                  bool      `json:"hasFile"`
		Monitored                bool      `json:"monitored"`
		UnverifiedSceneNumbering bool      `json:"unverifiedSceneNumbering"`
		Series                   struct {
			Title       string `json:"title"`
			SortTitle   string `json:"sortTitle"`
			SeasonCount int    `json:"seasonCount"`
			Status      string `json:"status"`
			Overview    string `json:"overview"`
			Network     string `json:"network"`
			AirTime     string `json:"airTime"`
			Images      []struct {
				CoverType string `json:"coverType"`
				URL       string `json:"url"`
			} `json:"images"`
			Seasons []struct {
				SeasonNumber int  `json:"seasonNumber"`
				Monitored    bool `json:"monitored"`
			} `json:"seasons"`
			Year              int           `json:"year"`
			Path              string        `json:"path"`
			ProfileID         int           `json:"profileId"`
			SeasonFolder      bool          `json:"seasonFolder"`
			Monitored         bool          `json:"monitored"`
			UseSceneNumbering bool          `json:"useSceneNumbering"`
			Runtime           int           `json:"runtime"`
			TvdbID            int           `json:"tvdbId"`
			TvRageID          int           `json:"tvRageId"`
			TvMazeID          int           `json:"tvMazeId"`
			FirstAired        time.Time     `json:"firstAired"`
			LastInfoSync      time.Time     `json:"lastInfoSync"`
			SeriesType        string        `json:"seriesType"`
			CleanTitle        string        `json:"cleanTitle"`
			TitleSlug         string        `json:"titleSlug"`
			Genres            []string      `json:"genres"`
			Tags              []interface{} `json:"tags"`
			Added             time.Time     `json:"added"`
			Ratings           struct {
				Votes int     `json:"votes"`
				Value float64 `json:"value"`
			} `json:"ratings"`
			QualityProfileID int `json:"qualityProfileId"`
			ID               int `json:"id"`
		} `json:"series"`
		ID int `json:"id"`
	} `json:"records"`
}

type sonarrEpisode struct {
	SeriesID                 int       `json:"seriesId"`
	EpisodeFileID            int       `json:"episodeFileId"`
	SeasonNumber             int       `json:"seasonNumber"`
	EpisodeNumber            int       `json:"episodeNumber"`
	Title                    string    `json:"title"`
	AirDate                  string    `json:"airDate"`
	AirDateUtc               time.Time `json:"airDateUtc"`
	Overview                 string    `json:"overview"`
	HasFile                  bool      `json:"hasFile"`
	Monitored                bool      `json:"monitored"`
	AbsoluteEpisodeNumber    int       `json:"absoluteEpisodeNumber"`
	UnverifiedSceneNumbering bool      `json:"unverifiedSceneNumbering"`
	Series                   struct {
		Title       string `json:"title"`
		SortTitle   string `json:"sortTitle"`
		SeasonCount int    `json:"seasonCount"`
		Status      string `json:"status"`
		Overview    string `json:"overview"`
		Network     string `json:"network"`
		Images      []struct {
			CoverType string `json:"coverType"`
			URL       string `json:"url"`
		} `json:"images"`
		Seasons []struct {
			SeasonNumber int  `json:"seasonNumber"`
			Monitored    bool `json:"monitored"`
		} `json:"seasons"`
		Year              int           `json:"year"`
		Path              string        `json:"path"`
		ProfileID         int           `json:"profileId"`
		SeasonFolder      bool          `json:"seasonFolder"`
		Monitored         bool          `json:"monitored"`
		UseSceneNumbering bool          `json:"useSceneNumbering"`
		Runtime           int           `json:"runtime"`
		TvdbID            int           `json:"tvdbId"`
		TvRageID          int           `json:"tvRageId"`
		TvMazeID          int           `json:"tvMazeId"`
		FirstAired        time.Time     `json:"firstAired"`
		LastInfoSync      time.Time     `json:"lastInfoSync"`
		SeriesType        string        `json:"seriesType"`
		CleanTitle        string        `json:"cleanTitle"`
		TitleSlug         string        `json:"titleSlug"`
		Certification     string        `json:"certification"`
		Genres            []string      `json:"genres"`
		Tags              []interface{} `json:"tags"`
		Added             time.Time     `json:"added"`
		Ratings           struct {
			Votes int     `json:"votes"`
			Value float64 `json:"value"`
		} `json:"ratings"`
		QualityProfileID int `json:"qualityProfileId"`
		ID               int `json:"id"`
	} `json:"series"`
	ID int `json:"id"`
}

type sonarrSearchResult struct {
	Name string `json:"name"`
	Body struct {
		EpisodeIds          []int  `json:"episodeIds"`
		SendUpdatesToClient bool   `json:"sendUpdatesToClient"`
		UpdateScheduledTask bool   `json:"updateScheduledTask"`
		CompletionMessage   string `json:"completionMessage"`
		Name                string `json:"name"`
		Trigger             string `json:"trigger"`
	} `json:"body"`
	Priority            string    `json:"priority"`
	Status              string    `json:"status"`
	Queued              time.Time `json:"queued"`
	Trigger             string    `json:"trigger"`
	State               string    `json:"state"`
	Manual              bool      `json:"manual"`
	StartedOn           time.Time `json:"startedOn"`
	SendUpdatesToClient bool      `json:"sendUpdatesToClient"`
	UpdateScheduledTask bool      `json:"updateScheduledTask"`
	ID                  int       `json:"id"`
}

type sonarrSearchEpisodeCommand struct {
	Name       string `json:"name"`
	EpisodeIds []int  `json:"episodeIds"`
}

type configuration struct {
	LogLocation string `json:"loglocation"`
	HostName    string `json:"hostname"`
	BaseURL     string `json:"baseurl"`
	HostPort    int    `json:"port"`
	APIKey      string `json:"apikey"`
}

var (
	flagIndex = flag.Int("i", -1, "index to search for. (if used, skips the random index and uses defined)")
	cfg       configuration
)

func main() {
	flag.Parse()

	cfgpath, _ := osext.ExecutableFolder()
	file, _ := os.Open(cfgpath + "/sre.json")
	decoder := json.NewDecoder(file)
	cfg = configuration{}
	err := decoder.Decode(&cfg)
	if err != nil {
		fmt.Println("error:", err)
	}

	if cfg.HostName == "" {
		cfg.HostName = "http://localhost"
	}
	if cfg.HostPort == 0 {
		cfg.HostPort = 8989
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = ""
	}
	if cfg.APIKey == "" {
		writeToLog("ERROR: Invalid ApiKey")
		return
	}

	urlRoot := cfg.HostName + ":" + strconv.Itoa(cfg.HostPort)
	if cfg.BaseURL != "" {
		if strings.Index(cfg.BaseURL, "/") == -1 {
			urlRoot += "/"
		}
		urlRoot += cfg.BaseURL
	}
	urlRoot += "/api"
	apiKey := cfg.APIKey

	randomEpisode := getRandomSonarrEpisode(urlRoot, apiKey)
	writeToLog(fmt.Sprintf("Random Episode ID: %d", randomEpisode))
	fmt.Printf("Random Episode ID: %d\n", randomEpisode)

	// Get Episode info
	episode := getSonarrEpisodeInfo(urlRoot, apiKey, randomEpisode)
	writeToLog(fmt.Sprintf("Searching: %s - S%dE%d - %s", episode.Series.Title, episode.SeasonNumber, episode.EpisodeNumber, episode.Title))
	fmt.Printf("Searching: %s - S%dE%d - %s", episode.Series.Title, episode.SeasonNumber, episode.EpisodeNumber, episode.Title)

	episodesURL := fmt.Sprintf("%s/command", urlRoot)
	writeToLog(episodesURL)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(struct {
		Name       string `json:"name"`
		EpisodeIds []int
	}{
		"episodeSearch",
		[]int{randomEpisode},
	})

	writeToLog(b.String())

	client := &http.Client{}
	req, err := http.NewRequest("POST", episodesURL, b)
	req.Header.Add("X-Api-Key", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	writeToLog(fmt.Sprintf("Resp Status: %s", resp.Status))

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		writeToLog(fmt.Sprintf("%s", err.Error()))
		os.Exit(1)
	}
	writeToLog("Response:\n")
	writeToLog(string(body))
	writeToLog("\n___________________________\n")
}

func getSonarrTotalRecords(urlRoot, apiKey string) int {
	page := 1
	pageSize := 1
	urlPath := fmt.Sprintf("%s/wanted/missing/?apikey=%s&page=%d&pageSize=%d&sortKey=airDateUtc&sortDir=asc", urlRoot, apiKey, page, pageSize)

	res, err := http.Get(urlPath)

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err.Error())
	}

	var data sonarrMissingEpisodes
	json.Unmarshal(body, &data)

	writeToLog(fmt.Sprintf("Total Records: %d", data.TotalRecords))
	return data.TotalRecords
}

func getRandomSonarrEpisode(urlRoot, apiKey string) int {
	records := getSonarrTotalRecords(urlRoot, apiKey)

	page := flagInt(flagIndex)

	if flagInt(flagIndex) == -1 {
		rand.Seed(time.Now().Unix())
		page = rand.Intn(records-1) + 1
	}

	writeToLog(fmt.Sprintf("Rand Record: %d", page))
	pageSize := 1
	urlPath := fmt.Sprintf("%s/wanted/missing/?apikey=%s&page=%d&pageSize=%d&sortKey=airDateUtc&sortDir=asc", urlRoot, apiKey, page, pageSize)

	res, err := http.Get(urlPath)

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err.Error())
	}

	var data sonarrMissingEpisodes
	json.Unmarshal(body, &data)

	return data.Records[0].ID
}

func getSonarrEpisodeInfo(urlRoot, apiKey string, episodeId int) (episode sonarrEpisode) {
	urlPath := fmt.Sprintf("%s/episode/%d?apikey=%s", urlRoot, episodeId, apiKey)

	res, err := http.Get(urlPath)

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err.Error())
	}

	var data sonarrEpisode
	json.Unmarshal(body, &data)

	return data
}

func flagInt(fs *int) int {
	return *fs
}

func writeToLog(str string) {
	filename := cfg.LogLocation
	if cfg.LogLocation == "" {
		filename, _ = osext.ExecutableFolder()
	}

	f, err := os.OpenFile(filename+"/sre.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(str)
}
