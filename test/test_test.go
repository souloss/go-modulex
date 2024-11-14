package test

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/souloss/go-clean-arch/pkg/logger"
	"go.uber.org/zap"
)

func TestLogrus(t *testing.T) {
	mlogger := logger.NewLogrusAdapter(logrus.New())
	nlogger := mlogger.Named("name1")
	n2logger := nlogger.Named("name2")
	mlogger.Info("hhhh")
	nlogger.Info("hhhh")
	n2logger.Info("hhhh")
}

func TestZap(t *testing.T) {
	logger.InitZapGlobalLogger(context.Background(), &logger.LoggerConfig{
		Level:   "info",
		Console: true,
	})
	mlogger := logger.NewZapAdapter(zap.L())
	nlogger := mlogger.Named("name1")
	n2logger := nlogger.Named("name2")
	mlogger.Info("no name")
	nlogger.Info("name1")
	n2logger.Info("name2")

}

func TestTz(t *testing.T) {
	fmt.Println(time.Now().Local())
}

func TestXXX(t *testing.T) {
	var val uint64 = math.MaxUint64
	fmt.Println(val)
	zf := zap.Uint64("haha", math.MaxUint64)
	fmt.Println(zf)
	fmt.Println(uint64(zf.Integer))
}
