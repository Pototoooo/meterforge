package httpdriver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Pototoooo/meterforge/meterforge/entitlement"
	"github.com/Pototoooo/meterforge/meterforge/namespace/namespacedriver"
	"github.com/Pototoooo/meterforge/meterforge/subject"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	SubjectHandler
}

type SubjectHandler interface {
	GetSubject() GetSubjectHandler
	ListSubjects() ListSubjectsHandler
	UpsertSubject() UpsertSubjectHandler
	DeleteSubject() DeleteSubjectHandler
}

var _ Handler = (*handler)(nil)

type handler struct {
	namespaceDecoder     namespacedriver.NamespaceDecoder
	options              []httptransport.HandlerOption
	logger               *slog.Logger
	subjectService       subject.Service
	entitlementConnector entitlement.Service
}

func (h *handler) resolveNamespace(ctx context.Context) (string, error) {
	ns, ok := h.namespaceDecoder.GetNamespace(ctx)
	if !ok {
		return "", commonhttp.NewHTTPError(http.StatusInternalServerError, errors.New("internal server error"))
	}

	return ns, nil
}

func New(
	namespaceDecoder namespacedriver.NamespaceDecoder,
	logger *slog.Logger,
	subjectService subject.Service,
	entitlementConnector entitlement.Service,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		namespaceDecoder:     namespaceDecoder,
		options:              options,
		logger:               logger,
		subjectService:       subjectService,
		entitlementConnector: entitlementConnector,
	}
}
