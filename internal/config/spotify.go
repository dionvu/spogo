package config

type Spotify struct {
	clientID     string `yaml:"client_id"`
	clientSecret string `yaml:"client_secret"`
}

func (s *Spotify) ClientID() string {
	return s.clientID
}

func (s *Spotify) ClientSecret() string {
	return s.clientSecret
}

func (s *Spotify) setID(str string) {
	s.clientID = str
}

func (s *Spotify) setSecret(str string) {
	s.clientSecret = str
}
