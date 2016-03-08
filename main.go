package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

const (
	headerEvent     = "X-GitHub-Event"  // HTTP header where the webhook event is stored
	headerSignature = "X-Hub-Signature" // HTTP header where the sha1 signature of the payload is stored
)

var (
	config            = tomlConfig{}                          // the program config
	verbose           = false                                 // weither we should log the output of the command
	gitHubSecretToken = os.Getenv("GITHUB_HOOK_SECRET_TOKEN") // the webhook secret token, used to verify signature
)

// HookHandler receive hooks from GitHub.
func HookHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// read the HTTP request body
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, color.RedString("Error: "+err.Error()))
		BadRequestHandler(w, r)
		return
	}

	// validate signature
	if gitHubSecretToken != "" {
		sign := r.Header.Get(headerSignature)

		// to compute the HMAC in order to check for equality with what has been sent by GitHub
		mac := hmac.New(sha1.New, []byte(gitHubSecretToken))
		mac.Write(payload)
		expectedHash := hex.EncodeToString(mac.Sum(nil))
		receivedHash := sign[5:] // remove 'sha1='

		// signature mismatch, do not process
		if !hmac.Equal([]byte(receivedHash), []byte(expectedHash)) {
			color.Set(color.FgRed)
			fmt.Fprintf(os.Stderr, "Mismatch between expected (%s) and received (%s) hash.", expectedHash, receivedHash)
			color.Set(color.Faint)
			BadRequestHandler(w, r)
			return
		}
	}

	var eventPayload HookWithRepository
	json.Unmarshal(payload, &eventPayload)

	// verify that this is an event that we should process
	event := r.Header.Get(headerEvent)
	if event == "ping" {
		return // always respond 200 to pings
	}

	// check weither we're interested in that event for that ref
	if _, ok := config.Events[event+":"+eventPayload.Repository.FullName+":"+eventPayload.Ref]; !ok {
		if verbose {
			color.Set(color.FgRed)
			fmt.Fprintf(os.Stderr, "Discarding %s on %s with ref %s.\n",
				color.CyanString(event), color.YellowString(eventPayload.Repository.FullName), color.YellowString(eventPayload.Ref))
			color.Set(color.Faint)
			BadRequestHandler(w, r)
			return // 400 Bad Request
		}
	}

	handleEvent(event, eventPayload, []byte(payload))
}

// handleEvent handles any event.
func handleEvent(event string, hook HookWithRepository, payload []byte) {
	// show related commits if push event
	if event == "push" {
		var pushEvent HookPush
		json.Unmarshal(payload, &pushEvent)
		fmt.Println(event, "detected on", color.YellowString(hook.Repository.FullName),
			"with ref", color.YellowString(hook.Ref), "with the following commits:")
		for _, commit := range pushEvent.Commits {
			fmt.Printf("\t%s - %s by %s\n", commit.Timestamp, color.CyanString(commit.Message), color.BlueString(commit.Author.Name))
		}
	}

	// prepare the command
	eventKey := event + ":" + hook.Repository.FullName + ":" + hook.Ref
	cmd := exec.Command(config.Events[eventKey].Cmd,
		strings.Split(config.Events[eventKey].Args, " ")...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		return
	}

	// in case of -verbose we log the output of the executed command
	if verbose {
		scanner := bufio.NewScanner(cmdReader)
		go func() {
			for scanner.Scan() {
				color.White("> " + scanner.Text() + "\n")
			}
		}()
	}

	// launch it
	err = cmd.Start()
	if err != nil {
		color.Set(color.FgRed)
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		color.Set(color.Faint)
		return
	}
}

// BadRequestHandler handles bad requests. Status 400 and JSON error message.
func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(`{"message": "I don't know what you're talking about"}`))
}

func main() {
	verbosePtr := flag.Bool("v", false, "Weither we output stuff.")
	flag.Parse()
	verbose = *verbosePtr

	// load the config.toml
	config = loadConfig()
	addr := config.Addr + ":" + strconv.Itoa(config.Port)
	color.White(`    __                                     
   / /_  ____ __________  ____  ____  ____ 
  / __ \/ __ ` + "`" + `/ ___/ __ \/ __ \/ __ \/ __ \
 / / / / /_/ / /  / /_/ / /_/ / /_/ / / / /
/_/ /_/\__,_/_/  / .___/\____/\____/_/ /_/ 
                /_/                        
`)
	color.White("\tListening on " + addr)

	// router & server
	r := mux.NewRouter()
	r.HandleFunc("/", HookHandler).Methods("POST")
	http.ListenAndServe(addr, r)
}
