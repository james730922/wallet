package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"

	"github.com/james730922/wallet/service/internal/pb/zqbapis"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/ctxs"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

var handlerTestLoggerOnce sync.Once

func initHandlerTestLogger() {
	handlerTestLoggerOnce.Do(logger.TestMock)
}

func TestBindProtoBufSupportsLegacyGeneratedMessages(t *testing.T) {
	initHandlerTestLogger()
	gin.SetMode(gin.TestMode)

	want := &zqbapis.LoginIdentifierPasswdReq{
		CountryCode: "886",
		Mobile:      "0912345678",
		Passwd:      "secret",
	}
	body, err := proto.Marshal(want)
	if err != nil {
		t.Fatalf("proto.Marshal() error = %v", err)
	}

	recorder := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(recorder)
	ginCtx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(body))
	self := ctxs.SetSessionID(context.Background(), "test-session")
	ctxs.SetSelfContext(self, ginCtx)

	got := &zqbapis.LoginIdentifierPasswdReq{}
	if err := (CtxHandler{}).BindProtoBuf(ginCtx, got); err != nil {
		t.Fatalf("BindProtoBuf() error = %v", err)
	}
	if !proto.Equal(got, want) {
		t.Fatalf("BindProtoBuf() = %v, want %v", got, want)
	}
}

func TestResponseProtoBufSupportsLegacyGeneratedMessages(t *testing.T) {
	initHandlerTestLogger()
	gin.SetMode(gin.TestMode)

	want := &zqbapis.LoginIdentifierPasswdResp{
		Token:       "test-token",
		CountryCode: "86",
		Mobile:      "13800000000",
	}
	recorder := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(recorder)
	self := ctxs.SetSessionID(context.Background(), "test-session")
	ctxs.SetSelfContext(self, ginCtx)

	if err := (CtxHandler{}).ResponseProtoBufStatusOK(ginCtx, want); err != nil {
		t.Fatalf("ResponseProtoBufStatusOK() error = %v", err)
	}
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Header().Get("Content-Type"); got != ContentTypeProtoBuf {
		t.Fatalf("Content-Type = %q, want %q", got, ContentTypeProtoBuf)
	}

	got := &zqbapis.LoginIdentifierPasswdResp{}
	if err := proto.Unmarshal(recorder.Body.Bytes(), got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !proto.Equal(got, want) {
		t.Fatalf("response = %v, want %v", got, want)
	}
}

func TestResponseProtoBufBadRequestSupportsLegacyGeneratedMessages(t *testing.T) {
	initHandlerTestLogger()
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(recorder)
	self := ctxs.SetSessionID(context.Background(), "test-session")
	ctxs.SetSelfContext(self, ginCtx)

	if err := (CtxHandler{}).ResponseProtoBufBadRequest(ginCtx, errs.MobileEmpty); err != nil {
		t.Fatalf("ResponseProtoBufBadRequest() error = %v", err)
	}
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}

	got := &zqbapis.Error{}
	if err := proto.Unmarshal(recorder.Body.Bytes(), got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Code == "" || got.Message == "" {
		t.Fatalf("error response is empty: %v", got)
	}
}

func TestResponseProtoBufBadRequestPreservesDefinedHTTPStatus(t *testing.T) {
	initHandlerTestLogger()
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(recorder)
	self := ctxs.SetSessionID(context.Background(), "test-session")
	ctxs.SetSelfContext(self, ginCtx)

	if err := (CtxHandler{}).ResponseProtoBufBadRequest(ginCtx, errs.AuthOperationForbidden); err != nil {
		t.Fatalf("ResponseProtoBufBadRequest() error = %v", err)
	}
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
	}
}

func TestSanitizeLogPayloadRedactsPasswords(t *testing.T) {
	payload := map[string]interface{}{
		"passwd":         "login-secret",
		"new_passwd":     "new-secret",
		"login_password": "reset-secret",
		"token":          "token-is-outside-this-change",
		"nested": map[string]interface{}{
			"securityPasswd": "payment-secret",
		},
	}

	got := sanitizeLogPayload(payload)
	for _, secret := range []string{"login-secret", "new-secret", "reset-secret", "payment-secret"} {
		if strings.Contains(got, secret) {
			t.Fatalf("sanitizeLogPayload() leaked %q in %s", secret, got)
		}
	}
	if strings.Count(got, "[REDACTED]") != 4 {
		t.Fatalf("sanitizeLogPayload() = %s, want four redacted fields", got)
	}
	if !strings.Contains(got, "token-is-outside-this-change") {
		t.Fatalf("sanitizeLogPayload() unexpectedly changed token: %s", got)
	}
}
