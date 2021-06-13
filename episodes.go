package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

type episode struct {
	ID          int64  `json:"id"`
	PodcastID   int64  `json:"podcast_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	SeasonNo    int    `json:"season_no"`
	EpisodeNo   int    `json:"episode_no"`
	// 1 = Full | 2 = Bonus | 3 = trailer
	TypeOfEpisode      int       `json:"type_of_episode"`
	IsExplicit         bool      `json:"isExplicit"`
	EpisodeArtID       int64     `json:"episode_art_id,omitempty"`
	EpisodeArtPath     int64     `json:"episode_art_path,omitempty"`
	EpisodeContentID   int64     `json:"episode_content_id,omitempty"`
	EpisodeContentPath int64     `json:"episode_content_path,omitempty"`
	Published          bool      `json:"published"`
	PublishedAt        time.Time `json:"published_at,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func createEpisodes(c echo.Context) (err error) {

	ep := &episode{}

	if err = c.Bind(ep); err != nil {
		return
	}

	if ep.PodcastID == 0 || ep.Title == "" || ep.Description == "" {
		return c.String(http.StatusBadRequest, "Insufficient fields")
	}

	ep.CreatedAt = time.Now()
	ep.UpdatedAt = time.Now()

	q := `
		INSERT INTO episodes (
			podcast_id, title, description, season_no, episode_no, type_of_episode, 
			is_explicit, episode_art_id, episode_content_id, 
			published, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id
	`
	err = db.QueryRow(q, ep.PodcastID, ep.Title, ep.Description, ep.SeasonNo, ep.EpisodeNo, ep.TypeOfEpisode, ep.IsExplicit, ep.EpisodeArtID, ep.EpisodeContentID, ep.Published, ep.CreatedAt, ep.UpdatedAt).Scan(&ep.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	return c.JSON(http.StatusOK, ep)
}

func getEpisodes(c echo.Context) (err error) {

	podcastID, err := strconv.Atoi(c.QueryParam("podcast_id"))
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, "podcast_id is required")
	}

	eps, err := _getPodcastEpisodes(int64(podcastID))
	if err != nil {
		fmt.Println(err)
		return
	}

	return c.JSON(http.StatusOK, eps)
}

func _getPodcastEpisodes(podcastID int64) (eps []episode, err error) {
	q := `SELECT id, podcast_id, title, description, season_no, episode_no,
	type_of_episode, is_explicit, episode_art_id, episode_content_id,
	published, created_at, updated_at FROM episodes WHERE podcast_id = $1
`
	rows, err := db.Query(q, podcastID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var ep episode

		err = rows.Scan(&ep.ID, &ep.PodcastID, &ep.Title, &ep.Description, &ep.SeasonNo, &ep.EpisodeNo, &ep.TypeOfEpisode, &ep.IsExplicit, &ep.EpisodeArtID, &ep.EpisodeContentID, &ep.Published, &ep.CreatedAt, &ep.UpdatedAt)

		if err != nil {
			fmt.Println(err)
			continue
		}

		eps = append(eps, ep)
	}
	return eps, nil
}
