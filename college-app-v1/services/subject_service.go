// services/subject_service.go
package services

import (
	"college-app-v1/models"
	"college-app-v1/repositories" // Para verificar sql.ErrNoRows
	"errors"                      // Para criar erros personalizados
	"fmt"                         // Para formatar mensagens de erro
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
	// --- MUDANÇA AQUI: Remover a validação de subject.ID para criação ---
	if subject.Name == "" || subject.Year == 0 { // ID não é mais verificado aqui
		return errors.New("nome e ano da matéria são obrigatórios") // Mensagem de erro atualizada
	}

	// Exemplo de validação: Matéria com o mesmo ID já existe
	// NOTA: Esta validação só faz sentido SE você permitir que o usuário forneça o ID,
	// ou se você gerar um ID único antes de chamar GetSubjectByID (o que não é o caso aqui).
	// Se o ID é gerado pelo repositório, esta verificação pode ser removida ou adaptada.
	// Por enquanto, vou comentá-la, pois o repositório gerará um UUID e garantirá unicidade.
	/*
		existingSubject, err := s.repo.GetSubjectByID(subject.ID) // subject.ID estaria vazio aqui
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("erro ao verificar matéria existente: %w", err)
		}
		if existingSubject != nil {
			return errors.New("matéria com este ID já existe")
		}
	*/

	// A geração do ID DEVE ocorrer no repositório antes de criar no DB
	return s.repo.CreateSubject(subject)
}

// GetSubjectByID busca uma matéria pelo ID.
func (s *SubjectService) GetSubjectByID(id string) (*models.Subject, error) {
	subject, err := s.repo.GetSubjectByID(id)
	if err != nil {
		// Encapsular erros do repositório para a camada de serviço
		if errors.Is(err, errors.New("matéria não encontrada")) { // Supondo que o repositório retorna este erro
			return nil, errors.New("matéria não encontrada")
		}
		return nil, fmt.Errorf("erro ao buscar matéria: %w", err)
	}
	// Se o repositório retornar (nil, nil) para não encontrado, esta verificação pega
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
	// Adicionar validação de nome e ano também para atualização
	if subject.Name == "" || subject.Year == 0 {
		return errors.New("nome e ano da matéria são obrigatórios para atualização")
	}

	// Validação: a matéria deve existir para ser atualizada
	existingSubject, err := s.repo.GetSubjectByID(subject.ID)
	if err != nil {
		return fmt.Errorf("erro ao verificar matéria para atualização: %w", err)
	}
	if existingSubject == nil {
		return errors.New("matéria não encontrada para atualização")
	}

	// Copiar os campos atualizáveis (se necessário, para não sobrescrever o que não deve)
	// existingSubject.Name = subject.Name
	// existingSubject.Year = subject.Year
	// existingSubject.Credits = subject.Credits // Se créditos forem atualizáveis

	return s.repo.UpdateSubject(subject) // Passe o subject recebido que já tem o ID
}

// DeleteSubject deleta uma matéria pelo ID.
func (s *SubjectService) DeleteSubject(id string) error {
	if id == "" {
		return errors.New("ID da matéria é obrigatório para exclusão")
	}
	// Validação: a matéria deve existir para ser deletada
	existingSubject, err := s.repo.GetSubjectByID(id)
	if err != nil {
		// Se o erro do repositório for "matéria não encontrada", encapsule.
		if errors.Is(err, errors.New("matéria não encontrada")) {
			return errors.New("matéria não encontrada para exclusão")
		}
		return fmt.Errorf("erro ao verificar matéria para exclusão: %w", err)
	}
	// Se existingSubject for nil (significa que GetSubjectByID retornou nil, nil),
	// então a matéria não foi encontrada.
	if existingSubject == nil {
		return errors.New("matéria não encontrada para exclusão")
	}

	return s.repo.DeleteSubject(id)
}
