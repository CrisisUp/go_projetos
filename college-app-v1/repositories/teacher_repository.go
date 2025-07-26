package repositories

import (
	"college-app-v1/models"
	"database/sql"
	"log"
	// Certifique-se de que este import está presente se você gerar UUIDs aqui
)

// TeacherRepository define as operações de CRUD para professores.
type TeacherRepository struct {
	db *sql.DB
}

// NewTeacherRepository cria uma nova instância de TeacherRepository.
func NewTeacherRepository(db *sql.DB) *TeacherRepository {
	return &TeacherRepository{db: db}
}

// CreateTeacher insere um novo professor no banco de dados.
// O ID, Registry e Email já devem vir preenchidos do Service.
func (r *TeacherRepository) CreateTeacher(teacher *models.Teacher) error {
	// Incluindo 'email' na query de INSERT
	query := `INSERT INTO teachers (id, registry, name, email, department) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, teacher.ID, teacher.Registry, teacher.Name, teacher.Email, teacher.Department)
	if err != nil {
		log.Printf("CreateTeacher: Erro ao executar INSERT para professor %s: %v", teacher.Name, err)
		return err
	}
	log.Printf("CreateTeacher: Professor '%s' (ID: %s, Registro: %s) criado com sucesso.", teacher.Name, teacher.ID, teacher.Registry)
	return nil
}

// GetTeacherByID busca um professor pelo ID e carrega suas matérias.
func (r *TeacherRepository) GetTeacherByID(id string) (*models.Teacher, error) {
	teacher := &models.Teacher{}
	// Incluindo 'email' na query de SELECT
	query := `SELECT id, registry, name, email, department FROM teachers WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&teacher.ID, &teacher.Registry, &teacher.Name, &teacher.Email, &teacher.Department)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("GetTeacherByID: Professor com ID %s não encontrado no DB.", id)
			return nil, nil
		}
		log.Printf("GetTeacherByID: Erro ao escanear dados do professor com ID %s: %v", id, err)
		return nil, err
	}

	// Carregar matérias associadas ao professor
	subjects, err := r.GetSubjectsByTeacherID(teacher.ID)
	if err != nil {
		log.Printf("GetTeacherByID: Erro ao buscar matérias para o professor %s (ID: %s): %v", teacher.Name, teacher.ID, err)
		return nil, err
	}
	if subjects == nil {
		teacher.Subjects = []models.Subject{} // Inicializa como slice vazio se não houver matérias
	} else {
		teacher.Subjects = subjects
	}

	log.Printf("GetTeacherByID: Professor '%s' (ID: %s) encontrado com sucesso, %d matérias.", teacher.Name, teacher.ID, len(teacher.Subjects))
	return teacher, nil
}

// GetAllTeachers busca todos os professores e carrega suas matérias.
func (r *TeacherRepository) GetAllTeachers() ([]models.Teacher, error) {
	// Incluindo 'email' na query de SELECT
	rows, err := r.db.Query(`SELECT id, registry, name, email, department FROM teachers`)
	if err != nil {
		log.Printf("GetAllTeachers: Erro ao executar SELECT ALL FROM teachers: %v", err)
		return nil, err
	}
	defer rows.Close()

	var teachers []models.Teacher
	for rows.Next() {
		teacher := models.Teacher{}
		// Incluindo 'email' no Scan
		if err := rows.Scan(&teacher.ID, &teacher.Registry, &teacher.Name, &teacher.Email, &teacher.Department); err != nil {
			log.Printf("GetAllTeachers: Erro ao escanear linha de professor do DB: %v", err)
			return nil, err
		}
		// Carregar matérias associadas a cada professor
		subjects, err := r.GetSubjectsByTeacherID(teacher.ID)
		if err != nil {
			log.Printf("GetAllTeachers: Erro ao buscar matérias para professor %s (ID: %s) durante iteração: %v", teacher.Name, teacher.ID, err)
			return nil, err
		}
		if subjects == nil {
			teacher.Subjects = []models.Subject{}
		} else {
			teacher.Subjects = subjects
		}
		teachers = append(teachers, teacher)
	}
	log.Printf("GetAllTeachers: %d professores encontrados e processados.", len(teachers))
	return teachers, nil
}

// UpdateTeacher atualiza um professor existente.
func (r *TeacherRepository) UpdateTeacher(teacher *models.Teacher) error {
	// Incluindo 'email' na query de UPDATE
	query := `UPDATE teachers SET registry = $1, name = $2, email = $3, department = $4 WHERE id = $5`
	result, err := r.db.Exec(query, teacher.Registry, teacher.Name, teacher.Email, teacher.Department, teacher.ID)
	if err != nil {
		log.Printf("UpdateTeacher: Erro ao atualizar professor: %v", err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("UpdateTeacher: Nenhum professor encontrado para atualizar com ID %s.", teacher.ID)
		return sql.ErrNoRows
	}
	log.Printf("UpdateTeacher: Professor '%s' (ID: %s) atualizado com sucesso.", teacher.Name, teacher.ID)
	return nil
}

// DeleteTeacher deleta um professor pelo ID.
func (r *TeacherRepository) DeleteTeacher(id string) error {
	query := `DELETE FROM teachers WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("DeleteTeacher: Erro ao deletar professor: %v", err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("DeleteTeacher: Nenhum professor encontrado para deletar com ID %s.", id)
		return sql.ErrNoRows
	}
	log.Printf("DeleteTeacher: Professor com ID %s deletado com sucesso.", id)
	return nil
}

// GetLastRegistryForDepartment busca o maior número de registro para o departamento especificado.
// Retorna o registro como string e um erro, se houver.
// Retorna "" e nil se não houver registros para o departamento.
func (r *TeacherRepository) GetLastRegistryForDepartment(departmentCode string) (string, error) {
	var lastRegistry sql.NullString
	query := `
		SELECT registry FROM teachers
		WHERE registry LIKE $1 || '-%'
		ORDER BY registry DESC
		LIMIT 1
	`
	err := r.db.QueryRow(query, departmentCode).Scan(&lastRegistry)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Nenhum registro encontrado para este departamento
		}
		log.Printf("GetLastRegistryForDepartment: Erro ao buscar último registro para o departamento %s: %v", departmentCode, err)
		return "", err
	}

	if lastRegistry.Valid {
		return lastRegistry.String, nil
	}
	return "", nil
}

// AddSubjectToTeacher associa uma matéria a um professor.
func (r *TeacherRepository) AddSubjectToTeacher(teacherID, subjectID string) error {
	query := `INSERT INTO teacher_subjects (teacher_id, subject_id) VALUES ($1, $2) ON CONFLICT (teacher_id, subject_id) DO NOTHING`
	_, err := r.db.Exec(query, teacherID, subjectID)
	if err != nil {
		log.Printf("AddSubjectToTeacher: Erro ao executar INSERT para associação professor %s - matéria %s: %v", teacherID, subjectID, err)
		return err
	}
	log.Printf("AddSubjectToTeacher: Associação professor %s - matéria %s criada/existente.", teacherID, subjectID)
	return nil
}

// RemoveSubjectFromTeacher desassocia uma matéria de um professor.
func (r *TeacherRepository) RemoveSubjectFromTeacher(teacherID, subjectID string) error {
	query := `DELETE FROM teacher_subjects WHERE teacher_id = $1 AND subject_id = $2`
	result, err := r.db.Exec(query, teacherID, subjectID)
	if err != nil {
		log.Printf("RemoveSubjectFromTeacher: Erro ao executar DELETE para associação professor %s - matéria %s: %v", teacherID, subjectID, err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("RemoveSubjectFromTeacher: Associação professor %s - matéria %s não encontrada para deletar.", teacherID, subjectID)
		return sql.ErrNoRows
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
		return nil, err
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		subject := models.Subject{}
		if err := rows.Scan(&subject.ID, &subject.Name, &subject.Year, &subject.Credits); err != nil {
			log.Printf("GetSubjectsByTeacherID: Erro ao escanear matéria do professor ID %s: %v", teacherID, err)
			return nil, err
		}
		subjects = append(subjects, subject)
	}
	if subjects == nil {
		log.Printf("GetSubjectsByTeacherID: Nenhuma matéria encontrada para professor ID %s (retornando slice vazio).", teacherID)
		return []models.Subject{}, nil
	}
	log.Printf("GetSubjectsByTeacherID: %d matérias encontradas para professor ID %s.", len(subjects), teacherID)
	return subjects, nil
}
