package domain

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Algorithm uint8

const (
	ECC Algorithm = iota + 1
	RSA
)

var algorithmName = map[uint8]string{
	1: "ECC",
	2: "RSA",
}

var algorithmValue = map[string]uint8{
	"ecc": 1,
	"rsa": 2,
}

func (algorithm Algorithm) String() string {
	return algorithmName[uint8(algorithm)]
}

func ParseAlgorithm(s string) (Algorithm, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	value, found := algorithmValue[s]
	if !found {
		return Algorithm(0), fmt.Errorf("%q is not a valid algorithm", s)
	}
	return Algorithm(value), nil
}

func (algorithm Algorithm) MarshalJSON() ([]byte, error) {
	return json.Marshal(algorithm.String())
}

func (algorithm *Algorithm) UnmarshalJSON(data []byte) (err error) {
	var stringValue string
	if err = json.Unmarshal(data, &stringValue); err != nil {
		return err
	}
	if *algorithm, err = ParseAlgorithm(stringValue); err != nil {
		return err
	}
	return nil
}
