// handlers/student_handler.go
package handlers

import (
	"college-app-v1/models"   // Certifique-se de que este caminho está correto
	"college-app-v1/services" // Certifique-se de que este caminho está correto
	"encoding/json"
	"log"
	"net/http"
	"strconv" // Adicionado para converter string de query param para int

	"github.com/gorilla/mux"
)

// StudentHandler gerencia as requisições HTTP para alunos.
type StudentHandler struct {
	service *services.StudentService // Ponteiro para o serviço, conforme sua definição original
}

// NewStudentHandler cria uma nova instância de StudentHandler.
func NewStudentHandler(s *services.StudentService) *StudentHandler {
	return &StudentHandler{service: s}
}

// CreateStudentHandler lida com a criação de um novo aluno.
// POST /students
func (h *StudentHandler) CreateStudentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Define o Content-Type no início

	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, `{"message": "Requisição inválida: corpo JSON malformado."}`, http.StatusBadRequest)
		return
	}

	if err := h.service.CreateStudent(&student); err != nil {
		log.Printf("CreateStudentHandler: Erro ao criar aluno no serviço: %v", err)
		// Você pode adicionar tratamento de erro mais granular aqui com base no tipo de erro retornado pelo serviço.
		// Ex: if strings.Contains(err.Error(), "turno inválido") { http.Error(w, err.Error(), http.StatusBadRequest) }
		http.Error(w, `{"message": "Erro ao criar aluno: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

// GetStudentByIDHandler lida com a busca de um aluno por ID.
// GET /students/{id}
func (h *StudentHandler) GetStudentByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Define o Content-Type no início

	vars := mux.Vars(r)
	id := vars["id"]

	student, err := h.service.GetStudentByID(id)
	if err != nil {
		// Use errors.Is para verificar tipos de erro específicos, se seus erros forem tratados assim
		if err.Error() == "aluno com ID "+id+" não encontrado" { // Mensagem de erro específica do serviço
			http.Error(w, `{"message": "Aluno não encontrado."}`, http.StatusNotFound)
			return
		}
		log.Printf("GetStudentByIDHandler: Erro ao buscar aluno no serviço: %v", err)
		http.Error(w, `{"message": "Erro ao buscar aluno: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(student)
}

// GetAllStudentsHandler lida com a busca de todos os alunos, com filtros opcionais.
// GET /students?current_year=X&shift=Y
func (h *StudentHandler) GetAllStudentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Define o Content-Type no início

	// Extrair parâmetros de consulta (query parameters)
	query := r.URL.Query()
	yearStr := query.Get("current_year")
	shiftFilter := query.Get("shift") // Renomeado para evitar conflito com 'shift' do modelo/frontend se houvesse

	var yearFilter *int // Usar ponteiro para diferenciar 0 de não fornecido
	if yearStr != "" {
		parsedYear, err := strconv.Atoi(yearStr)
		if err != nil {
			http.Error(w, `{"message": "Ano inválido fornecido. Deve ser um número inteiro."}`, http.StatusBadRequest)
			return
		}
		yearFilter = &parsedYear
	}

	// Chamar o serviço com os filtros
	students, err := h.service.GetAllStudents(yearFilter, shiftFilter)
	if err != nil {
		log.Printf("GetAllStudentsHandler: Erro ao buscar todos os alunos no serviço: %v", err)
		// Aqui, você pode adicionar tratamento mais específico para erros do serviço (ex: turno inválido no filtro)
		http.Error(w, `{"message": "Erro ao buscar alunos: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Se tudo correr bem, retorna a lista (pode ser vazia se nenhum aluno corresponder ao filtro)
	json.NewEncoder(w).Encode(students)
}

// UpdateStudentHandler lida com a atualização de um aluno existente.
// PUT /students/{id}
func (h *StudentHandler) UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Define o Content-Type no início

	vars := mux.Vars(r)
	id := vars["id"]

	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, `{"message": "Requisição inválida: corpo JSON malformado."}`, http.StatusBadRequest)
		return
	}

	student.ID = id // Garante que o ID da URL seja usado para a atualização

	if err := h.service.UpdateStudent(&student); err != nil {
		if err.Error() == "aluno não encontrado para atualização" || err.Error() == "ID do aluno é obrigatório para atualização" {
			http.Error(w, `{"message": "`+err.Error()+`"}`, http.StatusNotFound) // Use 404 para não encontrado
			return
		}
		// Se o erro for de validação (ex: nome, ano, turno obrigatórios), você pode retornar 400 Bad Request
		// if strings.Contains(err.Error(), "obrigatório") || strings.Contains(err.Error(), "turno inválido") {
		//    http.Error(w, `{"message": "`+err.Error()+`"}`, http.StatusBadRequest)
		//    return
		// }
		log.Printf("UpdateStudentHandler: Erro ao atualizar aluno no serviço: %v", err)
		http.Error(w, `{"message": "Erro ao atualizar aluno: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(student) // Retorna o aluno atualizado
}

// DeleteStudentHandler lida com a exclusão de um aluno por ID.
// DELETE /students/{id}
func (h *StudentHandler) DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Define o Content-Type no início

	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.DeleteStudent(id); err != nil {
		if err.Error() == "aluno não encontrado para exclusão" {
			http.Error(w, `{"message": "`+err.Error()+`"}`, http.StatusNotFound)
			return
		}
		log.Printf("DeleteStudentHandler: Erro ao deletar aluno no serviço: %v", err)
		http.Error(w, `{"message": "Erro ao deletar aluno: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// AddSubjectToStudentHandler lida com a adição de uma matéria a um aluno.
// POST /students/{studentID}/subjects/{subjectID}
func (h *StudentHandler) AddSubjectToStudentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Define o Content-Type no início

	vars := mux.Vars(r)
	studentID := vars["studentID"]
	subjectID := vars["subjectID"]

	if err := h.service.AddSubjectToStudent(studentID, subjectID); err != nil {
		// Use um switch ou if-else if para erros mais específicos do serviço
		errorMessage := err.Error()
		if errorMessage == "aluno com ID "+studentID+" não encontrado para associação" || errorMessage == "matéria com ID "+subjectID+" não encontrada para associação" {
			http.Error(w, `{"message": "`+errorMessage+`"}`, http.StatusNotFound) // 404 Not Found
			return
		}
		// Se o erro for "matéria já associada a este aluno", pode ser 409 Conflict ou 400 Bad Request
		if errorMessage == "matéria já associada a este aluno" {
			http.Error(w, `{"message": "`+errorMessage+`"}`, http.StatusConflict) // 409 Conflict
			return
		}
		log.Printf("AddSubjectToStudentHandler: Erro ao adicionar matéria ao aluno no serviço: %v", err)
		http.Error(w, `{"message": "Erro ao adicionar matéria ao aluno: `+errorMessage+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Matéria adicionada ao aluno com sucesso."})
}

// RemoveSubjectFromStudentHandler lida com a remoção de uma matéria de um aluno.
// DELETE /students/{studentID}/subjects/{subjectID}
func (h *StudentHandler) RemoveSubjectFromStudentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Define o Content-Type no início

	vars := mux.Vars(r)
	studentID := vars["studentID"]
	subjectID := vars["subjectID"]

	if err := h.service.RemoveSubjectFromStudent(studentID, subjectID); err != nil {
		errorMessage := err.Error()
		if errorMessage == "aluno com ID "+studentID+" não encontrado para desassociação" || errorMessage == "associação entre aluno "+studentID+" e matéria "+subjectID+" não encontrada para desassociação" {
			http.Error(w, `{"message": "`+errorMessage+`"}`, http.StatusNotFound) // 404 Not Found
			return
		}
		log.Printf("RemoveSubjectFromStudentHandler: Erro ao remover matéria do aluno no serviço: %v", err)
		http.Error(w, `{"message": "Erro ao remover matéria do aluno: `+errorMessage+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
