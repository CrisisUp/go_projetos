// handlers/student_handler.go
package handlers

import (
	"college-app-v1/models"
	"college-app-v1/services"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// StudentHandler gerencia as requisições HTTP para alunos.
type StudentHandler struct {
	service *services.StudentService
}

// NewStudentHandler cria uma nova instância de StudentHandler.
func NewStudentHandler(s *services.StudentService) *StudentHandler {
	return &StudentHandler{service: s}
}

// CreateStudentHandler lida com a criação de um novo aluno.
// POST /students
func (h *StudentHandler) CreateStudentHandler(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	// Não precisamos mais da matrícula no JSON de entrada, o serviço a gerará.
	// No entanto, precisamos do nome, ano e AGORA O TURNO (shift).
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Requisição inválida: "+err.Error(), http.StatusBadRequest)
		return
	}

	// O serviço agora validará e gerará a matrícula.
	if err := h.service.CreateStudent(&student); err != nil {
		log.Printf("Erro ao criar aluno no serviço: %v", err)
		// Aqui, você pode refinar o tratamento de erro para retornar 400 Bad Request
		// se o erro for de validação (ex: turno inválido).
		// Por enquanto, manteremos 500 para simplicidade, mas um tratamento mais granular seria melhor.
		http.Error(w, "Erro ao criar aluno: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

// GetStudentByIDHandler lida com a busca de um aluno por ID.
// GET /students/{id}
func (h *StudentHandler) GetStudentByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	student, err := h.service.GetStudentByID(id)
	if err != nil {
		if err.Error() == "aluno não encontrado" { // Erro personalizado do serviço
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Erro ao buscar aluno no serviço: %v", err)
		http.Error(w, "Erro ao buscar aluno: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

// GetAllStudentsHandler lida com a busca de todos os alunos.
// GET /students
func (h *StudentHandler) GetAllStudentsHandler(w http.ResponseWriter, r *http.Request) {
	students, err := h.service.GetAllStudents()
	if err != nil {
		log.Printf("Erro ao buscar todos os alunos no serviço: %v", err)
		http.Error(w, "Erro ao buscar alunos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

// UpdateStudentHandler lida com a atualização de um aluno existente.
// PUT /students/{id}
func (h *StudentHandler) UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Requisição inválida: "+err.Error(), http.StatusBadRequest)
		return
	}

	student.ID = id // Garante que o ID da URL seja usado

	if err := h.service.UpdateStudent(&student); err != nil {
		if err.Error() == "aluno não encontrado para atualização" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Erro ao atualizar aluno no serviço: %v", err)
		http.Error(w, "Erro ao atualizar aluno: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(student)
}

// DeleteStudentHandler lida com a exclusão de um aluno por ID.
// DELETE /students/{id}
func (h *StudentHandler) DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.DeleteStudent(id); err != nil {
		if err.Error() == "aluno não encontrado para exclusão" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Erro ao deletar aluno no serviço: %v", err)
		http.Error(w, "Erro ao deletar aluno: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// AddSubjectToStudentHandler lida com a adição de uma matéria a um aluno.
// POST /students/{studentID}/subjects/{subjectID}
func (h *StudentHandler) AddSubjectToStudentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["studentID"]
	subjectID := vars["subjectID"]

	if err := h.service.AddSubjectToStudent(studentID, subjectID); err != nil {
		if err.Error() == "aluno não encontrado" || err.Error() == "matéria não encontrada" || err.Error() == "matéria já associada a este aluno" {
			http.Error(w, err.Error(), http.StatusBadRequest) // 400 Bad Request para erros de validação
			return
		}
		log.Printf("Erro ao adicionar matéria ao aluno no serviço: %v", err)
		http.Error(w, "Erro ao adicionar matéria ao aluno: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Matéria adicionada ao aluno com sucesso."})
}

// RemoveSubjectFromStudentHandler lida com a remoção de uma matéria de um aluno.
// DELETE /students/{studentID}/subjects/{subjectID}
func (h *StudentHandler) RemoveSubjectFromStudentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["studentID"]
	subjectID := vars["subjectID"]

	if err := h.service.RemoveSubjectFromStudent(studentID, subjectID); err != nil {
		if err.Error() == "aluno não encontrado" || err.Error() == "matéria não associada a este aluno" {
			http.Error(w, err.Error(), http.StatusNotFound) // 404 Not Found para associação inexistente
			return
		}
		log.Printf("Erro ao remover matéria do aluno no serviço: %v", err)
		http.Error(w, "Erro ao remover matéria do aluno: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
