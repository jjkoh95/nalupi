package recipi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// GetLastComputedPI retrieves last computed PI from permanent storage
func GetLastComputedPI() {

}

// SaveComputedPI append an entry to permanent storage of PI
func SaveComputedPI(precision, pi string) error {
	reqBody := [][]string{{precision, pi}}
	url := "https://nalusi-b235sdkoha-de.a.run.app/nalupi?spreadsheetID=1YnXZwX5ABPmBUFhktGVLDVnmgluVgSMFjIkMyIJ8Lt0&a1Range=Data"
	b, _ := json.Marshal(reqBody)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	if response.StatusCode == 200 {
		return nil
	}
	return errors.New("Server error")
}

// GetSnapshot retrieves last snapshot of PI computation
// to proceed further in case there's any failure of server
func GetSnapshot() {}

// SaveSnapshot saves snapshot - temporary metadata required to calculate PI
func SaveSnapshot(k, Lk, Xk, Mk string) error {
	reqBody := [][]string{{k, Lk, Xk, Mk}}
	url := "https://nalusi-b235sdkoha-de.a.run.app/nalupi?spreadsheetID=1FMUFV2z_MaccKswNLh3-x2vDeBY3RRNNzzAusjh848c&a1Range=Data"
	b, _ := json.Marshal(reqBody)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	if response.StatusCode == 200 {
		return nil
	}
	return errors.New("Server error")
}
