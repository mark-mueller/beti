/*
:beti on

*/

package beti

import (
    "fmt"
    "strings"
    "strconv"
)

var    Version             string
var    Copyright           string

///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///


// Output document to stdout using hard-coded template.
//
//
func Output () {

    make_index()

    switch Parameters["index"] {
        case "top":
            document["body"] = document["index"] + document["body"]
        case "none":
        default:
            document["body"] = strings.Replace(document["body"], "[_INDEX_DEFAULT_]", document["index"], 1)
    }
    document["body"] = strings.Replace(document["body"], "[_INDEX_DEFAULT_]", "", 1)
    document["stylesheet"] = Parameters["stylesheet"]
    document["title"] = compiledRegex["LeadingDecimalHeading"].ReplaceAllString(Parameters["title"], "$2")
    document["title"] = capitalize_string(document["title"])

    fmt.Print ("<!DOCTYPE html>\n<HTML id=\"_top\">\n<HEAD>\n")
    fmt.Printf("<TITLE>%v</TITLE>\n", document["title"])
    fmt.Print ("<META charset=\"utf-8\" />\n")
    fmt.Printf("<META name=\"generator\" content=\"%v\" />\n", Version)
    fmt.Print ("<!-- - - - - - - - - - - - - - - - - - - - - - - - - - -\n")
    fmt.Printf("\n%v\n%v\n\n", Version, Copyright)
    fmt.Print ("- - - - - - - - - - - - - - - - - - - - - - - - - - - -->\n")
    fmt.Printf("<LINK type=\"text/css\" rel=\"stylesheet\" href=\"%v\" />\n", document["stylesheet"])
    fmt.Print (document["head"])
    fmt.Print ("</HEAD>\n<BODY>\n")
    fmt.Print (document["pageheading"])
    fmt.Print (document["body"])

    // add a footer later
    //fmt.Print (document["footer"])

    fmt.Print ("<a id=\"topbuttton\" href=\"#_top\" title=\"Top of document\"></a>\n")
    fmt.Print ("</BODY>\n</HTML>")
}

func make_index () {
    table_guts := "\n"
    rowno := 1

    for _, item := range index_arr {
        text := strings.Trim(item.Text, ":")
        num  := ""
        m := compiledRegex["LeadingDecimalHeading"].FindStringSubmatch(text)
        if m != nil {
            num  = m[1]
            text = m[2]
        }
        text = capitalize_string(text)
        text  = Make_link("#"+item.Name, text, "link-level-"+item.Level)
        index_item := Make_tag("div", num, "cell index_num")
        index_item += Make_tag("div", text, "cell index_text")
        table_guts += Make_tag("div", index_item, "row rowno_"+strconv.Itoa(rowno)) + "\n"
        rowno++
    }

    document["index"] = Make_a_heading("1", "Index", "") +
                        Make_tag("div", table_guts, "index table") + "\n"
}


func AppendBody (p string) {
    if p != "" {
        document["body"] += p + "\n"
    }
}


// Create table HTML from an array. Return html.
//
func Make_table (rawtable string, class string) (html string) {
    html = ""
    class = strings.Trim(class, ":")
    ndx1 := []int{0}
    ok := false
    even := true
    arr := strings.Split(rawtable, "\n")
    // find all column index positions in first row
    m := compiledRegex["tableColumn"].FindAllStringIndex(arr[0], -1)
    for n := range m {
        ok = true
        ndx1 = append(ndx1, m[n][0]+2)
    }
    if ok {
        rowno := 0;
        // shift ndx1 into ndx2
        ndx2 := append(ndx1[1:], 0)
        // iterate for each table row
        for _, s := range arr {
            colno := 0;
            row := ""
            ndx2[len(ndx2)-1] = len(s)
            // iterate for each cell
            for n, i := range ndx1 {
                r := ndx2[n]
                str := strings.Trim(s[i:r], " ")
                row += "  " + Make_tag("div", str, "cell colno_"+strconv.Itoa(colno)) + "\n"
                colno++
            }
            html += Make_tag("div", "\n"+row, "row rowno_"+strconv.Itoa(rowno)) + "\n"
            rowno++
            if even {

            }
        }

        if class != "" {
            class += " "
        }
        class += "table"
        html = Make_tag("div", "\n"+html, class) + "\n"
    }
    return
}

// Make a bullet list from an array. Return html.
//
//
func Make_bullet_list (s string, tag string) (html string) {
    html = ""
    _ = tag
    arr := strings.Split(s, "\n")
    for _, s = range arr {
        m := compiledRegex["bulletItem"].FindStringSubmatch(s)
        if m != nil {
            html += " " + Make_tag("li", m[2], "") + "\n"
        } else {
            if s != "" {
                s = strings.Trim(s, " ")
                html += "     " + s + "\n"
            }
        }
    }
    html = Make_tag(tag, "\n"+html, "") + "\n"
    return
}

// Create an html tag. Return html.
//
func Make_tag (tag, text, class string) string {
    attr := ""
    if class != "" {
        attr = " class=\"" + class + "\""
    }
    tag = strings.ToUpper(tag)
    return fmt.Sprintf("<%v%v>%v</%v>", tag, attr, text, tag)
}

// Create a link tag from href and text strings. Return html.
//
func Make_link (href, text, class string) string {
    attr := ""
    attr += " href=\"" + href + "\""
    if class != "" {
        attr += " class=\"" + class + "\""
    }
    return fmt.Sprintf("<A%v>%v</A>", attr, text)
}


// Create a paragraph
//
func Make_paragraph (content string, class string) string {
    if len(content) == 0 {
        return ""
    }
    return Make_tag("p", content, class) + "\n"
}

//  Create an anchored heading
func Make_a_heading (level, heading, section_heading string) (html string) {
    if section_heading != "" {
        section_heading += "-"
    }
    name := Make_anchor_name(section_heading+heading)
    index_arr = append(index_arr, index_struct{Text:heading, Name:name, Level:level})
    //html = "\n" + Make_tag("h"+level, capitalize_string(heading), "") + "\n"
    html = fmt.Sprintf("<H%v id=\"%v\">%v</H%v>", level, name, capitalize_string(heading), level)
    return
}

//  Create a heading
func Make_heading (level string, heading string) (html string) {
    html = Make_tag("h"+level, capitalize_string(heading), "") + "\n"
    return
}

// Create an anchor name
// Remove leading numbers
//
func Make_anchor_name (heading string) (name string) {
    name = strings.ToLower(heading)
    name = compiledRegex["NotAnchorNameChars"].ReplaceAllLiteralString(name, " ")
    name = compiledRegex["LeadingDecimalHeading"].ReplaceAllString(name, "$2")
    name = strings.Trim(name, ": ")
    name = strings.Replace(name, " ", "-", -1)
    return
}

//eof//