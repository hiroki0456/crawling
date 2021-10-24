package healthcheck

import (
	"github.com/slack-go/slack"
)

func NoticeSlack(illegalCheck error) {
	// divのワークスペースで使用できるアクセストークン
	tkn := "xoxb-81666833746-2632126536469-qcIWScF3h2ba75TwH87XoQpw"
	c := slack.New(tkn)

	_, _, err := c.PostMessage("#example_test", slack.MsgOptionText(illegalCheck.Error(), true))
	if err != nil {
		panic(err)
	}
}
