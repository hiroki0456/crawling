package crawlingrepository

import (
	"context"
	"testing"

	"cloud.google.com/go/spanner"
	pb "upsider.crawling/crawlingproto"
	"upsider.crawling/testhelper"
)

var userData1 User = User{
	Id:               "test_data1",
	UserIdOfficeName: "hiroki_test",
	UserId:           "hiroki",
	OfficeName:       "株式会社test",
	LastId:           "123456",
}

func OfficeRead(t *testing.T) {
	cases := map[string]struct {
		schema string
		data   []*spanner.Mutation
		cardID string
		ctx    context.Context

		// output
		want *pb.Office
	}{
		"Get physical card": {
			data: []*spanner.Mutation{
				testhelper.InsertStruct(t, "Cards", cardData1),
				testhelper.InsertStruct(t, "Cards", cardData2),
				testhelper.InsertStruct(t, "Cards", cardData2Latest),
			},
			cardID: "test_card1",
			ctx:    context.Background(),
			want: &cardreader.ReadCardResponse{
				Card: []*cardpb.Card{
					cardData1.ConvertToPB(),
					cardData2.ConvertToPB(),
					cardData2Latest.ConvertToPB(),
				},
			},
		},
	}
}
