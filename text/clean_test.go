package text_test

import (
	"testing"

	"github.com/konstellation-io/kli/text"

	"github.com/stretchr/testify/require"
)

func TestText(t *testing.T) {
	t.Run("Sanitize", func(t *testing.T) {
		str := " some     long     string    "
		expected := "some long string"
		require.Equal(t, text.Sanitize(str), expected)
	})

	t.Run("Normalize", func(t *testing.T) {
		str := " SoME     lONg     STring    "
		expected := "some long string"
		require.Equal(t, text.Normalize(str), expected)
	})

	t.Run("LinesTrim", func(t *testing.T) {
		str := "String \t More   \n    Test   \n   New"
		expected := "String More\nTest\nNew"
		require.Equal(t, text.LinesTrim(str), expected)
	})
}
