package models

import "time"

type Project struct {
	Name        string      `json:"name"`
	FilePath    string      `json:"filePath"`
	VolumeID    string      `json:"volumeId"`
	Entries     []FileEntry `json:"entries"`
	ISOOptions  ISOOptions  `json:"isoOptions"`
	BurnOptions BurnOptions `json:"burnOptions"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

type FileEntry struct {
	SourcePath string `json:"sourcePath"`
	DestPath   string `json:"destPath"`
	Name       string `json:"name"`
	IsDir      bool   `json:"isDir"`
	Size       int64  `json:"size"`
	ModTime    int64  `json:"modTime"` // Unix timestamp в миллисекундах
}

type ISOOptions struct {
	UDF        bool `json:"udf"`
	ISOLevel   int  `json:"isoLevel"`
	RockRidge  bool `json:"rockRidge"`
	Joliet     bool `json:"joliet"`
	HFSPlus    bool `json:"hfsPlus"`
	Zisofs     bool `json:"zisofs"`
	MD5        bool `json:"md5"`
	BackupMode bool `json:"backupMode"`
}

type BurnOptions struct {
	Speed           string `json:"speed"`
	DummyMode       bool   `json:"dummyMode"`
	Verify          bool   `json:"verify"`
	CloseDisc       bool   `json:"closeDisc"`
	StreamRecording bool   `json:"streamRecording"`
	Eject           bool   `json:"eject"`
	BurnMode        string `json:"burnMode"`
	Padding         int    `json:"padding"`
	Multisession    bool   `json:"multisession"`
	CleanupISO      bool   `json:"cleanupIso"`
}
