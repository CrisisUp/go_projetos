// services/subject_service.go
package services

import (
	"college-app-v1/models"
	"college-app-v1/repositories"
	"database/sql" // Para verificar sql.ErrNoRows
	"errors"       // Para criar erros personalizados
	"fmt"          // Para formatar mensagens de erro
)

// SubjectService define a interface para as operações de serviço de matérias.
type SubjectService struct {
	repo *repositories.SubjectRepository
}

// NewSubjectService cria uma nova instância de SubjectService.
func NewSubjectService(repo *repositories.SubjectRepository) *SubjectService {
	return &SubjectService{repo: repo}
}

// CreateSubject adiciona uma nova matéria após validações.
func (s *SubjectService) CreateSubject(subject *models.Subject) error {
	// Exemplo de validação: ID, nome e ano são obrigatórios
	if subject.ID == "" || subject.Name == "" || subject.Year == 0 {
		return errors.New("ID, nome e ano da matéria são obrigatórios")
	}

	// Exemplo de validação: Matéria com o mesmo ID já existe
	existingSubject, err := s.repo.GetSubjectByID(subject.ID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("erro ao verificar matéria existente: %w", err)
	}
	if existingSubject != nil {
		return errors.New("matéria com este ID já existe")
	}

	return s.repo.CreateSubject(subject)
}

// GetSubjectByID busca uma matéria pelo ID.
func (s *SubjectService) GetSubjectByID(id string) (*models.Subject, error) {
	subject, err := s.repo.GetSubjectByID(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar matéria: %w", err)
	}
	if subject == nil {
		return nil, errors.New("matéria não encontrada")
	}
	return subject, nil
}

// GetAllSubjects busca todas as matérias.
func (s *SubjectService) GetAllSubjects() ([]models.Subject, error) {
	subjects, err := s.repo.GetAllSubjects()
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar todas as matérias: %w", err)
	}
	return subjects, nil
}

// UpdateSubject atualiza uma matéria existente após validações.
func (s *SubjectService) UpdateSubject(subject *models.Subject) error {
	if subject.ID == "" {
		return errors.New("ID da matéria é obrigatório para atualização")
	}
	// Validação: a matéria deve existir para ser atualizada
	existingSubject, err := s.repo.GetSubjectByID(subject.ID)
	if err != nil {
		return fmt.Errorf("erro ao verificar matéria para atualização: %w", err)
	}
	if existingSubject == nil {
		return errors.New("matéria não encontrada para atualização")
	}

	return s.repo.UpdateSubject(subject)
}

// DeleteSubject deleta uma matéria pelo ID.
func (s *SubjectService) DeleteSubject(id string) error {
	if id == "" {
		return errors.New("ID da matéria é obrigatório para exclusão")
	}
	// Validação: a matéria deve existir para ser deletada
	existingSubject, err := s.repo.GetSubjectByID(id)
	if err != nil {
		return fmt.Errorf("erro ao verificar matéria para exclusão: %w", err)
	}
	if existingSubject == nil {
		return errors.New("matéria não encontrada para exclusão")
	}

	return s.repo.DeleteSubject(id)
}
