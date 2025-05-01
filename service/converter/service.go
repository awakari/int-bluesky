package converter

import (
	"context"
	"github.com/awakari/int-bluesky/model"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/microcosm-cc/bluemonday"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

type Service interface {
	EventToPost(ctx context.Context, evt *pb.CloudEvent, interestId string, t *time.Time) (post *bsky.FeedPost, err error)
}

type service struct {
	fmtLenMaxBodyTxt int
	htmlStripTags    *bluemonday.Policy
}

var reMultiSpace = regexp.MustCompile(`\s+`)

func NewService(
	fmtLenMaxBodyTxt int,
	addSpaceWhenStripTags bool,
) Service {
	return service{
		fmtLenMaxBodyTxt: fmtLenMaxBodyTxt,
		htmlStripTags: bluemonday.
			StrictPolicy().
			AddSpaceWhenStrippingTag(addSpaceWhenStripTags),
	}
}

func (s service) EventToPost(ctx context.Context, evt *pb.CloudEvent, interestId string, t *time.Time) (post *bsky.FeedPost, err error) {

	post = &bsky.FeedPost{}
	switch t {
	case nil:
		post.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	default:
		post.CreatedAt = t.UTC().Format(time.RFC3339)
	}
	attrLang, langPresent := evt.Attributes[model.CeKeyLanguage]
	if langPresent {
		post.Langs = append(post.Langs, attrLang.GetCeString())
	}

	var addrOrigin string
	attrObjUrl, attrObjUrlPresent := evt.Attributes[model.CeKeyObjectUrl]
	if attrObjUrlPresent {
		addrOrigin = attrObjUrl.GetCeString()
		if addrOrigin == "" {
			addrOrigin = attrObjUrl.GetCeUri()
		}
	}
	if addrOrigin == "" {
		addrOrigin = evt.Source
	}
	if strings.HasPrefix(addrOrigin, "@") { // telegram source
		addrOrigin = "https://t.me/" + addrOrigin[1:]
	}

	post.Embed = &bsky.FeedPost_Embed{
		EmbedExternal: &bsky.EmbedExternal{
			External: &bsky.EmbedExternal_External{
				Description: "Origin",
				Uri:         addrOrigin,
			},
		},
	}

	attrCats, _ := evt.Attributes[model.CeKeyCategories]
	cats := strings.Split(attrCats.GetCeString(), " ")
	for _, cat := range cats {
		var tagName string
		switch strings.HasPrefix(cat, "#") {
		case true:
			tagName = cat[1:]
		default:
			tagName = cat
		}
		if len(tagName) > 0 {
			post.Tags = append(post.Tags, tagName)
		}
	}

	post.Text = eventSummaryText(evt)
	post.Text = s.htmlStripTags.Sanitize(post.Text)
	post.Text = reMultiSpace.ReplaceAllString(post.Text, " ")
	post.Text = truncateStringUtf8(post.Text, s.fmtLenMaxBodyTxt)

	addrInterest := "https://awakari.com/sub-details.html?id=" + interestId
	post.Text += "\n<a href=\"" + addrInterest + "\">Interest</a>\n"

	addrEvtAttrs := "https://awakari.com/pub-msg.html?id=" + evt.Id + "&interestId=" + interestId
	post.Text += "\n<a href=\"" + addrEvtAttrs + "\">Result Details</a>"

	return
}

func eventSummaryText(evt *pb.CloudEvent) (txt string) {

	attrHead, headPresent := evt.Attributes[model.CeKeyHeadline]
	if headPresent {
		txt = strings.TrimSpace(attrHead.GetCeString())
	}

	attrTitle, titlePresent := evt.Attributes[model.CeKeyTitle]
	if titlePresent {
		if txt != "" {
			txt += " "
		}
		txt += strings.TrimSpace(attrTitle.GetCeString())
	}

	attrDescr, descrPresent := evt.Attributes[model.CeKeyDescription]
	if descrPresent {
		if txt != "" {
			txt += " "
		}
		txt += strings.TrimSpace(attrDescr.GetCeString())
	}

	attrSummary, summaryPresent := evt.Attributes[model.CeKeySummary]
	if summaryPresent {
		if txt != "" {
			txt += " "
		}
		txt += strings.TrimSpace(attrSummary.GetCeString())
	}

	if evt.GetTextData() != "" {
		if txt != "" {
			txt += " "
		}
		txt += strings.TrimSpace(evt.GetTextData())
	}
	if txt == "" {
		attrName, namePresent := evt.Attributes[model.CeKeyName]
		if namePresent {
			txt = strings.TrimSpace(attrName.GetCeString()) + "<br/>"
		}
	}

	return
}

func truncateStringUtf8(s string, lenMax int) string {
	if len(s) <= lenMax {
		return s
	}
	// Ensure we don't split a UTF-8 character in the middle.
	for i := lenMax - 3; i > 0; i-- {
		if utf8.RuneStart(s[i]) {
			return s[:i] + "..."
		}
	}
	return ""
}
