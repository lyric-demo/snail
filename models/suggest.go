package models

import (
	"encoding/json"

	"github.com/tidwall/buntdb"
)

// Suggest 反馈建议
type Suggest struct {
	ID         string // 唯一ID
	AccessIP   string // 访问IP
	Email      string // 电子邮箱
	Comment    string // 意见
	CreateTime string // 创建时间
}

// Create 增加反馈建议
func (s *Suggest) Create() (err error) {
	buf, err := json.Marshal(s)
	if err != nil {
		return
	}
	err = SuggestDB.Update(func(tx *buntdb.Tx) (err error) {
		tx.Set(s.ID, string(buf), nil)
		return
	})
	return
}
