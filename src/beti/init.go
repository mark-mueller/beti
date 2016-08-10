package beti

import (
	"regexp"
)

const VERSION           = "0.0.0"
const DEBUG             = true


//  Mode is used to indicate the parsing mode.
//
// 0    default
// 1    .go
// 2    .beti, .txt
//
var Mode                int
const MODE_DEFAULT =    0
const MODE_GO =         1
const MODE_TXT =        2

type assoc_array        map[string]string

var Parameters          assoc_array
var Flags               assoc_array

var (
    data_ar             []string
    document            assoc_array
    index_arr           []index_struct
    paragraphRegexKeys  []string
    headingRegexKeys    []string
    compiledRegex       map[string]*regexp.Regexp
)

type index_struct struct {
	Text                string
	Name                string
    Level               string
}

type row_               []string
var  index_table        [][]string




///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///

func init () {

    Mode = 0

    Flags = assoc_array{
        "stylesheet":       "",
        "title":            "",
        "index":            "",
    }

    Parameters = assoc_array{
        "beti":             "on",
        "stylesheet":       "beti.css",
        "title":            "",
        "index":            "default",
        "nheadings":        "",
    }

    document = assoc_array{
        "title":        "",
        "head":         "",
        "pageheading":  "",
        "index":        "",
        "body":         "",
    }

	re := assoc_array{
		"beginComment":     "/\\*(.*)",
		"endComment":       "(.*)\\*/",
		"leadingSplats":    "\\*{1,} *(.*)",
		"lineComment":      "//(.*)",
		"pComment":         "#(.*)",
		"funcComment":      "^// *(.*)",
		"func":             "^(func *([A-Z][\\w_]+) *([^{]+))",
		"gofunc":           "^(func *([a-z][\\w_]+) *([^{]+))",
		"nfunction":         "^(function *([\\w_]+) *([^{]*))",
		"directive":        "^:(\\w+) ?(.*)([\\n \\S]*)",
        "package":          "^package +\"(\\w+)\"",
        "beti":             "^:beti *(\\w+)",

		"h1":               "^\\. *([\\S ]{1,60})([\\n \\S]*)",                    // Non-numeric heading. Paragraph ok.
		"nheading":         "^\\d{1,2}\\. *([\\S ]{1,60})([\\n \\S]*)",            // Numeric with decimal only. Paragraph ok.
		"h2":               "^\\d{0,2}\\.\\d{1,3} *([\\S ]{1,60})([\\n \\S]*)",    // Numeric, decimal, numeric. Paragraph ok.
		"h3a":              "^([\\d\\w ]{1,40}:)()$",                              // 40 chars and a colon. Paragraph ok.
		"h3b":              "^([\\d\\w ]{1,40}:)\\n([\\n \\S]+)$",                 // 40 chars and a colon. Paragraph ok.

		"example1":         "^(?i)(examples?:?)()$",
		"example2":         "^(?i)(examples?:?)\\n([\\n \\S]+)$",
		"table":            "^(?i)(table ?\\w{0,3}:)\\n([\\n \\S]+)$",

        // paragraphRegexKeys
		"bulletList":       "^() *([-\\*] +[\\n \\S]*)",
		"orderedList":      "^() *(\\d\\. +[\\n \\S]*)",
		"bulletItem":       "^ *([-\\*\\d])\\.? +(.*)",
		"indent":           "^()(  +\\S[\\n\\s\\S]+)",
		"oneDent":          "^()( +\\S[\\n\\s\\S]+)",


        "expandTabs":       "^([^\t]*)(\t+)",
		"leadingSpaces":    "^( +)\\S",
		"tableColumn":      "  \\b",
    	"textSpan":         "(.*)(\\b[ihub]{1,3})\\|([^\\|.]+)\\|(.*)",
    	"fullURL":          "(.*)(https*://[\\S]+)(.*)",
    	"textURL":          "(.*)\\|(.*)\\|(https*://[\\S]+)(.*)",
    	"textRelativeURL":  "(.*)\\|(.*)\\|(\\./[\\S]+)( *.*)",
    	"pageAnchor":       "(#[^\\d][-\\w]+)",
    	"emailAddress":     "([^\\d]+@[-\\w\\.]+\\.[a-z]{2,15})",
		"firstWordChar":    "\\b[a-z]",
        "escapedChar":      "\\\\(.)",
        "hexifiedChars":    "\\[0x([0-9a-fA-F]{2,})\\]",
        "NotAnchorNameChars": "[^a-z0-9]+",
        "LeadingDecimalHeading": "^([\\d\\.]+) +(.+)",

    }
	compiledRegex = make(map[string]*regexp.Regexp)
	for k, r := range re {
		compiledRegex[k] = regexp.MustCompile(r)
	}

	paragraphRegexKeys = []string{
        "bulletList",
        "orderedList",
        "bulletItem",
//        "orderedItem",
        "indent",
        "oneDent",
	}

    headingRegexKeys = []string{
        "example2",
        "example1",
        "table",
        "h3a",
        "h2",
        "h1",
        "nheading",
    }
}

//eof//