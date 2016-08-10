

package beti

import (
	"strings"
	"io/ioutil"
)



///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///



//  Examine data to determine what kind of script we have.
//
func determineParseMode () (mode int) {

	// loop for each line
    // Modes:
    // 1   Comments, package and func
    // 2   Comments and function
    // 3   :beti
    //

    item := map[string]bool{
        "lineComment":  false,
        "pComment":     false,
        "beginComment": false,
        "endComment":   false,
        "package":      false,
        "nfunction":    false,
        "gofunc":       false,
        "beti":         false,
    }

    // look for items
    for _, s := range data_ar {
        for key, _ := range item {
            m := compiledRegex[key].FindStringSubmatch(s)
            if m != nil {
                item[key] = true
            }
        }
    }

    // set the mood
    mode = 0
    if (item["lineComment"] || (item["beginComment"] && item["endComment"])) && item["package"] && item["gofunc"] {
        mode = 1
    }
    if (item["lineComment"] || item["pComment"] || (item["beginComment"] && item["endComment"])) && item["nfunction"] {
        mode = 2
    }
    if (item["lineComment"] || item["pComment"] || (item["beginComment"] && item["endComment"])) && item["nfunction"] {
        mode = 3
    }

    return
}

///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///

func Parse_paragraph (s string) (ptype, head, paragraph string) {
	ptype = ""
	head = ""
	paragraph = s
    var m []string

    NOTALOOP:
    for true {

    	for _, key := range headingRegexKeys {
    		m = compiledRegex[key].FindStringSubmatch(s)
    		if m != nil {
        		ptype = key
        		head = m[1]
                if len(m) > 1 {
        		    paragraph = m[2]
        	    }
        		break NOTALOOP
    		}
    	}

    	for _, key := range paragraphRegexKeys {
    		m = compiledRegex[key].FindStringSubmatch(s)
    		if m != nil {
        		ptype = key
        		paragraph = m[2]
        		break NOTALOOP
    		}
    	}

    	m = compiledRegex["directive"].FindStringSubmatch(s)
    	if m != nil {
            ptype = "directive"
            value := "on"
            key := m[1]
    		if m[2] != "" {
                value = strings.Trim(m[2], " ")
    		}
            key = strings.ToLower(key)
            Parameters[key] = strings.ToLower(value)
      		break NOTALOOP
        }

        if Is_GoFile {
        	m = compiledRegex["func"].FindStringSubmatch(s)
        	if m != nil {
        		ptype = "func"
        		head = m[2]
        		paragraph = m[3]
            }
     		break NOTALOOP
        }
   		break NOTALOOP
    }

    paragraph = trimParagraphLeadingSpaces(paragraph)
	arr := strings.Split(paragraph, "\n")
    for i, s := range arr {
        s = hexifyEscapedChars(s)
		s = parse_spans(s)
		s = parse_text_url(s)
		s = parse_url(s)
		s = parse_all_page_anchors(s)
		s = parse_all_email_addresses(s)
		s = parse_text_relative_url(s)
		s = unHexifyChars(s)
        arr[i] = s
    }
	paragraph = strings.Join(arr, "\n")

	return
}




func parse_spans (s string) string {
	m := compiledRegex["textSpan"].FindStringSubmatch(s)
	if m != nil {
        class := m[2]
        if len(class) > 1 {
            arr := strings.Split(class, "")
            class = strings.Join(arr, " ")
        }
        st := m[1] + Make_tag("span", m[3], class) + m[4]
		return st
	}
	return s
}

func parse_text_url (s string) string {
	m := compiledRegex["textURL"].FindStringSubmatch(s)
	if m != nil {
		// hexify the protocol so the url will not be picked up again
		href := strings.Replace(m[3], "http", "[0x68747470]", -1)
        st := m[1] + Make_link(href, m[2], "") + m[4]
		return st
	}
	return s
}

func parse_url (s string) string {
	m := compiledRegex["fullURL"].FindStringSubmatch(s)
	if m != nil {
		// need to hexify the protocol so the url will not be picked up again
		href := strings.Replace(m[2], "http", "[0x68747470]", -1)
        st := m[1] + Make_link(href, href, "") + m[3]
		return st
	}
	return s
}

func parse_text_relative_url (s string) string {
	m := compiledRegex["textRelativeURL"].FindStringSubmatch(s)
	if m != nil {
        href := m[3]
        text := m[2]
        st := m[1] + Make_link(href, text, "") + m[4]
		return st
	}
	return s
}


func parse_all_page_anchors (s string) string {
    for true {
        m := compiledRegex["pageAnchor"].FindStringIndex(s)
    	if m == nil {
            break
        }
        p1 := m[0]
        p2 := m[1]
        text := strings.Replace(s[p1:p2], "#", "", 1)
        // substitute the symbol
        href := strings.Replace(s[p1:p2], "#", "[0x23]", 1)
        text  = strings.Replace(capitalize_string(text), "-", " ", -1)
        s = s[0:p1] + Make_link(href, text, "") + s[p2:]
    }
    return s
}

func parse_all_email_addresses (s string) string {
    for true {
        m := compiledRegex["emailAddress"].FindStringIndex(s)
    	if m == nil {
            break
        }
        p1 := m[0]
        p2 := m[1]
        // substitute the @ symbol
        href := strings.Replace(s[p1:p2], "@", "[0x40]", 1)
        s = s[0:p1] + Make_link("mailto:"+href, href, "") + s[p2:]
    }
    return s
}


func ReadFileGuts (filename string) (ok bool) {
    ok = false
	data, err := ioutil.ReadFile(filename)
	if err != nil {
        return
	}
    ok = true
	// split file guts into an array
	data_ar = strings.Split(string(data), "\n")
    return
}

//eof//