package db

import (
	"encoding/base64"
	"fmt"
	"strings"
)

type PaginationOptions struct {
	Direction string
	CursorID  string
	PageSize  int
}

func (p *PaginationOptions) Cursor() (*PaginationCursor, error) {
	return DecodePaginationCursor(p.CursorID)
}

type PaginationCursor struct {
	LaunchDate string
	ID         string
}

func EncodePaginationCursor(p *PaginationCursor) string {
	if p == nil {
		return ""
	}
	s := fmt.Sprintf("%s|%s", p.ID, p.LaunchDate)
	return base64.URLEncoding.EncodeToString([]byte(s))
}

func DecodePaginationCursor(coded string) (c *PaginationCursor, err error) {
	if coded == "" {
		c = &PaginationCursor{}
		return
	}
	s, err := base64.URLEncoding.DecodeString(coded)
	if err != nil {
		return
	}

	splits := strings.Split(string(s), "|")

	if len(splits) != 2 {
		err = fmt.Errorf("invalid cursor id, cannot decode")
		return
	}

	c = &PaginationCursor{
		ID:         splits[0],
		LaunchDate: splits[1],
	}

	return
}
