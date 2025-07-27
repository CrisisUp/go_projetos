// src/components/TeacherManagement.jsx

import React, { useState, useEffect } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { 
  faUserTie, faEdit, faTrashAlt, faPlus, faSave, faTimes, faFilter,
  faLink, faUnlink, faBook, faInfoCircle
} from '@fortawesome/free-solid-svg-icons';

// O componente TeacherManagement será responsável por toda a lógica e UI dos professores
function TeacherManagement() {
  const [teachers, setTeachers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Estados para o formulário de novo professor
  const [newTeacherName, setNewTeacherName] = useState('');
  const [newTeacherDepartment, setNewTeacherDepartment] = useState('');
  const [newTeacherEmail, setNewTeacherEmail] = useState(''); // <-- NOVO ESTADO PARA EMAIL
  const [formMessage, setFormMessage] = useState('');

  // Estados para edição de professor
  const [editingTeacherId, setEditingTeacherId] = useState(null);
  const [editTeacherName, setEditTeacherName] = useState('');
  const [editTeacherDepartment, setEditTeacherDepartment] = useState('');
  const [editTeacherEmail, setEditTeacherEmail] = useState(''); // <-- NOVO ESTADO PARA EMAIL DE EDIÇÃO
  const [editMessage, setEditMessage] = useState('');

  // Estados para filtro de professor (ex: por departamento)
  const [filterDepartment, setFilterDepartment] = useState('');
  const [hasFiltered, setHasFiltered] = useState(false);

  // --- Estados para Gerenciamento de Matérias do Professor ---
  const [allSubjectsForTeacherAssignment, setAllSubjectsForTeacherAssignment] = useState([]);
  const [selectedSubjectIdForTeacherAssignment, setSelectedSubjectIdForTeacherAssignment] = useState('');
  const [assignmentMessage, setAssignmentMessage] = useState('');
  const [assignedSubjectsOfSelectedTeacher, setAssignedSubjectsOfSelectedTeacher] = useState([]);


  // Função para buscar os professores (com filtros, se houver)
  const fetchTeachers = async (department) => {
    setLoading(true);
    setError(null);
    setTeachers([]);

    let url = '/api/teachers';
    const params = new URLSearchParams();

    if (department) {
      params.append('department', department); // Assumindo que o backend suporta filtro por departamento
    }

    if (params.toString()) {
      url += `?${params.toString()}`;
    }
    
    try {
      const response = await fetch(url);
      if (!response.ok) {
        const err = await response.json();
        throw new Error(err.message || 'Erro ao buscar professores');
      }
      const data = await response.json();
      setTeachers(Array.isArray(data) ? data : []);
    } catch (err) {
      setError(err);
      console.error("Erro ao buscar professores:", err);
    } finally {
      setLoading(false);
    }
  };

  // --- Funções para Gerenciamento de Matérias do Professor ---
  const fetchAllSubjectsForTeacherAssignment = async () => {
    try {
      const response = await fetch('/api/subjects');
      if (!response.ok) {
        const err = await response.json();
        throw new Error(err.message || 'Erro ao buscar matérias para atribuição');
      }
      const data = await response.json();
      setAllSubjectsForTeacherAssignment(Array.isArray(data) ? data : []);
      if (Array.isArray(data) && data.length > 0) {
        setSelectedSubjectIdForTeacherAssignment(data[0].id);
      } else {
        setSelectedSubjectIdForTeacherAssignment('');
      }
    } catch (err) {
      console.error("Erro ao buscar matérias para atribuição:", err);
      setAssignmentMessage(`Erro ao carregar matérias: ${err.message}`);
    }
  };

  const fetchAssignedSubjectsForTeacher = async (teacherId) => {
    if (!teacherId) {
      setAssignedSubjectsOfSelectedTeacher([]);
      return;
    }
    try {
      // Sua API GetAllTeachersHandler já retorna as matérias associadas se o modelo de professor tiver 'Subjects'
      // Se não, você precisará de uma rota GET /api/teachers/{id}/subjects
      const response = await fetch(`/api/teachers/${teacherId}`); // Buscar professor por ID, assumindo que ele tem subjects
      if (!response.ok) {
        const err = await response.json();
        throw new Error(err.message || 'Erro ao buscar matérias atribuídas ao professor');
      }
      const teacherData = await response.json();
      setAssignedSubjectsOfSelectedTeacher(teacherData.subjects || []); // Assegura que é um array
    } catch (err) {
      console.error(`Erro ao buscar matérias atribuídas ao professor ${teacherId}:`, err);
      setAssignedSubjectsOfSelectedTeacher([]);
    }
  };

  const handleAssignSubjectToTeacher = async (teacherId, subjectId) => {
    setAssignmentMessage('');
    if (!teacherId || !subjectId) {
      setAssignmentMessage('Erro: Selecione um professor e uma matéria.');
      return;
    }

    const teacherToUpdate = teachers.find(t => t.id === teacherId);
    const subjectToAdd = allSubjectsForTeacherAssignment.find(subj => subj.id === subjectId);

    if (teacherToUpdate && subjectToAdd) {
        const isAlreadyAssigned = assignedSubjectsOfSelectedTeacher.some(s => s.id === subjectToAdd.id);
        if (isAlreadyAssigned) {
            setAssignmentMessage('Erro: Matéria já está atribuída a este professor.');
            return;
        }
    }

    try {
      const response = await fetch(`/api/teachers/${teacherId}/subjects/${subjectId}`, {
        method: 'POST',
      });
      const result = await response.json();
      if (!response.ok) {
        if (result.message && result.message.includes("já associada")) {
            throw new Error("Matéria já está atribuída a este professor.");
        }
        throw new Error(result.message || 'Erro ao atribuir matéria ao professor.');
      }

      setAssignmentMessage('Matéria atribuída com sucesso ao professor!');
      // Atualização otimista do estado
      setAssignedSubjectsOfSelectedTeacher(prev => [...prev, subjectToAdd]);

    } catch (err) {
      setAssignmentMessage(`Erro ao atribuir: ${err.message}`);
      console.error("Erro ao atribuir matéria ao professor:", err);
    }
  };

  const handleRemoveSubjectFromTeacher = async (teacherId, subjectId) => {
    setAssignmentMessage('');
    if (window.confirm('Tem certeza que deseja remover esta matéria do professor?')) {
      try {
        const response = await fetch(`/api/teachers/${teacherId}/subjects/${subjectId}`, {
          method: 'DELETE',
        });
        if (!response.ok) {
          const err = await response.json();
          throw new Error(err.message || 'Erro ao remover matéria do professor.');
        }

        setAssignmentMessage('Matéria removida com sucesso do professor!');
        // Atualização otimista do estado
        setAssignedSubjectsOfSelectedTeacher(prev => prev.filter(s => s.id !== subjectId));
      } catch (err) {
        setAssignmentMessage(`Erro ao remover: ${err.message}`);
        console.error("Erro ao remover matéria do professor:", err);
      }
    }
  };
  // --- FIM Funções para Gerenciamento de Matérias do Professor ---


  // Handler para criar um novo professor
  const handleCreateTeacher = async (e) => {
    e.preventDefault();
    setFormMessage('');
    // <-- VALIDAÇÃO ATUALIZADA COM EMAIL
    if (!newTeacherName || !newTeacherDepartment || !newTeacherEmail) {
      setFormMessage('Erro: Nome, Departamento e Email são obrigatórios!');
      return;
    }
    const teacherData = {
      name: newTeacherName,
      department: newTeacherDepartment,
      email: newTeacherEmail, // <-- INCLUÍDO NO PAYLOAD
    };
    try {
      const response = await fetch('/api/teachers', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(teacherData),
      });
      const result = await response.json();
      if (!response.ok) { throw new Error(result.message || 'Erro ao criar professor'); }
      setFormMessage('Sucesso: Professor criado com sucesso!');
      setNewTeacherName('');
      setNewTeacherDepartment('');
      setNewTeacherEmail(''); // <-- LIMPAR O CAMPO EMAIL
      fetchTeachers(filterDepartment); // Recarrega a lista
    } catch (err) {
      setFormMessage(`Erro: ${err.message}`);
      console.error("Erro ao criar professor:", err);
    }
  };

  // Handler para iniciar edição
  const handleEditTeacher = (teacher) => {
    setEditingTeacherId(teacher.id);
    setEditTeacherName(teacher.name);
    setEditTeacherDepartment(teacher.department);
    setEditTeacherEmail(teacher.email); // <-- CARREGAR EMAIL EXISTENTE
    setEditMessage('');
    setAssignmentMessage(''); // Limpa mensagens de atribuição
    fetchAllSubjectsForTeacherAssignment(); // Busca matérias para atribuição
    fetchAssignedSubjectsForTeacher(teacher.id); // Busca matérias atribuídas ao professor
  };

  // Handler para salvar edição
  const handleSaveEditTeacher = async (e) => {
    e.preventDefault();
    setEditMessage('');
    // <-- VALIDAÇÃO ATUALIZADA COM EMAIL
    if (!editTeacherName || !editTeacherDepartment || !editTeacherEmail) {
      setEditMessage('Erro: Nome, Departamento e Email são obrigatórios!');
      return;
    }
    const updatedTeacherData = {
      id: editingTeacherId,
      name: editTeacherName,
      department: editTeacherDepartment,
      email: editTeacherEmail, // <-- INCLUÍDO NO PAYLOAD
    };
    try {
      const response = await fetch(`/api/teachers/${editingTeacherId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updatedTeacherData),
      });
      const result = await response.json();
      if (!response.ok) { throw new Error(result.message || 'Erro ao atualizar professor'); }
      setEditMessage('Sucesso: Professor atualizado com sucesso!');
      setEditingTeacherId(null);
      fetchTeachers(filterDepartment); // Recarrega a lista
    } catch (err) {
      setEditMessage(`Erro: ${err.message}`);
      console.error("Erro ao atualizar professor:", err);
    }
  };

  // Handler para cancelar edição
  const handleCancelEditTeacher = () => {
    setEditingTeacherId(null);
    setEditMessage('');
    setAssignmentMessage('');
  };

  // Handler para deletar professor
  const handleDeleteTeacher = async (id) => {
    if (window.confirm('Tem certeza que deseja excluir este professor?')) {
      try {
        const response = await fetch(`/api/teachers/${id}`, { method: 'DELETE' });
        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.message || `Erro ao excluir professor com ID: ${id}`);
        }
        setFormMessage('Sucesso: Professor excluído com sucesso!');
        fetchTeachers(filterDepartment); // Recarrega a lista
      } catch (err) {
        setFormMessage(`Erro ao excluir: ${err.message}`);
        console.error("Erro ao excluir professor:", err);
      }
    }
  };

  // Handler para filtro
  const handleFilterTeachers = (e) => {
    e.preventDefault();
    setHasFiltered(true);
    fetchTeachers(filterDepartment);
  };

  // Carrega professores ao montar (ou após filtro)
  useEffect(() => {
    fetchTeachers(''); // Carrega todos os professores inicialmente
  }, []); // Array de dependências vazio para carregar apenas uma vez


  return (
    <div className="teacher-management-page">
      {/* Formulário de Criação de Professor */}
      <h2>Adicionar Novo Professor</h2>
      <form onSubmit={handleCreateTeacher}>
        <div>
          <input type="text" value={newTeacherName} onChange={(e) => setNewTeacherName(e.target.value)} placeholder="Nome do Professor:" required />
        </div>
        <div>
          <input type="text" value={newTeacherDepartment} onChange={(e) => setNewTeacherDepartment(e.target.value)} placeholder="Departamento:" required />
        </div>
        <div> {/* NOVO CAMPO DE EMAIL PARA CRIAÇÃO */}
          <input type="email" value={newTeacherEmail} onChange={(e) => setNewTeacherEmail(e.target.value)} placeholder="Email do Professor:" required />
        </div>
        <button type="submit" className="create-button">
          <FontAwesomeIcon icon={faPlus} /> Criar Professor
        </button>
      </form>

      {formMessage && (
        <p className={formMessage.startsWith('Erro') ? 'error-message' : 'success-message'}>
          {formMessage}
        </p>
      )}

      <hr />

      {/* Formulário de Filtro de Professores */}
      <h2>Filtrar Professores</h2>
      <form onSubmit={handleFilterTeachers} className="filter-form">
        <div className="form-group-inline">
          <div id="filter-department-group">
            <input type="text" value={filterDepartment} onChange={(e) => setFilterDepartment(e.target.value)} placeholder="Filtrar por Departamento:" />
          </div>
          <div className="form-button-container">
            <button type="submit" className="filter-button">
              <FontAwesomeIcon icon={faFilter} /> Aplicar Filtros
            </button>
          </div>
        </div>
      </form>

      <hr />

      {/* Lista de Professores Existentes */}
      <h2>Lista de Professores</h2>
      {loading ? (
        <div>Carregando professores...</div>
      ) : error ? (
        <div>Erro ao carregar professores: {error.message}</div>
      ) : teachers.length === 0 && hasFiltered ? (
        <p>Nenhum professor encontrado com os filtros aplicados.</p>
      ) : teachers.length === 0 && !hasFiltered ? (
        <p>Nenhum professor cadastrado. Crie um acima!</p>
      ) : (
        <ul>
          {teachers.map(teacher => (
            <li key={teacher.id}>
              {editingTeacherId === teacher.id ? (
                // Formulário de Edição de Professor
                <form onSubmit={handleSaveEditTeacher}>
                  <h3>Editando Professor: {teacher.name}</h3>
                  <div>
                    <input type="text" value={editTeacherName} onChange={(e) => setEditTeacherName(e.target.value)} placeholder="Nome:" required />
                  </div>
                  <div>
                    <input type="text" value={editTeacherDepartment} onChange={(e) => setEditTeacherDepartment(e.target.value)} placeholder="Departamento:" required />
                  </div>
                  <div> {/* NOVO CAMPO DE EMAIL PARA EDIÇÃO */}
                    <input type="email" value={editTeacherEmail} onChange={(e) => setEditTeacherEmail(e.target.value)} placeholder="Email:" required />
                  </div>
                  <div className="form-button-container">
                    <button type="submit" className="edit-button">
                        <FontAwesomeIcon icon={faSave} />
                    </button>
                    <button type="button" className="cancel-button" onClick={handleCancelEditTeacher}>
                        <FontAwesomeIcon icon={faTimes} />
                    </button>
                  </div>
                  {editMessage && (
                    <p className={editMessage.startsWith('Erro') ? 'error-message' : 'success-message'}>
                      {editMessage}
                    </p>
                  )}

                  {/* --- SEÇÃO: Gerenciar Matérias do Professor --- */}
                  <div className="subject-management-section">
                    <h4>Gerenciar Matérias do Professor</h4>
                    {assignmentMessage && (
                        <p className={assignmentMessage.startsWith('Erro') ? 'error-message' : 'success-message'}>
                            {assignmentMessage}
                        </p>
                    )}
                    <div className="subject-add-form">
                        <select
                            value={selectedSubjectIdForTeacherAssignment}
                            onChange={(e) => setSelectedSubjectIdForTeacherAssignment(e.target.value)}
                            className="subject-select"
                        >
                            {allSubjectsForTeacherAssignment.length > 0 ? (
                                allSubjectsForTeacherAssignment.map(subject => (
                                    <option key={subject.id} value={subject.id}>
                                        {subject.name} (Ano: {subject.year}, Créditos: {subject.credits})
                                    </option>
                                ))
                            ) : (
                                <option value="">Carregando matérias...</option>
                            )}
                        </select>
                        <button type="button" className="add-subject-button" onClick={() => handleAssignSubjectToTeacher(teacher.id, selectedSubjectIdForTeacherAssignment)}>
                            <FontAwesomeIcon icon={faLink} /> Atribuir
                        </button>
                    </div>
                    
                    <h5>Matérias Atribuídas:</h5>
                    {assignedSubjectsOfSelectedTeacher && assignedSubjectsOfSelectedTeacher.length > 0 ? (
                        <ul className="assigned-subjects-list">
                            {assignedSubjectsOfSelectedTeacher.map(subject => (
                                <li key={subject.id}>
                                    <FontAwesomeIcon icon={faBook} /> {subject.name} (Ano: {subject.year})
                                    <button
                                        type="button"
                                        className="remove-subject-button"
                                        onClick={() => handleRemoveSubjectFromTeacher(teacher.id, subject.id)}
                                    >
                                        <FontAwesomeIcon icon={faUnlink} /> Remover
                                    </button>
                                </li>
                            ))}
                        </ul>
                    ) : (
                        <p>Nenhuma matéria atribuída a este professor.</p>
                    )}
                  </div>
                  {/* --- FIM SEÇÃO --- */}

                </form>
              ) : (
                // Visualização Normal do Professor
                <>
                  <FontAwesomeIcon icon={faUserTie} style={{ marginRight: '10px' }} />
                  <strong>{teacher.name}</strong> ({teacher.department}) - {teacher.email} {/* <-- EXIBIR EMAIL */}
                  
                  <div className="list-buttons-container">
                      <button className="details-button" onClick={() => alert(`Detalhes de ${teacher.name}:\nID: ${teacher.id}\nDepartamento: ${teacher.department}\nEmail: ${teacher.email}\nMatérias: ${teacher.subjects ? teacher.subjects.map(s => s.name).join(', ') : 'Nenhuma'}`)}> {/* <-- ATUALIZAR DETALHES */}
                          <FontAwesomeIcon icon={faInfoCircle} />
                      </button>
                      <button className="edit-button" onClick={() => handleEditTeacher(teacher)}>
                          <FontAwesomeIcon icon={faEdit} />
                      </button>
                      <button className="delete-button" onClick={() => handleDeleteTeacher(teacher.id)}>
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

export default TeacherManagement;