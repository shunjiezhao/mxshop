package initialize

import (
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"web-api/shop_cart-web/global"
	"web-api/shop_cart-web/utils/token"
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
