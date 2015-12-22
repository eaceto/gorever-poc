package gorever

import (
	"fmt"
	"github.com/inconshreveable/go-update"
	"net/http"
	"time"
)

type Updater struct {
	newVersion bool
	ch         chan bool
}

func NewUpdater(ch chan bool) (*Updater, error) {
	u := &Updater{
		newVersion: true,
		ch:         ch,
	}

	// TEST
	go func(updater *Updater) {
		time.Sleep(6*time.Second)
		updater.ch <- true
	}(u)

	return u, nil
}

func (u *Updater) HasNewVersion() bool {
	return u.newVersion
}

func (u *Updater) Update() error {
	u.newVersion = false
	// get url of new file
	url := "http://localhost:8078/gorever-poc" //TEST: do http-serve -p 8078
	return u.doUpdate(url)
}

func (u *Updater) doUpdate(url string) error {
	// request the new file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		if rerr := update.RollbackError(err); rerr != nil {
			return fmt.Errorf("Failed to rollback from bad update: %v", rerr)
		}
	}
	return err
}
