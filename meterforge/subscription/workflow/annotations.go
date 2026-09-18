package subscriptionworkflow

import (
	"github.com/oklog/ulid/v2"

	"github.com/Pototoooo/meterforge/pkg/models"
)

const (
	AnnotationEditUniqueKey = "subscription.workflow.patchid"
)

type annotationParser struct{}

var AnnotationParser = annotationParser{}

func (a annotationParser) SetUniquePatchID(annotations models.Annotations) models.Annotations {
	id := ulid.Make().String()

	if annotations == nil {
		annotations = models.Annotations{}
	}

	annotations[AnnotationEditUniqueKey] = id

	return annotations
}
