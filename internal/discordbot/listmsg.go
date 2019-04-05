package discordbot

import (
	"fmt"
	"strings"

	"github.com/zekroTJA/yuri2/internal/static"

	"github.com/zekroTJA/discordgo"
)

const (
	emojiForward = "👉"
	emojiBack    = "👈"
)

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

func NewListMessage(s *discordgo.Session, chanID, title, header string, pageCont []string, maxPageSize int) (*ListMessage, error) {
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

	// lm.emb = &discordgo.MessageEmbed{
	// 	Title: title,

	// }
	return nil, nil
}

func (lm *ListMessage) setPageEmbed(page int) {
	if page >= len(lm.pages) {
		page = len(lm.pages) - 1
	}
	lm.emb.Description = lm.header + "\n\n" + lm.pages[page]
	lm.emb.Footer.Text = fmt.Sprintf("Page %d / %d", page+1, len(lm.pages))
}

func (lm *ListMessage) updateMessage() error {
	var err error

	if lm.Message == nil {
		lm.Message, err = lm.session.ChannelMessageSendEmbed(lm.ChannelID, lm.emb)
		return err
	}

	lm.Message, err = lm.session.ChannelMessageEditEmbed(lm.ChannelID, lm.ID, lm.emb)
	return err
}
