package utils

import (
	"github.com/google/uuid"
)

func UUID() string {
	u := uuid.New()
	return u.String()
}
