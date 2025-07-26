// config/database.go
package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

var DB *sql.DB

func InitDB() {
	var err error
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("Variável de ambiente DATABASE_URL não definida. Por favor, defina-a.")
	}

	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Erro ao abrir o banco de dados PostgreSQL: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados PostgreSQL: %v", err)
	}

	log.Println("Conexão com o banco de dados PostgreSQL estabelecida com sucesso!")
	createTables()
}

func createTables() {
	// ATUALIZADO: Adicionada a coluna 'shift' e removido 'UNIQUE' de 'enrollment' temporariamente
	// para permitir a geração de matrículas mais flexíveis antes de definir a unicidade composta.
	// A unicidade será garantida pela lógica de geração no serviço.
	createStudentsTableSQL := `
    CREATE TABLE IF NOT EXISTS students (
        id TEXT PRIMARY KEY,
        enrollment TEXT NOT NULL UNIQUE, -- Re-adicionando UNIQUE após ajustar a lógica
        name TEXT NOT NULL,
        current_year INTEGER NOT NULL,
        shift TEXT NOT NULL -- Nova coluna para o turno
    );`

	createSubjectsTableSQL := `
    CREATE TABLE IF NOT EXISTS subjects (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        year INTEGER NOT NULL,
        credits INTEGER NOT NULL
    );`

	createTeachersTableSQL := `
    CREATE TABLE IF NOT EXISTS teachers (
        id TEXT PRIMARY KEY,
        registry TEXT NOT NULL UNIQUE,
        name TEXT NOT NULL,
        department TEXT NOT NULL
    );`

	createStudentSubjectsTableSQL := `
    CREATE TABLE IF NOT EXISTS student_subjects (
        student_id TEXT NOT NULL,
        subject_id TEXT NOT NULL,
        PRIMARY KEY (student_id, subject_id),
        FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
        FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE
    );`

	_, err := DB.Exec(createStudentsTableSQL)
	if err != nil {
		log.Fatalf("Erro ao criar tabela students: %v", err)
	}
	_, err = DB.Exec(createSubjectsTableSQL)
	if err != nil {
		log.Fatalf("Erro ao criar tabela subjects: %v", err)
	}
	_, err = DB.Exec(createTeachersTableSQL)
	if err != nil {
		log.Fatalf("Erro ao criar tabela teachers: %v", err)
	}
	_, err = DB.Exec(createStudentSubjectsTableSQL)
	if err != nil {
		log.Fatalf("Erro ao criar tabela student_subjects: %v", err)
	}

	log.Println("Tabelas verificadas/criadas com sucesso!")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Conexão com o banco de dados PostgreSQL fechada.")
	}
}
