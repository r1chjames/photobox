package database

import (
	"fmt"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"gorm.io/gorm"
	"testing"
)

var testDb *gorm.DB

func newAlbumRecord() Album {
	return Album{
		ID:          "a",
		Name:        "b",
		Description: "c",
	}
}

func albumMocker(n int64) []Album {
	var offset int64
	testDb.Model(&Album{}).Count(&offset)
	var ret []Album
	for i := offset + 1; i <= offset+n; i++ {
		userModel := Album{
			ID:          fmt.Sprintf("user%v", i),
			Name:        fmt.Sprintf("user%v@linkedin.com", i),
			Description: fmt.Sprintf("bio%v", i),
		}
		testDb.Create(&userModel)
		ret = append(ret, userModel)
	}
	return ret
}
