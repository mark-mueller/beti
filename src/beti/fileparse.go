

package beti

import (
	"strings"
)




//   Read comments and function definitions. Return array of paragraphs.
//
func FileParse () (paragraphs []string) {

    // initialize
	paragraphs = []string{""}
	ndx := 0
	block_comment_on := false
	p_flag := false
	beti_on := true
    func_comment_on := false
    func_comments := []string{""}

    //
    // 1   Comments, package and func
    // 2   Comments and function
    // 3   :beti
    Mode = determineParseMode()

	// loop for each line
	for _, s := range data_ar {

		s = strings.TrimRight(s, " \t\r")
		s = expand_tabs(s)
        var m []string

		if block_comment_on || Is_BetiFile {

    		if !Is_BetiFile {
    			m := compiledRegex["endComment"].FindStringSubmatch(s)
    			if m != nil {
    				if len(paragraphs[ndx]) > 0 {
    					paragraphs[ndx] += "\n"
    				}
    				paragraphs[ndx] += m[1]
    				paragraphs = append(paragraphs, "")
    				ndx += 1
    				block_comment_on = false
    				continue
    			}
    			m = compiledRegex["leadingSplats"].FindStringSubmatch(s)
    		    if m != nil {
                    s = m[1]
                }
            }

			m = compiledRegex["directive"].FindStringSubmatch(s)
			if m != nil {
                s = strings.ToLower(s)
                if m[2] == "" {
                    m[2] = "on"
                    //s = "$"+m[1]+" on"
                }
                if m[1] == "beti" {
                    if m[2] == "on" {
                        beti_on = true
                        s = m[3]
                    }
                    if m[2] == "off" {
                        beti_on = false
                        continue
                    }
                }
                if !beti_on {
                    continue
                }
				paragraphs = append(paragraphs, s)
				ndx += 1
				p_flag = false
                continue
			}
            if !beti_on {
                continue
            }

			if s == "" {
				if p_flag {
					continue
				}
				paragraphs = append(paragraphs, "")
				ndx += 1
				p_flag = true
				continue
			}

			p_flag = false
			if len(paragraphs[ndx]) > 0 {
				paragraphs[ndx] += "\n"
			}
			paragraphs[ndx] += s

		} else {

           if Is_GoFile {
    			m = compiledRegex["func"].FindStringSubmatch(s)
    			if m != nil {
                    fc := ""
                    if func_comment_on {
                        fc = strings.Join(func_comments, "\n")
                        func_comments = []string{""}
                        func_comment_on = false
                    }
    				paragraphs = append(paragraphs, m[1])
    				paragraphs = append(paragraphs, fc)
    				ndx += 1
    				continue
    			}
    			m = compiledRegex["funcComment"].FindStringSubmatch(s)
    			if m != nil {
                    if !func_comment_on {
                        func_comments = []string{""}
                    }
                    func_comment_on = true
    				s = strings.TrimRight(m[1], " \t")
                    func_comments = append(func_comments, s)
                    continue
    			}
                func_comment_on = false
            }

			m = compiledRegex["beginComment"].FindStringSubmatch(s)
			if m != nil {
				block_comment_on = true
                func_comment_on = false
				continue
			}

		}
	}
	return
}

