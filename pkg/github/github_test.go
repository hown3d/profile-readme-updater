package github

import (
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
)

func newGithubMockClient(options ...mock.MockBackendOption) *github.Client {
	mockedHTTPClient := mock.NewMockedHTTPClient(options...)
	return github.NewClient(mockedHTTPClient)
}

func Test_splitRepoName(t *testing.T) {
	type args struct {
		repoName string
	}
	tests := []struct {
		name      string
		args      args
		wantOwner string
		wantName  string
		wantErr   bool
	}{
		{
			name: "valid repo Name",
			args: args{
				repoName: "test/foo",
			},
			wantOwner: "test",
			wantName:  "foo",
		},
		{
			name: "2 slashes repo Name",
			args: args{
				repoName: "test/foo/",
			},
			wantErr: true,
		},

		{
			name: "no slash repo Name",
			args: args{
				repoName: "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOwner, gotName, err := splitRepoName(tt.args.repoName)
			if tt.wantErr {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.wantOwner, gotOwner)
			assert.Equal(t, tt.wantName, gotName)
		})
	}
}
