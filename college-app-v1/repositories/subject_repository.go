// repositories/subject_repository.go
package repositories

import (
	"college-app-v1/models"
	"database/sql"
	"fmt" // Importar fmt para usar fmt.Errorf
	"log"

	"github.com/google/uuid" // <-- Adicionar este import!
)

type SubjectRepository struct {
	db *sql.DB
}

func NewSubjectRepository(db *sql.DB) *SubjectRepository {
	return &SubjectRepository{db: db}
}

// CreateSubject insere uma nova matéria no banco de dados.
func (r *SubjectRepository) CreateSubject(subject *models.Subject) error {
	// --- MUDANÇA CRÍTICA AQUI: Gerar o UUID para o ID da matéria ---
	subject.ID = uuid.New().String() // Gera um ID único para a matéria

	query := `INSERT INTO subjects (id, name, year, credits) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, subject.ID, subject.Name, subject.Year, subject.Credits)
	if err != nil {
		log.Printf("CreateSubject: Erro ao executar INSERT para matéria %s (Name: %s, Year: %d): %v", subject.ID, subject.Name, subject.Year, err)
		return fmt.Errorf("falha ao criar matéria no DB: %w", err) // Encapsular o erro
	}
	log.Printf("CreateSubject: Matéria '%s' (ID: %s, Ano: %d) criada com sucesso.", subject.Name, subject.ID, subject.Year)
	return nil
}

// GetSubjectByID busca uma matéria pelo ID.
func (r *SubjectRepository) GetSubjectByID(id string) (*models.Subject, error) {
	subject := &models.Subject{}
	query := `SELECT id, name, year, credits FROM subjects WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&subject.ID, &subject.Name, &subject.Year, &subject.Credits)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("GetSubjectByID: Matéria com ID %s não encontrada no DB.", id)
			return nil, fmt.Errorf("matéria não encontrada") // Retorna erro mais específico
		}
		log.Printf("GetSubjectByID: Erro ao buscar matéria por ID %s: %v", id, err)
		return nil, fmt.Errorf("falha ao buscar matéria por ID: %w", err)
	}
	log.Printf("GetSubjectByID: Matéria '%s' (ID: %s) encontrada.", subject.Name, subject.ID)
	return subject, nil
}

// GetAllSubjects busca todas as matérias.
func (r *SubjectRepository) GetAllSubjects() ([]models.Subject, error) {
	rows, err := r.db.Query(`SELECT id, name, year, credits FROM subjects`)
	if err != nil {
		log.Printf("GetAllSubjects: Erro ao buscar todas as matérias: %v", err)
		return nil, fmt.Errorf("falha ao buscar todas as matérias: %w", err)
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		subject := models.Subject{}
		if err := rows.Scan(&subject.ID, &subject.Name, &subject.Year, &subject.Credits); err != nil {
			log.Printf("GetAllSubjects: Erro ao escanear matéria: %v", err)
			return nil, fmt.Errorf("falha ao escanear dados da matéria: %w", err)
		}
		subjects = append(subjects, subject)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração de matérias: %w", err)
	}
	log.Printf("GetAllSubjects: %d matérias encontradas.", len(subjects))
	return subjects, nil
}

// UpdateSubject atualiza uma matéria existente.
func (r *SubjectRepository) UpdateSubject(subject *models.Subject) error {
	query := `UPDATE subjects SET name = $1, year = $2, credits = $3 WHERE id = $4`
	result, err := r.db.Exec(query, subject.Name, subject.Year, subject.Credits, subject.ID)
	if err != nil {
		log.Printf("UpdateSubject: Erro ao atualizar matéria %s (ID: %s): %v", subject.Name, subject.ID, err)
		return fmt.Errorf("falha ao atualizar matéria: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("falha ao verificar linhas afetadas após atualização: %w", err)
	}
	if rowsAffected == 0 {
		log.Printf("UpdateSubject: Nenhuma matéria encontrada para atualizar com ID %s.", subject.ID)
		return fmt.Errorf("matéria não encontrada para atualização")
	}
	log.Printf("UpdateSubject: Matéria '%s' (ID: %s) atualizada com sucesso.", subject.Name, subject.ID)
	return nil
}

// DeleteSubject deleta uma matéria pelo ID.
func (r *SubjectRepository) DeleteSubject(id string) error {
	query := `DELETE FROM subjects WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("DeleteSubject: Erro ao deletar matéria ID %s: %v", id, err)
		return fmt.Errorf("falha ao deletar matéria: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("falha ao verificar linhas afetadas após exclusão: %w", err)
	}
	if rowsAffected == 0 {
		log.Printf("DeleteSubject: Nenhuma matéria encontrada para deletar com ID %s.", id)
		return fmt.Errorf("matéria não encontrada para exclusão")
	}
	log.Printf("DeleteSubject: Matéria com ID %s deletada com sucesso.", id)
	return nil
}
