package types

import (
	"time"
)

// MiniCommit mini-commitのデータ構造
type MiniCommit struct {
	ID        string    `json:"id"`        // SHA1ハッシュ
	Message   string    `json:"message"`   // コミットメッセージ
	CreatedAt time.Time `json:"createdAt"` // 作成日時
	Patch     string    `json:"patch"`     // 差分（patch形式）
}

// MiniCommitList mini-commitの一覧
type MiniCommitList []MiniCommit
