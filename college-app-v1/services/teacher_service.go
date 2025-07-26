package services

import (
	"college-app-v1/models"
	"college-app-v1/repositories"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// TeacherService define as operações de negócio para professores.
type TeacherService struct {
	teacherRepo *repositories.TeacherRepository // Renomeei 'repo' para 'teacherRepo' para clareza
	subjectRepo *repositories.SubjectRepository // <-- Novo campo para o repositório de matérias
}

// NewTeacherService cria uma nova instância de TeacherService.
// Agora recebe o subjectRepo para poder verificar matérias.
func NewTeacherService(tr *repositories.TeacherRepository, subR *repositories.SubjectRepository) *TeacherService {
	return &TeacherService{
		teacherRepo: tr,
		subjectRepo: subR, // <-- Atribui o subjectRepo recebido
	}
}

// CreateTeacher adiciona um novo professor com registro gerado automaticamente.
func (s *TeacherService) CreateTeacher(teacher *models.Teacher) error {
	// 1. Validações básicas
	if teacher.Name == "" || teacher.Department == "" || teacher.Email == "" {
		return errors.New("nome, departamento e email do professor são obrigatórios")
	}

	// 2. Verificar se o email já existe
	// (Você precisaria de um método GetTeacherByEmail no repositório para isso.
	// Por enquanto, vamos assumir que o DB cuidará da restrição UNIQUE no email.)

	// 3. Gerar ID único para o professor
	teacher.ID = uuid.New().String()

	// 4. Gerar o registro do professor (ex: "COMP-0001", "MATH-0002")
	departmentCode := strings.ToUpper(teacher.Department[:4]) // Ex: "COMP", "MATH"
	if len(teacher.Department) < 4 {
		departmentCode = strings.ToUpper(teacher.Department)
	}

	lastRegistry, err := s.teacherRepo.GetLastRegistryForDepartment(departmentCode)
	if err != nil {
		return fmt.Errorf("erro ao buscar último registro para o departamento: %w", err)
	}

	newSequence := 1
	if lastRegistry != "" {
		parts := strings.Split(lastRegistry, "-")
		if len(parts) == 2 {
			lastSequence, err := strconv.Atoi(parts[1])
			if err == nil {
				newSequence = lastSequence + 1
			} else {
				log.Printf("Aviso: Não foi possível converter sequência '%s' para int no registro. Reiniciando sequência para 1. Erro: %v", parts[1], err)
			}
		}
	}
	teacher.Registry = fmt.Sprintf("%s-%04d", departmentCode, newSequence)

	return s.teacherRepo.CreateTeacher(teacher)
}

// GetTeacherByID busca um professor pelo ID.
func (s *TeacherService) GetTeacherByID(id string) (*models.Teacher, error) {
	teacher, err := s.teacherRepo.GetTeacherByID(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar professor: %w", err)
	}
	if teacher == nil {
		return nil, errors.New("professor não encontrado")
	}
	// TODO: Carregar matérias associadas ao professor se necessário
	return teacher, nil
}

// GetAllTeachers busca todos os professores.
func (s *TeacherService) GetAllTeachers() ([]models.Teacher, error) {
	teachers, err := s.teacherRepo.GetAllTeachers()
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar todos os professores: %w", err)
	}
	// TODO: Carregar matérias associadas a cada professor se necessário
	return teachers, nil
}

// UpdateTeacher atualiza um professor existente.
func (s *TeacherService) UpdateTeacher(teacher *models.Teacher) error {
	if teacher.ID == "" {
		return errors.New("ID do professor é obrigatório para atualização")
	}
	if teacher.Name == "" || teacher.Department == "" || teacher.Email == "" {
		return errors.New("nome, departamento e email do professor são obrigatórios para atualização")
	}

	existingTeacher, err := s.teacherRepo.GetTeacherByID(teacher.ID)
	if err != nil {
		return fmt.Errorf("erro ao buscar professor existente para atualização: %w", err)
	}
	if existingTeacher == nil {
		return fmt.Errorf("professor com ID %s não encontrado para atualização", teacher.ID)
	}

	// Atualiza apenas os campos permitidos
	existingTeacher.Name = teacher.Name
	existingTeacher.Department = teacher.Department
	existingTeacher.Email = teacher.Email
	// O registro (Registry) não deve ser alterado após a criação, então não o copiamos do 'teacher' de entrada.
	// O existingTeacher já tem o Registry correto do DB.

	return s.teacherRepo.UpdateTeacher(existingTeacher)
}

// DeleteTeacher deleta um professor pelo ID.
func (s *TeacherService) DeleteTeacher(id string) error {
	if id == "" {
		return errors.New("ID do professor é obrigatório para exclusão")
	}
	existingTeacher, err := s.teacherRepo.GetTeacherByID(id)
	if err != nil {
		return fmt.Errorf("erro ao verificar professor para exclusão: %w", err)
	}
	if existingTeacher == nil {
		return errors.New("professor não encontrado para exclusão")
	}
	return s.teacherRepo.DeleteTeacher(id)
}

// AddSubjectToTeacherService associa uma matéria a um professor após validações.
func (s *TeacherService) AddSubjectToTeacherService(teacherID, subjectID string) error {
	// 1. Verificar se o professor existe
	teacher, err := s.teacherRepo.GetTeacherByID(teacherID)
	if err != nil {
		log.Printf("AddSubjectToTeacherService: Erro ao buscar professor %s: %v", teacherID, err)
		return err
	}
	if teacher == nil {
		return errors.New("professor não encontrado")
	}

	// 2. Verificar se a matéria existe (usa o subjectRepo)
	subject, err := s.subjectRepo.GetSubjectByID(subjectID)
	if err != nil {
		log.Printf("AddSubjectToTeacherService: Erro ao buscar matéria %s: %v", subjectID, err)
		return err
	}
	if subject == nil {
		return errors.New("matéria não encontrada")
	}

	// 3. Chamar o repositório para fazer a associação
	err = s.teacherRepo.AddSubjectToTeacher(teacherID, subjectID)
	if err != nil {
		log.Printf("AddSubjectToTeacherService: Erro ao adicionar matéria %s ao professor %s no DB: %v", subjectID, teacherID, err)
		return err
	}
	log.Printf("AddSubjectToTeacherService: Matéria %s adicionada ao professor %s com sucesso.", subjectID, teacherID)
	return nil
}

// RemoveSubjectFromTeacherService desassocia uma matéria de um professor após validações.
func (s *TeacherService) RemoveSubjectFromTeacherService(teacherID, subjectID string) error {
	// 1. Verificar se o professor existe
	teacher, err := s.teacherRepo.GetTeacherByID(teacherID)
	if err != nil {
		log.Printf("RemoveSubjectFromTeacherService: Erro ao buscar professor %s: %v", teacherID, err)
		return err
	}
	if teacher == nil {
		return errors.New("professor não encontrado")
	}

	// 2. Verificar se a matéria existe (opcional, mas boa prática)
	subject, err := s.subjectRepo.GetSubjectByID(subjectID)
	if err != nil {
		log.Printf("RemoveSubjectFromTeacherService: Erro ao buscar matéria %s: %v", subjectID, err)
		return err
	}
	if subject == nil {
		return errors.New("matéria não encontrada")
	}

	// 3. Chamar o repositório para remover a associação
	err = s.teacherRepo.RemoveSubjectFromTeacher(teacherID, subjectID)
	if err != nil {
		log.Printf("RemoveSubjectFromTeacherService: Erro ao remover matéria %s do professor %s no DB: %v", subjectID, teacherID, err)
		return err
	}
	log.Printf("RemoveSubjectFromTeacherService: Matéria %s removida do professor %s com sucesso.", subjectID, teacherID)
	return nil
}
