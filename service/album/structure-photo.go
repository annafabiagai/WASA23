package album

import (
	"fmt"
	"path/filepath"
)

type Photo struct {
	ID      uint64 `json:"photoID"`
	OwnerID uint64 `json:"ownerID"`
	Format  string `json:"format"`
	Date    string `json:"date"`
}

func (p *Photo) Path() (photoPath string) {
	return filepath.Join(Root, fmt.Sprintf("%d.%s", p.ID, p.Format))
}
