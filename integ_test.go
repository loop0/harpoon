package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// switch to testdata dir to access dummy config.toml
	os.Chdir("./testdata")

	// set up verbose flags for testing
	os.Args = append(os.Args, "-test.v")

	go main()

	// make sure the server is up before running tests
	tries := 10
	for i := 0; i < tries; i++ {
		if conn, err := net.Dial("tcp", "localhost:9001"); err == nil {
			conn.Close()
			fmt.Println("Server up! Waiting another second.")
			<-time.After(1 * time.Second)
			break
		}
		if i < tries-1 {
			fmt.Printf("Server not up yet. Try %d\n", i)
		} else {
			fmt.Printf("Server not up after %d tries. Failing tests.\n", tries)
			os.Exit(-1)
		}
		<-time.After(1 * time.Second)
	}

	os.Exit(m.Run())
}

func TestPushEventExplicit(t *testing.T) {
	body, err := os.Open("./push_develop.json")
	if err != nil {
		t.Fatalf("Failed to open payload file. Error: %v", err)
	}

	client := http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:9001", body)
	req.Header.Add("X-GitHub-Event", "push")

	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error posting to server. Error: %v", err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Bad response from server. Response code: %v", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Failed to read in response body. Error: %v", err)
		t.FailNow()
	}
	res.Body.Close()

	t.Logf("Data returned: %v", data)
}

func TestPushEventAll(t *testing.T) {
	body, err := os.Open("./push_master.json")
	if err != nil {
		t.Fatalf("Failed to open payload file. Error: %v", err)
	}

	client := http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:9001", body)
	req.Header.Add("X-GitHub-Event", "push")

	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error posting to server. Error: %v", err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Bad response from server. Response code: %v", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Failed to read in response body. Error: %v", err)
		t.FailNow()
	}
	res.Body.Close()

	t.Logf("Data returned: %v", data)
}

func TestPushEventBadRepo(t *testing.T) {
	body, err := os.Open("./push_differentrepo.json")
	if err != nil {
		t.Fatalf("Failed to open payload file. Error: %v", err)
	}

	client := http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:9001", body)
	req.Header.Add("X-GitHub-Event", "push")

	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error posting to server. Error: %v", err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Bad response from server. Response code: %v", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Failed to read in response body. Error: %v", err)
		t.FailNow()
	}
	res.Body.Close()

	t.Logf("Data returned: %v", data)
}
