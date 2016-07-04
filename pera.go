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
	var defaultEncoding string
	var text string
	var err error

	//get flags
	if runtime.GOOS == "windows" {
		defaultEncoding = "sjis"
	} else {
		defaultEncoding = "utf-8"
	}
	flag.StringVar(&encode, "e", defaultEncoding, "encoding")
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
		clearTerm(out)
		termbox.Close()
	}()
	width, heigt, err := tty.Size()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	if interval == 0 {
		KeyEvent(out, width, heigt, text, loop, gravity)
	} else {
		AutoPlay(out, width, heigt, interval, text, loop, gravity)
	}
	clearTerm(out)
}

// KeyEvent :listen key event
func KeyEvent(out io.Writer, width int, heigt int, text string, loop bool, gravity bool) {
	rep := regexp.MustCompile(`(?m:\n^---\n)`)
	strs := rep.Split(text, -1)
	i := 0

	Draw(out, width, heigt, strs[i], gravity)
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				return
			case termbox.KeySpace, termbox.KeyArrowRight, termbox.KeyEnter:
				i++
				if i >= len(strs) {
					if loop {
						i = 0
					} else {
						return
					}
				}
			case termbox.KeyArrowLeft, termbox.KeyBackspace, termbox.KeyDelete, termbox.KeyBackspace2:
				i--
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
		Draw(out, width, heigt, strs[i], gravity)
	}
}

//AutoPlay :auto
func AutoPlay(out io.Writer, width int, heigt int, interval int, text string, loop bool, gravity bool) {
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
				return
			}
		}
		Draw(out, width, heigt, strs[i], gravity)
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}

//Draw :draw
func Draw(out io.Writer, width int, heigt int, str string, gravity bool) {
	clearTerm(out)
	strs := strings.Split(str, "\n")
	clr := "\x1b[2K"

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

func clearTerm(out io.Writer) {
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
