package msort

import "testing"

func TestSortedByKeys(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		out     string
		wantErr bool
	}{
		{
			name: "same",
			in:   "abc: 1\ndef: 2\n",
			out:  "abc: 1\ndef: 2\n",
		},
		{
			name: "missing linebreak",
			in:   "abc: 1\ndef: 2",
			out:  "abc: 1\ndef: 2\n",
		},
		{
			name: "starting doc",
			in:   "---\nabc: 1\ndef: 2\n",
			out:  "abc: 1\ndef: 2\n",
		},
		{
			name: "sort map",
			in:   "def: 2\nabc: 1\n",
			out:  "abc: 1\ndef: 2\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SortedByKeys(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("SortedByKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.out {
				t.Errorf("SortedByKeys() = %v, want %v", got, tt.out)
			}
		})
	}
}
