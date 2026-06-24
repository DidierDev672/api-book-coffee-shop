package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Equipment struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	Status          string `json:"status"`
	LastMaintenance string `json:"last_maintenance"`
}

type store struct {
	mu    sync.RWMutex
	items map[string]Equipment
}

func newStore() *store {
	return &store{items: make(map[string]Equipment)}
}

func (s *store) create(e Equipment) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[e.ID] = e
}

func (s *store) get(id string) (Equipment, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.items[id]
	return e, ok
}

func (s *store) getAll() []Equipment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]Equipment, 0, len(s.items))
	for _, e := range s.items {
		list = append(list, e)
	}
	return list
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	http.Error(w, `{"error":"`+msg+`"}`, status)
}

func equipmentHandler(s *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/equipment")
		id = strings.TrimPrefix(id, "/")
		id = strings.TrimSpace(id)

		switch r.Method {
		case http.MethodPost:
			createEquipment(w, r, s)
		case http.MethodGet:
			if id != "" {
				getEquipment(w, r, s, id)
			} else {
				listEquipment(w, r, s)
			}
		case http.MethodPut:
			if id == "" {
				writeError(w, "id is required in URL", http.StatusBadRequest)
				return
			}
			updateEquipment(w, r, s, id)
		case http.MethodDelete:
			if id == "" {
				writeError(w, "id is required in URL", http.StatusBadRequest)
				return
			}
			deleteEquipment(w, r, s, id)
		default:
			writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func createEquipment(w http.ResponseWriter, r *http.Request, s *store) {
	var e Equipment
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(e.Name) == "" {
		writeError(w, "name is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(e.Type) == "" {
		writeError(w, "type is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(e.Status) == "" {
		writeError(w, "status is required", http.StatusBadRequest)
		return
	}

	if e.ID == "" {
		e.ID = time.Now().Format("20060102150405.000000000")
	}
	if e.LastMaintenance == "" {
		e.LastMaintenance = time.Now().Format("2006-01-02")
	}

	s.create(e)
	writeJSON(w, http.StatusCreated, e)
}

func getEquipment(w http.ResponseWriter, _ *http.Request, s *store, id string) {
	e, ok := s.get(id)
	if !ok {
		writeError(w, "equipment not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, e)
}

func listEquipment(w http.ResponseWriter, _ *http.Request, s *store) {
	writeJSON(w, http.StatusOK, s.getAll())
}

func updateEquipment(w http.ResponseWriter, r *http.Request, s *store, id string) {
	existing, ok := s.get(id)
	if !ok {
		writeError(w, "equipment not found", http.StatusNotFound)
		return
	}

	var incoming Equipment
	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(incoming.Name) != "" {
		existing.Name = incoming.Name
	}
	if strings.TrimSpace(incoming.Type) != "" {
		existing.Type = incoming.Type
	}
	if strings.TrimSpace(incoming.Status) != "" {
		existing.Status = incoming.Status
	}
	if strings.TrimSpace(incoming.LastMaintenance) != "" {
		existing.LastMaintenance = incoming.LastMaintenance
	}

	s.create(existing)
	writeJSON(w, http.StatusOK, existing)
}

func deleteEquipment(w http.ResponseWriter, _ *http.Request, s *store, id string) {
	if _, ok := s.get(id); !ok {
		writeError(w, "equipment not found", http.StatusNotFound)
		return
	}
	s.mu.Lock()
	delete(s.items, id)
	s.mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func main() {
	s := newStore()

	http.HandleFunc("/equipment", equipmentHandler(s))
	http.HandleFunc("/equipment/", equipmentHandler(s))

	log.Println("Server listening at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
