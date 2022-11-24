package token

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JWTTokenVerifier struct {
	PublicKey *rsa.PublicKey
}
type CustomClaim struct {
	Nickname string `json:"nickname"`
	GoodsId  uint32 `json:"Goods_id"`
	Role     uint32 `json:"role,omitempty"`
	jwt.StandardClaims
}

// 验证 token 返回 account ID
func (v *JWTTokenVerifier) Verify(token string) (*CustomClaim, error) {
	// key func 就是 返回解密用的 key
	t, err := jwt.ParseWithClaims(token, &CustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		return v.PublicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("can not parse token %v ", err)
	}

	// 验证 方法是否正确
	if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
	}
	// 验证 是否签名合法
	if !t.Valid {
		return nil, fmt.Errorf("token not valid")
	}

	clm, ok := t.Claims.(*CustomClaim)
	if !ok {
		return nil, fmt.Errorf("claim is not StandardClaims")
	}
	// 验证 data 的内容 如过期时间exp
	if err = clm.Valid(); err != nil {
		return nil, fmt.Errorf("claim is valid")
	}
	// 返回 account ID
	return clm, nil
}

// Generate TOken

type JWTokenGen struct {
	privateKey *rsa.PrivateKey
	ExpiresAt  time.Duration
	Issue      string
	nowFunc    func() time.Time
}

func NewJWTokenGen(issue string, ExpiresAt time.Duration, privateKey *rsa.PrivateKey) *JWTokenGen {
	return &JWTokenGen{
		Issue:      issue,
		nowFunc:    time.Now,
		ExpiresAt:  (ExpiresAt),
		privateKey: privateKey,
	}
}
