package models

import (
	"encoding/json"

	"github.com/tidwall/buntdb"
)

// History 下载历史
type History struct {
	ID         string // 唯一ID
	AccessIP   string // 访问IP
	FileLink   string // 下载链接
	FileName   string // 文件名
	FileSize   int    // 文件大小
	FileType   int    // 文件类型（0：普通文件，1：视频文件）
	CreateTime string // 创建时间
}

// Create 增加下载历史
func (h *History) Create() (err error) {
	buf, err := json.Marshal(h)
	if err != nil {
		return
	}
	err = HistoryDB.Update(func(tx *buntdb.Tx) (err error) {
		tx.Set(h.ID, string(buf), nil)
		return
	})
	return
}
