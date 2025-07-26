import { useState, useEffect } from 'react';
import './App.css'; // Importa o ficheiro CSS para estilização - Descomente esta linha

function App() {
  const [students, setStudents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Estados para o formulário de novo aluno
  const [newName, setNewName] = useState('');
  const [newEnrollment, setNewEnrollment] = useState('');
  const [newCurrentYear, setNewCurrentYear] = useState('');
  const [newShift, setNewShift] = useState('');
  const [formMessage, setFormMessage] = useState(''); // Para mensagens de sucesso/erro do formulário de criação

  // Estados para edição
  const [editingStudentId, setEditingStudentId] = useState(null); // ID do aluno a ser editado
  const [editName, setEditName] = useState('');
  const [editEnrollment, setEditEnrollment] = useState('');
  const [editCurrentYear, setEditCurrentYear] = useState('');
  const [editShift, setEditShift] = useState('');
  const [editMessage, setEditMessage] = useState(''); // Mensagens para o formulário de edição

  // Função para buscar os alunos (reutilizada)
  const fetchStudents = () => {
    setLoading(true); // Ativa o carregamento ao buscar
    fetch('/api/students') // Alterado de volta para URL relativo
      .then(response => {
        if (!response.ok) {
          // Tenta ler a mensagem de erro da API se a resposta não for OK
          return response.json().then(err => { throw new Error(err.message || 'Erro desconhecido da API'); });
        }
        return response.json();
      })
      .then(data => {
        setStudents(Array.isArray(data) ? data : []);
        setLoading(false);
      })
      .catch(err => {
        setError(err);
        setLoading(false);
        console.error("Erro ao buscar alunos:", err);
      });
  };

  // Chama fetchStudents ao montar o componente
  useEffect(() => {
    fetchStudents();
  }, []);

  // Handler para submissão do formulário de CRIAÇÃO
  const handleSubmit = async (e) => {
    e.preventDefault(); // Previne o comportamento padrão de recarregar a página

    setFormMessage(''); // Limpa mensagens anteriores

    // Validação básica do formulário
    if (!newName || !newEnrollment || !newCurrentYear || !newShift) {
      setFormMessage('Erro: Todos os campos são obrigatórios!'); // Mensagem de erro
      return;
    }

    const studentData = {
      name: newName,
      enrollment: newEnrollment,
      current_year: parseInt(newCurrentYear, 10), // Converte para número inteiro
      shift: newShift,
    };

    try {
      const response = await fetch('/api/students', { // Alterado de volta para URL relativo
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(studentData),
      });

      const result = await response.json();

      if (!response.ok) {
        // Se a resposta não for OK (ex: 400 Bad Request, 500 Internal Server Error)
        throw new Error(result.message || 'Erro ao criar aluno');
      }

      setFormMessage('Sucesso: Aluno criado com sucesso!'); // Mensagem de sucesso
      setNewName('');
      setNewEnrollment('');
      setNewCurrentYear('');
      setNewShift('');
      fetchStudents(); // Recarrega a lista de alunos após a criação
    } catch (err) {
      setFormMessage(`Erro: ${err.message}`); // Mensagem de erro
      console.error("Erro ao enviar formulário:", err);
    }
  };

  // Handler para DELETAR um aluno
  const handleDeleteStudent = async (id) => {
    if (window.confirm('Tem certeza que deseja excluir este aluno?')) {
      try {
        const response = await fetch(`/api/students/${id}`, { // Alterado de volta para URL relativo
          method: 'DELETE',
        });

        if (!response.ok) {
          const errorData = await response.json(); // Tenta ler a mensagem de erro da API
          throw new Error(errorData.message || `Erro ao excluir aluno com ID: ${id}`);
        }

        setFormMessage('Sucesso: Aluno excluído com sucesso!'); // Mensagem de sucesso
        fetchStudents(); // Recarrega a lista
      } catch (err) {
        setFormMessage(`Erro ao excluir: ${err.message}`); // Mensagem de erro
        console.error("Erro ao excluir aluno:", err);
      }
    }
  };

  // Handler para iniciar a EDIÇÃO de um aluno
  const handleEditStudent = (student) => {
    setEditingStudentId(student.id);
    setEditName(student.name);
    setEditEnrollment(student.enrollment);
    setEditCurrentYear(student.current_year.toString()); // Converte para string para o input
    setEditShift(student.shift);
    setEditMessage(''); // Limpa mensagens anteriores
  };

  // Handler para SALVAR a EDIÇÃO de um aluno
  const handleSaveEdit = async (e) => {
    e.preventDefault();
    setEditMessage('');

    if (!editName || !editEnrollment || !editCurrentYear || !editShift) {
      setEditMessage('Erro: Todos os campos são obrigatórios!'); // Mensagem de erro
      return;
    }

    const updatedStudentData = {
      id: editingStudentId, // O ID é importante para a requisição PUT
      name: editName,
      enrollment: editEnrollment,
      current_year: parseInt(editCurrentYear, 10),
      shift: editShift,
    };

    try {
      const response = await fetch(`/api/students/${editingStudentId}`, { // Alterado de volta para URL relativo
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(updatedStudentData),
      });

      const result = await response.json();

      if (!response.ok) {
        throw new Error(result.message || 'Erro ao atualizar aluno');
      }

      setEditMessage('Sucesso: Aluno atualizado com sucesso!'); // Mensagem de sucesso
      setEditingStudentId(null); // Sai do modo de edição
      fetchStudents(); // Recarrega a lista
    } catch (err) {
      setEditMessage(`Erro: ${err.message}`); // Mensagem de erro
      console.error("Erro ao atualizar aluno:", err);
    }
  };

  // Handler para CANCELAR a edição
  const handleCancelEdit = () => {
    setEditingStudentId(null);
    setEditMessage('');
  };

  if (loading) {
    return <div>Carregando alunos...</div>;
  }

  if (error) {
    return <div>Erro ao carregar alunos: {error.message}</div>;
  }

  return (
    <div className="App">
      <h1>Gerenciamento de Alunos</h1>

      {/* Formulário de Criação de Alunos */}
      <h2>Adicionar Novo Aluno</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="name">Nome:</label>
          <input
            type="text"
            id="name"
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
            required
          />
        </div>
        {/* Novos campos agrupados para ficarem na mesma linha */}
        <div className="form-group-inline">
            <div>
                <label htmlFor="enrollment">Matrícula:</label>
                <input
                    type="text"
                    id="enrollment"
                    value={newEnrollment}
                    onChange={(e) => setNewEnrollment(e.target.value)}
                    required
                />
            </div>
            {/* Adicionado um ID específico para o div do Ano Atual para controle de largura */}
            <div id="new-current-year-group">
                <label htmlFor="new_current_year">Ano Atual:</label>
                <input
                    type="number"
                    id="new_current_year"
                    value={newCurrentYear}
                    onChange={(e) => setNewCurrentYear(e.target.value)}
                    required
                />
            </div>
            {/* Adicionado um ID específico para o div do Turno para controle de largura */}
            <div id="new-shift-group">
                <label htmlFor="new_shift">Turno (M/T/N):</label>
                <input
                    type="text"
                    id="new_shift"
                    value={newShift}
                    onChange={(e) => setNewShift(e.target.value)}
                    maxLength="1"
                    required
                />
            </div>
            {/* Botão "Criar Aluno" movido para dentro do form-group-inline */}
            <div className="form-button-container">
                <button type="submit" className="create-button">Criar Aluno</button>
            </div>
        </div>
      </form>

      {/* Exibe mensagens do formulário de criação */}
      {formMessage && (
        <p className={formMessage.startsWith('Erro') ? 'error-message' : 'success-message'}>
          {formMessage}
        </p>
      )}

      <hr />

      {/* Lista de Alunos Existentes */}
      <h2>Lista de Alunos</h2>
      {students.length === 0 ? (
        <p>Nenhum aluno encontrado.</p>
      ) : (
        <ul>
          {students.map(student => (
            <li key={student.id}>
              {editingStudentId === student.id ? (
                // Formulário de Edição
                <form onSubmit={handleSaveEdit}>
                  <h3>Editando Aluno: {student.name}</h3>
                  <div>
                    <label htmlFor={`edit-name-${student.id}`}>Nome:</label>
                    <input
                      type="text"
                      id={`edit-name-${student.id}`}
                      value={editName}
                      onChange={(e) => setEditName(e.target.value)}
                      required
                    />
                  </div>
                  {/* Campos de edição agrupados para ficarem na mesma linha */}
                  <div className="form-group-inline">
                    <div>
                        <label htmlFor={`edit-enrollment-${student.id}`}>Matrícula:</label>
                        <input
                            type="text"
                            id={`edit-enrollment-${student.id}`}
                            value={editEnrollment}
                            onChange={(e) => setEditEnrollment(e.target.value)}
                            required
                        />
                    </div>
                    {/* Adicionado um ID específico para o div do Ano Atual para controle de largura */}
                    <div id={`edit-current-year-group-${student.id}`}>
                        <label htmlFor={`edit-current_year-${student.id}`}>Ano Atual:</label>
                        <input
                            type="number"
                            id={`edit-current_year-${student.id}`}
                            value={editCurrentYear}
                            onChange={(e) => setEditCurrentYear(e.target.value)}
                            required
                        />
                    </div>
                    {/* Adicionado um ID específico para o div do Turno para controle de largura */}
                    <div id={`edit-shift-group-${student.id}`}>
                        <label htmlFor={`edit-shift-${student.id}`}>Turno (M/T/N):</label>
                        <input
                            type="text"
                            id={`edit-shift-${student.id}`}
                            value={editShift}
                            onChange={(e) => setEditShift(e.target.value)}
                            maxLength="1"
                            required
                        />
                    </div>
                    {/* Botões de edição e cancelamento movidos para dentro do form-group-inline */}
                    <div className="form-button-container">
                        <button type="submit" className="edit-button">Salvar Edição</button>
                        <button type="button" className="cancel-button" onClick={handleCancelEdit}>Cancelar</button>
                    </div>
                  </div>
                  {editMessage && (
                    <p className={editMessage.startsWith('Erro') ? 'error-message' : 'success-message'}>
                      {editMessage}
                    </p>
                  )}
                </form>
              ) : (
                // Visualização Normal do Aluno
                <>
                  <strong>{student.name}</strong> ({student.enrollment}) - Ano: {student.current_year}, Turno: {student.shift}
                  <button className="details-button" onClick={() => alert(`Detalhes de ${student.name}:\nID: ${student.id}\nMatrícula: ${student.enrollment}\nAno: ${student.current_year}\nTurno: ${student.shift}`)}>Detalhes</button>
                  <button className="edit-button" onClick={() => handleEditStudent(student)}>Editar</button>
                  <button className="delete-button" onClick={() => handleDeleteStudent(student.id)}>Excluir</button>
                </>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

export default App;