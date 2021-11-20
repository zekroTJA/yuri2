package discordbot

import (
	"fmt"
	"strings"
	"time"

	"github.com/zekroTJA/yuri2/internal/static"

	"github.com/bwmarrin/discordgo"
)

const (
	emojiForward = "👉"
	emojiBack    = "👈"
)

// A ListMessage is a chat message which line content can
// be browsed by using reactions to turn "pages".
type ListMessage struct {
	*discordgo.Message

	header      string
	maxPageSize int
	chanID      string

	session  *discordgo.Session
	currPage int
	pages    []string
	emb      *discordgo.MessageEmbed
	unhandle func()
}

// NewListMessage creates the ListMessage and sends it to the
// specified channel.
//   s           : Discord bot session
//   chanID      : The text channel ID where the message will be posted in
//   title       : The embed title
//   header      : A text block which will be at the top of every page
//   pageCont    : Content of the pages. Elements will be joined with a line break ('\n')
//   maxPageSize : the maximum ammount of content lines per page
//   startPage   : the page, which should be shown on first send
func NewListMessage(s *discordgo.Session, chanID, title, header string, pageCont []string, maxPageSize int, startPage int) (*ListMessage, error) {
	var err error

	lm := &ListMessage{
		header:      header,
		maxPageSize: maxPageSize,
		chanID:      chanID,
		session:     s,
		emb: &discordgo.MessageEmbed{
			Color: static.ColorDefault,
			Title: title,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Page 1",
			},
		},
	}

	_pageCount := float64(len(pageCont)) / float64(maxPageSize)
	pageCount := int(_pageCount)
	if _pageCount > float64(pageCount) {
		pageCount++
	}

	lm.pages = make([]string, pageCount)

	for i := range lm.pages {
		from := i * maxPageSize
		to := from + maxPageSize
		if to > len(pageCont) {
			to = len(pageCont)
		}

		lm.pages[i] = strings.Join(pageCont[from:to], "\n")
	}

	lm.setPageEmbed(startPage)
	if err = lm.updateMessage(); err != nil {
		return nil, err
	}

	if err = lm.session.MessageReactionAdd(lm.chanID, lm.ID, emojiBack); err != nil {
		return nil, err
	}
	if err = lm.session.MessageReactionAdd(lm.chanID, lm.ID, emojiForward); err != nil {
		return nil, err
	}

	lm.unhandle = s.AddHandler(lm.reactionHandler)

	return lm, err
}

// setEmbed sets the embed content for
// the specified page.
func (lm *ListMessage) setPageEmbed(page int) {
	if page >= len(lm.pages) && len(lm.pages) > 1 {
		page = len(lm.pages) - 1
	}

	lm.emb.Description = lm.header + "\n\n" + lm.pages[page]
	lm.emb.Footer.Text = fmt.Sprintf("Page %d / %d", page+1, len(lm.pages))
}

// updateMessage actually updates the message
// in the text channel.
func (lm *ListMessage) updateMessage() error {
	var err error

	if lm.Message == nil {
		lm.Message, err = lm.session.ChannelMessageSendEmbed(lm.chanID, lm.emb)
		return err
	}

	lm.Message, err = lm.session.ChannelMessageEditEmbed(lm.chanID, lm.ID, lm.emb)
	return err
}

// turnForward is shorthand for increasing
// the page value by one (if larger or equal
// to page size -> 0), setting the embeds
// content and updating the message.
func (lm *ListMessage) turnForward() error {
	lm.currPage++
	if lm.currPage >= len(lm.pages) {
		lm.currPage = 0
	}
	lm.setPageEmbed(lm.currPage)
	return lm.updateMessage()
}

// turnBack is shorthand for decreasing
// the page value by one (if smaller or equal
// to 0 -> pages count - 1), setting the embeds
// content and updating the message.
func (lm *ListMessage) turnBack() error {
	lm.currPage--
	if lm.currPage < 0 {
		lm.currPage = len(lm.pages) - 1
	}
	lm.setPageEmbed(lm.currPage)
	return lm.updateMessage()
}

// reactionHandler is the handler function for the
// message reaction add event which will be registered
// to the bot session.
func (lm *ListMessage) reactionHandler(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
	if e.MessageID != lm.ID || e.UserID == lm.session.State.User.ID {
		return
	}

	switch e.Emoji.Name {
	case emojiBack:
		lm.turnBack()
	case emojiForward:
		lm.turnForward()
	}

	s.MessageReactionRemove(lm.chanID, lm.ID, e.Emoji.Name, e.UserID)
}

// Delete will delete the message from
// text channel and unregisters the
// reaction event handler.
func (lm *ListMessage) Delete() error {
	lm.unhandle()
	return lm.session.ChannelMessageDelete(lm.chanID, lm.ID)
}

// DeleteAfter executes Delete after the
// specified duration.
func (lm *ListMessage) DeleteAfter(d time.Duration) {
	time.AfterFunc(d, func() {
		lm.Delete()
	})
}
