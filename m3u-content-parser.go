package gom3u_content_parser

import (
	"regexp"
)

type M3UContentParser struct {
	m3uFileContent string
	dirtyItems     []string
	Items          []M3UItem
	CountItems     int

	TvgUrl  string
	Cache   int
	Refresh int

	Offsets int
	Limits  int
}

func NewM3UContentParser() *M3UContentParser {
	return &M3UContentParser{}
}

func (parser *M3UContentParser) LoadSource(source string, fromFile bool) *M3UContentParser {
	if fromFile {
		parser.m3uFileContent = ReadStringContentFromFile(source)
	} else {
		parser.m3uFileContent = ReadStringContentFromRemote(source)
	}

	return parser
}

func (parser *M3UContentParser) Parse() *M3UContentParser {

	a := regexp.MustCompile(`(#EXTINF:0|#EXTINF:-1|#EXTINF:-1,)`)
	parser.dirtyItems = a.Split(parser.m3uFileContent, -1)

	// first is #EXTM3U tag
	parser.parseAndSetTvgUrl(parser.dirtyItems[0])
	parser.dirtyItems = parser.dirtyItems[1:]

	for _, item := range parser.dirtyItems {
		parser.Items = append(parser.Items, *NewM3UItem(item))
		parser.CountItems++
	}

	return parser
}

func (parser *M3UContentParser) parseAndSetTvgUrl(url string) *M3UContentParser {
	re := regexp.MustCompile(`"([^"]+)"`)
	matches := re.FindAllString(url, -1)

	// if tv program exist
	if len(matches) > 0 {
		parser.TvgUrl = matches[0]
	}

	return parser
}

func (parser *M3UContentParser) GetTvgUrl() string {
	return parser.TvgUrl
}

func (parser *M3UContentParser) GetM3UContent() string {
	return parser.m3uFileContent
}

func (parser *M3UContentParser) GetDirtyItems() []string {
	return parser.dirtyItems
}

func (parser *M3UContentParser) GetItems() []M3UItem {
	return parser.Items
}

func (parser *M3UContentParser) Offset(offset int) *M3UContentParser {
	parser.Offsets = offset

	return parser
}

func (parser *M3UContentParser) Limit(limit int) *M3UContentParser {
	parser.Limits = limit

	return parser
}

func (parser *M3UContentParser) All() []M3UItem {
	if parser.Limits <= 0 {
		parser.Limits = parser.CountItems
	}

	return parser.Items[parser.Offsets:parser.Limits]
}
