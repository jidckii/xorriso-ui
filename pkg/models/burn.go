package models

import "time"

type BurnState string

const (
	BurnStatePending   BurnState = "pending"
	BurnStateFormat    BurnState = "formatting"
	BurnStateWriting   BurnState = "writing"
	BurnStateVerifying BurnState = "verifying"
	BurnStateDone      BurnState = "done"
	BurnStateError     BurnState = "error"
	BurnStateCancelled BurnState = "cancelled"
)

type BurnJob struct {
	ID         string       `json:"id"`
	State      BurnState    `json:"state"`
	Progress   BurnProgress `json:"progress"`
	Result     *BurnResult  `json:"result,omitempty"`
	StartedAt  time.Time    `json:"startedAt"`
	FinishedAt time.Time    `json:"finishedAt,omitempty"`
	Error      string       `json:"error,omitempty"`
}

type BurnProgress struct {
	Phase        string  `json:"phase"`
	Percent      float64 `json:"percent"`
	Speed        string  `json:"speed"`
	BytesWritten int64   `json:"bytesWritten"`
	BytesTotal   int64   `json:"bytesTotal"`
	ETA          string  `json:"eta"`
	FIFOFill     int     `json:"fifoFill"`
}

type BurnResult struct {
	Success      bool   `json:"success"`
	BytesWritten int64  `json:"bytesWritten"`
	Duration     string `json:"duration"`
	AverageSpeed string `json:"averageSpeed"`
	MD5Match     bool   `json:"md5Match"`
	VerifyErrors int    `json:"verifyErrors"`
}
