package logs

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aaryansinhaa/panes/utils/services/storage"
)

func ListLogHandler(w http.ResponseWriter, r *http.Request, s *storage.SQLite) {
	limit, err := strconv.Atoi(r.PathValue("limit"))
	if err != nil || limit <= 0 {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}
	logEntries, err := s.GetLogEntries(limit)
	if err != nil {
		http.Error(w, "Error retrieving log entries", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(logEntries)
	if err != nil {
		http.Error(w, "Error marshalling log entries", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}

func DeleteLogEntryHandler(w http.ResponseWriter, r *http.Request, s *storage.SQLite) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "No log entry ID provided", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid log entry ID", http.StatusBadRequest)
		return
	}

	err = s.DeleteLogEntry(id)
	if err != nil {
		http.Error(w, "Error deleting log entry", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteAllLogEntriesHandler(w http.ResponseWriter, r *http.Request, s *storage.SQLite) {
	rowsAffected, err := s.DeleteAllLogEntries()
	if err != nil {
		http.Error(w, "Error deleting all log entries", http.StatusInternalServerError)
		return
	}

	response := map[string]int64{"rowsAffected": rowsAffected}
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
