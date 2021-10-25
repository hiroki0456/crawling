package crawlingrepository

import (
	"context"
	"testing"

	"cloud.google.com/go/spanner"
	"google.golang.org/protobuf/proto"
	"upsider.crawling/config"
	"upsider.crawling/crawlingproto"
)

// var today = time.Now()

var userData1 User = User{
	Id:               "test_data1",
	UserIdOfficeName: "hiroki_株式会社test1",
	UserId:           "hiroki",
	OfficeName:       "株式会社test1",
	LastId:           "123456",
}

var userData2 User = User{
	Id:               "test_data1",
	UserIdOfficeName: "hiroki_株式会社test2",
	UserId:           "hiroki", 
	OfficeName:       "株式会社test2",
	LastId:           "12345678",
}

// var bankData1 CreateBank = CreateBank{
// 	Id         "dkfalhgarhdgdffgs"
//     UserId     "hiroki"
//     BankId     "12345678"
//     lastCommit string     `spanner:"LastCommit"`
//     OfficeName string     `spanner:"OfficeName"`
//     BankName   string     `spanner:"BankName"`
//     Amount     int64      `spanner:"Amount"`
//     UpdatedAt  *time.Time `spanner:"updatedAt"`
// }

func TestOfficeRead(t *testing.T) {
	cases := map[string]struct {
		shema string
		data  []*spanner.Mutation
		ctx   context.Context

		want *crawlingproto.Office
	}{
		"Get Office Name": {
			data: []*spanner.Mutation{
				config.InsertStruct(t, "Users", userData1),
				config.InsertStruct(t, "Users", userData2),
			},
			ctx: context.Background(),
			want: &crawlingproto.Office{
				OfficeName: "株式会社test1",
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			schema := tc.shema
			if schema == "" {
				schema = config.UsersTable
			}
			cfg := &config.MockSpanner{
				Dbschemas: []string{schema},
				Mutations: tc.data,
			}
			client := config.SetUpTestDatabase(t, cfg)
			s := &CrawlingRead{client: client}
			req := &crawlingproto.FreeeRequest{UserInput: &crawlingproto.UserInput{UserId: "hiroki"}}
			res, err := s.OfficeRead(tc.ctx, req)
			if err != nil {
				return
			}
			for _, office := range res {
				proto.Equal(tc.want, office)
			}
		})
	}
}

// func TestGetBankCount(t *testing.T) {
// 	cases := map[string]struct {
// 		shema string
// 		data  []*spanner.Mutation
// 		ctx   context.Context

// 		want *crawlingproto.Office
// 	}{
// 		"Get Office Name": {
// 			data: []*spanner.Mutation{
// 				config.InsertStruct(t, "Users", userData1),
// 				config.InsertStruct(t, "Users", userData2),
// 			},
// 			ctx: context.Background(),
// 			want: &crawlingproto.Office{
// 				OfficeName: "株式会社test1",
// 			},
// 		},
// 	}
// }
