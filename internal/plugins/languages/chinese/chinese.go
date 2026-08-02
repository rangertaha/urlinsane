package chinese

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "zh"

type Chinese struct {
	code         string
	name         string
	description  string
	numerals     map[string][]string
	graphemes    []string
	vowels       []string
	misspellings [][]string
	homophones   [][]string
	antonyms     map[string][]string
	homoglyphs   map[string][]string
}

func (l *Chinese) Id() string                      { return l.code }
func (l *Chinese) Name() string                    { return l.name }
func (l *Chinese) Description() string             { return l.description }
func (l *Chinese) Numerals() map[string][]string   { return l.numerals }
func (l *Chinese) Cardinal() map[string]string     { return languages.NumeralMap(l.numerals, 0) }
func (l *Chinese) Ordinal() map[string]string      { return languages.NumeralMap(l.numerals, 1) }
func (l *Chinese) Graphemes() []string             { return l.graphemes }
func (l *Chinese) Vowels() []string                { return l.vowels }
func (l *Chinese) Misspellings() [][]string        { return l.misspellings }
func (l *Chinese) Homophones() [][]string          { return l.homophones }
func (l *Chinese) Antonyms() map[string][]string   { return l.antonyms }
func (l *Chinese) Homoglyphs() map[string][]string { return l.homoglyphs }
func (l *Chinese) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}
func (l *Chinese) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}
func (l *Chinese) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}
func (l *Chinese) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	zhMisspellings = [][]string{
		{"布署", "部署"},
		{"竞然", "竟然"},
		{"既使", "即使"},
		{"迫不急待", "迫不及待"},
		{"一如继往", "一如既往"},
		{"甘败下风", "甘拜下风"},
		{"走头无路", "走投无路"},
		{"变本加利", "变本加厉"},
		{"再接再励", "再接再厉"},
		{"穿流不息", "川流不息"},
		{"弱不经风", "弱不禁风"},
		{"出奇不意", "出其不意"},
		{"各抒已见", "各抒己见"},
		{"谈笑风声", "谈笑风生"},
		{"世外桃园", "世外桃源"},
		{"张灯结采", "张灯结彩"},
		{"破斧沉舟", "破釜沉舟"},
		{"不加思索", "不假思索"},
		{"迭挫迭战", "屡挫屡战"},
		{"名燥一时", "名噪一时"},
		{"计日成功", "计日程功"},
		{"悬梁刺骨", "悬梁刺股"},
		{"一愁莫展", "一筹莫展"},
		{"黄梁美梦", "黄粱美梦"},
		{"不径而走", "不胫而走"},
		{"纷至踏来", "纷至沓来"},
		{"如愿以尝", "如愿以偿"},
		{"心心相映", "心心相印"},
		{"仗义直言", "仗义执言"},
		{"振聋发馈", "振聋发聩"},
		{"纵横弛骋", "纵横驰骋"},
		{"不落巢臼", "不落窠臼"},
		{"和霭可亲", "和蔼可亲"},
		{"挖墙角", "挖墙脚"},
	}

	zhHomophones = [][]string{
		{"句号", "。", "."},
		{"艾特", "@"},
		{"横杠", "-"},
		{"斜杠", "/"},

		{"制定", "制订"},
		{"辩证", "辨证"},
		{"秉性", "禀性"},
		{"反应", "反映"},
		{"其他", "其它"},
		{"账户", "帐户"},
		{"成分", "成份"},
		{"坐落", "座落"},
		{"必须", "必需"},
		{"以至", "以致"},
		{"权力", "权利"},
		{"起用", "启用"},
		{"终生", "终身"},
		{"交待", "交代"},
		{"定金", "订金"},
		{"法治", "法制"},
		{"分辩", "分辨"},
		{"妨害", "妨碍"},
		{"勾通", "沟通"},
		{"含糊", "含混"},
		{"决不", "绝不"},
		{"连接", "联接"},
		{"灌注", "贯注"},
		{"需要", "须要"},
		{"学历", "学力"},
		{"意气", "义气"},
		{"预定", "预订"},
		{"原型", "原形"},
		{"中止", "终止"},
	}
	zhAntonyms = map[string][]string{
		"好": {"坏"},
		"坏": {"好"},
		"大": {"小"},
		"小": {"大"},
		"高": {"低"},
		"低": {"高"},
		"长": {"短"},
		"短": {"长"},
		"多": {"少"},
		"少": {"多"},
		"新": {"旧"},
		"旧": {"新"},
		"早": {"晚"},
		"晚": {"早"},
		"快": {"慢"},
		"慢": {"快"},
		"热": {"冷"},
		"冷": {"热"},
		"明": {"暗"},
		"暗": {"明"},
		"强": {"弱"},
		"弱": {"强"},
		"胖": {"瘦"},
		"瘦": {"胖"},
		"轻": {"重"},
		"重": {"轻"},
		"远": {"近"},
		"近": {"远"},
		"上": {"下"},
		"下": {"上"},
		"前": {"后"},
		"后": {"前"},
		"左": {"右"},
		"右": {"左"},
		"对": {"错"},
		"错": {"对"},
		"真": {"假"},
		"假": {"真"},
		"开": {"关"},
		"关": {"开"},
		"进": {"出"},
		"出": {"进"},
		"来": {"去"},
		"去": {"来"},
		"有": {"无"},
		"无": {"有"},
		"买": {"卖"},
		"卖": {"买"},
		"赢": {"输"},
		"输": {"赢"},
	}

	Language = Chinese{
		code:        LANGUAGE,
		name:        "Chinese",
		description: "Chinese is a group of languages written with Han characters (Simplified/Traditional)",
		numerals: map[string][]string{
			"0":          {"零"},
			"1":          {"一", "第一"},
			"2":          {"二", "第二"},
			"3":          {"三", "第三"},
			"4":          {"四", "第四"},
			"5":          {"五", "第五"},
			"6":          {"六", "第六"},
			"7":          {"七", "第七"},
			"8":          {"八", "第八"},
			"9":          {"九", "第九"},
			"10":         {"十", "第十"},
			"11":         {"十一", "第十一"},
			"12":         {"十二", "第十二"},
			"13":         {"十三", "第十三"},
			"14":         {"十四", "第十四"},
			"15":         {"十五", "第十五"},
			"16":         {"十六", "第十六"},
			"17":         {"十七", "第十七"},
			"18":         {"十八", "第十八"},
			"19":         {"十九", "第十九"},
			"20":         {"二十", "第二十"},
			"30":         {"三十"},
			"40":         {"四十"},
			"50":         {"五十"},
			"60":         {"六十"},
			"70":         {"七十"},
			"80":         {"八十"},
			"90":         {"九十"},
			"100":        {"一百", "百"},
			"1000":       {"一千", "千"},
			"10000":      {"一万", "万"},
			"1000000":    {"一百万", "百万"},
			"1000000000": {"十亿", "十億"},
		},
		// Graphemes for CJK are huge; keep this minimal until a better tokenizer exists.
		graphemes:    []string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九", "十"},
		vowels:       []string{},
		misspellings: zhMisspellings,
		homophones:   zhHomophones,
		antonyms:     zhAntonyms,
		homoglyphs:   languages.DefaultLatinHomoglyphs(),
	}
)

func init() { languages.AddLanguage(LANGUAGE, func() internal.Language { return &Language }) }
