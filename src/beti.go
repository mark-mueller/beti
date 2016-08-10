/*
*/

package main

import (
	"strconv"
	"./beti"
//	"os"
//	"fmt"
)

const (
    NAME          = "Beti"
    VERSION       = "0.0.0"             // Major.Minor.Patch
    BUILD_VERSION = "x.0.1"             // BuildStatus.BuildMajor.BuildMinor
    BUILD_NO      = "90"                // BuildNumber (x.0.1.90)
    BUILD_DATE    = "2015-12-16"
    AUTHOR        = "Mark K Mueller, mark@markmueller.com"
    COPYRIGHT     = "Copyright 2015 Mark K Mueller. All rights reserved.\nUse of this source code is governed by a BSD-style\nlicense that can be found in the LICENSE file."
)



///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///  ///

func main () {

    beti.Version = NAME+" v"+VERSION+"-"+BUILD_VERSION+"."+BUILD_NO
    beti.Copyright = COPYRIGHT

	hcount := 0
    section_heading := ""
    levels := []int{0, 0, 0, 0, 0, 0, 0}

    filename, exists := beti.GetFilename()
    if filename == "" {
        println("Usage: Beti [filename]")
        return
    }
    if !exists {
        println("File does not exist")
        return
    }
	// read file guts
    var ok bool
	ok = beti.ReadFileGuts(filename)
	if !ok {
        println("Error: Cannot read file")
        return
	}

    // Check file extension
    // Later I may need different parsing routines depending on the file type
    // For now, I am only checking for Go source files. All others are treated as plain text.
    switch beti.FileExtension(filename) {
        case "go":
            beti.Is_GoFile = true
        //case "c":
        //    beti.Is_CFile = true
    }

    // parse text into an array of paragraphs
	ar_paragraphs := beti.FileParse()

    // loop through each paragraph
	for _, str := range ar_paragraphs {

		ptype, heading, paragraph := beti.Parse_paragraph(str)

        if ptype == "directive" {
            continue
        }

		switch ptype {

            case "nheading":
                if( beti.Parameters["nheadings"] == "" ){
                    beti.Parameters["nheadings"] = "on"
                }
                fallthrough
            case "h1":
                if( beti.Parameters["nheadings"] == "" ){
                    beti.Parameters["nheadings"] = "off"
                }
                levels[1]++
                levels[2] = 0
                if( beti.Parameters["nheadings"] == "on" ){
                    heading = strconv.Itoa(levels[1]) +". "+ heading
                }
				paragraph = beti.Make_a_heading("1", heading, "") + beti.Make_paragraph(paragraph, "")
				hcount += 1
                section_heading = heading
                if levels[1] == 1 && beti.Parameters["title"] == "" {
                    beti.Parameters["title"] = heading
                }
                if levels[1] == 2 {
                    // Insert index token before the 2nd heading
                    paragraph = "[_INDEX_DEFAULT_]\n" + paragraph
                }

            case "h2":
                levels[2]++
                if( beti.Parameters["nheadings"] == "on" ){
                    heading = strconv.Itoa(levels[1]) +"."+ strconv.Itoa(levels[2]) +" "+ heading
	            }
				paragraph = beti.Make_a_heading("2", heading, "") + beti.Make_paragraph(paragraph, "")
				hcount += 1
                section_heading = heading

            case "h3a", "h3b":
				paragraph = beti.Make_a_heading("3", heading, section_heading) + beti.Make_paragraph(paragraph, "")
//				paragraph = beti.Make_heading("3", heading) + beti.Make_paragraph(paragraph, "")
                //section_heading = heading

            case "example1", "example2":
				paragraph = beti.Make_a_heading("3", heading, section_heading) + beti.Make_paragraph(paragraph, "pre example")

			case "bulletList":
				paragraph = beti.Make_bullet_list(paragraph, "ul")

			case "orderedList":
				paragraph = beti.Make_bullet_list(paragraph, "ol")

			case "indent":
				paragraph = beti.Make_paragraph(paragraph, "pre")

			case "oneDent":
				paragraph = beti.Make_paragraph(paragraph, "wrap")

			case "table":
				paragraph = beti.Make_a_heading("3", heading, section_heading) +"\n"+ beti.Make_table(paragraph, "")

			case "func":
                if !beti.Is_GoFile {
                    continue
                }
			    paragraph = beti.Make_a_heading("2", heading, section_heading) + beti.Make_paragraph(heading+" "+paragraph, "pre")

			default:
				paragraph = beti.Make_paragraph(paragraph, ptype)

		}

        // Append the paragraph
        if paragraph != "" {
			beti.AppendBody(paragraph)
        }

	}

    beti.SetParameters()

    beti.Output()

}


//eof//