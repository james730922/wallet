package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/ctxs"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/tools"
)

const (
	ContentTypeProtoBuf = "application/x-protobuf"
	ContentTypeJson     = "application/json; charset=utf-8"
)

func newCtxHandler() CtxHandler {
	return CtxHandler{}
}

type CtxHandler struct {
}

func (eh CtxHandler) ResponseProtoBufStatusOK(ctx interface{}, payload proto.Message) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%s\n %s", e, debug.Stack())
			return
		}
	}()

	switch t := ctx.(type) {
	//case *ydras.Context:
	//	respBuf, err := proto.Marshal(payload)
	//	if err != nil {
	//		return err
	//	}
	//
	//	t.Response(http.StatusOK, respBuf)
	case *gin.Context:
		// API messages are generated with the legacy protobuf runtime. Gin 1.12's
		// ProtoBuf renderer expects protoreflect.ProtoMessage and panics for these
		// messages, so marshal with the same runtime used by BindProtoBuf.
		respBuf, marshalErr := proto.Marshal(payload)
		if marshalErr != nil {
			return marshalErr
		}
		t.Data(http.StatusOK, ContentTypeProtoBuf, respBuf)
	default:
		m := fmt.Sprintf("%s, type: %v", errs.FrameworkContextErrorType, t)
		panic(m)
	}

	eh.loggerResponse(ctx, payload)

	return nil
}

func (eh CtxHandler) ResponseJsonStatusOK(ctx interface{}, payload interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%s\n %s", e, debug.Stack())
			return
		}
	}()

	respBuf, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	eh.ResponseStatusOK(ctx, ContentTypeJson, respBuf)

	return nil
}

func (eh CtxHandler) ResponseStatusOK(ctx interface{}, contentType string, payload []byte) error {
	return eh.Response(ctx, http.StatusOK, contentType, payload)
}

func (eh CtxHandler) Response(ctx interface{}, status int, contentType string, payload []byte) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%s\n %s", e, debug.Stack())
			return
		}
	}()

	switch t := ctx.(type) {
	//case *ydras.Context:
	//	t.Response(uint32(status), payload)
	case *gin.Context:
		t.Data(status, contentType, payload)
	default:
		m := fmt.Sprintf("%s, type: %v", errs.FrameworkContextErrorType, t)
		panic(m)
	}

	eh.loggerResponse(ctx, payload)

	return nil
}

func (eh CtxHandler) ResponseProtoBufBadRequest(ctx interface{}, errMsg error) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%s\n %s", e, debug.Stack())
			return
		}
	}()

	parsedError := ErrorHandler.Parse(errMsg)
	switch t := ctx.(type) {
	//case *ydras.Context:
	//	t.Response(http.StatusBadRequest, ErrorHandler.Marshal(errMsg))
	case *gin.Context:
		respBuf, marshalErr := proto.Marshal(parsedError.InnerError)
		if marshalErr != nil {
			return marshalErr
		}
		// Preserve the status carried by the domain error (for example 401,
		// 403, or 503) instead of flattening every API failure to HTTP 400.
		t.Data(parsedError.StatusCode, ContentTypeProtoBuf, respBuf)
	default:
		m := fmt.Sprintf("%s, type: %v", errs.FrameworkContextErrorType, t)
		panic(m)
	}

	eh.loggerResponse(ctx, errMsg)

	return nil
}

func (eh CtxHandler) ResponseJsonBadRequest(ctx interface{}, errMsg error) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%s\n %s", e, debug.Stack())
			return
		}
	}()

	respBuf, err := tools.JsonMarshal(ErrorHandler.Parse(errMsg).InnerError)
	if err != nil {
		return err
	}

	eh.Response(ctx, http.StatusBadRequest, ContentTypeJson, respBuf)

	return nil
}

func (eh CtxHandler) ResponseJsonWithError(ctx interface{}, errMsg error) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%s\n %s", e, debug.Stack())
			return
		}
	}()

	selfDefinedError := ErrorHandler.Parse(errMsg)
	respBuf, err := tools.JsonMarshal(selfDefinedError.InnerError)
	if err != nil {
		return err
	}

	eh.Response(ctx, selfDefinedError.StatusCode, ContentTypeJson, respBuf)

	return nil
}

func (eh CtxHandler) ResponseJsonNotFound(ctx interface{}, errMsg error) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%s\n %s", e, debug.Stack())
			return
		}
	}()

	respBuf, err := tools.JsonMarshal(ErrorHandler.Parse(errMsg).InnerError)
	if err != nil {
		return err
	}

	eh.Response(ctx, http.StatusNotFound, ContentTypeJson, respBuf)

	return nil
}

func (eh CtxHandler) ResponseJsonInternalServerError(ctx interface{}, errMsg error) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%s\n %s", e, debug.Stack())
			return
		}
	}()

	respBuf, err := tools.JsonMarshal(ErrorHandler.Parse(errMsg).InnerError)
	if err != nil {
		return err
	}

	eh.Response(ctx, http.StatusInternalServerError, ContentTypeJson, respBuf)

	return nil
}

func (eh CtxHandler) BindProtoBuf(ctx interface{}, obj interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%s\n %s", e, debug.Stack())
			return
		}
	}()

	pb, ok := obj.(proto.Message)
	if !ok {
		return errs.FrameworkNotProtoBuf
	}
	defer eh.loggerRequest(ctx, obj)

	switch t := ctx.(type) {
	case *gin.Context:
		switch t.Request.Method {
		case http.MethodGet:
			return eh.bindGinMethodGet(t, obj)
		default:
			body, err := io.ReadAll(t.Request.Body)
			if err != nil {
				logger.AccessLog().Errorf("%s: %s", errs.FrameworkRequestBodyParseError, err)
				return errs.FrameworkRequestBodyParseError
			}
			// Generated API messages still implement the legacy protobuf interface.
			// Gin 1.12's binder uses google.golang.org/protobuf and rejects them,
			// so decode through the protobuf runtime they were generated for.
			if err := proto.Unmarshal(body, pb); err != nil {
				logger.AccessLog().Errorf("%s: %s", errs.FrameworkRequestBodyParseError, err)
				return errs.FrameworkRequestBodyParseError
			}

		}
	default:
		m := fmt.Sprintf("%s, type: %v", errs.FrameworkContextErrorType, t)
		panic(m)
	}

	return nil
}

func (eh CtxHandler) BindJson(ctx interface{}, obj interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%s\n %s", e, debug.Stack())
			return
		}
	}()

	defer eh.loggerRequest(ctx, obj)

	switch t := ctx.(type) {
	case *gin.Context:
		switch t.Request.Method {
		case http.MethodGet:
			return eh.bindGinMethodGet(t, obj)
		default:
			if err := binding.JSON.Bind(t.Request, obj); err != nil {
				logger.AccessLog().Errorf("%s: %s", errs.FrameworkRequestBodyParseError, err)
				return errs.FrameworkRequestBodyParseError
			}

		}
	default:
		m := fmt.Sprintf("%s, type: %v", errs.FrameworkContextErrorType, t)
		panic(m)
	}

	return nil
}

func (eh CtxHandler) bindGinMethodGet(ctx *gin.Context, obj interface{}) (err error) {
	if err := binding.Query.Bind(ctx.Request, obj); err != nil {
		c := ctxs.GetSelfContext(ctx)
		sessionID, _ := ctxs.GetSessionID(c)
		logger.AccessLog().Errorf("Req: [sessionID: %s], %s payload: %s, err: %s",
			sessionID, errs.FrameworkRequestBodyParseError, sanitizeLogPayload(obj), err)
		return errs.FrameworkRequestBodyParseError
	}

	var paging zqbapis.Paging
	valPaging := reflect.ValueOf(&paging)

	// 當發現有zqbapis.Paging 時，解析並塞入paging
	val := reflect.Indirect(reflect.ValueOf(obj))
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if field.Type().Kind() != reflect.Ptr {
			continue
		}

		switch field.Type() {
		case valPaging.Type():
			// 解析paging
			if err := binding.Query.Bind(ctx.Request, &paging); err != nil {
				logger.AccessLog().Errorf("%s: %s", errs.FrameworkRequestBodyParseError, err)
				return errs.FrameworkRequestBodyParseError
			}

			field.Set(reflect.ValueOf(&paging))
		}

		break
	}

	return nil
}

func (eh CtxHandler) GetParamInt32(ctx interface{}, key string) (int32, error) {
	val, err := eh.getParam(ctx, key)
	if err != nil {
		return 0, nil
	}

	v, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return 0, errs.CommonRequestParamParseFailed
	}

	return int32(v), nil
}

func (eh CtxHandler) GetParamInt64(ctx interface{}, key string) (int64, error) {
	val, err := eh.getParam(ctx, key)
	if err != nil {
		return 0, nil
	}

	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, errs.CommonRequestParamParseFailed
	}

	return v, nil
}

func (eh CtxHandler) getParam(ctx interface{}, key string) (str string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%s\n %s", e, debug.Stack())
			return
		}
	}()

	tmpResult := ""

	switch t := ctx.(type) {
	case *gin.Context:
		tv := t.Param(key)
		if tv == "" {
			return "", errs.CommonRequestParamInvalid
		}

		tmpResult = tv
	default:
		m := fmt.Sprintf("%s, type: %v", errs.FrameworkContextErrorType, t)
		panic(m)
	}

	return tmpResult, nil
}

func (eh CtxHandler) loggerRequest(ctx interface{}, payload interface{}) {
	c := ctxs.GetSelfContext(ctx)
	sessionID, ok := ctxs.GetSessionID(c)
	if !ok {
		logger.AccessLog().Errorf("sessionID not found: %v, %v", tools.JsonMarshalString(ctx), sanitizeLogPayload(payload))
		return
	}
	logger.AccessLog().Infof("Req: [sessionID: %s], payload: %s", sessionID, sanitizeLogPayload(payload))
}

func (eh CtxHandler) loggerResponse(ctx interface{}, payload interface{}) {
	c := ctxs.GetSelfContext(ctx)
	sessionID, ok := ctxs.GetSessionID(c)

	logPayload := ""
	switch payload.(type) {
	case []byte:
		logPayload = sanitizeLogPayload(json.RawMessage(payload.([]byte)))
	default:
		logPayload = sanitizeLogPayload(payload)
	}

	if !ok {
		logger.AccessLog().Errorf("sessionID not found: %v, %v \n%s", tools.JsonMarshalString(ctx), logPayload, string(debug.Stack()))
		return
	}
	logger.AccessLog().Infof("Resp: [sessionID: %s], payload: %s", sessionID, logPayload)
}

func sanitizeLogPayload(payload interface{}) string {
	raw, err := json.Marshal(payload)
	if err != nil {
		return `"[UNSERIALIZABLE]"`
	}

	var value interface{}
	if err := json.Unmarshal(raw, &value); err != nil {
		return `"[UNSERIALIZABLE]"`
	}

	redactPasswordFields(value)
	redacted, err := json.Marshal(value)
	if err != nil {
		return `"[UNSERIALIZABLE]"`
	}
	return string(redacted)
}

func redactPasswordFields(value interface{}) {
	switch item := value.(type) {
	case map[string]interface{}:
		for key, field := range item {
			normalizedKey := strings.ToLower(strings.ReplaceAll(key, "-", "_"))
			if strings.Contains(normalizedKey, "passwd") || strings.Contains(normalizedKey, "password") {
				item[key] = "[REDACTED]"
				continue
			}
			redactPasswordFields(field)
		}
	case []interface{}:
		for _, field := range item {
			redactPasswordFields(field)
		}
	}
}
