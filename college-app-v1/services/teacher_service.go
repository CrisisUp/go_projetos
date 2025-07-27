// services/teacher_service.go

package services

import (
	"college-app-v1/models"       // Certifique-se de que este caminho está correto
	"college-app-v1/repositories" // Certifique-se de que este caminho está correto
	"errors"
	"fmt"
	"log" // <-- ADICIONADO: Importar log para depuração
	// Para usar strings.ToUpper (se necessário para normalização de filtros)
)

// TeacherService define a interface para operações de negócio de professor.
// Assinatura atualizada para GetAllTeachers.
type TeacherService struct {
	teacherRepo *repositories.TeacherRepository
	subjectRepo *repositories.SubjectRepository // Se o serviço precisar interagir com matérias
}

// NewTeacherService cria uma nova instância de TeacherService.
func NewTeacherService(tr *repositories.TeacherRepository, sr *repositories.SubjectRepository) *TeacherService {
	return &TeacherService{teacherRepo: tr, subjectRepo: sr}
}

// CreateTeacher implementa a criação de um novo professor.
func (s *TeacherService) CreateTeacher(teacher *models.Teacher) error {
	// Validações de negócio para criação (Name, Department, Email)
	if teacher.Name == "" || teacher.Department == "" || teacher.Email == "" {
		return errors.New("nome, departamento e email do professor são obrigatórios")
	}
	// Você pode adicionar validação de formato de email aqui.

	return s.teacherRepo.CreateTeacher(teacher)
}

// GetTeacherByID implementa a busca de professor por ID.
func (s *TeacherService) GetTeacherByID(id string) (*models.Teacher, error) {
	teacher, err := s.teacherRepo.GetTeacherByID(id)
	if err != nil {
		if errors.Is(err, errors.New("professor não encontrado")) {
			return nil, fmt.Errorf("professor com ID %s não encontrado", id)
		}
		return nil, fmt.Errorf("erro ao buscar professor por ID: %w", err)
	}
	return teacher, nil
}

// GetAllTeachers implementa a busca de todos os professores com filtros.
// Adicionado nameFilter e emailFilter.
func (s *TeacherService) GetAllTeachers(nameFilter, departmentFilter, emailFilter string) ([]models.Teacher, error) {
	// --- ADICIONADO: LOG TEMPORÁRIO AQUI ---
	log.Printf("Service: Chamando GetAllTeachers com filtros -> Nome: '%s', Departamento: '%s', Email: '%s'", nameFilter, departmentFilter, emailFilter)
	// --- FIM DO LOG TEMPORÁRIO ---

	// Exemplo de normalização do filtro (opcional, mas boa prática)
	// if departmentFilter != "" {
	//     departmentFilter = strings.ToUpper(departmentFilter)
	// }
	// if nameFilter != "" {
	//     nameFilter = strings.ToLower(nameFilter) // Para buscas case-insensitive no repositório
	// }

	teachers, err := s.teacherRepo.GetAllTeachers(nameFilter, departmentFilter, emailFilter)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar todos os professores com filtros: %w", err)
	}
	return teachers, nil
}

// UpdateTeacher implementa a atualização de um professor.
func (s *TeacherService) UpdateTeacher(teacher *models.Teacher) error {
	if teacher.ID == "" {
		return errors.New("ID do professor é obrigatório para atualização")
	}
	if teacher.Name == "" || teacher.Department == "" || teacher.Email == "" { // Validação completa
		return errors.New("nome, departamento e email do professor são obrigatórios para atualização")
	}

	existingTeacher, err := s.teacherRepo.GetTeacherByID(teacher.ID)
	if err != nil {
		if errors.Is(err, errors.New("professor não encontrado")) {
			return fmt.Errorf("professor com ID %s não encontrado para atualização", teacher.ID)
		}
		return fmt.Errorf("erro ao buscar professor existente para atualização: %w", err)
	}

	// Atualiza os campos do professor existente
	existingTeacher.Name = teacher.Name
	existingTeacher.Department = teacher.Department
	existingTeacher.Email = teacher.Email

	return s.teacherRepo.UpdateTeacher(existingTeacher)
}

// DeleteTeacher implementa a exclusão de um professor.
func (s *TeacherService) DeleteTeacher(id string) error {
	err := s.teacherRepo.DeleteTeacher(id)
	if err != nil {
		if errors.Is(err, errors.New("professor não encontrado para exclusão")) {
			return fmt.Errorf("professor com ID %s não encontrado para exclusão", id)
		}
		return fmt.Errorf("erro ao deletar professor: %w", err)
	}
	return nil
}

// AddSubjectToTeacher associa uma matéria a um professor.
func (s *TeacherService) AddSubjectToTeacher(teacherID, subjectID string) error {
	_, err := s.teacherRepo.GetTeacherByID(teacherID)
	if err != nil {
		if errors.Is(err, errors.New("professor não encontrado")) {
			return fmt.Errorf("professor com ID %s não encontrado para associação", teacherID)
		}
		return fmt.Errorf("erro ao buscar professor para associação: %w", err)
	}
	// Supondo que GetSubjectByID no subjectRepo retorna (nil, nil) se não encontrar
	subject, err := s.subjectRepo.GetSubjectByID(subjectID)
	if err != nil {
		return fmt.Errorf("erro ao buscar matéria para associação: %w", err)
	}
	if subject == nil {
		return fmt.Errorf("matéria com ID %s não encontrada para associação", subjectID)
	}

	return s.teacherRepo.AddSubjectToTeacher(teacherID, subjectID)
}

// RemoveSubjectFromTeacher desassocia uma matéria de um professor.
func (s *TeacherService) RemoveSubjectFromTeacher(teacherID, subjectID string) error {
	err := s.teacherRepo.RemoveSubjectFromTeacher(teacherID, subjectID)
	if err != nil {
		if errors.Is(err, errors.New("associação não encontrada para desassociação")) {
			return fmt.Errorf("associação entre professor %s e matéria %s não encontrada para desassociação", teacherID, subjectID)
		}
		return fmt.Errorf("erro ao desassociar matéria do professor: %w", err)
	}
	return nil
}
