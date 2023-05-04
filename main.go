package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/schollz/pianoai/ai2"
	"github.com/schollz/pianoai/player"
	"github.com/urfave/cli"
)

var (
	version  string
	BPM      = 120
	DEBUG    = false
	MANUAL   = false
	JAZZY    = false
	STACCATO = false
	CHORDS   = false
	FOLLOW   = false
)

func main() {
	menu()
	readSettings()
	pianoAI()
}

func copyOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

// Function that runs the Python code from Go and send the output to stdout and error to stderr
// Partially taken from here https://stackoverflow.com/questions/41415337/running-external-python-in-golang-catching-continuous-exec-command-stdout
// Specifically, the answer written by minimijoyo (https://stackoverflow.com/users/7357841/minamijoyo)
func menu() {
	defer fmt.Println("Running PianoAI...")
	// defer wg.Done()

	cmd := exec.Command("python3", "../../tuneHat/main.py")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	go copyOutput(stdout)
	go copyOutput(stderr)

	// err = cmd.Start()
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

func readSettings() {

	readFile, err := os.Open("../../tuneHat/Settings/user_settings.txt")
	defer func(readFile *os.File) {
		err := readFile.Close()
		if err != nil {

		}
	}(readFile)

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	for _, line := range fileLines {
		fmt.Println(line)
	}

	BPM, _ = strconv.Atoi(fileLines[0])
	DEBUG, _ = strconv.ParseBool(fileLines[1])
	MANUAL, _ = strconv.ParseBool(fileLines[2])
	JAZZY, _ = strconv.ParseBool(fileLines[3])
	STACCATO, _ = strconv.ParseBool(fileLines[4])
	CHORDS, _ = strconv.ParseBool(fileLines[5])
	FOLLOW, _ = strconv.ParseBool(fileLines[6])
}

// Also from minimijoyo's stackoverflow answer
func pianoAI() {

	app := cli.NewApp()
	app.Version = version
	app.Compiled = time.Now()
	app.Name = "pianoai"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "bpm",
			Value: 120,
			Usage: "BPM to use",
		},
		cli.IntFlag{
			Name:  "tick",
			Value: 500,
			Usage: "tick frequency in hertz",
		},
		cli.IntFlag{
			Name:  "hp",
			Value: 65,
			Usage: "high pass note threshold to use for leraning",
		},
		cli.IntFlag{
			Name:  "waits",
			Value: 2,
			Usage: "beats of silence before AI jumps in",
		},
		cli.IntFlag{
			Name:  "quantize",
			Value: 64,
			Usage: "1/quantize is shortest possible note",
		},
		cli.StringFlag{
			Name:  "file,f",
			Value: "music_history.json",
			Usage: "file save/load to when pressing bottom C",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "debug mode",
		},
		cli.BoolFlag{
			Name:  "manual",
			Usage: "AI is activated manually",
		},
		cli.IntFlag{
			Name:  "link",
			Value: 3,
			Usage: "AI LinkLength",
		},
		cli.BoolFlag{
			Name:  "jazzy",
			Usage: "AI Jazziness",
		},
		cli.BoolFlag{
			Name:  "stacatto",
			Usage: "AI Stacattoness",
		},
		cli.BoolFlag{
			Name:  "chords",
			Usage: "AI Allow chords",
		},
		cli.BoolFlag{
			Name:  "follow",
			Usage: "AI velocities follow the host",
		},
	}

	app.Action = func(c *cli.Context) (err error) {
		fmt.Println(`
		
		______ _____                   ___  _____ 
		| ___ \_   _|                 / _ \|_   _|
		| |_/ / | |  __ _ _ __   ___ / /_\ \ | |  
		|  __/  | | / _` + "`" + ` | '_ \ / _ \|  _  | | |  
		| |    _| || (_| | | | | (_) | | | |_| |_ 
		\_|    \___/\__,_|_| |_|\___/\_| |_/\___/ 
																							
																							
		 _______________________________________  
		|  | | | |  |  | | | | | |  |  | | | |  | 
		|  | | | |  |  | | | | | |  |  | | | |  | 
		|  | | | |  |  | | | | | |  |  | | | |  | 
		|  |_| |_|  |  |_| |_| |_|  |  |_| |_|  | 
		|   |   |   |   |   |   |   |   |   |   | 
		|   |   |   |   |   |   |   |   |   |   | 
		|___|___|___|___|___|___|___|___|___|___| 

	 Lets play some music!
											`)
		p, err := player.New(BPM, c.GlobalInt("tick"), DEBUG)
		if err != nil {
			return
		}
		p.HighPassFilter = c.GlobalInt("hp")
		p.AI = ai2.New(p.TicksPerBeat)
		p.AI.HighPassFilter = c.GlobalInt("hp")
		p.AI.LinkLength = c.GlobalInt("link")
		p.AI.Jazzy = JAZZY
		p.AI.Stacatto = STACCATO
		p.AI.DisallowChords = !CHORDS
		p.ManualAI = MANUAL
		p.UseHostVelocity = FOLLOW

		p.Start()

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Print(err)
	}
}
