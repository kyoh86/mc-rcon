package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kyoh86/mcrcon"
	"github.com/morikuni/aec"
	"golang.org/x/term"
)

var (
	flags struct {
		command  string
		host     string
		port     int
		password bool
	}
)

func init() {
	flag.StringVar(&flags.host, "host", "localhost", "Minecraft server host")
	flag.IntVar(&flags.port, "port", 25575, "Minecraft server port")
	flag.BoolVar(&flags.password, "password", false, "Set if you want to input Minecraft server RCON password as secret")
	flag.StringVar(&flags.command, "command", "", "Command line to execute")
	flag.Parse()
}

func promptCommand() (string, error) {
	fmt.Print(aec.Bold.String(), aec.CyanF.String(), "> ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	fmt.Print(aec.Reset)
	return input, nil
}

func main() {
	conn := new(mcrcon.MCConn)
	password := os.Getenv("MCRCON_PASSWORD")
	if flags.password {
		fmt.Print("Enter Minecraft server RCON password: ")
		buf, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatalln("Failed to read password", err)
		}
		password = string(buf)
		fmt.Println()
	}
	err := conn.Open(fmt.Sprintf("%s:%d", flags.host, flags.port), password)
	if err != nil {
		log.Fatalln("Open failed", err)
	}
	defer conn.Close()

	err = conn.Authenticate()
	if err != nil {
		log.Fatalln("Auth failed", err)
	}

	if flags.command == "" {
		for {
			cmd, err := promptCommand()
			if err != nil {
				log.Fatalln("Failed to read command: ", err)
			}
			words := strings.Split(strings.TrimSpace(cmd), " ")
			if len(words) > 0 && words[0] == "exit" {
				return
			}
			resp, err := conn.SendCommand(cmd)
			if err != nil {
				log.Fatalln("Command failed", err)
			}
			fmt.Println(resp)
			fmt.Println()
		}
	} else {
		resp, err := conn.SendCommand(flags.command)
		if err != nil {
			log.Fatalln("Command failed", err)
		}
		fmt.Println(resp)
	}
}
