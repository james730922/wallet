package auth

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"

	"github.com/james730922/wallet/service/internal/models/model"
	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/conf"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

type tokenTypeEnum int

const (
	tokenTypeEnum_Member tokenTypeEnum = iota
	tokenTypeEnum_MemberSecurityPasswd
	tokenTypeEnum_MemberForgetSecurityPasswd
)

const rotateTokenScript = `
local old_token = redis.call("GET", KEYS[2])
if old_token then
  redis.call("DEL", ARGV[1] .. old_token)
end
redis.call("SET", KEYS[1], ARGV[2], "EX", ARGV[4])
redis.call("SET", KEYS[2], ARGV[3], "EX", ARGV[4])
return 1`

const deleteTokenScript = `
local id = redis.call("GET", KEYS[1])
if not id then
  return 0
end
local id_key = ARGV[1] .. id
if redis.call("GET", id_key) == ARGV[2] then
  redis.call("DEL", id_key)
end
redis.call("DEL", KEYS[1])
return 1`

const deleteMemberTokenScript = `
local token = redis.call("GET", KEYS[1])
if token then
  redis.call("DEL", ARGV[1] .. token)
end
redis.call("DEL", KEYS[1])
return 1`

func (t tokenTypeEnum) Names() map[tokenTypeEnum]string {
	return map[tokenTypeEnum]string{
		tokenTypeEnum_Member:                     "member",                     // 會員
		tokenTypeEnum_MemberSecurityPasswd:       "memberSecurityPasswd",       // 會員安全密碼
		tokenTypeEnum_MemberForgetSecurityPasswd: "memberForgetSecurityPasswd", // 會員安全密碼忘記密碼修改
	}
}

func (t tokenTypeEnum) String() string {
	return t.Names()[t]
}

var token IToken

func newToken() IToken {
	token = &tokenUseCase{}
	return token
}

type IToken interface {
	GenMember(ctx context.Context, id int64) (string, error)
	DeleteMemberByToken(ctx context.Context, token string) error
	DeleteMemberByID(ctx context.Context, id int64) error
	AuthMember(ctx context.Context, token string) (*model.MemberRepository, error)
	GenMemberSecurityPasswd(ctx context.Context, memberID int64) (string, error)
	DeleteMemberSecurityPasswdByToken(ctx context.Context, token string) error
	DeleteMemberSecurityPasswdByID(ctx context.Context, id int64) error
	AuthMemberSecurityPasswd(ctx context.Context, token string) (*model.MemberRepository, error)

	GenMemberForgetSecurityPasswd(ctx context.Context, memberID int64) (string, error)
	DeleteMemberForgetSecurityPasswdByToken(ctx context.Context, token string) error
	DeleteMemberForgetSecurityPasswdByID(ctx context.Context, id int64) error
	AuthMemberForgetSecurityPasswd(ctx context.Context, token string) (*model.MemberRepository, error)
}

type tokenUseCase struct{}

func (tk *tokenUseCase) GenMember(ctx context.Context, id int64) (string, error) {
	return tk.gen(ctx, tokenTypeEnum_Member, id)
}

func (tk *tokenUseCase) GenMemberSecurityPasswd(ctx context.Context, id int64) (string, error) {
	return tk.gen(ctx, tokenTypeEnum_MemberSecurityPasswd, id)
}

func (tk *tokenUseCase) GenMemberForgetSecurityPasswd(ctx context.Context, id int64) (string, error) {
	return tk.gen(ctx, tokenTypeEnum_MemberForgetSecurityPasswd, id)
}

func (tk *tokenUseCase) gen(ctx context.Context, tokenType tokenTypeEnum, id int64) (string, error) {
	token := tk.genToken()
	idStr := strconv.FormatInt(id, 10)

	prefix := tokenType.String() + "_"
	keyByToken := prefix + token
	keyByID := prefix + idStr
	ttl, err := tk.getTTl(tokenType)
	if err != nil {
		logger.ApLog().Errorf("id: %d, err: %s", id, err)
		return "", errs.AuthTokenCreateFailed
	}

	if ttl <= 0 {
		return "", errs.AuthTokenCreateFailed
	}
	if err := packet.Redis.Eval(
		rotateTokenScript,
		[]string{keyByToken, keyByID},
		prefix,
		idStr,
		token,
		int64(ttl/time.Second),
	).Err(); err != nil {
		logger.ApLog().Errorf("id: %d, err: %s", id, err)
		return "", errs.AuthTokenCreateFailed
	}

	return token, nil
}

// gen token
func (tk *tokenUseCase) genToken() string {
	return uuid.New().String()
}

func (tk *tokenUseCase) DeleteMemberByToken(ctx context.Context, token string) error {
	return tk.deleteByToken(ctx, tokenTypeEnum_Member, token)
}

func (tk *tokenUseCase) DeleteMemberSecurityPasswdByToken(ctx context.Context, token string) error {
	return tk.deleteByToken(ctx, tokenTypeEnum_MemberSecurityPasswd, token)
}

func (tk *tokenUseCase) DeleteMemberForgetSecurityPasswdByToken(ctx context.Context, token string) error {
	return tk.deleteByToken(ctx, tokenTypeEnum_MemberForgetSecurityPasswd, token)
}

func (tk *tokenUseCase) deleteByToken(ctx context.Context, tokenType tokenTypeEnum, token string) error {
	prefix := tokenType.String() + "_"
	if err := packet.Redis.Eval(
		deleteTokenScript,
		[]string{prefix + token},
		prefix,
		token,
	).Err(); err != nil {
		logger.ApLog().Warnf("delete token mapping failed: tokenType=%s err=%v", tokenType.String(), err)
		return errs.AuthTokenDelFailed
	}

	return nil
}

func (tk *tokenUseCase) DeleteMemberByID(ctx context.Context, id int64) error {
	return tk.deleteByID(ctx, tokenTypeEnum_Member, id)
}

func (tk *tokenUseCase) DeleteMemberSecurityPasswdByID(ctx context.Context, id int64) error {
	return tk.deleteByID(ctx, tokenTypeEnum_MemberSecurityPasswd, id)
}

func (tk *tokenUseCase) DeleteMemberForgetSecurityPasswdByID(ctx context.Context, id int64) error {
	return tk.deleteByID(ctx, tokenTypeEnum_MemberForgetSecurityPasswd, id)
}

func (tk *tokenUseCase) deleteByID(ctx context.Context, tokenType tokenTypeEnum, id int64) error {
	prefix := tokenType.String() + "_"
	keyByID := prefix + strconv.FormatInt(id, 10)
	if err := packet.Redis.Eval(
		deleteMemberTokenScript,
		[]string{keyByID},
		prefix,
	).Err(); err != nil {
		logger.ApLog().Warnf("delete member token failed: tokenType=%s memberID=%d err=%v", tokenType.String(), id, err)
		return errs.AuthTokenDelFailed
	}

	return nil
}

func (tk *tokenUseCase) AuthMember(ctx context.Context, token string) (*model.MemberRepository, error) {
	id, err := tk.auth(ctx, tokenTypeEnum_Member, token)
	if err != nil {
		logger.ApLog().Debugf("auth member token failed: err=%s", err)
		return nil, errs.AuthTokenUnauthorized
	}

	memberRepository := &model.MemberRepository{
		Token: token,
		ID:    id,
	}

	return memberRepository, nil
}

func (tk *tokenUseCase) AuthMemberSecurityPasswd(ctx context.Context, token string) (*model.MemberRepository, error) {
	id, err := tk.auth(ctx, tokenTypeEnum_MemberSecurityPasswd, token)
	if err != nil {
		logger.ApLog().Warnf("auth memberSecurityPasswd token failed: err=%s", err)
		return nil, errs.AuthTokenUnauthorized
	}

	memberRepository := &model.MemberRepository{
		Token: token,
		ID:    id,
	}

	return memberRepository, nil
}

func (tk *tokenUseCase) AuthMemberForgetSecurityPasswd(ctx context.Context, token string) (*model.MemberRepository, error) {
	id, err := tk.auth(ctx, tokenTypeEnum_MemberForgetSecurityPasswd, token)
	if err != nil {
		logger.ApLog().Warnf("auth memberForgetSecurityPasswd token failed: err=%s", err)
		return nil, errs.AuthTokenUnauthorized
	}

	memberRepository := &model.MemberRepository{
		Token: token,
		ID:    id,
	}

	return memberRepository, nil
}

func (tk *tokenUseCase) auth(ctx context.Context, tokenType tokenTypeEnum, token string) (int64, error) {
	if token == "" {
		return 0, errs.AuthTokenUnauthorized
	}

	prefix := tokenType.String() + "_"
	keyByToken := prefix + token

	idStr, err := packet.Redis.Get(keyByToken).Result()
	if err != nil {
		if err != redis.Nil {
			logger.ApLog().Warnf("read token failed: tokenType=%s err=%v", tokenType.String(), err)
		}
		return 0, errs.AuthTokenUnauthorized
	}
	keyByID := prefix + idStr
	currentToken, err := packet.Redis.Get(keyByID).Result()
	if err != nil || currentToken != token {
		return 0, errs.AuthTokenUnauthorized
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.ApLog().Errorf("invalid member id in token mapping: err=%v", err)
		return 0, errs.AuthTokenUnauthorized
	}

	return id, nil
}

func (tk *tokenUseCase) getTTl(tokenType tokenTypeEnum) (time.Duration, error) {
	ttl := time.Duration(0)
	switch tokenType {
	case tokenTypeEnum_Member:
		ttl = conf.LoginMember().GetMemberTokenExpiration()
	case tokenTypeEnum_MemberSecurityPasswd:
		ttl = conf.LoginMember().GetMemberSecureTokenExpiration()
	case tokenTypeEnum_MemberForgetSecurityPasswd:
		ttl = conf.LoginMember().GetMemberForgetSecureTokenExpiration()
	default:
		return 0, errs.AuthTokenCreateFailed
	}

	return ttl, nil
}
