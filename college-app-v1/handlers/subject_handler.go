// handlers/subject_handler.go
package handlers

import (
	"college-app-v1/models"
	"college-app-v1/services"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// SubjectHandler gerencia as requisições HTTP para matérias.
type SubjectHandler struct {
	service *services.SubjectService
}

// NewSubjectHandler cria uma nova instância de SubjectHandler.
func NewSubjectHandler(s *services.SubjectService) *SubjectHandler {
	return &SubjectHandler{service: s}
}

// CreateSubjectHandler lida com a criação de uma nova matéria.
// POST /subjects
func (h *SubjectHandler) CreateSubjectHandler(w http.ResponseWriter, r *http.Request) {
	var subject models.Subject
	if err := json.NewDecoder(r.Body).Decode(&subject); err != nil {
		http.Error(w, "Requisição inválida: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateSubject(&subject); err != nil {
		log.Printf("Erro ao criar matéria no serviço: %v", err)
		http.Error(w, "Erro ao criar matéria: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subject)
}

// GetSubjectByIDHandler lida com a busca de uma matéria por ID.
// GET /subjects/{id}
func (h *SubjectHandler) GetSubjectByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	subject, err := h.service.GetSubjectByID(id)
	if err != nil {
		if err.Error() == "matéria não encontrada" { // Erro personalizado do serviço
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Erro ao buscar matéria no serviço: %v", err)
		http.Error(w, "Erro ao buscar matéria: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subject)
}

// GetAllSubjectsHandler lida com a busca de todas as matérias.
// GET /subjects
func (h *SubjectHandler) GetAllSubjectsHandler(w http.ResponseWriter, r *http.Request) {
	subjects, err := h.service.GetAllSubjects()
	if err != nil {
		log.Printf("Erro ao buscar todas as matérias no serviço: %v", err)
		http.Error(w, "Erro ao buscar matérias: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subjects)
}

// UpdateSubjectHandler lida com a atualização de uma matéria existente.
// PUT /subjects/{id}
func (h *SubjectHandler) UpdateSubjectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var subject models.Subject
	if err := json.NewDecoder(r.Body).Decode(&subject); err != nil {
		http.Error(w, "Requisição inválida: "+err.Error(), http.StatusBadRequest)
		return
	}

	subject.ID = id // Garante que o ID da URL seja usado

	if err := h.service.UpdateSubject(&subject); err != nil {
		if err.Error() == "matéria não encontrada para atualização" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Erro ao atualizar matéria no serviço: %v", err)
		http.Error(w, "Erro ao atualizar matéria: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(subject)
}

// DeleteSubjectHandler lida com a exclusão de uma matéria por ID.
// DELETE /subjects/{id}
func (h *SubjectHandler) DeleteSubjectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.DeleteSubject(id); err != nil {
		if err.Error() == "matéria não encontrada para exclusão" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Erro ao deletar matéria no serviço: %v", err)
		http.Error(w, "Erro ao deletar matéria: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent) // 204 No Content para deleção bem-sucedida
}
