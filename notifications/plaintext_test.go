package notifications_test

import (
	"github.com/bubulearn/bubucore/notifications"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlainText_Filter(t *testing.T) {
	values := map[string]string{
		"   ":     "",
		" test ":  "test",
		"Test  1": "Test  1",
	}

	for in, expected := range values {
		n := &notifications.PlainText{
			Text: in,
		}
		n.Filter()
		assert.Equal(t, expected, n.Text)
	}
}

func TestPlainText_Validate(t *testing.T) {
	valid := []string{
		"Test 1",
		"<pre>test 2</pre>",
	}
	invalid := []string{
		"",
		"    ",
		"\n\n",
	}

	for _, v := range valid {
		n := &notifications.PlainText{
			Text: v,
		}
		assert.NoError(t, n.Validate())
	}

	for _, v := range invalid {
		n := &notifications.PlainText{
			Text: v,
		}
		assert.Error(t, n.Validate())
	}
}
