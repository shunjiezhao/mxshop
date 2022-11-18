package utils

import (
	"crypto/md5"
	"encoding/hex"
	"go.uber.org/zap"
)

func Hash(b []byte) string {
	hash := md5.New()
	_, err := hash.Write(b)
	if err != nil {
		zap.L().Error("can not hash", zap.Error(err))
		return ""
	}
	return hex.EncodeToString(hash.Sum(nil))
}
