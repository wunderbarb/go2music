package audio

type AVTransport struct {
	CurrentTransportState string `json:"CurrentTransportState"`
	CurrentTrack          struct {
		TrackNumber string `json:"TrackNumber"`
		Title       string `json:"Title"`
		Artist      string `json:"Artist"`
	} `json:"CurrentTrack"`
	CurrentPosition string `json:"CurrentPosition"`
	TransportStatus string `json:"TransportStatus"`
}
