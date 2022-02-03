package github

import (
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/stretchr/testify/assert"
)

func TestClient_collect(t *testing.T) {
	type fields struct {
		infos *Infos
	}
	type args struct {
		key int64
		val IssueWithRepository
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[int64]IssueWithRepository
		collectFunc
	}{
		{
			name: "empty map",
			fields: fields{
				infos: &Infos{Issues: map[int64]IssueWithRepository{}},
			},
			args: args{
				key: 123,
				val: IssueWithRepository{},
			},
			want: map[int64]IssueWithRepository{
				123: {},
			},
		},
		{
			name: "key already exists",
			fields: fields{
				infos: &Infos{Issues: map[int64]IssueWithRepository{
					123: {Repo: &github.Repository{}},
				}},
			},
			args: args{
				key: 123,
				val: IssueWithRepository{},
			},
			want: map[int64]IssueWithRepository{
				123: {Repo: &github.Repository{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				infos: tt.fields.infos,
			}
			c.collectIssue(tt.args.key, tt.args.val)
			assert.Equal(t, tt.want, c.infos.Issues)
		})
	}
}
