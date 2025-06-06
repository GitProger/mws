// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"
	"net/http"
	"time"

	"github.com/go-faster/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/otelogen"
)

type codeRecorder struct {
	http.ResponseWriter
	status int
}

func (c *codeRecorder) WriteHeader(status int) {
	c.status = status
	c.ResponseWriter.WriteHeader(status)
}

// handleAddUserBookRequest handles addUserBook operation.
//
// Appends a book to user's list if the list doesn't contain it, if contains, an error returned.
//
// POST /users/{user_id}/books
func (s *Server) handleAddUserBookRequest(args [1]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	statusWriter := &codeRecorder{ResponseWriter: w}
	w = statusWriter
	otelAttrs := []attribute.KeyValue{
		otelogen.OperationID("addUserBook"),
		semconv.HTTPRequestMethodKey.String("POST"),
		semconv.HTTPRouteKey.String("/users/{user_id}/books"),
	}

	// Start a span for this request.
	ctx, span := s.cfg.Tracer.Start(r.Context(), AddUserBookOperation,
		trace.WithAttributes(otelAttrs...),
		serverSpanKind,
	)
	defer span.End()

	// Add Labeler to context.
	labeler := &Labeler{attrs: otelAttrs}
	ctx = contextWithLabeler(ctx, labeler)

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		elapsedDuration := time.Since(startTime)

		attrSet := labeler.AttributeSet()
		attrs := attrSet.ToSlice()
		code := statusWriter.status
		if code != 0 {
			codeAttr := semconv.HTTPResponseStatusCode(code)
			attrs = append(attrs, codeAttr)
			span.SetAttributes(codeAttr)
		}
		attrOpt := metric.WithAttributes(attrs...)

		// Increment request counter.
		s.requests.Add(ctx, 1, attrOpt)

		// Use floating point division here for higher precision (instead of Millisecond method).
		s.duration.Record(ctx, float64(elapsedDuration)/float64(time.Millisecond), attrOpt)
	}()

	var (
		recordError = func(stage string, err error) {
			span.RecordError(err)

			// https://opentelemetry.io/docs/specs/semconv/http/http-spans/#status
			// Span Status MUST be left unset if HTTP status code was in the 1xx, 2xx or 3xx ranges,
			// unless there was another error (e.g., network error receiving the response body; or 3xx codes with
			// max redirects exceeded), in which case status MUST be set to Error.
			code := statusWriter.status
			if code >= 100 && code < 500 {
				span.SetStatus(codes.Error, stage)
			}

			attrSet := labeler.AttributeSet()
			attrs := attrSet.ToSlice()
			if code != 0 {
				attrs = append(attrs, semconv.HTTPResponseStatusCode(code))
			}

			s.errors.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: AddUserBookOperation,
			ID:   "addUserBook",
		}
	)
	params, err := decodeAddUserBookParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}
	request, close, err := s.decodeAddUserBookRequest(r)
	if err != nil {
		err = &ogenerrors.DecodeRequestError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeRequest", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}
	defer func() {
		if err := close(); err != nil {
			recordError("CloseRequest", err)
		}
	}()

	var response AddUserBookRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    AddUserBookOperation,
			OperationSummary: "Add a new book for user",
			OperationID:      "addUserBook",
			Body:             request,
			Params: middleware.Parameters{
				{
					Name: "user_id",
					In:   "path",
				}: params.UserID,
			},
			Raw: r,
		}

		type (
			Request  = *Book
			Params   = AddUserBookParams
			Response = AddUserBookRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackAddUserBookParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.AddUserBook(ctx, request, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.AddUserBook(ctx, request, params)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeAddUserBookResponse(response, w, span); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleGetUserBookRequest handles getUserBook operation.
//
// Returns a book by user's and book's ids.
//
// GET /users/{user_id}/books/{book_id}
func (s *Server) handleGetUserBookRequest(args [2]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	statusWriter := &codeRecorder{ResponseWriter: w}
	w = statusWriter
	otelAttrs := []attribute.KeyValue{
		otelogen.OperationID("getUserBook"),
		semconv.HTTPRequestMethodKey.String("GET"),
		semconv.HTTPRouteKey.String("/users/{user_id}/books/{book_id}"),
	}

	// Start a span for this request.
	ctx, span := s.cfg.Tracer.Start(r.Context(), GetUserBookOperation,
		trace.WithAttributes(otelAttrs...),
		serverSpanKind,
	)
	defer span.End()

	// Add Labeler to context.
	labeler := &Labeler{attrs: otelAttrs}
	ctx = contextWithLabeler(ctx, labeler)

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		elapsedDuration := time.Since(startTime)

		attrSet := labeler.AttributeSet()
		attrs := attrSet.ToSlice()
		code := statusWriter.status
		if code != 0 {
			codeAttr := semconv.HTTPResponseStatusCode(code)
			attrs = append(attrs, codeAttr)
			span.SetAttributes(codeAttr)
		}
		attrOpt := metric.WithAttributes(attrs...)

		// Increment request counter.
		s.requests.Add(ctx, 1, attrOpt)

		// Use floating point division here for higher precision (instead of Millisecond method).
		s.duration.Record(ctx, float64(elapsedDuration)/float64(time.Millisecond), attrOpt)
	}()

	var (
		recordError = func(stage string, err error) {
			span.RecordError(err)

			// https://opentelemetry.io/docs/specs/semconv/http/http-spans/#status
			// Span Status MUST be left unset if HTTP status code was in the 1xx, 2xx or 3xx ranges,
			// unless there was another error (e.g., network error receiving the response body; or 3xx codes with
			// max redirects exceeded), in which case status MUST be set to Error.
			code := statusWriter.status
			if code >= 100 && code < 500 {
				span.SetStatus(codes.Error, stage)
			}

			attrSet := labeler.AttributeSet()
			attrs := attrSet.ToSlice()
			if code != 0 {
				attrs = append(attrs, semconv.HTTPResponseStatusCode(code))
			}

			s.errors.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: GetUserBookOperation,
			ID:   "getUserBook",
		}
	)
	params, err := decodeGetUserBookParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response GetUserBookRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    GetUserBookOperation,
			OperationSummary: "Get book by it's id",
			OperationID:      "getUserBook",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "user_id",
					In:   "path",
				}: params.UserID,
				{
					Name: "book_id",
					In:   "path",
				}: params.BookID,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = GetUserBookParams
			Response = GetUserBookRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackGetUserBookParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.GetUserBook(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.GetUserBook(ctx, params)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeGetUserBookResponse(response, w, span); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleGetUserBooksRequest handles getUserBooks operation.
//
// Returns list of user's books by their id.
//
// GET /users/{user_id}/books
func (s *Server) handleGetUserBooksRequest(args [1]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	statusWriter := &codeRecorder{ResponseWriter: w}
	w = statusWriter
	otelAttrs := []attribute.KeyValue{
		otelogen.OperationID("getUserBooks"),
		semconv.HTTPRequestMethodKey.String("GET"),
		semconv.HTTPRouteKey.String("/users/{user_id}/books"),
	}

	// Start a span for this request.
	ctx, span := s.cfg.Tracer.Start(r.Context(), GetUserBooksOperation,
		trace.WithAttributes(otelAttrs...),
		serverSpanKind,
	)
	defer span.End()

	// Add Labeler to context.
	labeler := &Labeler{attrs: otelAttrs}
	ctx = contextWithLabeler(ctx, labeler)

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		elapsedDuration := time.Since(startTime)

		attrSet := labeler.AttributeSet()
		attrs := attrSet.ToSlice()
		code := statusWriter.status
		if code != 0 {
			codeAttr := semconv.HTTPResponseStatusCode(code)
			attrs = append(attrs, codeAttr)
			span.SetAttributes(codeAttr)
		}
		attrOpt := metric.WithAttributes(attrs...)

		// Increment request counter.
		s.requests.Add(ctx, 1, attrOpt)

		// Use floating point division here for higher precision (instead of Millisecond method).
		s.duration.Record(ctx, float64(elapsedDuration)/float64(time.Millisecond), attrOpt)
	}()

	var (
		recordError = func(stage string, err error) {
			span.RecordError(err)

			// https://opentelemetry.io/docs/specs/semconv/http/http-spans/#status
			// Span Status MUST be left unset if HTTP status code was in the 1xx, 2xx or 3xx ranges,
			// unless there was another error (e.g., network error receiving the response body; or 3xx codes with
			// max redirects exceeded), in which case status MUST be set to Error.
			code := statusWriter.status
			if code >= 100 && code < 500 {
				span.SetStatus(codes.Error, stage)
			}

			attrSet := labeler.AttributeSet()
			attrs := attrSet.ToSlice()
			if code != 0 {
				attrs = append(attrs, semconv.HTTPResponseStatusCode(code))
			}

			s.errors.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: GetUserBooksOperation,
			ID:   "getUserBooks",
		}
	)
	params, err := decodeGetUserBooksParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response []Book
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    GetUserBooksOperation,
			OperationSummary: "Get all user's books with current progresses",
			OperationID:      "getUserBooks",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "user_id",
					In:   "path",
				}: params.UserID,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = GetUserBooksParams
			Response = []Book
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackGetUserBooksParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.GetUserBooks(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.GetUserBooks(ctx, params)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeGetUserBooksResponse(response, w, span); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleRemoveUserBookRequest handles removeUserBook operation.
//
// Removes a book by id if exists, otherwise an error returned.
//
// DELETE /users/{user_id}/books/{book_id}
func (s *Server) handleRemoveUserBookRequest(args [2]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	statusWriter := &codeRecorder{ResponseWriter: w}
	w = statusWriter
	otelAttrs := []attribute.KeyValue{
		otelogen.OperationID("removeUserBook"),
		semconv.HTTPRequestMethodKey.String("DELETE"),
		semconv.HTTPRouteKey.String("/users/{user_id}/books/{book_id}"),
	}

	// Start a span for this request.
	ctx, span := s.cfg.Tracer.Start(r.Context(), RemoveUserBookOperation,
		trace.WithAttributes(otelAttrs...),
		serverSpanKind,
	)
	defer span.End()

	// Add Labeler to context.
	labeler := &Labeler{attrs: otelAttrs}
	ctx = contextWithLabeler(ctx, labeler)

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		elapsedDuration := time.Since(startTime)

		attrSet := labeler.AttributeSet()
		attrs := attrSet.ToSlice()
		code := statusWriter.status
		if code != 0 {
			codeAttr := semconv.HTTPResponseStatusCode(code)
			attrs = append(attrs, codeAttr)
			span.SetAttributes(codeAttr)
		}
		attrOpt := metric.WithAttributes(attrs...)

		// Increment request counter.
		s.requests.Add(ctx, 1, attrOpt)

		// Use floating point division here for higher precision (instead of Millisecond method).
		s.duration.Record(ctx, float64(elapsedDuration)/float64(time.Millisecond), attrOpt)
	}()

	var (
		recordError = func(stage string, err error) {
			span.RecordError(err)

			// https://opentelemetry.io/docs/specs/semconv/http/http-spans/#status
			// Span Status MUST be left unset if HTTP status code was in the 1xx, 2xx or 3xx ranges,
			// unless there was another error (e.g., network error receiving the response body; or 3xx codes with
			// max redirects exceeded), in which case status MUST be set to Error.
			code := statusWriter.status
			if code >= 100 && code < 500 {
				span.SetStatus(codes.Error, stage)
			}

			attrSet := labeler.AttributeSet()
			attrs := attrSet.ToSlice()
			if code != 0 {
				attrs = append(attrs, semconv.HTTPResponseStatusCode(code))
			}

			s.errors.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: RemoveUserBookOperation,
			ID:   "removeUserBook",
		}
	)
	params, err := decodeRemoveUserBookParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response RemoveUserBookRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    RemoveUserBookOperation,
			OperationSummary: "Remove book (or complete reading)",
			OperationID:      "removeUserBook",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "user_id",
					In:   "path",
				}: params.UserID,
				{
					Name: "book_id",
					In:   "path",
				}: params.BookID,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = RemoveUserBookParams
			Response = RemoveUserBookRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackRemoveUserBookParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.RemoveUserBook(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.RemoveUserBook(ctx, params)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeRemoveUserBookResponse(response, w, span); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleUpdateReadingProgressRequest handles updateReadingProgress operation.
//
// Sets page value to a new one, returns an error if the book doesn't exist.
//
// PUT /users/{user_id}/books/{book_id}
func (s *Server) handleUpdateReadingProgressRequest(args [2]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	statusWriter := &codeRecorder{ResponseWriter: w}
	w = statusWriter
	otelAttrs := []attribute.KeyValue{
		otelogen.OperationID("updateReadingProgress"),
		semconv.HTTPRequestMethodKey.String("PUT"),
		semconv.HTTPRouteKey.String("/users/{user_id}/books/{book_id}"),
	}

	// Start a span for this request.
	ctx, span := s.cfg.Tracer.Start(r.Context(), UpdateReadingProgressOperation,
		trace.WithAttributes(otelAttrs...),
		serverSpanKind,
	)
	defer span.End()

	// Add Labeler to context.
	labeler := &Labeler{attrs: otelAttrs}
	ctx = contextWithLabeler(ctx, labeler)

	// Run stopwatch.
	startTime := time.Now()
	defer func() {
		elapsedDuration := time.Since(startTime)

		attrSet := labeler.AttributeSet()
		attrs := attrSet.ToSlice()
		code := statusWriter.status
		if code != 0 {
			codeAttr := semconv.HTTPResponseStatusCode(code)
			attrs = append(attrs, codeAttr)
			span.SetAttributes(codeAttr)
		}
		attrOpt := metric.WithAttributes(attrs...)

		// Increment request counter.
		s.requests.Add(ctx, 1, attrOpt)

		// Use floating point division here for higher precision (instead of Millisecond method).
		s.duration.Record(ctx, float64(elapsedDuration)/float64(time.Millisecond), attrOpt)
	}()

	var (
		recordError = func(stage string, err error) {
			span.RecordError(err)

			// https://opentelemetry.io/docs/specs/semconv/http/http-spans/#status
			// Span Status MUST be left unset if HTTP status code was in the 1xx, 2xx or 3xx ranges,
			// unless there was another error (e.g., network error receiving the response body; or 3xx codes with
			// max redirects exceeded), in which case status MUST be set to Error.
			code := statusWriter.status
			if code >= 100 && code < 500 {
				span.SetStatus(codes.Error, stage)
			}

			attrSet := labeler.AttributeSet()
			attrs := attrSet.ToSlice()
			if code != 0 {
				attrs = append(attrs, semconv.HTTPResponseStatusCode(code))
			}

			s.errors.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: UpdateReadingProgressOperation,
			ID:   "updateReadingProgress",
		}
	)
	params, err := decodeUpdateReadingProgressParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}
	request, close, err := s.decodeUpdateReadingProgressRequest(r)
	if err != nil {
		err = &ogenerrors.DecodeRequestError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeRequest", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}
	defer func() {
		if err := close(); err != nil {
			recordError("CloseRequest", err)
		}
	}()

	var response UpdateReadingProgressRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    UpdateReadingProgressOperation,
			OperationSummary: "Update reading progess with new current page",
			OperationID:      "updateReadingProgress",
			Body:             request,
			Params: middleware.Parameters{
				{
					Name: "user_id",
					In:   "path",
				}: params.UserID,
				{
					Name: "book_id",
					In:   "path",
				}: params.BookID,
			},
			Raw: r,
		}

		type (
			Request  = *UpdateReadingProgressReq
			Params   = UpdateReadingProgressParams
			Response = UpdateReadingProgressRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackUpdateReadingProgressParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.UpdateReadingProgress(ctx, request, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.UpdateReadingProgress(ctx, request, params)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeUpdateReadingProgressResponse(response, w, span); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}
