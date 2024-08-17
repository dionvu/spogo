package spotify

import "fmt"

type Show struct {
	Description   string  `json:"description"`
	ID            string  `json:"id"`
	Images        []Image `json:"images"`
	Name          string  `json:"name"`
	Uri           string  `json:"uri"`
	TotalEpisodes int     `json:"total_episodes"`
}

func (s *Show) String() string {
	return fmt.Sprintf("%s%v", s.Name+" | ", s.TotalEpisodes)
}

type Episode struct {
	Description string      `json:"description"`
	DurationMs  int         `json:"duration_ms"`
	ID          string      `json:"id"`
	Images      []Image     `json:"images"`
	Name        string      `json:"name"`
	ReleaseDate string      `json:"release_date"`
	ResumePoint ResumePoint `json:"resume_point"`
	Uri         string      `json:"uri"`
}

func (ep *Episode) String() string {
	return fmt.Sprintf("%v", ep.Name)
}

type ResumePoint struct {
	FullyPlayed      bool `json:"fully_played"`
	ResumePositionMs int  `json:"resume_position_ms"`
}
