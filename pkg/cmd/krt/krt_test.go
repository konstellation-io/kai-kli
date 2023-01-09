package krt_test

import (
	"testing"

	"github.com/konstellation-io/kli/internal/logger"
	"github.com/konstellation-io/kli/pkg/cmd/krt"

	"github.com/MakeNowJust/heredoc"
	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"

	"github.com/konstellation-io/kli/internal/testhelpers"
	"github.com/konstellation-io/kli/mocks"
)

func TestNewKRTCreateCmd(t *testing.T) {
	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		c := gomock.NewController(t)
		l := logger.NewDefaultLogger()
		krtTooler := mocks.NewMockKrtTooler(c)
		f.EXPECT().Krt().Return(krtTooler).AnyTimes()
		f.EXPECT().Logger().Return(l).AnyTimes()

		krtTooler.EXPECT().Build("/test/krt", "test.krt", "testversion").Return(nil)

		return krt.NewKRTCmd(f)
	})

	r.Run("krt create /test/krt test.krt -v testversion").
		Contains(heredoc.Doc(`
	  New KRT file created.
  `))
}

func TestNewKRTCreateCmdWithoutVersion(t *testing.T) {
	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		c := gomock.NewController(t)
		krtTooler := mocks.NewMockKrtTooler(c)
		f.EXPECT().Krt().Return(krtTooler)

		krtTooler.EXPECT().Build("/test/krt", "test.krt", "").Return(nil)

		return krt.NewKRTCmd(f)
	})

	r.Run("krt create /test/krt test.krt").
		Contains(heredoc.Doc(`
	  New KRT file created.
  `))
}

func TestNewValidateCmd(t *testing.T) {
	r := testhelpers.NewRunner(t, func(f *mocks.MockCmdFactory) *cobra.Command {
		c := gomock.NewController(t)
		krtTooler := mocks.NewMockKrtTooler(c)

		f.EXPECT().Krt().Return(krtTooler)
		krtTooler.EXPECT().Validate("/test/krt.yaml").Return(nil)

		return krt.NewKRTCmd(f)
	})
	r.Run("krt validate /test/krt.yaml").
		Contains(heredoc.Doc(`
	  Krt file is valid.
  `))
}
