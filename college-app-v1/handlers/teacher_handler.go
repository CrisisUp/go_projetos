package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings" // Para verificar erros específicos

	"college-app-v1/models"   // Importa o modelo de professor
	"college-app-v1/services" // Importa o serviço de professor

	"github.com/gorilla/mux"
)

// TeacherHandler gerencia requisições HTTP para professores.
type TeacherHandler struct {
	teacherService *services.TeacherService // Renomeei 'service' para 'teacherService' para clareza
}

// NewTeacherHandler cria uma nova instância de TeacherHandler.
func NewTeacherHandler(service *services.TeacherService) *TeacherHandler {
	return &TeacherHandler{teacherService: service}
}

// CreateTeacherHandler cria um novo professor via HTTP POST.
// Rota: /teachers (POST)
func (h *TeacherHandler) CreateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var teacher models.Teacher
	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		http.Error(w, "Requisição inválida: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.teacherService.CreateTeacher(&teacher); err != nil {
		if strings.Contains(err.Error(), "obrigatórios") || strings.Contains(err.Error(), "já existe") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Erro interno ao criar professor: %v", err)
		http.Error(w, "Erro interno do servidor ao criar professor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(teacher)
}

// GetTeacherByIDHandler busca um professor pelo ID via HTTP GET.
// Rota: /teachers/{id} (GET)
func (h *TeacherHandler) GetTeacherByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	teacher, err := h.teacherService.GetTeacherByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Erro interno ao buscar professor por ID %s: %v", id, err)
		http.Error(w, "Erro interno do servidor ao buscar professor", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(teacher)
}

// GetAllTeachersHandler busca todos os professores via HTTP GET.
// Rota: /teachers (GET)
func (h *TeacherHandler) GetAllTeachersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	teachers, err := h.teacherService.GetAllTeachers()
	if err != nil {
		log.Printf("Erro interno ao buscar todos os professores: %v", err)
		http.Error(w, "Erro interno do servidor ao buscar professores", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(teachers)
}

// UpdateTeacherHandler atualiza um professor existente via HTTP PUT.
// Rota: /teachers/{id} (PUT)
func (h *TeacherHandler) UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	var teacher models.Teacher
	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		http.Error(w, "Requisição inválida: "+err.Error(), http.StatusBadRequest)
		return
	}
	teacher.ID = id // Garante que o ID da URL seja usado para atualização

	if err := h.teacherService.UpdateTeacher(&teacher); err != nil {
		if strings.Contains(err.Error(), "obrigatório") || strings.Contains(err.Error(), "não encontrado") {
			http.Error(w, err.Error(), http.StatusNotFound) // Use 404 para não encontrado
			return
		}
		log.Printf("Erro interno ao atualizar professor %s: %v", id, err)
		http.Error(w, "Erro interno do servidor ao atualizar professor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(teacher) // Retorna o professor atualizado
}

// DeleteTeacherHandler deleta um professor pelo ID via HTTP DELETE.
// Rota: /teachers/{id} (DELETE)
func (h *TeacherHandler) DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	if err := h.teacherService.DeleteTeacher(id); err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Erro interno ao deletar professor %s: %v", id, err)
		http.Error(w, "Erro interno do servidor ao deletar professor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content para deleção bem-sucedida
}

// AddSubjectToTeacherHandler associa uma matéria a um professor via HTTP POST.
// Rota: /teachers/{teacherID}/subjects/{subjectID} (POST)
func (h *TeacherHandler) AddSubjectToTeacherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	teacherID := params["teacherID"]
	subjectID := params["subjectID"]

	if teacherID == "" || subjectID == "" {
		http.Error(w, "IDs de professor e matéria são obrigatórios", http.StatusBadRequest)
		return
	}

	err := h.teacherService.AddSubjectToTeacherService(teacherID, subjectID)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Erro ao adicionar matéria ao professor: %v", err)
		http.Error(w, "Erro interno do servidor ao adicionar matéria ao professor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Matéria adicionada ao professor com sucesso."})
}

// RemoveSubjectFromTeacherHandler desassocia uma matéria de um professor via HTTP DELETE.
// Rota: /teachers/{teacherID}/subjects/{subjectID} (DELETE)
func (h *TeacherHandler) RemoveSubjectFromTeacherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	teacherID := params["teacherID"]
	subjectID := params["subjectID"]

	if teacherID == "" || subjectID == "" {
		http.Error(w, "IDs de professor e matéria são obrigatórios", http.StatusBadRequest)
		return
	}

	err := h.teacherService.RemoveSubjectFromTeacherService(teacherID, subjectID)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Erro ao remover matéria do professor: %v", err)
		http.Error(w, "Erro interno do servidor ao remover matéria do professor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Matéria removida do professor com sucesso."})
}
