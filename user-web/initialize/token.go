package initialize

import (
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"time"
	"web-api/user-web/global"
	"web-api/user-web/utils/token"
)

func InitJwtVerifier() {
	pbBytes := readKey(global.ServerConfig.JwtInfo.PublicKeyPath)
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pbBytes)
	if err != nil {
		zap.L().Fatal("can not ParseRSAPublicKeyFromPEM", zap.Error(err))
	}
	global.JWTTokenVerifier = &token.JWTTokenVerifier{
		PublicKey: publicKey,
	}
	pbBytes = readKey(global.ServerConfig.JwtInfo.PrivateKeyPath)
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pbBytes)
	if err != nil {
		zap.L().Fatal("can not ParseRSAPrivateKeyFromPEM", zap.Error(err))
	}
	global.JwtTokenGen = token.NewJWTokenGen(global.ServerConfig.JwtInfo.Issuer,
		(global.ServerConfig.JwtInfo.Expire * time.Minute),
		privateKey)
}

func readKey(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		zap.L().Fatal("can not open file public key ", zap.Error(err))
	}
	defer file.Close()
	pbBytes, err := ioutil.ReadAll(file)
	if err != nil {
		zap.L().Fatal("can not read all bytes from public key", zap.Error(err))
	}
	return pbBytes
}
