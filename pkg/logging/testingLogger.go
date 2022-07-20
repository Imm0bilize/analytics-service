package logging

import "github.com/sirupsen/logrus"

type loggerStub struct {
	*logrus.Logger
}

func NewLoggerStub() *loggerStub {
	return &loggerStub{logrus.New()}
}

func (l loggerStub) Debugf(_ string, _ ...interface{}) {}

func (l loggerStub) Infof(_ string, _ ...interface{}) {}

func (l loggerStub) Printf(_ string, _ ...interface{}) {}

func (l loggerStub) Warnf(_ string, _ ...interface{}) {}

func (l loggerStub) Warningf(_ string, _ ...interface{}) {}

func (l loggerStub) Errorf(_ string, _ ...interface{}) {}

func (l loggerStub) Fatalf(_ string, _ ...interface{}) {}

func (l loggerStub) Panicf(_ string, _ ...interface{}) {}

func (l loggerStub) Debug(_ ...interface{}) {}

func (l loggerStub) Info(_ ...interface{}) {}

func (l loggerStub) Print(_ ...interface{}) {}

func (l loggerStub) Warn(_ ...interface{}) {}

func (l loggerStub) Warning(_ ...interface{}) {}

func (l loggerStub) Error(_ ...interface{}) {}

func (l loggerStub) Fatal(_ ...interface{}) {}

func (l loggerStub) Panic(_ ...interface{}) {}

func (l loggerStub) Debugln(_ ...interface{}) {}

func (l loggerStub) Infoln(_ ...interface{}) {}

func (l loggerStub) Println(_ ...interface{}) {}

func (l loggerStub) Warnln(_ ...interface{}) {}

func (l loggerStub) Warningln(_ ...interface{}) {}

func (l loggerStub) Errorln(_ ...interface{}) {}

func (l loggerStub) Fatalln(_ ...interface{}) {}

func (l loggerStub) Panicln(_ ...interface{}) {}
