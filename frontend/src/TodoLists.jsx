import { useEffect, useState } from 'react';
import { getLists, logout, getListById, updateList } from './api';

function TodoLists({ onLogout }) {
  const [lists, setLists] = useState(null);
  const [selectedList, setSelectedList] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newTitle, setNewTitle] = useState('');
  const [newDesc, setNewDesc] = useState('');
  const [editTitle, setEditTitle] = useState('');
  const [editDesc, setEditDesc] = useState('');

  useEffect(() => {
    getLists()
      .then(response => {
        const normalized = Array.isArray(response)
          ? response
          : Array.isArray(response?.data)
            ? response.data
            : [];

        setLists(normalized);
      })
      .catch(err => {
        console.error("Failed to fetch lists:", err);
        setLists([]);
      });
  }, []);

  const handleLogout = async () => {
    await logout();
    onLogout();
  };

  async function loadLists() {
    try {
      const data = await getLists();
      setLists(data || []);
    } catch (err) {
      console.error('Error loading lists:', err);
    }
  }
  const openEditModal = async (id) => {
    try {
      const list = await getListById(id);
      setSelectedList(list);
      setEditTitle(list.title);
      setEditDesc(list.description || '');
      setIsModalOpen(true);
    } catch (err) {
      console.error("Error loading list:", err);
    }
  };

  async function handleCreate() {
    try {
      if (!newTitle.trim()) {
        alert('Title is required');
        return;
      }
      const newId = await createList(newTitle.trim(), newDesc.trim());
      console.log('New list created with id:', newId);
      setShowModal(false);
      setNewTitle('');
      setNewDesc('');
      await loadLists();
    } catch (err) {
      console.error('Error creating list:', err);
    }
  }

  const handleUpdate = async () => {
    if (!selectedList) return;

    console.log("clicked update");

    try {
      const result = await updateList(selectedList.id, editTitle, editDesc);

      if (!result || !result.id) {
        throw new Error("Update returned invalid result");
      }

      console.log("Updated successfully");

      const updated = await getLists();
      const normalized = Array.isArray(updated)
        ? updated
        : Array.isArray(updated?.data)
          ? updated.data
          : [];

      setLists(normalized);
      setIsModalOpen(false);
      setSelectedList(null);
    } catch (err) {
      console.error("Error updating list:", err);
      alert("Update failed. Try again.");
    }
  };


  if (!Array.isArray(lists)) return <div>Loading...</div>;

  return (
    <div>
      <h1>Your Todo Lists</h1>
      <button onClick={handleLogout}>Logout</button>

      <ul>
        {lists.map(list => (
          <li key={list.id}>
            <button
              id={`list-${list.id}`}
              onClick={() => openEditModal(list.id)}
              style={{ display: 'block', margin: '10px 0' }}
            >
              {list.title}
            </button>
          </li>
        ))}
      </ul>

      {/* Modal */}
      {isModalOpen && (
        <div style={modalOverlay}>
          <div style={modalBox}>
            <button onClick={handleUpdate} style={updateButton}>Update</button>
            <h2>Edit Todo List</h2>
            <input
              type="text"
              value={editTitle}
              onChange={e => setEditTitle(e.target.value)}
              placeholder="Title"
              style={{ width: '100%', marginBottom: 10 }}
            />
            <textarea
              value={editDesc}
              onChange={e => setEditDesc(e.target.value)}
              placeholder="Description (optional)"
              rows={4}
              style={{ width: '100%' }}
            />
            <button onClick={() => setIsModalOpen(false)} style={{ marginTop: 10 }}>Cancel</button>
          </div>
        </div>
      )}
      <div style={{ marginBottom: '20px', border: '1px solid #ccc', padding: '10px' }}>
        <h3>Create New List</h3>
        <input
          type="text"
          placeholder="Title"
          value={newTitle}
          onChange={(e) => setNewTitle(e.target.value)}
          style={{ display: 'block', marginBottom: '8px' }}
        />
        <textarea
          placeholder="Description (optional)"
          value={newDesc}
          onChange={(e) => setNewDesc(e.target.value)}
          style={{ display: 'block', marginBottom: '8px' }}
        />
        <button onClick={handleCreate}>Create New Todo List</button>
      </div>
    </div>
  );
}

const modalOverlay = {
  position: 'fixed',
  top: 0, left: 0, right: 0, bottom: 0,
  backgroundColor: 'rgba(0, 0, 0, 0.5)',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  zIndex: 1000,
};

const modalBox = {
  background: 'white',
  padding: '20px',
  borderRadius: '8px',
  minWidth: '300px',
  maxWidth: '500px',
  position: 'relative',
};

const updateButton = {
  position: 'absolute',
  top: 10,
  right: 10,
};

export default TodoLists;