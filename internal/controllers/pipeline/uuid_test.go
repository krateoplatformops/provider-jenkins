//go:build integration
// +build integration

package pipeline

import "testing"

func TestUUID(t *testing.T) {
	str := `<aaa class="xxx">__UUID__</a>`

	exp, err := fillWithUUID(str, "__UUID__")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("\n%s\n", exp)
}
