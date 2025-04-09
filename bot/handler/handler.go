package handler

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"

	"reacjira/config"
	myjira "reacjira/jira"
	myslack "reacjira/slack"
)

// CommandHandler handle command.
type CommandHandler struct {
	slackMessenger *myslack.Messenger
	jira           *myjira.MyJira
	slackCtx       config.SlackCtx
	botProfile     config.BotProfile
	reacjiras      []config.Reacjira
}

// NewCommandHandler is ...
func NewCommandHandler(
	slack *slack.Client,
	slackCtx config.SlackCtx,
	jiraCtx config.JiraCtx,
	botProfile config.BotProfile,
	reacjiras []config.Reacjira,
) *CommandHandler {
	return &CommandHandler{
		myslack.New(slack),
		myjira.New(jiraCtx.Host, jiraCtx.Email, jiraCtx.Token),
		slackCtx,
		botProfile,
		reacjiras,
	}
}

// HandleCommand handle command
func (commandHandler *CommandHandler) HandleCommand(ev *slackevents.ReactionAddedEvent) {
	// 対応するreacjiraを探す
	var found config.Reacjira

	for _, r := range commandHandler.reacjiras {
		if r.Emoji == ev.Reaction {
			found = r
		}
	}

	if found.Emoji == "" {
		// 対応していないemojiなら無視する
		return
	}

	commandHandler.createTicket(ev, found)
}

func (commandHandler *CommandHandler) createTicket(ev *slackevents.ReactionAddedEvent, reacjira config.Reacjira) {
	log.Printf("A slack reaction was received: channel: %s, timestamp: %s, user: %s", ev.Item.Channel, ev.Item.Timestamp, ev.User)
	log.Printf("prepare to create a ticket %+v", reacjira)

	// reactionがついたメッセージを取得
	msg, err := findMessage(commandHandler.slackMessenger, ev)
	if err != nil {
		log.Printf("findMessage error: %s", err.Error())
		return
	}
	var replyTo = msg.ThreadTimestamp
	if msg.ThreadTimestamp == "" {
		replyTo = msg.Timestamp
	}

	// 何回もチケットを作らないように２個目以降の発言には反応しないようにしている。
	//（この方式だとemojiをつけてキャンセル後に、またemojiをつけると再度チケットが作成されてしまうので本当はチケットタイトルを検索しに行ったほうが良いが実装が面倒なので手抜き実装）
	for _, reaction := range msg.Reactions {
		if reaction.Name == reacjira.Emoji && reaction.Count > 1 {
			log.Printf("multiple reacjira reactions were found(%s, %d). skip", reaction.Name, reaction.Count)
			return
		}
	}

	// リアクションを行った人をreporterとする。reporterのuserIdをuserに変換
	reporter, err := commandHandler.slackMessenger.SearchUser(ev.User)
	if err != nil {
		log.Printf("searchUser error: %s", err.Error())
		_ = commandHandler.slackMessenger.SendMessages([]string{err.Error()}, ev.Item.Channel, replyTo)
		return
	}
	// reporterはjira ticketを作るときに必須なので空ならエラーにする
	if reporter.Profile.Email == "" {
		log.Printf("reporter email is empty (userId=%s)", ev.User)
		_ = commandHandler.slackMessenger.SendMessages([]string{fmt.Sprintf("reporter email is empty(userId=%s)", ev.User)}, ev.Item.Channel, replyTo)
		return
	}

	// channelIdをchannelに変換
	channel, err := commandHandler.slackMessenger.SearchChannel(ev.Item.Channel)
	if err != nil {
		log.Printf("searchChannel error: %s", err.Error())
		_ = commandHandler.slackMessenger.SendMessages([]string{err.Error()}, ev.Item.Channel, replyTo)
		return
	}

	// messageからpermlinkを取得
	link, err := getPermLink(commandHandler.slackMessenger, ev, msg)
	if err != nil {
		log.Printf("getPermLink error: %s", err.Error())
		_ = commandHandler.slackMessenger.SendMessages([]string{err.Error()}, ev.Item.Channel, replyTo)
		return
	}

	// チケットタイトルは発言 (改行を除去し、200文字までにする)
	title := strings.ReplaceAll(msg.Text, "\n", " ")
	limit := int(math.Min(float64(len(title)), 200))
	title = title[0:limit]

	// チケットの中身はslackリンク + 予め指定されたformat
	description := fmt.Sprintf(`auto generated
from: %s
at: %s
%s`, link, channel.Name, reacjira.Description)
	log.Printf("reporter:%s, channel: %s, title: %s", reporter.Name, channel.Name, title)

	log.Printf("attempt to create a ticket.")
	ticket, err := commandHandler.jira.CreateTicket(
		reacjira.Project,
		reacjira.IssueType,
		reacjira.EpicKey,
		reporter.Profile.Email,
		title,
		description,
	)
	if err != nil {
		log.Printf("createTicket error: %s", err.Error())
		_ = commandHandler.slackMessenger.SendMessages([]string{err.Error()}, ev.Item.Channel, replyTo)
		return
	}

	slackMessage := *ticket
	if reacjira.SlackMessage != "" {
		slackMessage += fmt.Sprintf("\n%s", reacjira.SlackMessage)
	}
	err = commandHandler.slackMessenger.SendMessages([]string{slackMessage}, ev.Item.Channel, replyTo)
	if err != nil {
		log.Printf("sendMessages error: %s", err.Error())
		return
	}
}

func findMessage(slackMessenger *myslack.Messenger, ev *slackevents.ReactionAddedEvent) (*slack.Message, error) {
	// reaction情報から元のメッセージを取得
	msg, err := slackMessenger.SearchMsg(ev.Item.Channel, ev.Item.Timestamp)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	log.Printf("slack message search finish: %s, %s, %s, reply:%d", msg.User, msg.Text, msg.Timestamp, msg.ReplyCount)
	return msg, nil
}

// getPermLink APIを使ってtimestampからpermlinkを取得する
func getPermLink(slackMessenger *myslack.Messenger, ev *slackevents.ReactionAddedEvent, msg *slack.Message) (string, error) {
	// msg.Channelは空文字が入っているので使えない。代わりにev.Item.Channelを使う
	link, err := slackMessenger.GetPermLink(ev.Item.Channel, msg.Timestamp)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	log.Printf("got slack permlink: %s", link)
	return link, nil
}
