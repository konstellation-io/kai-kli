package mocks

import (
	"github.com/golang/mock/gomock"
)

func AddLoggerExpects(logger *MockLogger) {
	logger.EXPECT().Debug(gomock.Any()).Return().AnyTimes()
	logger.EXPECT().Info(gomock.Any()).Return().AnyTimes()
	logger.EXPECT().Warn(gomock.Any()).Return().AnyTimes()
	logger.EXPECT().Error(gomock.Any()).Return().AnyTimes()
	logger.EXPECT().Success(gomock.Any()).Return().AnyTimes()
}
