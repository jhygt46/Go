package main

import (
	"fmt"
	"strings"
	"strconv"
	"unicode"
	"io/ioutil"
	"encoding/json"
	"github.com/valyala/fasthttp"
)

type MyHandler struct {}

func main() {

	/*

	//https://pkg.go.dev/unicode

	writeUnicode("Cc", unicode.Cc)
	writeUnicode("Cf", unicode.Cf)
	writeUnicode("Co", unicode.Co)
	writeUnicode("Cs", unicode.Cs)
	writeUnicode("Digit", unicode.Digit)
	writeUnicode("Nd", unicode.Nd)
	writeUnicode("Letter", unicode.Letter)
	writeUnicode("L", unicode.L)
	writeUnicode("Lm", unicode.Lm)
	writeUnicode("Lo", unicode.Lo)
	writeUnicode("Lower", unicode.Lower)
	writeUnicode("Ll", unicode.Ll)
	writeUnicode("Mark", unicode.Mark)
	writeUnicode("M", unicode.M)
	writeUnicode("Mc", unicode.Mc)
	writeUnicode("Me", unicode.Me)
	writeUnicode("Mn", unicode.Mn)
	writeUnicode("Nl", unicode.Nl)
	writeUnicode("No", unicode.No)
	writeUnicode("Number", unicode.Number)
	writeUnicode("N", unicode.N)
	writeUnicode("Other", unicode.Other)
	writeUnicode("C", unicode.C)
	writeUnicode("Pc", unicode.Pc)
	writeUnicode("Pd", unicode.Pd)
	writeUnicode("Pe", unicode.Pe)
	writeUnicode("Pf", unicode.Pf)
	writeUnicode("Pi", unicode.Pi)
	writeUnicode("Po", unicode.Po)
	writeUnicode("Ps", unicode.Ps)
	writeUnicode("Punct", unicode.Punct)
	writeUnicode("P", unicode.P)
	writeUnicode("Sc", unicode.Sc)
	writeUnicode("Sk", unicode.Sk)
	writeUnicode("Sm", unicode.Sm)
	writeUnicode("So", unicode.So)
	writeUnicode("Space", unicode.Space)
	writeUnicode("Z", unicode.Z)
	writeUnicode("Symbol", unicode.Symbol)
	writeUnicode("S", unicode.S)
	writeUnicode("Title", unicode.Title)
	writeUnicode("Lt", unicode.Lt)
	writeUnicode("Upper", unicode.Upper)
	writeUnicode("Lu", unicode.Lu)
	writeUnicode("Zl", unicode.Zl)
	writeUnicode("Zp", unicode.Zp)
	writeUnicode("Zs", unicode.Zs)



	writeUnicode("ASCII_Hex_Digit", unicode.ASCII_Hex_Digit)
	writeUnicode("Bidi_Control", unicode.Bidi_Control)
	writeUnicode("Dash", unicode.Dash)
	writeUnicode("Deprecated", unicode.Deprecated)
	writeUnicode("Diacritic", unicode.Diacritic)
	writeUnicode("Extender", unicode.Extender)
	writeUnicode("Hex_Digit", unicode.Hex_Digit)
	writeUnicode("Hyphen", unicode.Hyphen)
	writeUnicode("IDS_Binary_Operator", unicode.IDS_Binary_Operator)
	writeUnicode("IDS_Trinary_Operator", unicode.IDS_Trinary_Operator)
	writeUnicode("Ideographic", unicode.Ideographic)
	writeUnicode("Join_Control", unicode.Join_Control)
	writeUnicode("Logical_Order_Exception", unicode.Logical_Order_Exception)
	writeUnicode("Noncharacter_Code_Point", unicode.Noncharacter_Code_Point)
	writeUnicode("Other_Alphabetic", unicode.Other_Alphabetic)
	writeUnicode("Other_Default_Ignorable_Code_Point", unicode.Other_Default_Ignorable_Code_Point)
	writeUnicode("Other_Grapheme_Extend", unicode.Other_Grapheme_Extend)
	writeUnicode("Other_ID_Continue", unicode.Other_ID_Continue)
	writeUnicode("Other_ID_Start", unicode.Other_ID_Start)
	writeUnicode("Other_Lowercase", unicode.Other_Lowercase)
	writeUnicode("Other_Math", unicode.Other_Math)
	writeUnicode("Other_Uppercase", unicode.Other_Uppercase)
	writeUnicode("Pattern_Syntax", unicode.Pattern_Syntax)
	writeUnicode("Pattern_White_Space", unicode.Pattern_White_Space)
	writeUnicode("Prepended_Concatenation_Mark", unicode.Prepended_Concatenation_Mark)
	writeUnicode("Quotation_Mark", unicode.Quotation_Mark)
	writeUnicode("Radical", unicode.Radical)
	writeUnicode("Regional_Indicator", unicode.Regional_Indicator)
	writeUnicode("STerm", unicode.STerm)
	writeUnicode("Sentence_Terminal", unicode.Sentence_Terminal)
	writeUnicode("Soft_Dotted", unicode.Soft_Dotted)
	writeUnicode("Terminal_Punctuation", unicode.Terminal_Punctuation)
	writeUnicode("Unified_Ideograph", unicode.Unified_Ideograph)
	writeUnicode("Variation_Selector", unicode.Variation_Selector)
	writeUnicode("White_Space", unicode.White_Space)



	writeUnicode("Adlam", unicode.Adlam)
	writeUnicode("Ahom", unicode.Ahom)
	writeUnicode("Anatolian_Hieroglyphs", unicode.Anatolian_Hieroglyphs)
	writeUnicode("Arabic", unicode.Arabic)
	writeUnicode("Armenian", unicode.Armenian)
	writeUnicode("Avestan", unicode.Avestan)
	writeUnicode("Balinese", unicode.Balinese)
	writeUnicode("Bamum", unicode.Bamum)
	writeUnicode("Bassa_Vah", unicode.Bassa_Vah)
	writeUnicode("Batak", unicode.Batak)
	writeUnicode("Bengali", unicode.Bengali)
	writeUnicode("Bhaiksuki", unicode.Bhaiksuki)
	writeUnicode("Bopomofo", unicode.Bopomofo)
	writeUnicode("Brahmi", unicode.Brahmi)
	writeUnicode("Braille", unicode.Braille)
	writeUnicode("Buginese", unicode.Buginese)
	writeUnicode("Buhid", unicode.Buhid)
	writeUnicode("Canadian_Aboriginal", unicode.Canadian_Aboriginal)
	writeUnicode("Carian", unicode.Carian)
	writeUnicode("Caucasian_Albanian", unicode.Caucasian_Albanian)
	writeUnicode("Chakma", unicode.Chakma)
	writeUnicode("Cham", unicode.Cham)
	writeUnicode("Cherokee", unicode.Cherokee)
	writeUnicode("Chorasmian", unicode.Chorasmian)
	writeUnicode("Common", unicode.Common)
	writeUnicode("Coptic", unicode.Coptic)
	writeUnicode("Cuneiform", unicode.Cuneiform)
	writeUnicode("Cypriot", unicode.Cypriot)
	writeUnicode("Cyrillic", unicode.Cyrillic)
	writeUnicode("Deseret", unicode.Deseret)
	writeUnicode("Devanagari", unicode.Devanagari)
	writeUnicode("Dives_Akuru", unicode.Dives_Akuru)
	writeUnicode("Dogra", unicode.Dogra)
	writeUnicode("Duployan", unicode.Duployan)
	writeUnicode("Egyptian_Hieroglyphs", unicode.Egyptian_Hieroglyphs)
	writeUnicode("Elbasan", unicode.Elbasan)
	writeUnicode("Elymaic", unicode.Elymaic)
	writeUnicode("Ethiopic", unicode.Ethiopic)
	writeUnicode("Georgian", unicode.Georgian)
	writeUnicode("Glagolitic", unicode.Glagolitic)
	writeUnicode("Gothic", unicode.Gothic)
	writeUnicode("Grantha", unicode.Grantha)
	writeUnicode("Greek", unicode.Greek)
	writeUnicode("Gujarati", unicode.Gujarati)
	writeUnicode("Gunjala_Gondi", unicode.Gunjala_Gondi)
	writeUnicode("Gurmukhi", unicode.Gurmukhi)
	writeUnicode("Han", unicode.Han)
	writeUnicode("Hangul", unicode.Hangul)
	writeUnicode("Hanifi_Rohingya", unicode.Hanifi_Rohingya)
	writeUnicode("Hanunoo", unicode.Hanunoo)
	writeUnicode("Hatran", unicode.Hatran)
	writeUnicode("Hebrew", unicode.Hebrew)
	writeUnicode("Hiragana", unicode.Hiragana)
	writeUnicode("Imperial_Aramaic", unicode.Imperial_Aramaic)
	writeUnicode("Inherited", unicode.Inherited)
	writeUnicode("Inscriptional_Pahlavi", unicode.Inscriptional_Pahlavi)
	writeUnicode("Inscriptional_Parthian", unicode.Inscriptional_Parthian)
	writeUnicode("Javanese", unicode.Javanese)
	writeUnicode("Kaithi", unicode.Kaithi)
	writeUnicode("Kannada", unicode.Kannada)
	writeUnicode("Katakana", unicode.Katakana)
	writeUnicode("Kayah_Li", unicode.Kayah_Li)
	writeUnicode("Kharoshthi", unicode.Kharoshthi)
	writeUnicode("Khitan_Small_Script", unicode.Khitan_Small_Script)
	writeUnicode("Khmer", unicode.Khmer)
	writeUnicode("Khojki", unicode.Khojki)
	writeUnicode("Khudawadi", unicode.Khudawadi)
	writeUnicode("Lao", unicode.Lao)
	writeUnicode("Latin", unicode.Latin)
	writeUnicode("Lepcha", unicode.Lepcha)
	writeUnicode("Limbu", unicode.Limbu)
	writeUnicode("Linear_A", unicode.Linear_A)
	writeUnicode("Linear_B", unicode.Linear_B)
	writeUnicode("Lisu", unicode.Lisu)
	writeUnicode("Lycian", unicode.Lycian)
	writeUnicode("Lydian", unicode.Lydian)
	writeUnicode("Mahajani", unicode.Mahajani)
	writeUnicode("Makasar", unicode.Makasar)
	writeUnicode("Malayalam", unicode.Malayalam)
	writeUnicode("Mandaic", unicode.Mandaic)
	writeUnicode("Manichaean", unicode.Manichaean)
	writeUnicode("Marchen", unicode.Marchen)
	writeUnicode("Masaram_Gondi", unicode.Masaram_Gondi)
	writeUnicode("Medefaidrin", unicode.Medefaidrin)
	writeUnicode("Meetei_Mayek", unicode.Meetei_Mayek)
	writeUnicode("Mende_Kikakui", unicode.Mende_Kikakui)
	writeUnicode("Meroitic_Cursive", unicode.Meroitic_Cursive)
	writeUnicode("Meroitic_Hieroglyphs", unicode.Meroitic_Hieroglyphs)
	writeUnicode("Miao", unicode.Miao)
	writeUnicode("Modi", unicode.Modi)
	writeUnicode("Mongolian", unicode.Mongolian)
	writeUnicode("Mro", unicode.Mro)
	writeUnicode("Multani", unicode.Multani)
	writeUnicode("Myanmar", unicode.Myanmar)
	writeUnicode("Nabataean", unicode.Nabataean)
	writeUnicode("Nandinagari", unicode.Nandinagari)
	writeUnicode("New_Tai_Lue", unicode.New_Tai_Lue)
	writeUnicode("Newa", unicode.Newa)
	writeUnicode("Nko", unicode.Nko)
	writeUnicode("Nushu", unicode.Nushu)
	writeUnicode("Nyiakeng_Puachue_Hmong", unicode.Nyiakeng_Puachue_Hmong)
	writeUnicode("Ogham", unicode.Ogham)
	writeUnicode("Ol_Chiki", unicode.Ol_Chiki)
	writeUnicode("Old_Hungarian", unicode.Old_Hungarian)
	writeUnicode("Old_Italic", unicode.Old_Italic)
	writeUnicode("Old_North_Arabian", unicode.Old_North_Arabian)
	writeUnicode("Old_Permic", unicode.Old_Permic)
	writeUnicode("Old_Persian", unicode.Old_Persian)
	writeUnicode("Old_Sogdian", unicode.Old_Sogdian)
	writeUnicode("Old_South_Arabian", unicode.Old_South_Arabian)
	writeUnicode("Old_Turkic", unicode.Old_Turkic)
	writeUnicode("Oriya", unicode.Oriya)
	writeUnicode("Osage", unicode.Osage)
	writeUnicode("Osmanya", unicode.Osmanya)
	writeUnicode("Pahawh_Hmong", unicode.Pahawh_Hmong)
	writeUnicode("Palmyrene", unicode.Palmyrene)
	writeUnicode("Pau_Cin_Hau", unicode.Pau_Cin_Hau)
	writeUnicode("Phags_Pa", unicode.Phags_Pa)
	writeUnicode("Phoenician", unicode.Phoenician)
	writeUnicode("Psalter_Pahlavi", unicode.Psalter_Pahlavi)
	writeUnicode("Rejang", unicode.Rejang)
	writeUnicode("Runic", unicode.Runic)
	writeUnicode("Samaritan", unicode.Samaritan)
	writeUnicode("Saurashtra", unicode.Saurashtra)
	writeUnicode("Sharada", unicode.Sharada)
	writeUnicode("Shavian", unicode.Shavian)
	writeUnicode("Siddham", unicode.Siddham)
	writeUnicode("SignWriting", unicode.SignWriting)
	writeUnicode("Sinhala", unicode.Sinhala)
	writeUnicode("Sogdian", unicode.Sogdian)
	writeUnicode("Sora_Sompeng", unicode.Sora_Sompeng)
	writeUnicode("Soyombo", unicode.Soyombo)
	writeUnicode("Sundanese", unicode.Sundanese)
	writeUnicode("Syloti_Nagri", unicode.Syloti_Nagri)
	writeUnicode("Syriac", unicode.Syriac)
	writeUnicode("Tagalog", unicode.Tagalog)
	writeUnicode("Tagbanwa", unicode.Tagbanwa)
	writeUnicode("Tai_Le", unicode.Tai_Le)
	writeUnicode("Tai_Tham", unicode.Tai_Tham)
	writeUnicode("Tai_Viet", unicode.Tai_Viet)
	writeUnicode("Takri", unicode.Takri)
	writeUnicode("Tamil", unicode.Tamil)
	writeUnicode("Tangut", unicode.Tangut)
	writeUnicode("Telugu", unicode.Telugu)
	writeUnicode("Thaana", unicode.Thaana)
	writeUnicode("Thai", unicode.Thai)
	writeUnicode("Tibetan", unicode.Tibetan)
	writeUnicode("Tifinagh", unicode.Tifinagh)
	writeUnicode("Tirhuta", unicode.Tirhuta)
	writeUnicode("Ugaritic", unicode.Ugaritic)
	writeUnicode("Vai", unicode.Vai)
	writeUnicode("Wancho", unicode.Wancho)
	writeUnicode("Warang_Citi", unicode.Warang_Citi)
	writeUnicode("Yezidi", unicode.Yezidi)
	writeUnicode("Yi", unicode.Yi)
	*/

	//s := "Hello, 世界&¡¿🏳️‍🌈🇩🇪"
	//fmt.Println(runetodir(s))

	h := &MyHandler{}
	fasthttp.ListenAndServe(":81", h.HandleFastHTTP)

}

func writeUnicode(name string, uni *unicode.RangeTable){

	total := 1114111
	var r rune
	var runes []rune
	for i:=0; i < total; i++ {

		r = int32(i)
		if unicode.Is(uni, r) { runes = append(runes, r) }
	
	}
	
	file, _ := json.MarshalIndent(runes, "", " ")
	err := ioutil.WriteFile("lang/"+name, file, 0777)
	if err != nil { fmt.Println(err) }

}

func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {

	total := 122
	output := make([][]byte, total)
	for i:=32; i < total; i++ {
		output[i] = []byte(string(i))
	}

	for _, v := range output {
		fmt.Fprintf(ctx, string(v))
	}
	
}
func runetodir(s string) string {
	var b strings.Builder
	a := []rune(s)
	for _, v := range a {
		x := formatInt32(v)
		fmt.Fprintf(&b, "/%v", x)
	}
	return b.String()
}
func formatInt32(n int32) string {
    return strconv.FormatInt(int64(n), 10)
}
/*
s := "Hello, 世界&¡¿🏳️‍🌈🇩🇪"
a := []rune(s)
output := make([][]byte, len(a))
for i, v := range a {
	output[i] = []byte(string(v))
}
fmt.Println(output)
for _, v1 := range output {
	fmt.Println(v1)
	for _, v2 := range v1 {
		fmt.Println(v2)
	}
}
*/
