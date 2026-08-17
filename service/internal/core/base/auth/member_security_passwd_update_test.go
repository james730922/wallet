package auth

import (
	"errors"
	"testing"

	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func TestSecurityPasswdFirstToAddRejectsExistingPassword(t *testing.T) {
	handler := newSecurityPasswdUpdate()

	if err := handler.validateFirstToAdd(&model.MemberMapping{}); err != nil {
		t.Fatalf("empty security password rejected: %v", err)
	}

	err := handler.validateFirstToAdd(&model.MemberMapping{SecurityPasswd: "already-set"})
	if !errors.Is(err, errs.AuthOperationForbidden) {
		t.Fatalf("existing security password error = %v, want %v", err, errs.AuthOperationForbidden)
	}
}
