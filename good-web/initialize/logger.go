package initialize

import "go.uber.org/zap"

func InitLogger() {
	p, err := zap.NewProduction()
	if err != nil {
		return
	}
	zap.ReplaceGlobals(p)
}
