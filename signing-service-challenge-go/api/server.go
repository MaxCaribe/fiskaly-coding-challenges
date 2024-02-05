package api

import (
	"encoding/json"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/gorilla/mux"
	"net/http"
)

// Response is the generic API response container.
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse is the generic error API response container.
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress     string
	devicesRepository domain.DevicesRepository
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string, devicesRepository domain.DevicesRepository) *Server {
	return &Server{
		listenAddress:     listenAddress,
		devicesRepository: devicesRepository,
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	router := mux.NewRouter()

	router.Handle("/api/v0/health", http.HandlerFunc(s.Health))
	router.Handle("/api/v0/devices/{uuid}/sign", http.HandlerFunc(s.DeviceSign))
	router.Handle("/api/v0/devices/{uuid}", http.HandlerFunc(s.Device))
	router.Handle("/api/v0/devices", http.HandlerFunc(s.Devices))

	return http.ListenAndServe(s.listenAddress, router)
}

// WriteInternalError writes a default internal error message as an HTTP response.
func WriteInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, code int, errors []string) {
	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Errors: errors,
	}

	bytes, err := json.Marshal(errorResponse)
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteAPIResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)

	response := Response{
		Data: data,
	}

	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}
