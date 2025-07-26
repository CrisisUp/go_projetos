-- Remover tabelas existentes para garantir um estado limpo (apenas para desenvolvimento)
DROP TABLE IF EXISTS student_subjects;
DROP TABLE IF EXISTS teacher_subjects;
DROP TABLE IF EXISTS students;
DROP TABLE IF EXISTS subjects;
DROP TABLE IF EXISTS teachers;

-- Tabela de Estudantes
CREATE TABLE students (
    id VARCHAR(255) PRIMARY KEY,
    enrollment VARCHAR(255) UNIQUE NOT NULL, -- Matrícula do aluno, única
    name VARCHAR(255) NOT NULL,
    current_year INT NOT NULL, -- Ano atual do curso (ex: 1, 2, 3)
    shift VARCHAR(50) NOT NULL -- Turno (ex: 'Manhã', 'Tarde', 'Noite')
);

-- Tabela de Matérias
CREATE TABLE subjects (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    year INT NOT NULL, -- Ano em que a matéria é oferecida (ex: 1º ano, 2º ano)
    credits INT NOT NULL
);

-- Tabela de Professores
CREATE TABLE teachers (
    id VARCHAR(255) PRIMARY KEY,
    registry VARCHAR(255) UNIQUE NOT NULL, -- <-- ADICIONADO: Registro único do professor
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    department VARCHAR(255) -- Departamento do professor (ex: "Ciência da Computação")
);

-- Tabela de associação Aluno-Matéria (muitos-para-muitos)
-- Um aluno pode ter várias matérias e uma matéria pode ter vários alunos
CREATE TABLE student_subjects (
    student_id VARCHAR(255) REFERENCES students(id) ON DELETE CASCADE,
    subject_id VARCHAR(255) REFERENCES subjects(id) ON DELETE CASCADE,
    PRIMARY KEY (student_id, subject_id)
);

-- Tabela de associação Professor-Matéria (muitos-para-muitos)
-- Um professor pode lecionar várias matérias e uma matéria pode ser lecionada por vários professores
CREATE TABLE teacher_subjects (
    teacher_id VARCHAR(255) REFERENCES teachers(id) ON DELETE CASCADE,
    subject_id VARCHAR(255) REFERENCES subjects(id) ON DELETE CASCADE,
    PRIMARY KEY (teacher_id, subject_id)
);

-- Índices para melhor performance em colunas frequentemente usadas em buscas ou junções
CREATE INDEX idx_students_enrollment ON students(enrollment);
CREATE INDEX idx_subjects_name ON subjects(name);
CREATE INDEX idx_teachers_email ON teachers(email);
CREATE INDEX idx_student_subjects_student_id ON student_subjects(student_id);
CREATE INDEX idx_student_subjects_subject_id ON student_subjects(subject_id);
CREATE INDEX idx_teacher_subjects_teacher_id ON teacher_subjects(teacher_id);
CREATE INDEX idx_teacher_subjects_subject_id ON teacher_subjects(subject_id);