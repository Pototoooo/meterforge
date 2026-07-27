package apiconverter

import (
	"github.com/Pototoooo/meterforge/api"
	"github.com/Pototoooo/meterforge/pkg/pagination/v2"
)

func ConvertCursor(s api.CursorPaginationCursor) (*pagination.Cursor, error) {
	return pagination.DecodeCursor(s)
}

func ConvertCursorPtr(s *api.CursorPaginationCursor) (*pagination.Cursor, error) {
	if s == nil {
		return nil, nil
	}

	return ConvertCursor(*s)
}
