import { useState, useEffect } from 'react';
import './App.css'; 
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faInfoCircle, faEdit, faTrashAlt, faPlus, faSave, faTimes } from '@fortawesome/free-solid-svg-icons'; // Importe os ícones que você vai usar

function App() {
  const [students, setStudents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Estados para o formulário de novo aluno
  const [newName, setNewName] = useState('');
  const [newEnrollment, setNewEnrollment] = useState('');
  const [newCurrentYear, setNewCurrentYear] = useState('');
  const [newShift, setNewShift] = useState('');
  const [formMessage, setFormMessage] = useState('');

  // Estados para edição
  const [editingStudentId, setEditingStudentId] = useState(null);
  const [editName, setEditName] = useState('');
  const [editEnrollment, setEditEnrollment] = useState('');
  const [editCurrentYear, setEditCurrentYear] = useState('');
  const [editShift, setEditShift] = useState('');
  const [editMessage, setEditMessage] = useState('');

  const fetchStudents = () => {
    setLoading(true);
    fetch('/api/students')
      .then(response => {
        if (!response.ok) {
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

  useEffect(() => {
    fetchStudents();
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();

    setFormMessage('');

    if (!newName || !newEnrollment || !newCurrentYear || !newShift) {
      setFormMessage('Erro: Todos os campos são obrigatórios!');
      return;
    }

    const studentData = {
      name: newName,
      enrollment: newEnrollment,
      current_year: parseInt(newCurrentYear, 10),
      shift: newShift,
    };

    try {
      const response = await fetch('/api/students', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(studentData),
      });

      const result = await response.json();

      if (!response.ok) {
        throw new Error(result.message || 'Erro ao criar aluno');
      }

      setFormMessage('Sucesso: Aluno criado com sucesso!');
      setNewName('');
      setNewEnrollment('');
      setNewCurrentYear('');
      setNewShift('');
      fetchStudents();
    } catch (err) {
      setFormMessage(`Erro: ${err.message}`);
      console.error("Erro ao enviar formulário:", err);
    }
  };

  const handleDeleteStudent = async (id) => {
    if (window.confirm('Tem certeza que deseja excluir este aluno?')) {
      try {
        const response = await fetch(`/api/students/${id}`, {
          method: 'DELETE',
        });

        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.message || `Erro ao excluir aluno com ID: ${id}`);
        }

        setFormMessage('Sucesso: Aluno excluído com sucesso!');
        fetchStudents();
      } catch (err) {
        setFormMessage(`Erro ao excluir: ${err.message}`);
        console.error("Erro ao excluir aluno:", err);
      }
    }
  };

  const handleEditStudent = (student) => {
    setEditingStudentId(student.id);
    setEditName(student.name);
    setEditEnrollment(student.enrollment);
    setEditCurrentYear(student.current_year.toString());
    setEditShift(student.shift);
    setEditMessage('');
  };

  const handleSaveEdit = async (e) => {
    e.preventDefault();
    setEditMessage('');

    if (!editName || !editEnrollment || !editCurrentYear || !editShift) {
      setEditMessage('Erro: Todos os campos são obrigatórios!');
      return;
    }

    const updatedStudentData = {
      id: editingStudentId,
      name: editName,
      enrollment: editEnrollment,
      current_year: parseInt(editCurrentYear, 10),
      shift: editShift,
    };

    try {
      const response = await fetch(`/api/students/${editingStudentId}`, {
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

      setEditMessage('Sucesso: Aluno atualizado com sucesso!');
      setEditingStudentId(null);
      fetchStudents();
    } catch (err) {
      setEditMessage(`Erro: ${err.message}`);
      console.error("Erro ao atualizar aluno:", err);
    }
  };

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
          <input
            type="text"
            id="name"
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
            placeholder="Nome:"
            required
          />
        </div>
        <div className="form-group-inline">
            <div>
                <input
                    type="text"
                    id="enrollment"
                    value={newEnrollment}
                    onChange={(e) => setNewEnrollment(e.target.value)}
                    placeholder="Matrícula:"
                    required
                />
            </div>
            <div id="new-current-year-group">
                <input
                    type="number"
                    id="new_current_year"
                    value={newCurrentYear}
                    onChange={(e) => setNewCurrentYear(e.target.value)}
                    placeholder="Ano Atual:"
                    required
                />
            </div>
            <div id="new-shift-group">
                <input
                    type="text"
                    id="new_shift"
                    value={newShift}
                    onChange={(e) => setNewShift(e.target.value)}
                    placeholder="Turno (M/T/N):"
                    maxLength="1"
                    required
                />
            </div>
            <div className="form-button-container">
                <button type="submit" className="create-button">
                    <FontAwesomeIcon icon={faPlus} /> Criar Aluno {/* Ícone para Criar */}
                </button>
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
                    <input
                      type="text"
                      id={`edit-name-${student.id}`}
                      value={editName}
                      onChange={(e) => setEditName(e.target.value)}
                      placeholder="Nome:"
                      required
                    />
                  </div>
                  <div className="form-group-inline">
                    <div>
                        <input
                            type="text"
                            id={`edit-enrollment-${student.id}`}
                            value={editEnrollment}
                            onChange={(e) => setEditEnrollment(e.target.value)}
                            placeholder="Matrícula:"
                            required
                        />
                    </div>
                    <div id={`edit-current-year-group-${student.id}`}>
                        <input
                            type="number"
                            id={`edit-current_year-${student.id}`}
                            value={editCurrentYear}
                            onChange={(e) => setEditCurrentYear(e.target.value)}
                            placeholder="Ano Atual:"
                            required
                        />
                    </div>
                    <div id={`edit-shift-group-${student.id}`}>
                        <input
                            type="text"
                            id={`edit-shift-${student.id}`}
                            value={editShift}
                            onChange={(e) => setEditShift(e.target.value)}
                            placeholder="Turno (M/T/N):"
                            maxLength="1"
                            required
                        />
                    </div>
                    <div className="form-button-container">
                        <button type="submit" className="edit-button">
                            <FontAwesomeIcon icon={faSave} /> {/* Ícone para Salvar */}
                        </button>
                        <button type="button" className="cancel-button" onClick={handleCancelEdit}>
                            <FontAwesomeIcon icon={faTimes} /> {/* Ícone para Cancelar */}
                        </button>
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
                  {/* Informações do aluno */}
                  <strong>{student.name}</strong> ({student.enrollment}) - Ano: {student.current_year}, Turno: {student.shift}
                  
                  {/* Novo div para agrupar e alinhar os botões de ação */}
                  <div className="list-buttons-container">
                      <button className="details-button" onClick={() => alert(`Detalhes de ${student.name}:\nID: ${student.id}\nMatrícula: ${student.enrollment}\nAno: ${student.current_year}\nTurno: ${student.shift}`)}>
                          <FontAwesomeIcon icon={faInfoCircle} /> {/* Ícone para Detalhes */}
                      </button>
                      <button className="edit-button" onClick={() => handleEditStudent(student)}>
                          <FontAwesomeIcon icon={faEdit} /> {/* Ícone para Editar */}
                      </button>
                      <button className="delete-button" onClick={() => handleDeleteStudent(student.id)}>
                          <FontAwesomeIcon icon={faTrashAlt} /> {/* Ícone para Excluir */}
                      </button>
                  </div>
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