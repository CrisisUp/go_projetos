// services/student_service.go
package services

import (
	"college-app-v1/models"
	"college-app-v1/repositories"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// StudentService define as operações de negócio para alunos.
type StudentService struct {
	studentRepo *repositories.StudentRepository
	subjectRepo *repositories.SubjectRepository
}

// NewStudentService cria uma nova instância de StudentService.
func NewStudentService(sr *repositories.StudentRepository, subR *repositories.SubjectRepository) *StudentService {
	return &StudentService{studentRepo: sr, subjectRepo: subR}
}

// CreateStudent cria um novo aluno com matrícula gerada automaticamente.
func (s *StudentService) CreateStudent(student *models.Student) error {
	// 1. Validar o turno (Shift)
	student.Shift = strings.ToUpper(student.Shift)
	if student.Shift != "M" && student.Shift != "T" && student.Shift != "N" {
		return fmt.Errorf("turno inválido: %s. Deve ser 'M' (Manhã), 'T' (Tarde) ou 'N' (Noite)", student.Shift)
	}

	// 2. Obter o ano atual
	currentYear := time.Now().Year()

	// 3. Buscar a última matrícula para o ano e turno atuais
	lastEnrollment, err := s.studentRepo.GetLastEnrollmentForYearAndShift(currentYear, student.Shift)
	if err != nil {
		return fmt.Errorf("erro ao buscar última matrícula: %v", err)
	}

	// 4. Gerar a nova matrícula
	newSequence := 1
	if lastEnrollment != "" {
		if len(lastEnrollment) >= 5 {
			seqStr := lastEnrollment[len(lastEnrollment)-4:]
			lastSequence, err := strconv.Atoi(seqStr)
			if err == nil {
				newSequence = lastSequence + 1
			} else {
				log.Printf("Aviso: Não foi possível converter sequência '%s' para int. Reiniciando sequência para 1. Erro: %v", seqStr, err)
			}
		}
	}

	// Formata a nova matrícula (ex: 2025M0001)
	student.Enrollment = fmt.Sprintf("%d%s%04d", currentYear, student.Shift, newSequence)

	// Define o CurrentYear como o ano atual (pode ser ajustado depois pelo frontend/admin)
	if student.CurrentYear == 0 {
		student.CurrentYear = 1
	}

	return s.studentRepo.CreateStudent(student)
}

// GetStudentByID busca um aluno pelo ID.
func (s *StudentService) GetStudentByID(id string) (*models.Student, error) {
	return s.studentRepo.GetStudentByID(id)
}

// GetAllStudents busca todos os alunos.
func (s *StudentService) GetAllStudents() ([]models.Student, error) {
	return s.studentRepo.GetAllStudents()
}

// UpdateStudent atualiza um aluno existente.
func (s *StudentService) UpdateStudent(student *models.Student) error {
	if student.ID == "" {
		return errors.New("ID do aluno é obrigatório para atualização")
	}
	if student.Name == "" || student.CurrentYear == 0 || student.Shift == "" { // CORRIGIDO: Valida o Shift aqui também
		return errors.New("nome, ano atual e turno do aluno são obrigatórios para atualização")
	}

	existingStudent, err := s.studentRepo.GetStudentByID(student.ID)
	if err != nil {
		return fmt.Errorf("erro ao buscar aluno existente para atualização: %w", err)
	}
	if existingStudent == nil {
		return fmt.Errorf("aluno com ID %s não encontrado para atualização", student.ID)
	}

	// ATUALIZADO: Copia os campos atualizáveis do 'student' (DTO de entrada) para 'existingStudent'
	existingStudent.Name = student.Name
	existingStudent.CurrentYear = student.CurrentYear
	existingStudent.Shift = strings.ToUpper(student.Shift) // CORRIGIDO: Copia o turno e garante maiúscula
	// A matrícula (Enrollment) é gerada na criação e não deve ser alterada aqui.
	// Ela já é parte do 'existingStudent' buscado do DB.

	return s.studentRepo.UpdateStudent(existingStudent)
}

// DeleteStudent deleta um aluno pelo ID.
func (s *StudentService) DeleteStudent(id string) error {
	return s.studentRepo.DeleteStudent(id)
}

// AddSubjectToStudent associa uma matéria a um aluno.
func (s *StudentService) AddSubjectToStudent(studentID, subjectID string) error {
	student, err := s.studentRepo.GetStudentByID(studentID)
	if err != nil {
		return fmt.Errorf("erro ao buscar aluno: %w", err)
	}
	if student == nil {
		return fmt.Errorf("aluno com ID %s não encontrado", studentID)
	}

	subject, err := s.subjectRepo.GetSubjectByID(subjectID)
	if err != nil {
		return fmt.Errorf("erro ao buscar matéria: %w", err)
	}
	if subject == nil {
		return fmt.Errorf("matéria com ID %s não encontrada", subjectID)
	}

	return s.studentRepo.AddSubjectToStudent(studentID, subjectID)
}

// RemoveSubjectFromStudent desassocia uma matéria de um aluno.
func (s *StudentService) RemoveSubjectFromStudent(studentID, subjectID string) error {
	return s.studentRepo.RemoveSubjectFromStudent(studentID, subjectID)
}
