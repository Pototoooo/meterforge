package http

import (
	"github.com/samber/lo"

	"github.com/Pototoooo/meterforge/api"
	"github.com/Pototoooo/meterforge/pkg/models"
)

func FromAnnotations(annotations models.Annotations) *api.Annotations {
	return lo.ToPtr((api.Annotations)(annotations))
}

func AsAnnotations(annotations *api.Annotations) models.Annotations {
	if annotations == nil {
		return nil
	}

	return (models.Annotations)(*annotations)
}
