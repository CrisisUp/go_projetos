import { useState, useEffect } from 'react';
import './App.css'; 
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faInfoCircle, faEdit, faTrashAlt, faPlus, faSave, faTimes, faFilter } from '@fortawesome/free-solid-svg-icons'; // Adicionado faFilter

function App() {
  const [students, setStudents] = useState([]);
  const [loading, setLoading] = useState(false); // Inicialmente false, pois a lista só carrega após o filtro
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

  // --- Novos Estados para o Filtro ---
  const [filterYear, setFilterYear] = useState('');
  const [filterShift, setFilterShift] = useState('');
  const [hasFiltered, setHasFiltered] = useState(false); // Para controlar se o filtro já foi aplicado ao menos uma vez

  // Função para buscar os alunos (agora com filtros)
  const fetchStudents = async (year, shift) => {
    setLoading(true);
    setError(null); // Limpa erros anteriores
    setStudents([]); // Limpa a lista enquanto carrega

    let url = '/api/students';
    const params = new URLSearchParams();

    if (year) {
      params.append('current_year', year);
    }
    if (shift) {
      params.append('shift', shift);
    }

    // Se houver parâmetros, adiciona-os à URL
    if (params.toString()) {
      url += `?${params.toString()}`;
    }
    
    try {
      const response = await fetch(url);
      if (!response.ok) {
        const err = await response.json();
        throw new Error(err.message || 'Erro ao buscar alunos');
      }
      const data = await response.json();
      setStudents(Array.isArray(data) ? data : []);
    } catch (err) {
      setError(err);
      console.error("Erro ao buscar alunos:", err);
    } finally {
      setLoading(false);
    }
  };

  // NENHUM useEffect inicial chamando fetchStudents() ao montar.
  // A lista só será carregada após o usuário clicar em "Aplicar Filtros".

  // Handler para submissão do formulário de CRIAÇÃO
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
      // Após criar, recarrega a lista de acordo com os filtros ATUAIS
      fetchStudents(filterYear, filterShift); 
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
        // Após excluir, recarrega a lista de acordo com os filtros ATUAIS
        fetchStudents(filterYear, filterShift);
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
      // Após editar, recarrega a lista de acordo com os filtros ATUAIS
      fetchStudents(filterYear, filterShift);
    } catch (err) {
      setEditMessage(`Erro: ${err.message}`);
      console.error("Erro ao atualizar aluno:", err);
    }
  };

  const handleCancelEdit = () => {
    setEditingStudentId(null);
    setEditMessage('');
  };

  // --- Handler para o formulário de filtro ---
  const handleFilterSubmit = (e) => {
    e.preventDefault();
    setHasFiltered(true); // Indica que o filtro já foi aplicado
    fetchStudents(filterYear, filterShift);
  };

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
                    <FontAwesomeIcon icon={faPlus} /> Criar Aluno
                </button>
            </div>
        </div>
      </form>

      {formMessage && (
        <p className={formMessage.startsWith('Erro') ? 'error-message' : 'success-message'}>
          {formMessage}
        </p>
      )}

      <hr />

      {/* --- Formulário de Filtro de Alunos --- */}
      <h2>Filtrar Alunos</h2>
      <form onSubmit={handleFilterSubmit} className="filter-form">
        <div className="form-group-inline">
          <div id="filter-year-group">
            <input
              type="number"
              id="filter_year"
              value={filterYear}
              onChange={(e) => setFilterYear(e.target.value)}
              placeholder="Ano do Aluno:"
            />
          </div>
          <div id="filter-shift-group">
            <input
              type="text"
              id="filter_shift"
              value={filterShift}
              onChange={(e) => setFilterShift(e.target.value)}
              placeholder="Turno (M/T/N):"
              maxLength="1"
            />
          </div>
          <div className="form-button-container">
            <button type="submit" className="filter-button">
              <FontAwesomeIcon icon={faFilter} /> Aplicar Filtros
            </button>
          </div>
        </div>
      </form>

      <hr />

      {/* Lista de Alunos Existentes */}
      <h2>Lista de Alunos</h2>
      {loading ? (
        <div>Carregando alunos...</div>
      ) : error ? (
        <div>Erro ao carregar alunos: {error.message}</div>
      ) : hasFiltered && students.length === 0 ? (
        <p>Nenhum aluno encontrado com os filtros aplicados.</p>
      ) : !hasFiltered ? (
        <p>Use o formulário acima para aplicar filtros e carregar a lista de alunos.</p>
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
                            <FontAwesomeIcon icon={faSave} />
                        </button>
                        <button type="button" className="cancel-button" onClick={handleCancelEdit}>
                            <FontAwesomeIcon icon={faTimes} />
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
                          <FontAwesomeIcon icon={faInfoCircle} />
                      </button>
                      <button className="edit-button" onClick={() => handleEditStudent(student)}>
                          <FontAwesomeIcon icon={faEdit} />
                      </button>
                      <button className="delete-button" onClick={() => handleDeleteStudent(student.id)}>
                          <FontAwesomeIcon icon={faTrashAlt} />
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