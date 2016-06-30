package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-tty"

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
	text, err = transEnc(text, encode)
	text = strings.Replace(text, "\r\n", "\n", -1)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	if text == "" {
		text = usage
	}
	//init tty
	tty, err := tty.Open()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	defer tty.Close()
	out := colorable.NewColorable(tty.Output())

	//term init
	err = termbox.Init()
	if err != nil {
		panic(err)
	}

	defer func() {
		clearTerm(out, tty)
		termbox.Close()
	}()

	if interval == 0 {
		keyEvent(out, tty, text, loop, gravity)
	} else {
		autoPlay(out, tty, interval, text, loop, gravity)
	}
}

// listen key event
func keyEvent(out io.Writer, tty *tty.TTY, text string, loop bool, gravity bool) {
	rep := regexp.MustCompile(`(?m:\n^---\n)`)
	strs := rep.Split(text, -1)
	i := 0
	draw(out, tty, strs[i], gravity)
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				clearTerm(out, tty)
				return
			case termbox.KeySpace, termbox.KeyArrowRight, termbox.KeyEnter:
				i += 1
				if i >= len(strs) {
					if loop {
						i = 0
					} else {
						return
					}
				}
			case termbox.KeyArrowLeft, termbox.KeyBackspace, termbox.KeyDelete, termbox.KeyBackspace2:
				i -= 1
				if i < 0 {
					if loop {
						i = len(strs) - 1
					} else {
						return
					}
				}
			default:
			}
		default:
		}
		draw(out, tty, strs[i], gravity)
	}
}

//auto
func autoPlay(out io.Writer, tty *tty.TTY, interval int, text string, loop bool, gravity bool) {
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
					clearTerm(out, tty)
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
				clearTerm(out, tty)
				return
			}
		}
		draw(out, tty, strs[i], gravity)
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}

//draw
func draw(out io.Writer, tty *tty.TTY, str string, gravity bool) {
	clearTerm(out, tty)
	strs := strings.Split(str, "\n")
	clr := "\x1b[2K"
	_, heigt, err := tty.Size()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	if gravity == false {
		for i := 0; i < len(strs) && i < heigt; i++ {
			if i != heigt-1 {
				out.Write([]byte(clr + strs[i] + "\r\n"))
			} else {
				out.Write([]byte(clr + strs[i]))
			}
		}
	} else {
		if len(strs) < heigt {
			strlen := len(strs)
			for i := 0; i < (heigt - strlen); i++ {
				strs = append([]string{" "}, strs...)
			}
		}
		for i := 0; i < len(strs) && i < heigt; i++ {
			if i != heigt-1 {
				out.Write([]byte(clr + strs[i] + "\r\n"))
			} else {
				out.Write([]byte(clr + strs[i]))
			}
		}
	}
}

func clearTerm(out io.Writer, tty *tty.TTY) {
	out.Write([]byte(fmt.Sprintf("\x1b[2J")))
	out.Write([]byte(fmt.Sprintf("\x1b[0;0H")))

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
                                                            .__________________
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
