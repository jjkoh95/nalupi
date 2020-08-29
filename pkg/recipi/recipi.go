package recipi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
)

// GetLastComputedPI retrieves last computed PI from permanent storage
func GetLastComputedPI() (precision int64, pi *big.Int, err error) {
	url := "https://nalusi-b235sdkoha-de.a.run.app/nalupi?spreadsheetID=1YnXZwX5ABPmBUFhktGVLDVnmgluVgSMFjIkMyIJ8Lt0&a1Range=Data!A2:B2"
	response, err := http.Get(url)
	if err != nil {
		return 0, nil, err
	}
	var respBody [][]string
	body, err := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return 0, nil, err
	}
	precision, err = strconv.ParseInt(respBody[0][0], 10, 64)
	if err != nil {
		return 0, nil, errors.New("Unable to convert string to big.Int")
	}
	var ok bool
	pi, ok = big.NewInt(0).SetString(respBody[0][1], 10)
	if !ok {
		return 0, nil, errors.New("Unable to convert string to big.Int")
	}
	return precision, pi, nil
}

// SaveComputedPI append an entry to permanent storage of PI
func SaveComputedPI(precision, pi string) error {
	client := &http.Client{}
	reqBody := [][]string{{precision, pi}}
	url := "https://nalusi-b235sdkoha-de.a.run.app/nalupi?spreadsheetID=1YnXZwX5ABPmBUFhktGVLDVnmgluVgSMFjIkMyIJ8Lt0&a1Range=Data!A2:B2"
	b, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
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
func GetSnapshot() (k, Lk, Xk, Mk, Kk *big.Int, err error) {
	url := "https://nalusi-b235sdkoha-de.a.run.app/nalupi?spreadsheetID=1FMUFV2z_MaccKswNLh3-x2vDeBY3RRNNzzAusjh848c&a1Range=Data!A:E"
	response, err := http.Get(url)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	var respBody [][]string
	body, err := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	var ok bool
	k, ok = big.NewInt(0).SetString(respBody[1][0], 10)
	if !ok {
		return nil, nil, nil, nil, nil, errors.New("Unable to convert string to big.Int")
	}
	Lk, ok = big.NewInt(0).SetString(respBody[1][1], 10)
	if !ok {
		return nil, nil, nil, nil, nil, errors.New("Unable to convert string to big.Int")
	}
	Xk, ok = big.NewInt(0).SetString(respBody[1][2], 10)
	if !ok {
		return nil, nil, nil, nil, nil, errors.New("Unable to convert string to big.Int")
	}
	Mk, ok = big.NewInt(0).SetString(respBody[1][3], 10)
	if !ok {
		return nil, nil, nil, nil, nil, errors.New("Unable to convert string to big.Int")
	}
	Kk, ok = big.NewInt(0).SetString(respBody[1][4], 10)
	if !ok {
		return nil, nil, nil, nil, nil, errors.New("Unable to convert string to big.Int")
	}
	return k, Lk, Xk, Mk, Kk, nil
}

// SaveSnapshot saves snapshot - temporary metadata required to calculate PI
func SaveSnapshot(k, Lk, Xk, Kk, Mk string) error {
	// we need to wrap PUT request our own
	client := &http.Client{}

	reqBody := [][]string{{k, Lk, Xk, Kk, Mk}}
	url := "https://nalusi-b235sdkoha-de.a.run.app/nalupi?spreadsheetID=1FMUFV2z_MaccKswNLh3-x2vDeBY3RRNNzzAusjh848c&a1Range=Data!A2:E2"
	b, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	if response.StatusCode == 200 {
		return nil
	}
	return errors.New("Server error")
}

// SaveFractionMeta saves fraction data (Mk*Lk) and Xk
func SaveFractionMeta(numerator, denominator string) error {
	reqBody := [][]string{{numerator, denominator}}
	url := "https://nalusi-b235sdkoha-de.a.run.app/nalupi?spreadsheetID=1EUh-4d-Xx1YQmfgdNaNsIWE5B1Q8s6WTm50UlUtHfIQ&a1Range=Data"
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

// LoadFractionMeta sums up the fraction data (Mk*Lk)/Xk with precision as defined
func LoadFractionMeta(precision int64) (*big.Int, error) {
	url := "https://nalusi-b235sdkoha-de.a.run.app/nalupi?spreadsheetID=1EUh-4d-Xx1YQmfgdNaNsIWE5B1Q8s6WTm50UlUtHfIQ&a1Range=Data!A:B"
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	var respBody [][]string
	body, err := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return nil, err
	}

	pow := big.NewInt(0)
	pow.Exp(big.NewInt(10), big.NewInt(precision), nil)
	res := big.NewInt(0)
	for i := 1; i < len(respBody); i++ {
		numerator, ok := big.NewInt(0).SetString(respBody[i][0], 10)
		if !ok {
			return nil, errors.New("Unable to convert string to big.Int")
		}
		denominator, ok := big.NewInt(0).SetString(respBody[i][1], 10)
		if !ok {
			return nil, errors.New("Unable to convert string to big.Int")
		}
		temp := big.NewInt(0)
		temp.Mul(numerator, pow)
		temp.Quo(temp, denominator)
		res.Add(res, temp)
	}
	return res, nil
}
