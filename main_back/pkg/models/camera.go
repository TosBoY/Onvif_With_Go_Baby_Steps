package models

type Camera struct {
	ID       string `json:"id"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	IsFake   bool   `json:"isFake"`
}

type EncoderConfig struct {
	Resolution Resolution
	Quality    int
	FPS        int
	Bitrate    int
	Encoding   string
}

type EncoderOption struct {
	Resolutions []Resolution
	Quality     []int
	FPSOptions  []int
	Bitrate     []int
}

type Resolution struct {
	Width  int
	Height int
}
