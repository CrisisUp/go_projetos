// services/student_service.go

package services

import (
	"college-app-v1/models"       // Ajuste o caminho do import
	"college-app-v1/repositories" // Ajuste o caminho do import
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time" // Necessário para time.Now().Year()
)

// StudentService representa as operações de negócio para alunos.
// Ajustado para receber ponteiros para os repositórios.
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

	// 2. Obter o ano atual (para a matrícula)
	currentYearForEnrollment := time.Now().Year()

	// 3. Buscar a última matrícula para o ano e turno atuais
	lastEnrollment, err := s.studentRepo.GetLastEnrollmentForYearAndShift(currentYearForEnrollment, student.Shift)
	if err != nil {
		return fmt.Errorf("erro ao buscar última matrícula para geração automática: %w", err)
	}

	// 4. Gerar a nova matrícula
	newSequence := 1
	if lastEnrollment != "" {
		if len(lastEnrollment) >= 5 { // Verifica se há caracteres suficientes para a sequência
			// Assume que a sequência são os últimos 4 caracteres
			seqStr := lastEnrollment[len(lastEnrollment)-4:]
			lastSequence, err := strconv.Atoi(seqStr)
			if err == nil {
				newSequence = lastSequence + 1
			} else {
				log.Printf("Aviso: CreateStudent: Não foi possível converter sequência '%s' da última matrícula para int. Reiniciando sequência para 1. Erro: %v", seqStr, err)
			}
		} else {
			log.Printf("Aviso: CreateStudent: Matrícula '%s' tem formato inesperado. Reiniciando sequência para 1.", lastEnrollment)
		}
	}

	// Formata a nova matrícula (ex: 2025M0001)
	student.Enrollment = fmt.Sprintf("%d%s%04d", currentYearForEnrollment, student.Shift, newSequence)

	// O `CurrentYear` do aluno pode vir do frontend ou ser padronizado.
	// Se `student.CurrentYear` vier do frontend e for válido, use-o.
	// Se for 0 (não fornecido pelo frontend ou inválido), padronize para 1.
	if student.CurrentYear == 0 {
		student.CurrentYear = 1 // Padrão para o primeiro ano se não especificado ou for 0
	}
	// Adicionar validação se o CurrentYear vindo do frontend for um valor futuro absurdo, etc.

	return s.studentRepo.CreateStudent(student)
}

// GetStudentByID busca um aluno pelo ID.
func (s *StudentService) GetStudentByID(id string) (*models.Student, error) {
	student, err := s.studentRepo.GetStudentByID(id)
	if err != nil {
		if errors.Is(err, errors.New("aluno não encontrado")) { // Verifique se é o erro de 'não encontrado' do repositório
			return nil, fmt.Errorf("aluno com ID %s não encontrado", id)
		}
		return nil, fmt.Errorf("erro ao buscar aluno por ID: %w", err)
	}
	return student, nil
}

// GetAllStudents busca todos os alunos, com opções de filtro.
// year: ponteiro para int para permitir nil (sem filtro de ano)
// shift: string para o turno (vazio significa sem filtro de turno)
func (s *StudentService) GetAllStudents(year *int, shift string) ([]models.Student, error) {
	// Aqui você pode adicionar lógica de negócio adicional ou validações para os filtros, se necessário.
	// Por exemplo, validar se o ano é um número razoável, ou se o turno é "M", "T", "N".
	if shift != "" {
		shift = strings.ToUpper(shift)
		if shift != "M" && shift != "T" && shift != "N" {
			// Se o frontend enviar um turno inválido, podemos retornar um erro aqui
			return nil, fmt.Errorf("turno inválido no filtro: %s. Deve ser 'M', 'T' ou 'N'", shift)
		}
	}

	// Delega a chamada para o repositório com os filtros
	students, err := s.studentRepo.GetAllStudents(year, shift)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar todos os alunos com filtros: %w", err)
	}
	return students, nil
}

// UpdateStudent atualiza um aluno existente.
func (s *StudentService) UpdateStudent(student *models.Student) error {
	if student.ID == "" {
		return errors.New("ID do aluno é obrigatório para atualização")
	}
	if student.Name == "" || student.CurrentYear == 0 || student.Shift == "" {
		return errors.New("nome, ano atual e turno do aluno são obrigatórios para atualização")
	}

	// Normalizar o turno para maiúsculas antes de usar
	student.Shift = strings.ToUpper(student.Shift)
	if student.Shift != "M" && student.Shift != "T" && student.Shift != "N" {
		return fmt.Errorf("turno inválido para atualização: %s. Deve ser 'M' (Manhã), 'T' (Tarde) ou 'N' (Noite)", student.Shift)
	}

	existingStudent, err := s.studentRepo.GetStudentByID(student.ID)
	if err != nil {
		if errors.Is(err, errors.New("aluno não encontrado")) { // Verifique se é o erro de 'não encontrado' do repositório
			return fmt.Errorf("aluno com ID %s não encontrado para atualização", student.ID)
		}
		return fmt.Errorf("erro ao buscar aluno existente para atualização: %w", err)
	}
	// `existingStudent` já não será nil aqui se o erro for tratado acima.

	// Copia os campos atualizáveis do 'student' (DTO de entrada) para 'existingStudent'
	existingStudent.Name = student.Name
	existingStudent.CurrentYear = student.CurrentYear
	existingStudent.Shift = student.Shift // Já normalizado para maiúscula acima

	// A matrícula (Enrollment) é gerada na criação e não deve ser alterada aqui.
	// Ela já é parte do 'existingStudent' buscado do DB.

	return s.studentRepo.UpdateStudent(existingStudent)
}

// DeleteStudent deleta um aluno pelo ID.
func (s *StudentService) DeleteStudent(id string) error {
	err := s.studentRepo.DeleteStudent(id)
	if err != nil {
		if errors.Is(err, errors.New("aluno não encontrado para exclusão")) { // Verifique o erro específico do repositório
			return fmt.Errorf("aluno com ID %s não encontrado para exclusão", id)
		}
		return fmt.Errorf("erro ao deletar aluno: %w", err)
	}
	return nil
}

// AddSubjectToStudent associa uma matéria a um aluno.
func (s *StudentService) AddSubjectToStudent(studentID, subjectID string) error {
	student, err := s.studentRepo.GetStudentByID(studentID)
	if err != nil {
		if err.Error() == "aluno não encontrado" {
			return fmt.Errorf("aluno com ID %s não encontrado para associação", studentID)
		}
		return fmt.Errorf("erro ao buscar aluno para associação: %w", err)
	}
	if student == nil {
		return fmt.Errorf("aluno com ID %s não encontrado para associação", studentID)
	}
	subject, err := s.subjectRepo.GetSubjectByID(subjectID)
	if err != nil {
		if err.Error() == "matéria não encontrada" {
			return fmt.Errorf("matéria com ID %s não encontrada para associação", subjectID)
		}
		return fmt.Errorf("erro ao buscar matéria para associação: %w", err)
	}
	if subject == nil {
		return fmt.Errorf("matéria com ID %s não encontrada para associação", subjectID)
	}

	return s.studentRepo.AddSubjectToStudent(studentID, subjectID)
}

// RemoveSubjectFromStudent desassocia uma matéria de um aluno.
func (s *StudentService) RemoveSubjectFromStudent(studentID, subjectID string) error {
	// Verifica se o aluno existe
	student, err := s.studentRepo.GetStudentByID(studentID)
	if err != nil {
		if errors.Is(err, errors.New("aluno não encontrado")) {
			return fmt.Errorf("aluno com ID %s não encontrado para desassociação", studentID)
		}
		return fmt.Errorf("erro ao buscar aluno para desassociação: %w", err)
	}
	if student == nil {
		return fmt.Errorf("aluno com ID %s não encontrado para desassociação", studentID)
	}

	// Verifica se a matéria existe
	subject, err := s.subjectRepo.GetSubjectByID(subjectID)
	if err != nil {
		return fmt.Errorf("erro ao buscar matéria para desassociação: %w", err)
	}
	if subject == nil {
		return fmt.Errorf("matéria com ID %s não encontrada para desassociação", subjectID)
	}

	// Tenta remover a associação
	err = s.studentRepo.RemoveSubjectFromStudent(studentID, subjectID)
	if err != nil {
		if err.Error() == "associação não encontrada para desassociação" {
			return fmt.Errorf("associação entre aluno %s e matéria %s não encontrada para desassociação", studentID, subjectID)
		}
		return fmt.Errorf("erro ao remover associação entre aluno e matéria: %w", err)
	}
	return nil
}
