package helpers

import (
	"database/sql"
	"errors"
	"time"
)

func Nullable(data any) (any, error) {
	// Handle nil input
	if data == nil {
		return nil, nil
	}

	// Check for string type
	if v, ok := data.(string); ok {
		// Handle string
		if v == "" {
			return sql.NullString{Valid: false}, nil
		}
		return sql.NullString{String: v, Valid: true}, nil
	} else if v, ok := data.([]byte); ok {
		// Handle byte slice
		if len(v) == 0 {
			return sql.NullString{Valid: false}, nil
		}
		return sql.NullString{String: string(v), Valid: true}, nil
	} else if v, ok := data.(int); ok {
		// Handle int
		return sql.NullInt64{Int64: int64(v), Valid: true}, nil
	} else if v, ok := data.(int64); ok {
		// Handle int64
		return sql.NullInt64{Int64: v, Valid: true}, nil
	} else if v, ok := data.(float64); ok {
		// Handle float64
		return sql.NullFloat64{Float64: v, Valid: true}, nil
	} else if v, ok := data.(bool); ok {
		// Handle bool
		return sql.NullBool{Bool: v, Valid: true}, nil
	} else if v, ok := data.(time.Time); ok {
		// Handle time.Time
		return sql.NullTime{Time: v, Valid: true}, nil
	} else {
		// Handle unsupported types
		return nil, errors.New("unsupported type data")
	}
}
