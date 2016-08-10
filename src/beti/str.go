package beti

import (
   	"os"
    "flag"
    "fmt"
	"math"
	"strings"
	"encoding/hex"
)

var Is_GoFile = false



///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///


func GetFilename () (filename string, exists bool) {
    filename = CmdLineParse()
    exists = false
    if filename != "" {
        // check if file exists
        if _, err := os.Stat(filename); err != nil {
            return
        }
        exists = true
    }
    return
}


func CmdLineParse () (file string) {
    file = ""
    var stylesheet, title, index string
    flag.StringVar(&stylesheet, "stylesheet", "", "Filename of stylesheet")
    flag.StringVar(&title,      "title",      "", "Title of document")
    flag.StringVar(&index,      "index",      "", "top|bottom|none")

    flag.Parse()
    if len(flag.Args()) > 0 {
        file = flag.Args()[0]
    }
    if stylesheet != "" {
        Flags["stylesheet"] = stylesheet
    }
    if title != "" {
        Flags["title"] = title
    }
    if index != "" {
        Flags["index"] = index
    }
    return
}

func SetParameters () {
    if Flags["stylesheet"] != "" {
        Parameters["stylesheet"] = Flags["stylesheet"]
    }
    if Flags["title"] != "" {
        Parameters["title"] = Flags["title"]
    }
    if Flags["index"] != "" {
        Parameters["index"] = Flags["index"]
    }
}

func FileExtension (filename string) (ext string) {
    ext = ""
    i := strings.LastIndex(filename, ".")
    if i > -1 {
        ext = filename[i+1:]
    }
    return
}


func hexifyEscapedChars (s string) string {
    for true {
    	m := compiledRegex["escapedChar"].FindStringSubmatch(s)
	    if m == nil {
            break
        }
        hex := fmt.Sprintf("%x", m[1])
        s = strings.Replace(s, m[0], "[0x"+ hex +"]", -1)
    }
    return s
}

func unHexifyChars (s string) string {
    for true {
    	m := compiledRegex["hexifiedChars"].FindStringSubmatch(s)
	    if m == nil {
            break
        }
        chr, err := hex.DecodeString(m[1])
	    if err != nil {
            println("Error:DecodeString,", err)
            break
        }
        s = strings.Replace(s, m[0], string(chr), -1)
    }
    return s
}

func capitalize_string (s string) string {
    s  = strings.ToLower(s)
    m := compiledRegex["firstWordChar"].FindAllStringIndex(s, -1)
	if m != nil {
        for _, n := range m {
            p1 := n[0]
            p2 := p1+1
            u := strings.ToUpper(s[p1:p2])
            s = s[0:p1] + u + s[p2:]
        }
    }
    return s
}


func expand_tabs (s string) string {
	tab_stop := float64(4)
	for 1 == 1 {
		var ln1 float64
		var ln2 float64
		m := compiledRegex["expandTabs"].FindStringSubmatch(s)
		if m == nil {
			break
		}
		ln1 = float64(len(m[1]))
		ln2 = float64(len(m[2]))
		mo := math.Mod(ln1, tab_stop)
		x := ln2*tab_stop - mo
		s = strings.Replace(s, m[0], m[1]+strings.Repeat(" ", int(x)), 1)
	}
	return s
}

func trimParagraphLeadingSpaces (p string) string {
	m := compiledRegex["leadingSpaces"].FindStringSubmatch(p)
    if m != nil {
        n := len(m[1])
        p = p[n:]
        sp := m[1]
        p = strings.Replace(p, "\n"+sp, "\n", -1)
    }
    return p
}

//eof//