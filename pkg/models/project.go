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
	SourcePath string      `json:"sourcePath"`
	DestPath   string      `json:"destPath"`
	Name       string      `json:"name"`
	IsDir      bool        `json:"isDir"`
	Size       int64       `json:"size"`
	Children   []FileEntry `json:"children,omitempty"`
}

type ISOOptions struct {
	RockRidge    bool   `json:"rockRidge"`
	Joliet       bool   `json:"joliet"`
	HFSPlus      bool   `json:"hfsPlus"`
	Zisofs       bool   `json:"zisofs"`
	MD5          bool   `json:"md5"`
	BackupMode   bool   `json:"backupMode"`
	BootImage    string `json:"bootImage,omitempty"`
	BootCatalog  string `json:"bootCatalog,omitempty"`
	EFIBootImage string `json:"efiBootImage,omitempty"`
	BootMode     string `json:"bootMode,omitempty"`
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
}
