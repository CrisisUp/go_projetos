package main // Mudar para 'main' para ser um executável

import (
	"log"
	"net/http"
	"os" // Adicionar para obter a porta do ambiente

	// Corrigir os caminhos dos imports para o nome exato do seu módulo
	"college-app-v1/config"
	"college-app-v1/handlers"
	"college-app-v1/repositories"
	"college-app-v1/services"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// O roteador Mux precisa ser uma variável global ou ser inicializado uma vez
// para que não seja re-inicializado em cada invocação da função serverless.
var router *mux.Router
var initOnce bool = false // Flag para garantir que a inicialização ocorra apenas uma vez

// Handler é a função de entrada para a Vercel Function.
func Handler(w http.ResponseWriter, r *http.Request) {
	// Inicializa o roteador e as dependências APENAS UMA VEZ
	if !initOnce {
		initAPI()
		initOnce = true
	}
	// Servir a requisição usando o roteador inicializado
	router.ServeHTTP(w, r)
}

// initAPI inicializa todas as dependências da aplicação
func initAPI() {
	// A DATABASE_URL será definida via variável de ambiente da Vercel.
	// NOTA: Certifique-se de que config.InitDB() lida com a conexão ao banco de dados.
	config.InitDB()
	// NOTE: defer config.CloseDB() não é usado em Serverless Functions
	// A conexão é mantida viva pela plataforma entre invocações.

	log.Println("Backend da universidade inicializando para Vercel Function...")

	// --- Inicializando Repositórios e Serviços ---
	// Certifique-se de que estas funções existem nos seus respectivos pacotes
	// e que elas aceitam as dependências corretas (ex: conexão DB)
	subjectRepo := repositories.NewSubjectRepository(config.DB) // Exemplo: passando a conexão do DB
	studentRepo := repositories.NewStudentRepository(config.DB)
	teacherRepo := repositories.NewTeacherRepository(config.DB)

	subjectService := services.NewSubjectService(subjectRepo)
	studentService := services.NewStudentService(studentRepo, subjectRepo)
	teacherService := services.NewTeacherService(teacherRepo, subjectRepo)

	// --- Inicializando Handlers ---
	subjectHandler := handlers.NewSubjectHandler(subjectService)
	studentHandler := handlers.NewStudentHandler(studentService)
	teacherHandler := handlers.NewTeacherHandler(teacherService)

	// --- Configurando o Roteador Mux ---
	router = mux.NewRouter()

	// Rotas para Matérias
	router.HandleFunc("/subjects", subjectHandler.CreateSubjectHandler).Methods("POST")
	router.HandleFunc("/subjects", subjectHandler.GetAllSubjectsHandler).Methods("GET")
	router.HandleFunc("/subjects/{id}", subjectHandler.GetSubjectByIDHandler).Methods("GET")
	router.HandleFunc("/subjects/{id}", subjectHandler.UpdateSubjectHandler).Methods("PUT")
	router.HandleFunc("/subjects/{id}", subjectHandler.DeleteSubjectHandler).Methods("DELETE")

	// Rotas para Alunos
	router.HandleFunc("/students", studentHandler.CreateStudentHandler).Methods("POST")
	router.HandleFunc("/students", studentHandler.GetAllStudentsHandler).Methods("GET")
	router.HandleFunc("/students/{id}", studentHandler.GetStudentByIDHandler).Methods("GET")
	router.HandleFunc("/students/{id}", studentHandler.UpdateStudentHandler).Methods("PUT")
	router.HandleFunc("/students/{id}", studentHandler.DeleteStudentHandler).Methods("DELETE")

	// Rotas para associação Aluno-Matéria
	router.HandleFunc("/students/{studentID}/subjects/{subjectID}", studentHandler.AddSubjectToStudentHandler).Methods("POST")
	router.HandleFunc("/students/{studentID}/subjects/{subjectID}", studentHandler.RemoveSubjectFromStudentHandler).Methods("DELETE")

	// --- ROTAS PARA PROFESSORES ---
	router.HandleFunc("/teachers", teacherHandler.CreateTeacherHandler).Methods("POST")
	router.HandleFunc("/teachers", teacherHandler.GetAllTeachersHandler).Methods("GET")
	router.HandleFunc("/teachers/{id}", teacherHandler.GetTeacherByIDHandler).Methods("GET")
	router.HandleFunc("/teachers/{id}", teacherHandler.UpdateTeacherHandler).Methods("PUT")
	router.HandleFunc("/teachers/{id}", teacherHandler.DeleteTeacherHandler).Methods("DELETE")

	// Rotas para associação Professor-Matéria
	router.HandleFunc("/teachers/{teacherID}/subjects/{subjectID}", teacherHandler.AddSubjectToTeacherHandler).Methods("POST")
	router.HandleFunc("/teachers/{teacherID}/subjects/{subjectID}", teacherHandler.RemoveSubjectFromTeacherHandler).Methods("DELETE")

	// --- Configuração do CORS ---
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Em produção, ajuste para domínios específicos
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"}, // Adicione outros headers se necessário
		AllowCredentials: true,
		Debug:            false, // Defina como false em produção
	})

	// Aplica o middleware CORS ao seu roteador
	router.Use(corsHandler.Handler) // Use diretamente corsHandler.Handler

	log.Println("Backend da universidade inicializado com sucesso para Vercel Function!")
}

// Adicionando uma função main() para testar localmente (opcional)
// Esta função NÃO será executada pela Vercel. A Vercel executará 'Handler'.
func main() {
	initAPI() // Inicializa a API

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Porta padrão se não estiver definida
	}
	addr := ":" + port

	log.Printf("Servidor da College App iniciado localmente na porta %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
