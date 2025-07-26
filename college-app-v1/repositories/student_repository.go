// repositories/student_repository.go
package repositories

import (
	"college-app-v1/models" // Certifique-se de que este caminho está correto
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
)

// StudentRepository define as operações de CRUD para alunos.
// A interface será ajustada no service para ser mais clara.
type StudentRepository struct {
	db *sql.DB
}

// NewStudentRepository cria uma nova instância de StudentRepository.
func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

// CreateStudent insere um novo aluno no banco de dados.
func (r *StudentRepository) CreateStudent(student *models.Student) error {
	student.ID = uuid.New().String() // Gera um ID único para o aluno
	query := `INSERT INTO students (id, enrollment, name, current_year, shift) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, student.ID, student.Enrollment, student.Name, student.CurrentYear, student.Shift)
	if err != nil {
		log.Printf("CreateStudent: Erro ao executar INSERT para aluno %s: %v", student.Name, err)
		return fmt.Errorf("falha ao criar aluno: %w", err) // Retorna erro encapsulado
	}

	// Insere as matérias do aluno na tabela de relacionamento
	if student.Subjects != nil {
		for _, subject := range student.Subjects {
			err := r.AddSubjectToStudent(student.ID, subject.ID) // student.ID é string, subject.ID é string
			if err != nil {
				log.Printf("CreateStudent: Aviso - Erro ao adicionar matéria %s ao aluno %s: %v", subject.ID, student.ID, err)
			}
		}
	}
	log.Printf("CreateStudent: Aluno %s (%s) criado com sucesso. Matrícula: %s", student.Name, student.ID, student.Enrollment)
	return nil
}

// GetStudentByID busca um aluno pelo ID.
func (r *StudentRepository) GetStudentByID(id string) (*models.Student, error) {
	student := &models.Student{}
	query := `SELECT id, enrollment, name, current_year, shift FROM students WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&student.ID, &student.Enrollment, &student.Name, &student.CurrentYear, &student.Shift)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("GetStudentByID: Aluno com ID %s não encontrado no DB.", id)
			return nil, fmt.Errorf("aluno não encontrado") // Retorna erro específico
		}
		log.Printf("GetStudentByID: Erro ao escanear dados do aluno com ID %s: %v", id, err)
		return nil, fmt.Errorf("falha ao buscar aluno por ID: %w", err) // Retorna erro encapsulado
	}

	subjects, err := r.GetSubjectsByStudentID(student.ID)
	if err != nil {
		log.Printf("GetStudentByID: Erro ao buscar matérias para o aluno %s (ID: %s): %v", student.Name, student.ID, err)
		return nil, fmt.Errorf("falha ao buscar matérias associadas: %w", err)
	}
	if subjects == nil {
		student.Subjects = []models.Subject{}
	} else {
		student.Subjects = subjects
	}
	log.Printf("GetStudentByID: Aluno '%s' (ID: %s) encontrado com sucesso, %d matérias.", student.Name, student.ID, len(student.Subjects))
	return student, nil
}

// GetAllStudents busca todos os alunos, com opções de filtro.
// year: ponteiro para int para permitir nil (sem filtro de ano)
// shift: string para o turno (vazio significa sem filtro de turno)
func (r *StudentRepository) GetAllStudents(year *int, shift string) ([]models.Student, error) {
	baseQuery := `SELECT id, enrollment, name, current_year, shift FROM students WHERE 1=1`
	args := []interface{}{}
	argCounter := 1

	// Adiciona filtro por ano
	if year != nil {
		baseQuery += fmt.Sprintf(" AND current_year = $%d", argCounter)
		args = append(args, *year)
		argCounter++
	}

	// Adiciona filtro por turno (case-insensitive)
	if shift != "" {
		baseQuery += fmt.Sprintf(" AND LOWER(shift) = LOWER($%d)", argCounter)
		args = append(args, shift)
		argCounter++
	}

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		log.Printf("GetAllStudents: Erro ao executar query com filtros '%s' %v: %v", baseQuery, args, err)
		return nil, fmt.Errorf("falha ao buscar alunos com filtros: %w", err)
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		student := models.Student{}
		if err := rows.Scan(&student.ID, &student.Enrollment, &student.Name, &student.CurrentYear, &student.Shift); err != nil {
			log.Printf("GetAllStudents: Erro ao escanear linha de aluno do DB: %v", err)
			return nil, fmt.Errorf("falha ao escanear dados do aluno: %w", err)
		}
		// Buscar matérias para cada aluno (mantendo o N+1 por enquanto)
		subjects, err := r.GetSubjectsByStudentID(student.ID)
		if err != nil {
			log.Printf("GetAllStudents: Erro ao buscar matérias para aluno %s (ID: %s) durante iteração: %v", student.Name, student.ID, err)
			// Decide como lidar com este erro. Pode ser fatal ou apenas logar e continuar.
			// Por enquanto, vamos retornar o erro, mas em um cenário real você poderia logar
			// e talvez retornar um erro específico para o cliente se matérias forem cruciais.
			return nil, fmt.Errorf("falha ao buscar matérias associadas ao aluno %s: %w", student.ID, err)
		}
		if subjects == nil {
			student.Subjects = []models.Subject{}
		} else {
			student.Subjects = subjects
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração de alunos: %w", err)
	}

	log.Printf("GetAllStudents: %d alunos encontrados e processados com filtros (Ano: %v, Turno: %s).", len(students), year, shift)
	return students, nil
}

// UpdateStudent atualiza um aluno existente.
func (r *StudentRepository) UpdateStudent(student *models.Student) error {
	query := `UPDATE students SET enrollment = $1, name = $2, current_year = $3, shift = $4 WHERE id = $5`
	result, err := r.db.Exec(query, student.Enrollment, student.Name, student.CurrentYear, student.Shift, student.ID)
	if err != nil {
		log.Printf("UpdateStudent: Erro ao executar UPDATE para aluno %s (ID: %s): %v", student.Name, student.ID, err)
		return fmt.Errorf("falha ao atualizar aluno: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("falha ao verificar linhas afetadas após atualização: %w", err)
	}
	if rowsAffected == 0 {
		log.Printf("UpdateStudent: Nenhum aluno encontrado para atualizar com ID %s.", student.ID)
		return fmt.Errorf("aluno não encontrado para atualização") // Retorna erro mais descritivo
	}
	log.Printf("UpdateStudent: Aluno %s (ID: %s) atualizado com sucesso.", student.Name, student.ID)
	return nil
}

// DeleteStudent deleta um aluno pelo ID.
func (r *StudentRepository) DeleteStudent(id string) error {
	query := `DELETE FROM students WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("DeleteStudent: Erro ao executar DELETE para aluno ID %s: %v", id, err)
		return fmt.Errorf("falha ao deletar aluno: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("falha ao verificar linhas afetadas após exclusão: %w", err)
	}
	if rowsAffected == 0 {
		log.Printf("DeleteStudent: Nenhum aluno encontrado para deletar com ID %s.", id)
		return fmt.Errorf("aluno não encontrado para exclusão") // Retorna erro mais descritivo
	}
	log.Printf("DeleteStudent: Aluno com ID %s deletado com sucesso.", id)
	return nil
}

// AddSubjectToStudent associa uma matéria a um aluno.
func (r *StudentRepository) AddSubjectToStudent(studentID, subjectID string) error {
	query := `INSERT INTO student_subjects (student_id, subject_id) VALUES ($1, $2) ON CONFLICT (student_id, subject_id) DO NOTHING`
	_, err := r.db.Exec(query, studentID, subjectID)
	if err != nil {
		log.Printf("AddSubjectToStudent: Erro ao executar INSERT para associação aluno %s - matéria %s: %v", studentID, subjectID, err)
		return fmt.Errorf("falha ao associar matéria ao aluno: %w", err)
	}
	log.Printf("AddSubjectToStudent: Associação aluno %s - matéria %s criada/existente.", studentID, subjectID)
	return nil
}

// RemoveSubjectFromStudent desassocia uma matéria de um aluno.
func (r *StudentRepository) RemoveSubjectFromStudent(studentID, subjectID string) error {
	query := `DELETE FROM student_subjects WHERE student_id = $1 AND subject_id = $2`
	result, err := r.db.Exec(query, studentID, subjectID)
	if err != nil {
		log.Printf("RemoveSubjectFromStudent: Erro ao executar DELETE para associação aluno %s - matéria %s: %v", studentID, subjectID, err)
		return fmt.Errorf("falha ao desassociar matéria do aluno: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("falha ao verificar linhas afetadas após desassociação: %w", err)
	}
	if rowsAffected == 0 {
		log.Printf("RemoveSubjectFromStudent: Associação aluno %s - matéria %s não encontrada para deletar.", studentID, subjectID)
		return fmt.Errorf("associação não encontrada para desassociação") // Retorna erro mais descritivo
	}
	log.Printf("RemoveSubjectFromStudent: Associação aluno %s - matéria %s deletada com sucesso.", studentID, subjectID)
	return nil
}

// GetLastEnrollmentForYearAndShift busca a maior matrícula para o ano e turno especificados.
func (r *StudentRepository) GetLastEnrollmentForYearAndShift(year int, studentShift string) (string, error) {
	var lastEnrollment sql.NullString // Usar sql.NullString para lidar com NULL do DB
	query := `
		SELECT enrollment FROM students
		WHERE enrollment LIKE $1 || $2 || '%'
		ORDER BY enrollment DESC
		LIMIT 1
	`
	err := r.db.QueryRow(query, fmt.Sprintf("%d", year), studentShift).Scan(&lastEnrollment)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("GetLastEnrollmentForYearAndShift: Nenhuma matrícula encontrada para o ano %d e turno %s.", year, studentShift)
			return "", nil // Nenhuma matrícula encontrada para este ano e turno
		}
		log.Printf("GetLastEnrollmentForYearAndShift: Erro ao buscar última matrícula para o ano %d e turno %s: %v", year, studentShift, err)
		return "", fmt.Errorf("falha ao buscar última matrícula: %w", err)
	}

	if lastEnrollment.Valid {
		log.Printf("GetLastEnrollmentForYearAndShift: Última matrícula encontrada para %d%s: %s", year, studentShift, lastEnrollment.String)
		return lastEnrollment.String, nil
	}
	log.Printf("GetLastEnrollmentForYearAndShift: Matrícula nula/inválida para o ano %d e turno %s.", year, studentShift)
	return "", nil // Caso a string seja nula (não deveria acontecer com LIMIT 1)
}

// GetSubjectsByStudentID busca todas as matérias associadas a um aluno.
func (r *StudentRepository) GetSubjectsByStudentID(studentID string) ([]models.Subject, error) {
	query := `
	SELECT s.id, s.name, s.year, s.credits
	FROM subjects s
	JOIN student_subjects ss ON s.id = ss.subject_id
	WHERE ss.student_id = $1`
	rows, err := r.db.Query(query, studentID)
	if err != nil {
		log.Printf("GetSubjectsByStudentID: Erro ao executar query para aluno ID %s: %v", studentID, err)
		return nil, fmt.Errorf("falha ao buscar matérias por aluno: %w", err)
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		subject := models.Subject{}
		if err := rows.Scan(&subject.ID, &subject.Name, &subject.Year, &subject.Credits); err != nil {
			log.Printf("GetSubjectsByStudentID: Erro ao escanear matéria do aluno ID %s: %v", studentID, err)
			return nil, fmt.Errorf("falha ao escanear matéria: %w", err)
		}
		subjects = append(subjects, subject)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração de matérias: %w", err)
	}

	if subjects == nil {
		log.Printf("GetSubjectsByStudentID: Nenhuma matéria encontrada para aluno ID %s (retornando slice vazio).", studentID)
		return []models.Subject{}, nil
	}
	log.Printf("GetSubjectsByStudentID: %d matérias encontradas para aluno ID %s.", len(subjects), studentID)
	return subjects, nil
}
