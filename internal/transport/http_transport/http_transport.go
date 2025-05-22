package http_transport

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yash3004/config_server/configurations"
	"github.com/yash3004/config_server/users"
)

// Server implements the HTTP server
type Server struct {
	userManager   *users.UserManager
	configManager *configurations.ConfigManager
	router        *mux.Router
}

// NewServer creates a new HTTP server
func NewServer(userManager *users.UserManager, configManager *configurations.ConfigManager) *Server {
	s := &Server{
		userManager:   userManager,
		configManager: configManager,
		router:        mux.NewRouter(),
	}
	s.setupRoutes()
	return s
}

// setupRoutes configures the HTTP routes
func (s *Server) setupRoutes() {
	// Configuration routes
	s.router.HandleFunc("/config", s.addConfig).Methods("POST")
	s.router.HandleFunc("/config", s.updateConfig).Methods("PUT")
	s.router.HandleFunc("/config", s.deleteConfig).Methods("DELETE")
	s.router.HandleFunc("/config", s.getConfig).Methods("GET")

	// User routes
	s.router.HandleFunc("/user", s.addUser).Methods("POST")
	s.router.HandleFunc("/user", s.updateUser).Methods("PUT")
	s.router.HandleFunc("/user/{userID}", s.deleteUser).Methods("DELETE")
}

// StartHTTPServer starts the HTTP server
func (s *Server) StartHTTPServer(address string) error {
	return http.ListenAndServe(address, s.router)
}

// Request and response types
type AuthRequest struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

type ConfigRequest struct {
	AuthRequest
	Filename string `json:"filename"`
	FileType int    `json:"file_type"`
	Data     []byte `json:"data"`
}

type UserRequest struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type ConfigResponse struct {
	UserID   string `json:"user_id"`
	Filename string `json:"filename"`
	FileType int    `json:"file_type"`
	Data     []byte `json:"data"`
}

// addConfig handles POST /config
func (s *Server) addConfig(w http.ResponseWriter, r *http.Request) {
	var req ConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Authenticate user
	authenticated, err := s.userManager.AuthenticateUser(r.Context(), req.UserID, req.Password)
	if err != nil || !authenticated {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Add config
	err = s.configManager.AddConfig(r.Context(), req.UserID, req.Filename, req.FileType, req.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// updateConfig handles PUT /config
func (s *Server) updateConfig(w http.ResponseWriter, r *http.Request) {
	var req ConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Authenticate user
	authenticated, err := s.userManager.AuthenticateUser(r.Context(), req.UserID, req.Password)
	if err != nil || !authenticated {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Update config
	err = s.configManager.UpdateConfig(r.Context(), req.UserID, req.Filename, req.FileType, req.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// deleteConfig handles DELETE /config
func (s *Server) deleteConfig(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}

	// Authenticate user
	authenticated, err := s.userManager.AuthenticateUser(r.Context(), req.UserID, req.Password)
	if err != nil || !authenticated {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Delete config
	err = s.configManager.DeleteConfig(r.Context(), req.UserID, filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// getConfig handles GET /config
func (s *Server) getConfig(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	password := r.URL.Query().Get("password")
	filename := r.URL.Query().Get("filename")

	if userID == "" || password == "" || filename == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Authenticate user
	authenticated, err := s.userManager.AuthenticateUser(r.Context(), userID, password)
	if err != nil || !authenticated {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Get config
	data, fileType, err := s.configManager.GetConfig(r.Context(), userID, filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ConfigResponse{
		UserID:   userID,
		Filename: filename,
		FileType: fileType,
		Data:     data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// addUser handles POST /user
func (s *Server) addUser(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := s.userManager.AddUser(r.Context(), req.UserID, req.Email, req.Name, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// updateUser handles PUT /user
func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := s.userManager.UpdateUser(r.Context(), req.UserID, req.Email, req.Name, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// deleteUser handles DELETE /user/{userID}
func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	err := s.userManager.DeleteUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}