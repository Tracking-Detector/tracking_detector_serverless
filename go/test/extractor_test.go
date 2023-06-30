package test

import (
	"reflect"
	"tds/shared/extractor"
	"testing"
)

func TestUrlExtractor(t *testing.T) {
	t.Run("should_extract_url_exactly_as_in_node_when_smaller", func(t *testing.T) {
		// given
		url := "https://bat.bing.com/actionp/0?ti=17550024&Ver=2&mid" +
			"=08bdac16-6066-4298-be0e-2fc4477e74e6&sid=f99ac5901a2711ed" +
			"8deb51e4b6913687&vid=f99b7ed01a2711edab50617d53bfba7b&v" +
			"ids=1&msclkid=N&evt=pageHide"
		// when
		got := extractor.URL_EXTRACTOR(url)

		// then
		want := []int{0, 0, 0, 0, 0, 0, 0, 16, 28, 28, 24, 27,
			59, 48, 48, 10, 9, 28, 47, 10, 17, 22, 15, 47, 11,
			23, 21, 48, 9, 11, 28, 17, 23, 22, 24, 48, 49, 64,
			28, 17, 62, 50, 56, 54, 54, 49, 49, 51, 53, 39, 87,
			13, 26, 62, 51, 39, 21, 17, 12, 62, 49, 57, 10, 12,
			9, 11, 50, 55, 46, 55, 49, 55, 55, 46, 53, 51, 58,
			57, 46, 10, 13, 49, 13, 46, 51, 14, 11, 53, 53, 56,
			56, 13, 56, 53, 13, 55, 39, 27, 17, 12, 62, 14, 58,
			58, 9, 11, 54, 58, 49, 50, 9, 51, 56, 50, 50, 13,
			12, 57, 12, 13, 10, 54, 50, 13, 53, 10, 55, 58, 50,
			52, 55, 57, 56, 39, 30, 17, 12, 62, 14, 58, 58, 10,
			56, 13, 12, 49, 50, 9, 51, 56, 50, 50, 13, 12, 9, 10,
			54, 49, 55, 50, 56, 12, 54, 52, 10, 14, 10, 9, 56, 10,
			39, 30, 17, 12, 27, 62, 50, 39, 21, 27, 11, 20, 19, 17,
			12, 62, 79, 39, 13, 30, 28, 62, 24, 9, 15, 13, 73, 17, 12, 13}

		if len(got) != len(want) {
			t.Errorf("got len %d wanted len %d", len(got), len(want))
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v wanted %v", got, want)
		}
	})
	t.Run("should_extract_url_exactly_as_in_node_when_bigger", func(t *testing.T) {
		// given
		url := "https://region1.google-analytics.com/g/collect?" +
			"v=2&tid=G-0P6ZGEWQVE&gtm=2oe880&_p=1729724280&cid=939202449" +
			".1660299502&ul=en-us&sr=1920x1200&_z=ccd.v9B&_s=1&sid" +
			"=1660299502&sct=1&seg=0&dl=https%3A%2F%2Fwww.scientificamerican." +
			"com%2F&dt=Scientific%20American%3A%20Science%20News%2C%20Expert%20" +
			"Analysis%2C%20Health%20Research%20-%20Scientific%20American&en=" +
			"page_view&_fv=1&_nsi=1&_ss=1&_ee=1"
		// when
		got := extractor.URL_EXTRACTOR(url)
		// then
		want := []int{16, 28, 28, 24, 27, 38, 52, 66, 38, 51, 71, 38, 51, 71, 31, 31,
			31, 47, 27, 11, 17, 13, 22, 28, 17, 14, 17, 11, 9, 21, 13, 26, 17, 11, 9,
			22, 47, 11, 23, 21, 38, 51, 71, 39, 12, 28, 62, 84, 11, 17, 13, 22, 28, 17,
			14, 17, 11, 38, 51, 49, 66, 21, 13, 26, 17, 11, 9, 22, 38, 52, 66, 38, 51,
			49, 84, 11, 17, 13, 22, 11, 13, 38, 51, 49, 79, 13, 31, 27, 38, 51, 68, 38,
			51, 49, 70, 32, 24, 13, 26, 28, 38, 51, 49, 66, 22, 9, 20, 33, 27, 17, 27,
			38, 51, 68, 38, 51, 49, 73, 13, 9, 20, 28, 16, 38, 51, 49, 83, 13, 27, 13,
			9, 26, 11, 16, 38, 51, 49, 46, 38, 51, 49, 84, 11, 17, 13, 22, 28, 17, 14,
			17, 11, 38, 51, 49, 66, 21, 13, 26, 17, 11, 9, 22, 39, 13, 22, 62, 24, 9,
			15, 13, 7, 30, 17, 13, 31, 39, 7, 14, 30, 62, 50, 39, 7, 22, 27, 17, 62,
			50, 39, 7, 27, 27, 62, 50, 39, 7, 13, 13, 62, 50}
		if len(got) != len(want) {
			t.Errorf("got len %d wanted len %d", len(got), len(want))
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v wanted %v", got, want)
		}
	})

}
