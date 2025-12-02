package parse

import (
	"reflect"
	"testing"
)

func Test_findGoraddElements(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		want    []GoraddElement
		wantErr bool
	}{
		{
			name: "basic single element",
			html: `<!doctype html>
<html>
<body>
  <div id="outer" data-goradd="true">Hello</div>
</body>
</html>`,
			want:    []GoraddElement{{"div", "outer", 32, 78}},
			wantErr: false,
		},
		{
			name: "nested single element",
			html: `
<div id="outer" data-goradd="1">
  <span id="inner">Text</span>
</div>`,
			want:    []GoraddElement{{"div", "outer", 1, 71}},
			wantErr: false,
		},
		{
			name: "self closing",
			html: `<input id="field1" data-goradd="yes" type="text">
<br id="br1" data-goradd="line-break" />
<img id="pic1" data-goradd="pic" src="x.png">`,
			want: []GoraddElement{
				{"input", "field1", 0, 49},
				{"br", "br1", 50, 90},
				{"img", "pic1", 91, 136},
			},
			wantErr: false,
		},
		{
			name: "tricky",
			html: `
<div id="outer" data-goradd="1" title="<span data-goradd='no'> &lt;fake&gt; tag">
  Content here
</div>`,
			want: []GoraddElement{
				{"div", "outer", 1, 104},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findGoraddElements([]byte(tt.html))
			if (err != nil) != tt.wantErr {
				t.Errorf("findGoraddElements() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findGoraddElements() got = %v, want %v", got, tt.want)
			}
		})
	}
}
