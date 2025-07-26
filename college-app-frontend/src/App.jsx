// src/App.jsx (arquivo principal)

import { useState, useEffect } from 'react';
import './App.css'; 
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { 
  faInfoCircle, faEdit, faTrashAlt, faPlus, faSave, faTimes, faFilter,
  faLink, faUnlink, faUserGraduate, faUserTie, faBook
} from '@fortawesome/free-solid-svg-icons'; 

// Importa o novo componente de gerenciamento de professores
import TeacherManagement from './components/TeacherManagement';


function App() {
  // --- Estado para controlar qual tela está visível ---
  const [currentView, setCurrentView] = useState('students'); // 'students' ou 'teachers'

  // O restante dos estados e funções do App.jsx (alunos) permanecem inalterados,
  // pois eles serão renderizados apenas quando currentView for 'students'.
  
  const [students, setStudents] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const [newName, setNewName] = useState('');
  const [newEnrollment, setNewEnrollment] = useState('');
  const [newCurrentYear, setNewCurrentYear] = useState('');
  const [newShift, setNewShift] = useState('');
  const [formMessage, setFormMessage] = useState('');

  const [editingStudentId, setEditingStudentId] = useState(null);
  const [editName, setEditName] = useState('');
  const [editEnrollment, setEditEnrollment] = useState('');
  const [editCurrentYear, setEditCurrentYear] = useState('');
  const [editShift, setEditShift] = useState('');
  const [editMessage, setEditMessage] = useState('');

  const [filterYear, setFilterYear] = useState('');
  const [filterShift, setFilterShift] = useState('');
  const [hasFiltered, setHasFiltered] = useState(false);

  const [allStudentsForSubjectAssignment, setAllStudentsForSubjectAssignment] = useState([]);
  const [allSubjectsForAssignment, setAllSubjectsForAssignment] = useState([]);
  const [selectedStudentIdForAssignment, setSelectedStudentIdForAssignment] = useState('');
  const [selectedSubjectIdForAssignment, setSelectedSubjectIdForAssignment] = useState('');
  const [assignmentMessage, setAssignmentMessage] = useState('');
  const [assignedSubjectsOfSelectedStudent, setAssignedSubjectsOfSelectedStudent] = useState([]);

  // Seus métodos fetch, handle, etc.
  const fetchStudents = async (year, shift) => {
    setLoading(true);
    setError(null);
    setStudents([]);

    let url = '/api/students';
    const params = new URLSearchParams();

    if (year) {
      params.append('current_year', year);
    }
    if (shift) {
      params.append('shift', shift);
    }

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

  const fetchAllStudentsForAssignment = async () => {
    try {
      const response = await fetch('/api/students'); // Requisição sem filtros
      if (!response.ok) {
        const err = await response.json();
        throw new Error(err.message || 'Erro ao buscar todos os alunos para atribuição');
      }
      const data = await response.json();
      setAllStudentsForSubjectAssignment(Array.isArray(data) ? data : []);
      if (Array.isArray(data) && data.length > 0) {
        setSelectedStudentIdForAssignment(data[0].id); // Seleciona o primeiro aluno por padrão
        // Ao selecionar o primeiro aluno, já busca suas matérias
        fetchAssignedSubjectsForStudent(data[0].id);
      } else {
        setSelectedStudentIdForAssignment('');
        setAssignedSubjectsOfSelectedStudent([]);
      }
    } catch (err) {
      console.error("Erro ao buscar alunos para atribuição:", err);
      setAssignmentMessage(`Erro ao carregar alunos: ${err.message}`);
    }
  };

  const fetchAllSubjectsForAssignment = async () => {
    try {
      const response = await fetch('/api/subjects');
      if (!response.ok) {
        const err = await response.json();
        throw new Error(err.message || 'Erro ao buscar todas as matérias para atribuição');
      }
      const data = await response.json();
      setAllSubjectsForAssignment(Array.isArray(data) ? data : []);
      if (Array.isArray(data) && data.length > 0) {
        setSelectedSubjectIdForAssignment(data[0].id); // Seleciona a primeira matéria por padrão
      } else {
        setSelectedSubjectIdForAssignment('');
      }
    } catch (err) {
      console.error("Erro ao buscar matérias para atribuição:", err);
      setAssignmentMessage(`Erro ao carregar matérias: ${err.message}`);
    }
  };

  const fetchAssignedSubjectsForStudent = async (studentId) => {
    if (!studentId) {
      setAssignedSubjectsOfSelectedStudent([]);
      return;
    }
    try {
      const response = await fetch(`/api/students/${studentId}`);
      if (!response.ok) {
        const err = await response.json();
        throw new Error(err.message || 'Erro ao buscar matérias atribuídas ao aluno');
      }
      const studentData = await response.json();
      setAssignedSubjectsOfSelectedStudent(studentData.subjects || []); // Assegura que é um array
    } catch (err) {
      console.error(`Erro ao buscar matérias atribuídas ao aluno ${studentId}:`, err);
      setAssignedSubjectsOfSelectedStudent([]); // Limpa em caso de erro
    }
  };

  const handleAssignSubject = async () => {
    setAssignmentMessage('');
    if (!selectedStudentIdForAssignment || !selectedSubjectIdForAssignment) {
      setAssignmentMessage('Erro: Selecione um aluno e uma matéria.');
      return;
    }

    // Otimista: Adiciona no UI primeiro
    const studentToUpdate = allStudentsForSubjectAssignment.find(s => s.id === selectedStudentIdForAssignment);
    const subjectToAdd = allSubjectsForAssignment.find(subj => subj.id === selectedSubjectIdForAssignment);

    if (studentToUpdate && subjectToAdd) {
        // Verifica se a matéria já está atribuída para evitar duplicação visual
        const isAlreadyAssigned = assignedSubjectsOfSelectedStudent.some(s => s.id === subjectToAdd.id);
        if (isAlreadyAssigned) {
            setAssignmentMessage('Erro: Matéria já está atribuída a este aluno.');
            return;
        }
    }


    try {
      const response = await fetch(`/api/students/${selectedStudentIdForAssignment}/subjects/${selectedSubjectIdForAssignment}`, {
        method: 'POST',
      });
      const result = await response.json();
      if (!response.ok) {
        // Se a API retornar um erro de "já existe", trate aqui
        if (result.message && result.message.includes("já associada")) {
            throw new Error("Matéria já está atribuída a este aluno.");
        }
        throw new Error(result.message || 'Erro ao atribuir matéria.');
      }

      setAssignmentMessage('Matéria atribuída com sucesso!');
      // Atualiza as matérias atribuídas para o aluno selecionado
      // fetchAssignedSubjectsForStudent(selectedStudentIdForAssignment); // Recarrega do backend
      // Alternativa: Atualização otimista
      setAssignedSubjectsOfSelectedStudent(prev => [...prev, subjectToAdd]);

    } catch (err) {
      setAssignmentMessage(`Erro ao atribuir: ${err.message}`);
      console.error("Erro ao atribuir matéria:", err);
    }
  };

  const handleRemoveSubjectFromAssignment = async (studentId, subjectId) => {
    setAssignmentMessage('');
    if (window.confirm('Tem certeza que deseja remover esta matéria?')) {
      try {
        const response = await fetch(`/api/students/${studentId}/subjects/${subjectId}`, {
          method: 'DELETE',
        });
        if (!response.ok) {
          const err = await response.json();
          throw new Error(err.message || 'Erro ao remover matéria.');
        }

        setAssignmentMessage('Matéria removida com sucesso!');
        // Atualiza as matérias atribuídas para o aluno selecionado
        // fetchAssignedSubjectsForStudent(selectedStudentIdForAssignment); // Recarrega do backend
        // Alternativa: Atualização otimista
        setAssignedSubjectsOfSelectedStudent(prev => prev.filter(s => s.id !== subjectId));
      } catch (err) {
        setAssignmentMessage(`Erro ao remover: ${err.message}`);
        console.error("Erro ao remover matéria:", err);
      }
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault(); setFormMessage('');
    if (!newName || !newEnrollment || !newCurrentYear || !newShift) { setFormMessage('Erro: Todos os campos são obrigatórios!'); return; }
    const studentData = { name: newName, enrollment: newEnrollment, current_year: parseInt(newCurrentYear, 10), shift: newShift, };
    try {
      const response = await fetch('/api/students', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(studentData), });
      const result = await response.json(); if (!response.ok) { throw new Error(result.message || 'Erro ao criar aluno'); }
      setFormMessage('Sucesso: Aluno criado com sucesso!');
      setNewName(''); setNewEnrollment(''); setNewCurrentYear(''); setNewShift('');
      fetchStudents(filterYear, filterShift); fetchAllStudentsForAssignment();
    } catch (err) { setFormMessage(`Erro: ${err.message}`); console.error("Erro ao enviar formulário:", err); }
  };

  const handleDeleteStudent = async (id) => {
    if (window.confirm('Tem certeza que deseja excluir este aluno?')) {
      try {
        const response = await fetch(`/api/students/${id}`, { method: 'DELETE' });
        if (!response.ok) { const errorData = await response.json(); throw new Error(errorData.message || `Erro ao excluir aluno com ID: ${id}`); }
        setFormMessage('Sucesso: Aluno excluído com sucesso!');
        fetchStudents(filterYear, filterShift); fetchAllStudentsForAssignment();
      } catch (err) { setFormMessage(`Erro ao excluir: ${err.message}`); console.error("Erro ao excluir aluno:", err); }
    }
  };

  const handleEditStudent = (student) => {
    setEditingStudentId(student.id); setEditName(student.name); setEditEnrollment(student.enrollment); setEditCurrentYear(student.current_year.toString()); setEditShift(student.shift); setEditMessage('');
  };

  const handleSaveEdit = async (e) => {
    e.preventDefault(); setEditMessage('');
    if (!editName || !editEnrollment || !editCurrentYear || !editShift) { setEditMessage('Erro: Todos os campos são obrigatórios!'); return; }
    const updatedStudentData = { id: editingStudentId, name: editName, enrollment: editEnrollment, current_year: parseInt(editCurrentYear, 10), shift: editShift, };
    try {
      const response = await fetch(`/api/students/${editingStudentId}`, { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(updatedStudentData), });
      const result = await response.json(); if (!response.ok) { throw new Error(result.message || 'Erro ao atualizar aluno'); }
      setEditMessage('Sucesso: Aluno atualizado com sucesso!'); setEditingStudentId(null); fetchStudents(filterYear, filterShift);
    } catch (err) { setEditMessage(`Erro: ${err.message}`); console.error("Erro ao atualizar aluno:", err); }
  };

  const handleCancelEdit = () => {
    setEditingStudentId(null); setEditMessage('');
  };

  const handleFilterSubmit = (e) => {
    e.preventDefault(); setHasFiltered(true); fetchStudents(filterYear, filterShift);
  };

  // Carrega alunos e matérias para a seção de atribuição UMA VEZ ao montar o componente App
  // quando a view inicial é 'students'.
  useEffect(() => {
    if (currentView === 'students') {
      fetchStudents(''); // Carrega todos os alunos para a lista principal
      fetchAllStudentsForAssignment(); // Carrega para o dropdown de atribuição
      fetchAllSubjectsForAssignment(); // Carrega matérias para o dropdown de atribuição
    }
  }, [currentView]); // Depende da view para carregar os dados iniciais corretamente

  return (
    <div className="App">
      <div className="main-header">
        <h1>Gerenciamento da Universidade</h1>
        <div className="navigation-buttons">
          <button 
            className={`nav-button ${currentView === 'students' ? 'active' : ''}`} 
            onClick={() => setCurrentView('students')}
          >
            <FontAwesomeIcon icon={faUserGraduate} /> Alunos
          </button>
          <button 
            className={`nav-button ${currentView === 'teachers' ? 'active' : ''}`} 
            onClick={() => setCurrentView('teachers')}
          >
            <FontAwesomeIcon icon={faUserTie} /> Professores
          </button>
        </div>
      </div>

      <hr />

      {currentView === 'students' && (
        <div className="students-section">
          {/* O conteúdo original do seu App.jsx (formulários de aluno, filtro, lista e atribuição de matéria) */}
          
          {/* Formulário de Criação de Alunos */}
          <h2>Adicionar Novo Aluno</h2>
          <form onSubmit={handleSubmit}>
            <div>
              <input type="text" id="name" value={newName} onChange={(e) => setNewName(e.target.value)} placeholder="Nome:" required />
            </div>
            <div className="form-group-inline">
                <div>
                    <input type="text" id="enrollment" value={newEnrollment} onChange={(e) => setNewEnrollment(e.target.value)} placeholder="Matrícula:" required />
                </div>
                <div id="new-current-year-group">
                    <input type="number" id="new_current_year" value={newCurrentYear} onChange={(e) => setNewCurrentYear(e.target.value)} placeholder="Ano Atual:" required />
                </div>
                <div id="new-shift-group">
                    <input type="text" id="new_shift" value={newShift} onChange={(e) => setNewShift(e.target.value)} placeholder="Turno (M/T/N):" maxLength="1" required />
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

          {/* Formulário de Filtro */}
          <h2>Filtrar Alunos</h2>
          <form onSubmit={handleFilterSubmit} className="filter-form">
            <div className="form-group-inline">
              <div id="filter-year-group">
                <input type="number" id="filter_year" value={filterYear} onChange={(e) => setFilterYear(e.target.value)} placeholder="Ano do Aluno:" />
              </div>
              <div id="filter-shift-group">
                <input type="text" id="filter_shift" value={filterShift} onChange={(e) => setFilterShift(e.target.value)} placeholder="Turno (M/T/N):" maxLength="1" />
              </div>
              <div className="form-button-container">
                <button type="submit" className="filter-button">
                  <FontAwesomeIcon icon={faFilter} /> Aplicar Filtros
                </button>
              </div>
            </div>
          </form>

          <hr />

          {/* Nova Seção de Gerenciamento de Matérias (agora faz parte da view de alunos) */}
          <div className="subject-assignment-section">
            <h2>Atribuir Matérias a Alunos</h2>
            {assignmentMessage && (
              <p className={assignmentMessage.startsWith('Erro') ? 'error-message' : 'success-message'}>
                {assignmentMessage}
              </p>
            )}

            <div className="form-group-inline">
              <div className="select-container">
                <select
                    value={selectedStudentIdForAssignment}
                    onChange={(e) => {
                        setSelectedStudentIdForAssignment(e.target.value);
                        fetchAssignedSubjectsForStudent(e.target.value); // Ao mudar o aluno, atualiza a lista de matérias atribuídas
                        setAssignmentMessage(''); // Limpa mensagens anteriores
                    }}
                    className="select-student-assignment"
                >
                    {allStudentsForSubjectAssignment.length > 0 ? (
                        allStudentsForSubjectAssignment.map(student => (
                            <option key={student.id} value={student.id}>
                                {student.name} ({student.enrollment})
                            </option>
                        ))
                    ) : (
                        <option value="">Carregando alunos...</option>
                    )}
                </select>
              </div>

              <div className="select-container">
                <select
                    value={selectedSubjectIdForAssignment}
                    onChange={(e) => {
                        setSelectedSubjectIdForAssignment(e.target.value);
                        setAssignmentMessage(''); // Limpa mensagens
                    }}
                    className="subject-select-assignment"
                >
                    {allSubjectsForAssignment.length > 0 ? (
                        allSubjectsForAssignment.map(subject => (
                            <option key={subject.id} value={subject.id}>
                                {subject.name} (Ano: {subject.year}, Créditos: {subject.credits})
                            </option>
                        ))
                    ) : (
                        <option value="">Carregando matérias...</option>
                    )}
                </select>
              </div>

              <div className="form-button-container">
                <button type="button" className="add-subject-to-student-button" onClick={handleAssignSubject}>
                  <FontAwesomeIcon icon={faLink} /> Atribuir
                </button>
              </div>
            </div>

            {selectedStudentIdForAssignment && (
                <div className="assigned-subjects-view">
                    <h3>Matérias de {allStudentsForSubjectAssignment.find(s => s.id === selectedStudentIdForAssignment)?.name || 'Aluno Selecionado'}</h3>
                    {assignedSubjectsOfSelectedStudent.length > 0 ? (
                        <ul className="assigned-subjects-list">
                            {assignedSubjectsOfSelectedStudent.map(subject => (
                                <li key={subject.id}>
                                    <FontAwesomeIcon icon={faBook} /> {subject.name} (Ano: {subject.year}, Créditos: {subject.credits})
                                    <button
                                        type="button"
                                        className="remove-subject-button"
                                        onClick={() => handleRemoveSubjectFromAssignment(selectedStudentIdForAssignment, subject.id)}
                                    >
                                        <FontAwesomeIcon icon={faUnlink} /> Remover
                                    </button>
                                </li>
                            ))}
                        </ul>
                    ) : (
                        <p>Nenhuma matéria atribuída a este aluno.</p>
                    )}
                </div>
            )}
          </div>

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
                    // Formulário de Edição (SEM gerenciamento de matéria aqui agora)
                    <form onSubmit={handleSaveEdit}>
                      <h3>Editando Aluno: {student.name}</h3>
                      <div>
                        <input type="text" id={`edit-name-${student.id}`} value={editName} onChange={(e) => setEditName(e.target.value)} placeholder="Nome:" required />
                      </div>
                      <div className="form-group-inline">
                        <div>
                            <input type="text" id={`edit-enrollment-${student.id}`} value={editEnrollment} onChange={(e) => setEditEnrollment(e.target.value)} placeholder="Matrícula:" required />
                        </div>
                        <div id={`edit-current-year-group-${student.id}`}>
                            <input type="number" id={`edit-current_year-${student.id}`} value={editCurrentYear} onChange={(e) => setEditCurrentYear(e.target.value)} placeholder="Ano Atual:" required />
                        </div>
                        <div id={`edit-shift-group-${student.id}`}>
                            <input type="text" id={`edit-shift-${student.id}`} value={editShift} onChange={(e) => setEditShift(e.target.value)} placeholder="Turno (M/T/N):" maxLength="1" required />
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
                      <FontAwesomeIcon icon={faUserGraduate} style={{ marginRight: '10px' }} /> {/* Ícone para aluno */}
                      <strong>{student.name}</strong> ({student.enrollment}) - Ano: {student.current_year}, Turno: {student.shift}
                      
                      {/* Div para agrupar e alinhar os botões de ação */}
                      <div className="list-buttons-container">
                          <button className="details-button" onClick={() => alert(`Detalhes de ${student.name}:\nID: ${student.id}\nMatrícula: ${student.enrollment}\nAno: ${student.current_year}\nTurno: ${student.shift}\nMatérias: ${student.subjects ? student.subjects.map(s => s.name).join(', ') : 'Nenhuma'}`)}>
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
      )} {/* Fim da renderização condicional de alunos */}

      {currentView === 'teachers' && (
        <div className="teachers-section">
          <TeacherManagement />
        </div>
      )} {/* Fim da renderização condicional de professores */}

    </div>
  );
}

export default App;