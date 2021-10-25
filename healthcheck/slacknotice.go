package healthcheck

import (
	"fmt"

	"github.com/slack-go/slack"
)

func NoticeSlack(res error) {
	// divのワークスペースで使用できるアクセストークン
	tkn := "xoxb-81666833746-2632126536469-qcIWScF3h2ba75TwH87XoQpw"
	c := slack.New(tkn)
	attachment := slack.Attachment{
		Pretext: res.Error(),
	}
	_, _, err := c.PostMessage(
		"#example_test",
		slack.MsgOptionText("*【freee】のヘルスチェックが失敗しました。*", true),
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}
