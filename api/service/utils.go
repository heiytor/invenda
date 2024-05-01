package service

import (
	"github.com/heiytor/invenda/api/pkg/errors"
	"github.com/heiytor/invenda/api/store"
	"github.com/rs/zerolog/log"
)

// mapError maps an error "in" to "out", it receives an optional entity to
// write as an attribute "entity".
func mapError(in error, entity string) error {
	if in == nil {
		return nil
	}

	out := errors.New().Layer(errors.LayerService)

	if entity != "" {
		out.Attr("entity", entity)
	}

	switch {
	case errors.Is(in, store.ErrNotFound):
		return out.Code(404).Msg(errors.MsgNotFound)
	default:
		// default branch handle non-expected mongo errors.
		// TODO: send to sentry
		log.Error().Err(in).Msg(errors.MsgUnexpected)

		return out.Layer(errors.LayerStore).Code(500).Msg("internal server error")
	}
}
