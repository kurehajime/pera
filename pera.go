package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

//main
func main() {
	var encode string
	var loop bool
	var gravity bool
	var interval int
	var default_encoding string
	var text string = ""
	var err error = nil

	//get flags
	if runtime.GOOS == "windows" {
		default_encoding = "sjis"
	} else {
		default_encoding = "utf-8"
	}
	flag.StringVar(&encode, "e", default_encoding, "encoding")
	flag.BoolVar(&loop, "l", false, "loop (Exit:Press Esc Key)")
	flag.BoolVar(&gravity, "g", false, "gravity(bottom align)")
	flag.IntVar(&interval, "i", 0, "interval <= 10 millisecond (enable auto mode)")
	flag.Usage = func() {
		fmt.Printf("Usage:\n\t$pera [options]  filepath \n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	// get str
	if len(flag.Args()) == 0 {
		text, err = readPipe()
	} else if flag.Arg(0) == "-" {
		text, err = readStdin()
	} else {
		text, err = readFileByArg(flag.Arg(0))
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	if text == "" {
		text = usage
	}

	//term init
	err = termbox.Init()
	if err != nil {
		panic(err)
	}

	defer termbox.Close()
	if interval == 0 {
		keyEvent(text, loop, gravity)
	} else {
		autoPlay(interval, text, loop, gravity)
	}
}

// listen key event
func keyEvent(text string, loop bool, gravity bool) {
	rep := regexp.MustCompile(`(?m:\n^---\n)`)
	strs := rep.Split(text, -1)
	i := 0
	draw(strs[i], gravity)
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				clear()
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				return
			case termbox.KeySpace, termbox.KeyArrowRight, termbox.KeyEnter:
				i += 1
				if i >= len(strs) {
					if loop {
						i = 0
					} else {
						clear()
						termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
						return
					}
				}
			case termbox.KeyArrowLeft, termbox.KeyBackspace, termbox.KeyDelete, termbox.KeyBackspace2:
				i -= 1
				if i < 0 {
					if loop {
						i = len(strs) - 1
					} else {
						clear()
						termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
						return
					}
				}
			default:
			}
		default:
		}
		draw(strs[i], gravity)
	}
}

//auto
func autoPlay(interval int, text string, loop bool, gravity bool) {
	rep := regexp.MustCompile(`(?m:\n^---\n)`)
	strs := rep.Split(text, -1)
	stop := false
	if interval < 10 {
		interval = 10
	}
	go func() {
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc, termbox.KeyCtrlC:
					stop = true
					clear()
					termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
					return
				default:
				}
			default:
			}
		}

	}()
	for i := 0; stop == false; i++ {
		if i >= len(strs) {
			if loop == true {
				i = 0
			} else {
				clear()
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				return
			}
		}
		draw(strs[i], gravity)
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}

//draw
func draw(str string, gravity bool) {
	strs := strings.Split(str, "\n")
	clear()
	if gravity == false {
		for i := 0; i < len(strs); i++ {
			drawLine(i, strs[i])
		}
	} else {
		_, heigt := termbox.Size()
		for i, r := heigt-1, len(strs)-1; i >= 0 && r >= 0; i, r = i-1, r-1 {
			drawLine(i, strs[r])
		}
	}

	termbox.Flush()
}

//clear terminal
func clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for r := 0; r < 40; r++ {
		for c := 0; c < 80; c++ {
			termbox.SetCell(c, r, []rune(" ")[0], termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	termbox.Flush()
}

//draw line
func drawLine(line int, str string) {
	runes := []rune(str)
	b := 0
	for i := 0; i < len(runes); i++ {
		termbox.SetCell(b, line, runes[i], termbox.ColorDefault, termbox.ColorDefault)
		b += runewidth.RuneWidth(runes[i])
	}
}

var usage = `

    * pera command *





          -> Please press right allow key.
---


    OK! 

    You can back by pressing left allow key .



            
---


    Basic usage :

        $ pera <filepath>

   that's easy :) 

　　　
---


　　　Q. Where to page breaks ?


　　　A. ---



---

　　　Example:

         first page.
　　　    
         ---         <--this
                        
         second page.
         It will be a new page in three hyphens(---) on beginning of line.
---


     * Useful Options *






---


     -e

         Set encoding. (default "utf-8")




---


     -g  

         Gravity option. Set vertical-align to bottom. 
                                                            .____________________
                                                               \o/ 
                                                                |   Oh,
                                                               /)   Help me!
---


     -l

         Loop option. 
　　　　　 It ends when you press the Esc key .

         

---


     -i <number>

         Interval option. 
　　　　　 It will be auto-play when you specify this option.

         

---





               THE
              o/ 
             <|  
          D   /)   EN
---





               THE EN
                   o/ 
                  /|  
          D        /) 
---





               THE END 
                    \o/
                     |  
                     /)  `
