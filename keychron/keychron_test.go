package keychron

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"
)

var (
	sample1 []byte
	sample2 []byte
)

func setup() error {
	f, err := os.Open("testdata/sample.html")
	if err != nil {
		fmt.Println()
		return err
	}

	sample1, err = ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	f1, err := os.Open("testdata/sample1.html")
	if err != nil {
		return err
	}

	sample2, err = ioutil.ReadAll(f1)
	if err != nil {
		return err
	}
	return nil
}

func TestVerifyAlert(t *testing.T) {
	err := setup()
	if err != nil {
		panic(err)
	}

	messages := make([]string, 0)
	messageLock := sync.Mutex{}
	Alert = func(m string) error {
		messageLock.Lock()
		defer messageLock.Unlock()
		messages = append(messages, m)
		return nil
	}

	startTime := time.Now()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if startTime.Add(1 * time.Second).Before(time.Now()) {
			fmt.Fprintln(w, string(sample2))
			return
		}
		fmt.Fprintln(w, string(sample1))
	}))
	defer ts.Close()

	allKeyboards = []string{ts.URL}

	delay = 2 * time.Second
	go Start()
	time.Sleep(5 * time.Second)

	messageLock.Lock()
	defer messageLock.Unlock()
	if len(messages) != 1 {
		t.Fatal("No alerts found")
	}
}
