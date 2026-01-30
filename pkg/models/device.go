package models

type Device struct {
	Path        string         `json:"path"`
	LinkPath    string         `json:"linkPath"`
	Vendor      string         `json:"vendor"`
	Model       string         `json:"model"`
	Revision    string         `json:"revision"`
	MediaLoaded bool           `json:"mediaLoaded"`
	Profiles    []MediaProfile `json:"profiles"`
}

type MediaInfo struct {
	DevicePath    string            `json:"devicePath"`
	MediaType     string            `json:"mediaType"`
	MediaStatus   string            `json:"mediaStatus"`
	Sessions      int               `json:"sessions"`
	FreeSpace     int64             `json:"freeSpace"`
	UsedSpace     int64             `json:"usedSpace"`
	TotalCapacity int64             `json:"totalCapacity"`
	Speeds        []SpeedDescriptor `json:"speeds"`
	Erasable      bool              `json:"erasable"`
}

type SpeedDescriptor struct {
	WriteSpeed  float64 `json:"writeSpeed"`
	DisplayName string  `json:"displayName"`
}

type MediaProfile struct {
	Name    string `json:"name"`
	Current bool   `json:"current"`
}

type TOC struct {
	Sessions []Session `json:"sessions"`
}

type Session struct {
	Number   int    `json:"number"`
	StartLBA int64  `json:"startLba"`
	Size     int64  `json:"size"`
	VolumeID string `json:"volumeId"`
}
