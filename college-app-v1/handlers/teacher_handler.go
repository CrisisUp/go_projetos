// handlers/teacher_handler.go

package handlers

import (
	"college-app-v1/models"   // Certifique-se de que este caminho está correto
	"college-app-v1/services" // Certifique-se de que este caminho está correto
	"encoding/json"
	"log"
	"net/http"

	// "strconv" // Não é necessário aqui, pois os filtros são strings (ou ponteiro para int se fosse numérico)

	"github.com/gorilla/mux"
)

// SubjectHandler gerencia as requisições HTTP para matérias. (Mantenha o SubjectHandler se ele existir)
// ... (seu código SubjectHandler) ...

// StudentHandler gerencia as requisições HTTP para alunos. (Mantenha o StudentHandler se ele existir)
// ... (seu código StudentHandler) ...

// TeacherHandler gerencia as requisições HTTP para professores.
type TeacherHandler struct {
	service *services.TeacherService // Ponteiro para o serviço de professor
}

// NewTeacherHandler cria uma nova instância de TeacherHandler.
func NewTeacherHandler(s *services.TeacherService) *TeacherHandler {
	return &TeacherHandler{service: s}
}

// CreateTeacherHandler lida com a criação de um novo professor.
// POST /teachers
func (h *TeacherHandler) CreateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var teacher models.Teacher // Assumimos que o modelo Teacher tem Name, Department, Email
	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		http.Error(w, `{"message": "Requisição inválida: corpo JSON malformado."}`, http.StatusBadRequest)
		return
	}

	if err := h.service.CreateTeacher(&teacher); err != nil {
		log.Printf("CreateTeacherHandler: Erro ao criar professor no serviço: %v", err)
		http.Error(w, `{"message": "Erro ao criar professor: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(teacher)
}

// GetTeacherByIDHandler lida com a busca de um professor por ID.
// GET /teachers/{id}
func (h *TeacherHandler) GetTeacherByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	teacher, err := h.service.GetTeacherByID(id)
	if err != nil {
		if err.Error() == "professor não encontrado" { // Mensagem de erro específica do serviço
			http.Error(w, `{"message": "Professor não encontrado."}`, http.StatusNotFound)
			return
		}
		log.Printf("GetTeacherByIDHandler: Erro ao buscar professor no serviço: %v", err)
		http.Error(w, `{"message": "Erro ao buscar professor: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(teacher)
}

// GetAllTeachersHandler lida com a busca de todos os professores, com filtros opcionais.
// GET /teachers?name=X&department=Y&email=Z
func (h *TeacherHandler) GetAllTeachersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extrair parâmetros de consulta (query parameters)
	query := r.URL.Query()
	nameFilter := query.Get("name")
	departmentFilter := query.Get("department")
	emailFilter := query.Get("email")

	// --- ADICIONE ESTES LOGS TEMPORÁRIOS PARA DEPURAR (REMOVER EM PRODUÇÃO) ---
	log.Printf("GetAllTeachersHandler: Requisição recebida. URL: %s", r.URL.String())
	log.Printf("GetAllTeachersHandler: Filtro Nome: '%s', Departamento: '%s', Email: '%s'", nameFilter, departmentFilter, emailFilter)
	// --- FIM DOS LOGS TEMPORÁRIOS ---

	// Chamar o serviço com os filtros
	teachers, err := h.service.GetAllTeachers(nameFilter, departmentFilter, emailFilter) // <-- NOVA ASSINATURA
	if err != nil {
		log.Printf("GetAllTeachersHandler: Erro ao buscar professores no serviço: %v", err)
		http.Error(w, `{"message": "Erro ao buscar professores: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(teachers)
}

// UpdateTeacherHandler lida com a atualização de um professor existente.
// PUT /teachers/{id}
func (h *TeacherHandler) UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	var teacher models.Teacher
	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		http.Error(w, `{"message": "Requisição inválida: corpo JSON malformado."}`, http.StatusBadRequest)
		return
	}

	teacher.ID = id // Garante que o ID da URL seja usado
	if err := h.service.UpdateTeacher(&teacher); err != nil {
		if err.Error() == "professor não encontrado para atualização" {
			http.Error(w, `{"message": "`+err.Error()+`"}`, http.StatusNotFound)
			return
		}
		log.Printf("UpdateTeacherHandler: Erro ao atualizar professor no serviço: %v", err)
		http.Error(w, `{"message": "Erro ao atualizar professor: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(teacher)
}

// DeleteTeacherHandler lida com a exclusão de um professor por ID.
// DELETE /teachers/{id}
func (h *TeacherHandler) DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.DeleteTeacher(id); err != nil {
		if err.Error() == "professor não encontrado para exclusão" {
			http.Error(w, `{"message": "`+err.Error()+`"}`, http.StatusNotFound)
			return
		}
		log.Printf("DeleteTeacherHandler: Erro ao deletar professor no serviço: %v", err)
		http.Error(w, `{"message": "Erro ao deletar professor: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddSubjectToTeacherHandler lida com a adição de uma matéria a um professor.
// POST /teachers/{teacherID}/subjects/{subjectID}
func (h *TeacherHandler) AddSubjectToTeacherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	teacherID := vars["teacherID"]
	subjectID := vars["subjectID"]

	if err := h.service.AddSubjectToTeacher(teacherID, subjectID); err != nil {
		errorMessage := err.Error()
		if errorMessage == "professor não encontrado para associação" || errorMessage == "matéria não encontrada para associação" {
			http.Error(w, `{"message": "`+errorMessage+`"}`, http.StatusNotFound)
			return
		}
		if errorMessage == "matéria já associada a este professor" {
			http.Error(w, `{"message": "`+errorMessage+`"}`, http.StatusConflict)
			return
		}
		log.Printf("AddSubjectToTeacherHandler: Erro ao adicionar matéria ao professor no serviço: %v", err)
		http.Error(w, `{"message": "Erro ao adicionar matéria ao professor: `+errorMessage+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Matéria adicionada ao professor com sucesso."})
}

// RemoveSubjectFromTeacherHandler lida com a remoção de uma matéria de um professor.
// DELETE /teachers/{teacherID}/subjects/{subjectID}
func (h *TeacherHandler) RemoveSubjectFromTeacherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	teacherID := vars["teacherID"]
	subjectID := vars["subjectID"]

	if err := h.service.RemoveSubjectFromTeacher(teacherID, subjectID); err != nil {
		errorMessage := err.Error()
		if errorMessage == "professor não encontrado para desassociação" || errorMessage == "associação entre professor "+teacherID+" e matéria "+subjectID+" não encontrada para desassociação" {
			http.Error(w, `{"message": "`+errorMessage+`"}`, http.StatusNotFound)
			return
		}
		log.Printf("RemoveSubjectFromTeacherHandler: Erro ao remover matéria do professor no serviço: %v", err)
		http.Error(w, `{"message": "Erro ao remover matéria do professor: `+errorMessage+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
