package http

import (
	"bytes"
	"strings"
	"testing"
)

func TestSortCookiePairs(t *testing.T) {
	values := []string{"c=3; a=1; d=4; b=2"}
	order := []string{"a", "b", "c"}
	got := SortCookiePairs(values, order)
	want := []string{"a=1", "b=2", "c=3", "d=4"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("want %v, got %v", want, got)
	}
}

func TestSortCookiePairsNoOrder(t *testing.T) {
	values := []string{"c=3; a=1; b=2"}
	got := SortCookiePairs(values, nil)
	want := []string{"c=3", "a=1", "b=2"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("want %v, got %v", want, got)
	}
}

func TestCookieOrderWrite(t *testing.T) {
	h := Header{
		"Cookie":        {"c=3; a=1; d=4; b=2"},
		CookieOrderKey:  {"a", "b", "c"},
		HeaderOrderKey:  {"cookie"},
		PHeaderOrderKey: {":method"},
	}
	var buf bytes.Buffer
	if err := h.Write(&buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()

	if want := "Cookie: a=1; b=2; c=3; d=4\r\n"; !strings.Contains(out, want) {
		t.Fatalf("want to contain %q, got:\n%s", want, out)
	}
	// Magic keys must not be written to the wire.
	if strings.Contains(out, CookieOrderKey) {
		t.Fatalf("CookieOrderKey leaked into output:\n%s", out)
	}
}
