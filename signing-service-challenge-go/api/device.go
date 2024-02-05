package api

import (
	"encoding/json"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"io"
	"net/http"
)

// Devices handles api/v0/devices route
func (s *Server) Devices(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		s.getAllSignatureDevices(response, request)
	case "POST":
		s.createSignatureDevice(response, request)
	default:
		WriteErrorResponse(response, 404, []string{"not found"})
	}
}

// Device handles api/v0/device/{id} route
func (s *Server) Device(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		s.getAllSignatureDevices(response, request)
	default:
		WriteErrorResponse(response, 404, []string{"not found"})
	}
}

// DeviceSign handles api/v0/device/{id}/sign route
func (s *Server) DeviceSign(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		s.signDataWithDevice(response, request)
	default:
		WriteErrorResponse(response, 404, []string{"not found"})
	}
}

func (s *Server) getAllSignatureDevices(response http.ResponseWriter, request *http.Request) {
	devices := s.devicesRepository.GetAll()
	WriteAPIResponse(response, 200, devices)
}

type createSignatureDeviceParams struct {
	algorithm domain.Algorithm
	label     string
}

func (s *Server) createSignatureDevice(response http.ResponseWriter, request *http.Request) {
	var params createSignatureDeviceParams
	read, _ := io.ReadAll(request.Body)
	err := json.Unmarshal(read, &params)
	if err != nil {
		WriteErrorResponse(response, 400, []string{err.Error()})
		return
	}

	device, err := domain.CreateSignatureDevice(params.algorithm, params.label, s.devicesRepository)
	if err != nil {
		WriteErrorResponse(response, 400, []string{err.Error()})
		return
	}

	WriteAPIResponse(response, 200, device)
}

func (s *Server) getSignatureDevice(response http.ResponseWriter, request *http.Request) {
	id := request.URL.Path
	device, found := s.devicesRepository.Get(id)
	if !found {
		WriteErrorResponse(response, 404, []string{"not found"})
		return
	}

	WriteAPIResponse(response, 200, device)
}

type signDataWithDeviceParams struct {
	data string
}

func (s *Server) signDataWithDevice(response http.ResponseWriter, request *http.Request) {
	id := request.URL.Path
	var params signDataWithDeviceParams
	read, _ := io.ReadAll(request.Body)
	err := json.Unmarshal(read, &params)
	signedData, err := domain.SignTransaction(id, params.data, s.devicesRepository)
	if err != nil {
		WriteErrorResponse(response, 400, []string{err.Error()})
		return
	}

	WriteAPIResponse(response, 200, signedData)
}
