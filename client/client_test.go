package client

import (
	"testing"

	"github.com/go-ocf/kit/net/grpc"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	type args struct {
		host string
		opts []grpc.ClientConnOption
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid",
			want: true,
			args: args{
				opts: []grpc.ClientConnOption{grpc.WithInsecure()},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.host, tt.args.opts...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tt.want {
				assert.NotNil(t, got)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}
