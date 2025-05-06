package converter

import (
	"context"
	"errors"
	"fmt"
	"github.com/awakari/int-bluesky/model"
	"github.com/awakari/int-bluesky/service/firehose"
	"github.com/awakari/int-bluesky/util"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/microcosm-cc/bluemonday"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Service interface {
	EventToPost(ctx context.Context, evt *pb.CloudEvent, interestId string, t *time.Time) (post *bsky.FeedPost, err error)
	PostToEvent(ctx context.Context, src *pb.CloudEvent) (dst *pb.CloudEvent, userId string, err error)
}

type service struct {
	fmtLenMaxBodyTxt int
	htmlStripTags    *bluemonday.Policy
	didPlc           string
}

var reMultiSpace = regexp.MustCompile(`\s+`)

func NewService(
	fmtLenMaxBodyTxt int,
	addSpaceWhenStripTags bool,
	didPlc string,
) Service {
	return service{
		fmtLenMaxBodyTxt: fmtLenMaxBodyTxt,
		htmlStripTags: bluemonday.
			StrictPolicy().
			AddSpaceWhenStrippingTag(addSpaceWhenStripTags),
		didPlc: didPlc,
	}
}

func (s service) EventToPost(ctx context.Context, evt *pb.CloudEvent, interestId string, t *time.Time) (post *bsky.FeedPost, err error) {

	post = &bsky.FeedPost{
		Labels: &bsky.FeedPost_Labels{
			LabelDefs_SelfLabels: &atproto.LabelDefs_SelfLabels{
				Values: []*atproto.LabelDefs_SelfLabel{{
					Val: interestId,
				}},
			},
		},
	}

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
	post.Text = util.TruncateStringUtf8(post.Text, s.fmtLenMaxBodyTxt)
	post.Text += "\nResult Details"

	post.Embed = &bsky.FeedPost_Embed{
		EmbedExternal: &bsky.EmbedExternal{
			External: &bsky.EmbedExternal_External{
				Description: "Origin",
				Uri:         addrOrigin,
			},
		},
		EmbedRecord: &bsky.EmbedRecord{
			Record: &atproto.RepoStrongRef{
				Uri: fmt.Sprintf("https://bsky.app/profile/%s/feed/%s", s.didPlc, interestId),
			},
		},
	}

	startResultDetails := strings.LastIndex(post.Text, "Result Details")
	post.Facets = append(post.Facets,
		&bsky.RichtextFacet{
			Features: []*bsky.RichtextFacet_Features_Elem{{
				RichtextFacet_Link: &bsky.RichtextFacet_Link{
					Uri: "https://awakari.com/pub-msg.html?id=" + evt.Id + "&interestId=" + interestId,
				},
			}},
			Index: &bsky.RichtextFacet_ByteSlice{
				ByteStart: int64(startResultDetails),
				ByteEnd:   int64(startResultDetails + len("Result Details")),
			},
		},
	)

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

type ConvertFunc func(evt *pb.CloudEvent, v any) (err error)

var convSchema = map[string]any{
	"action":    toAttrStringFunc("action"),
	"createdAt": toAttrTimestampFunc("time"),    // bluesky
	"did":       toAttrStringFunc("blueskydid"), // bluesky
	"objecturl": toAttrStringFunc("objecturl"),  // as it is
	"rev":       toAttrStringFunc("blueskyrev"),
	"subject":   toAttrStringFunc("subject"),
	"text":      toTextDataFunc(),
	"time":      toAttrTimestampFunc("time"),
	"version":   toAttrInt32ElseStringFunc("blueskyversion"),
}

var ErrConversion = errors.New("conversion failure")

func (s service) PostToEvent(ctx context.Context, src *pb.CloudEvent) (dst *pb.CloudEvent, userId string, err error) {
	var raw map[string]any
	raw, userId, err = firehose.DecodePost(src.GetBinaryData())
	if raw != nil {
		dst = &pb.CloudEvent{
			Id:          src.Id,
			Source:      src.Source,
			SpecVersion: src.SpecVersion,
			Type:        "com_awakari_bluesky_v1",
			Attributes:  make(map[string]*pb.CloudEventAttributeValue),
			Data:        &pb.CloudEvent_TextData{},
		}
		err = convert(dst, raw, convSchema)
	}
	return
}

func convert(evt *pb.CloudEvent, node map[string]any, schema map[string]any) (err error) {
	for k, v := range node {
		schemaChild, schemaChildOk := schema[k]
		if schemaChildOk {
			switch schemaChildT := schemaChild.(type) {
			case ConvertFunc:
				err = errors.Join(err, schemaChildT(evt, v))
			case map[string]any:
				branch, branchOk := v.(map[string]any)
				if branchOk {
					err = errors.Join(convert(evt, branch, schemaChildT))
				}
			}
		}
	}
	return
}

func toAttrStringFunc(k string) ConvertFunc {
	return func(evt *pb.CloudEvent, v any) (err error) {
		var s string
		s, err = toString(k, v)
		if err == nil {
			evt.Attributes[k] = &pb.CloudEventAttributeValue{
				Attr: &pb.CloudEventAttributeValue_CeString{
					CeString: s,
				},
			}
		}
		return
	}
}

func toString(k string, v any) (str string, err error) {
	switch vt := v.(type) {
	case bool:
		str = strconv.FormatBool(vt)
	case int:
		str = strconv.Itoa(vt)
	case int8:
		str = strconv.Itoa(int(vt))
	case int16:
		str = strconv.Itoa(int(vt))
	case int32:
		str = strconv.Itoa(int(vt))
	case int64:
		str = strconv.FormatInt(vt, 10)
	case float32:
		switch float32(int(vt)) == vt {
		case true:
			str = strconv.Itoa(int(vt))
		default:
			str = fmt.Sprintf("%f", vt)
		}
	case float64:
		switch float64(int(vt)) == vt {
		case true:
			str = strconv.Itoa(int(vt))
		default:
			str = fmt.Sprintf("%f", vt)
		}
	case string:
		str = vt
	default:
		err = fmt.Errorf("%w: key: %s, value: %v, type: %s, expected: string/bool/int/float", ErrConversion, k, v, reflect.TypeOf(v))
	}
	return
}

func toAttrTimestampFunc(k string) ConvertFunc {
	return func(evt *pb.CloudEvent, v any) (err error) {
		switch vt := v.(type) {
		case int:
			if vt > 1e15 {
				// timestamp is unix micros
				vt /= 1_000_000
			}
			if vt > 1e12 {
				// timestamp is unix millis
				vt /= 1_000
			}
			evt.Attributes[k] = &pb.CloudEventAttributeValue{
				Attr: &pb.CloudEventAttributeValue_CeTimestamp{
					CeTimestamp: &timestamppb.Timestamp{
						Seconds: int64(vt),
					},
				},
			}
		case int32:
			evt.Attributes[k] = &pb.CloudEventAttributeValue{
				Attr: &pb.CloudEventAttributeValue_CeTimestamp{
					CeTimestamp: &timestamppb.Timestamp{
						Seconds: int64(vt),
					},
				},
			}
		case int64:
			if vt > 1e15 {
				// timestamp is unix micros
				vt /= 1_000_000
			}
			if vt > 1e12 {
				// timestamp is unix millis
				vt /= 1_000
			}
			evt.Attributes[k] = &pb.CloudEventAttributeValue{
				Attr: &pb.CloudEventAttributeValue_CeTimestamp{
					CeTimestamp: &timestamppb.Timestamp{
						Seconds: vt,
					},
				},
			}
		case float32:
			if vt > 1e15 {
				// timestamp is unix micros
				vt /= 1_000_000
			}
			if vt > 1e12 {
				// timestamp is unix millis
				vt /= 1_000
			}
			evt.Attributes[k] = &pb.CloudEventAttributeValue{
				Attr: &pb.CloudEventAttributeValue_CeTimestamp{
					CeTimestamp: &timestamppb.Timestamp{
						Seconds: int64(vt),
					},
				},
			}
		case float64:
			if vt > 1e15 {
				// timestamp is unix micros
				vt /= 1_000_000
			}
			if vt > 1e12 {
				// timestamp is unix millis
				vt /= 1_000
			}
			evt.Attributes[k] = &pb.CloudEventAttributeValue{
				Attr: &pb.CloudEventAttributeValue_CeTimestamp{
					CeTimestamp: &timestamppb.Timestamp{
						Seconds: int64(vt),
					},
				},
			}
		case string:
			var t time.Time
			t, err = time.Parse(time.RFC3339, vt)
			evt.Attributes[k] = &pb.CloudEventAttributeValue{
				Attr: &pb.CloudEventAttributeValue_CeTimestamp{
					CeTimestamp: timestamppb.New(t),
				},
			}
		default:
			err = fmt.Errorf("%w: key: %s, value %v, type: %s, expected timestamp in RFC3339 format", ErrConversion, k, v, reflect.TypeOf(k))
		}
		return
	}
}

func toAttrInt32ElseStringFunc(k string) ConvertFunc {
	return func(evt *pb.CloudEvent, v any) (err error) {
		i, ok := toInt32(v)
		switch ok {
		case true:
			evt.Attributes[k] = &pb.CloudEventAttributeValue{
				Attr: &pb.CloudEventAttributeValue_CeInteger{
					CeInteger: i,
				},
			}
		default:
			var s string
			s, err = toString(k, v)
			if err == nil {
				evt.Attributes[k] = &pb.CloudEventAttributeValue{
					Attr: &pb.CloudEventAttributeValue_CeString{
						CeString: s,
					},
				}
			}
		}
		return
	}
}

func toInt32(v any) (i int32, ok bool) {
	switch vt := v.(type) {
	case bool:
		if vt {
			i = 1
		}
		ok = true
	case int:
		if vt >= math.MinInt32 && vt <= math.MaxInt32 {
			i = int32(vt)
			ok = true
		}
	case int8:
		i = int32(vt)
		ok = true
	case int16:
		i = int32(vt)
		ok = true
	case int32:
		i = vt
		ok = true
	case int64:
		if vt >= math.MinInt32 && vt <= math.MaxInt32 {
			i = int32(vt)
			ok = true
		}
	case float32:
		if vt >= math.MinInt32 && vt <= math.MaxInt32 {
			i = int32(vt)
			ok = float32(i) == vt
		}
	case float64:
		if vt >= math.MinInt32 && vt <= math.MaxInt32 {
			i = int32(vt)
			ok = float64(i) == vt
		}
	case string:
		i64, err := strconv.ParseInt(vt, 10, 32)
		if err == nil && i64 >= math.MinInt32 && i64 <= math.MaxInt32 {
			i = int32(i64)
			ok = true
		}
	}
	return
}

func toTextDataFunc() ConvertFunc {
	return func(evt *pb.CloudEvent, v any) (err error) {
		vs, vsOk := v.(string)
		switch vsOk {
		case true:
			evt.Data = &pb.CloudEvent_TextData{
				TextData: vs,
			}
		default:
			err = fmt.Errorf("invalid value type %T", v)
		}
		return
	}
}
