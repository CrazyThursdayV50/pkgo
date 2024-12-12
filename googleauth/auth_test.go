package googleauth

import "testing"

func TestAuth(t *testing.T) {
	key, err := NewKey()
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	t.Logf("key: %s", key)

	url, err := NewKeyUrl("Alex", "Amazon")
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	t.Logf("url: %s", url)
	key = "E4I7XMPPXAUEXK4KJSBGAR74IAPR5ARH"
	var num int64 = 486324

	t.Logf("ok: %v", Validate(key, num))
}
