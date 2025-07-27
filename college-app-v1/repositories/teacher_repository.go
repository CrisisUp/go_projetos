// repositories/teacher_repository.go

package repositories

import (
	"college-app-v1/models" // Certifique-se de que este caminho está correto
	"database/sql"
	"fmt"
	"log"

	// Adicionado para strings.ToLower
	"github.com/google/uuid" // Adicionar este import se ainda não estiver
)

// TeacherRepository define a interface para operações de persistência de professor.
type TeacherRepository struct {
	db *sql.DB
}

// NewTeacherRepository cria uma nova instância de TeacherRepository.
func NewTeacherRepository(db *sql.DB) *TeacherRepository {
	return &TeacherRepository{db: db}
}

// CreateTeacher insere um novo professor no banco de dados.
// Assumimos que o ID é gerado aqui.
func (r *TeacherRepository) CreateTeacher(teacher *models.Teacher) error {
	teacher.ID = uuid.New().String() // Gera um ID único para o professor
	query := `INSERT INTO teachers (id, name, department, email) VALUES ($1, $2, $3, $4) RETURNING id`
	_, err := r.db.Exec(query, teacher.ID, teacher.Name, teacher.Department, teacher.Email)
	if err != nil {
		log.Printf("CreateTeacher: Erro ao executar INSERT para professor %s: %v", teacher.Name, err)
		return fmt.Errorf("falha ao criar professor no DB: %w", err)
	}
	log.Printf("CreateTeacher: Professor '%s' (ID: %s, Dept: %s, Email: %s) criado com sucesso.", teacher.Name, teacher.ID, teacher.Department, teacher.Email)
	return nil
}

// GetTeacherByID busca um professor pelo ID, incluindo matérias associadas.
func (r *TeacherRepository) GetTeacherByID(id string) (*models.Teacher, error) {
	var teacher models.Teacher
	query := `SELECT id, name, department, email FROM teachers WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&teacher.ID, &teacher.Name, &teacher.Department, &teacher.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("GetTeacherByID: Professor com ID %s não encontrado no DB.", id)
			return nil, fmt.Errorf("professor não encontrado")
		}
		log.Printf("GetTeacherByID: Erro ao buscar professor por ID %s: %v", id, err)
		return nil, fmt.Errorf("falha ao buscar professor por ID: %w", err)
	}

	// Buscar matérias para este professor
	subjects, err := r.GetSubjectsByTeacherID(teacher.ID) // Assumindo que você tem essa função
	if err != nil {
		log.Printf("GetTeacherByID: Erro ao buscar matérias para o professor %s (ID: %s): %v", teacher.Name, teacher.ID, err)
		return nil, fmt.Errorf("falha ao buscar matérias associadas: %w", err)
	}
	teacher.Subjects = subjects // Atribui as matérias
	log.Printf("GetTeacherByID: Professor '%s' (ID: %s) encontrado com sucesso, %d matérias.", teacher.Name, teacher.ID, len(teacher.Subjects))
	return &teacher, nil
}

// GetAllTeachers busca todos os professores com filtros.
// nameFilter, departmentFilter, emailFilter: strings vazias significam sem filtro.
func (r *TeacherRepository) GetAllTeachers(nameFilter, departmentFilter, emailFilter string) ([]models.Teacher, error) {
	baseQuery := `SELECT id, name, department, email FROM teachers WHERE 1=1`
	args := []interface{}{}
	argCounter := 1

	if nameFilter != "" {
		baseQuery += fmt.Sprintf(" AND LOWER(name) LIKE LOWER($%d)", argCounter)
		args = append(args, "%"+nameFilter+"%") // % para LIKE
		argCounter++
	}
	if departmentFilter != "" {
		baseQuery += fmt.Sprintf(" AND LOWER(department) LIKE LOWER($%d)", argCounter)
		args = append(args, "%"+departmentFilter+"%")
		argCounter++
	}
	if emailFilter != "" {
		baseQuery += fmt.Sprintf(" AND LOWER(email) LIKE LOWER($%d)", argCounter)
		args = append(args, "%"+emailFilter+"%")
		argCounter++
	}

	// --- ADICIONE ESTES LOGS TEMPORÁRIOS PARA DEPURAR (REMOVER EM PRODUÇÃO) ---
	log.Printf("GetAllTeachers Repository: Query SQL final: '%s'", baseQuery)
	log.Printf("GetAllTeachers Repository: Argumentos SQL: %v", args)
	// --- FIM DOS LOGS TEMPORÁRIOS ---

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		log.Printf("GetAllTeachers: Erro ao executar query com filtros: %v", err)
		return nil, fmt.Errorf("falha ao buscar professores com filtros: %w", err)
	}
	defer rows.Close()

	var teachers []models.Teacher
	for rows.Next() {
		var t models.Teacher
		if err := rows.Scan(&t.ID, &t.Name, &t.Department, &t.Email); err != nil {
			log.Printf("GetAllTeachers: Erro ao escanear professor: %v", err)
			return nil, fmt.Errorf("falha ao escanear dados do professor: %w", err)
		}
		// Buscar matérias para cada professor (abordagem N+1 - pode ser otimizada com JOINs)
		subjects, err := r.GetSubjectsByTeacherID(t.ID) // Assumindo que você tem essa função
		if err != nil {
			log.Printf("GetAllTeachers: Erro ao buscar matérias para professor %s (ID: %s) durante iteração: %v", t.Name, t.ID, err)
			return nil, fmt.Errorf("falha ao buscar matérias associadas ao professor %s: %w", t.ID, err)
		}
		t.Subjects = subjects
		teachers = append(teachers, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração de professores: %w", err)
	}
	log.Printf("GetAllTeachers: %d professores encontrados e processados com filtros (Nome: '%s', Dept: '%s', Email: '%s').", len(teachers), nameFilter, departmentFilter, emailFilter)
	return teachers, nil
}

// UpdateTeacher atualiza um professor existente.
func (r *TeacherRepository) UpdateTeacher(teacher *models.Teacher) error {
	query := `UPDATE teachers SET name = $1, department = $2, email = $3 WHERE id = $4`
	res, err := r.db.Exec(query, teacher.Name, teacher.Department, teacher.Email, teacher.ID)
	if err != nil {
		log.Printf("UpdateTeacher: Erro ao atualizar professor %s (ID: %s): %v", teacher.Name, teacher.ID, err)
		return fmt.Errorf("falha ao atualizar professor: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("falha ao verificar linhas afetadas após atualização: %w", err)
	}
	if rowsAffected == 0 {
		log.Printf("UpdateTeacher: Nenhum professor encontrado para atualizar com ID %s.", teacher.ID)
		return fmt.Errorf("professor não encontrado para atualização")
	}
	log.Printf("UpdateTeacher: Professor '%s' (ID: %s) atualizado com sucesso.", teacher.Name, teacher.ID)
	return nil
}

// DeleteTeacher deleta um professor pelo ID.
func (r *TeacherRepository) DeleteTeacher(id string) error {
	query := `DELETE FROM teachers WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("DeleteTeacher: Erro ao deletar professor ID %s: %v", id, err)
		return fmt.Errorf("falha ao deletar professor: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("falha ao verificar linhas afetadas após exclusão: %w", err)
	}
	if rowsAffected == 0 {
		log.Printf("DeleteTeacher: Nenhum professor encontrado para deletar com ID %s.", id)
		return fmt.Errorf("professor não encontrado para exclusão")
	}
	log.Printf("DeleteTeacher: Professor com ID %s deletado com sucesso.", id)
	return nil
}

// AddSubjectToTeacher associa uma matéria a um professor (tabela teacher_subjects).
func (r *TeacherRepository) AddSubjectToTeacher(teacherID, subjectID string) error {
	query := `INSERT INTO teacher_subjects (teacher_id, subject_id) VALUES ($1, $2) ON CONFLICT (teacher_id, subject_id) DO NOTHING`
	_, err := r.db.Exec(query, teacherID, subjectID)
	if err != nil {
		log.Printf("AddSubjectToTeacher: Erro ao executar INSERT para associação professor %s - matéria %s: %v", teacherID, subjectID, err)
		return fmt.Errorf("falha ao associar matéria ao professor: %w", err)
	}
	log.Printf("AddSubjectToTeacher: Associação professor %s - matéria %s criada/existente.", teacherID, subjectID)
	return nil
}

// RemoveSubjectFromTeacher desassocia uma matéria de um professor.
func (r *TeacherRepository) RemoveSubjectFromTeacher(teacherID, subjectID string) error {
	query := `DELETE FROM teacher_subjects WHERE teacher_id = $1 AND subject_id = $2`
	res, err := r.db.Exec(query, teacherID, subjectID)
	if err != nil {
		log.Printf("RemoveSubjectFromTeacher: Erro ao executar DELETE para associação professor %s - matéria %s: %v", teacherID, subjectID, err)
		return fmt.Errorf("falha ao desassociar matéria do professor: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("falha ao verificar linhas afetadas após desassociação: %w", err)
	}
	if rowsAffected == 0 {
		log.Printf("RemoveSubjectFromTeacher: Associação professor %s - matéria %s não encontrada para deletar.", teacherID, subjectID)
		return fmt.Errorf("associação não encontrada para desassociação")
	}
	log.Printf("RemoveSubjectFromTeacher: Associação professor %s - matéria %s deletada com sucesso.", teacherID, subjectID)
	return nil
}

// GetSubjectsByTeacherID busca todas as matérias associadas a um professor.
func (r *TeacherRepository) GetSubjectsByTeacherID(teacherID string) ([]models.Subject, error) {
	query := `
	SELECT s.id, s.name, s.year, s.credits
	FROM subjects s
	JOIN teacher_subjects ts ON s.id = ts.subject_id
	WHERE ts.teacher_id = $1`
	rows, err := r.db.Query(query, teacherID)
	if err != nil {
		log.Printf("GetSubjectsByTeacherID: Erro ao executar query para professor ID %s: %v", teacherID, err)
		return nil, fmt.Errorf("falha ao buscar matérias por professor: %w", err)
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		subject := models.Subject{}
		if err := rows.Scan(&subject.ID, &subject.Name, &subject.Year, &subject.Credits); err != nil {
			log.Printf("GetSubjectsByTeacherID: Erro ao escanear matéria do professor ID %s: %v", teacherID, err)
			return nil, fmt.Errorf("falha ao escanear matéria: %w", err)
		}
		subjects = append(subjects, subject)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração de matérias: %w", err)
	}

	if subjects == nil {
		log.Printf("GetSubjectsByTeacherID: Nenhuma matéria encontrada para professor ID %s (retornando slice vazio).", teacherID)
		return []models.Subject{}, nil
	}
	log.Printf("GetSubjectsByTeacherID: %d matérias encontradas para professor ID %s.", len(subjects), teacherID)
	return subjects, nil
}
