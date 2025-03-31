package lending

import (
	"context"
	"testing"
	"time"

	"github.com/ngoctd314/c72-api-server/pkg/dto"
	"github.com/ngoctd314/c72-api-server/pkg/helper"
	"github.com/stretchr/testify/assert"
)

func Test_list_Validate(t *testing.T) {
	type want struct {
		from int64
		to   int64
	}
	type args struct {
		ctx context.Context
		req *dto.ListLendingRequest
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "from, to both equal to nil",
			args: args{
				req: &dto.ListLendingRequest{
					From: nil,
					To:   nil,
				},
			},
			want: want{
				from: time.Now().Add(-weakDuration).Unix(),
				to:   time.Now().Unix(),
			},
			wantErr: false,
		},
		{
			name: "to equal to nil",
			args: args{
				req: &dto.ListLendingRequest{
					From: helper.Ptr[int64](1000),
					To:   nil,
				},
			},
			want: want{
				from: 1000,
				to:   time.Now().Unix(),
			},
			wantErr: false,
		},
		{
			name: "from equal to nil",
			args: args{
				req: &dto.ListLendingRequest{
					To: helper.Ptr(time.Now().Unix()),
				},
			},
			want: want{
				from: time.Now().Add(-weakDuration).Unix(),
				to:   time.Now().Unix(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			is := assert.New(t)

			uc := &list{}
			if err := uc.Validate(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("list.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			is.NotNil(tt.args.req.From)
			is.NotNil(tt.args.req.To)
			is.Equal(tt.want.from, *tt.args.req.From)
			is.Equal(tt.want.to, *tt.args.req.To)
		})
	}
}
