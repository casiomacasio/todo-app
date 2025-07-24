import { useEffect, useState } from 'react';
import './index.css'; 
import {createList,deleteListById, getLists, logout, getListById, updateList } from './api';

function TodoLists({ onLogout }) {
  const [lists, setLists] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [showModal, setShowModal] = useState(false);
  const [newTitle, setNewTitle] = useState('');
  const [newDesc, setNewDesc] = useState('');
  const [loadingList, setLoadingList] = useState(false);
  const [editTitle, setEditTitle] = useState('');
  const [editDesc, setEditDesc] = useState('');
  const [selectedList, setSelectedList] = useState(null);

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
    const confirmed = window.confirm("Are you sure you want to log out?");
    if (confirmed) {
      await logout();
      onLogout();
    }
  };

  const openEditModal = async (id) => {
    try {
      console.log('Clicked list ID:', id); 
      const list = await getListById(id);
      console.log('Fetched list data:', list); 
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

      if (!newId) {
        throw new Error("Create returned invalid id");
      }

      console.log("Created successfully");

      const updated = await getLists();
      const normalized = Array.isArray(updated)
        ? updated
        : Array.isArray(updated?.data)
          ? updated.data
          : [];

      setLists(normalized);
      setNewTitle('');
      setNewDesc('');
      setShowModal(false);
    } catch (err) {
      console.error('Error creating list:', err);
      alert('Create failed. Try again.');
    }
  }

  const handleDelete = async () => {
    if (!selectedList) return;

    if (!window.confirm("Are you sure you want to delete this list?")) return;

    try {
      const result = await deleteListById(selectedList.id);

      if (!result || result.status !== "deleted") {
        throw new Error("Delete failed");
      }

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
      console.error("Error deleting list:", err);
      alert("Failed to delete list. Try again.");
    }
  };

  const handleUpdate = async () => {
    if (!selectedList) return;

    console.log("clicked update");

    try {
      const result = await updateList(selectedList.id, editTitle, editDesc);

    if (!result || result.status !== "updated") {
      throw new Error("Update failed");
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
    <div className="min-h-screen bg-gray-100 py-10 px-4">
      <div className="max-w-2xl mx-auto bg-white shadow-md rounded-lg p-6">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-3xl font-bold text-gray-800">Your Todo Lists</h1>
          <button
            onClick={handleLogout}
            className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600"
          >
            Logout
          </button>
        </div>

        <button
          onClick={() => setShowModal(true)}
          className="bg-blue-600 text-white px-4 py-2 rounded mb-6 hover:bg-blue-700"
        >
          + Create New Todo List
        </button>

        <ul className="space-y-3">
          {lists.map((list) => (
            <li key={list.id}>
              <button
                id={`list-${list.id}`}
                onClick={() => openEditModal(list.id)}
                disabled={loadingList}
                className="w-full text-left px-4 py-3 bg-gray-200 rounded hover:bg-gray-300 disabled:opacity-50"
              >
                {list.title}
              </button>
            </li>
          ))}
        </ul>
      </div>

      {isModalOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h2 className="text-xl font-semibold mb-4">Edit Todo List</h2>

            <input
              type="text"
              value={editTitle}
              onChange={(e) => setEditTitle(e.target.value)}
              placeholder="Title"
              className="w-full mb-3 px-3 py-2 border rounded"
            />
            <textarea
              value={editDesc}
              onChange={(e) => setEditDesc(e.target.value)}
              placeholder="Description (optional)"
              rows={4}
              className="w-full mb-4 px-3 py-2 border rounded"
            />

            <div className="flex justify-end space-x-3">
              <button
                onClick={handleUpdate}
                className="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600"
              >
                Update
              </button>
              <button
                onClick={handleDelete}
                className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600"
              >
                Delete
              </button>
              <button
                onClick={() => setIsModalOpen(false)}
                className="bg-gray-300 px-4 py-2 rounded hover:bg-gray-400"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}

      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h2 className="text-xl font-semibold mb-4">Create New Todo List</h2>

            <input
              type="text"
              placeholder="Title"
              value={newTitle}
              onChange={(e) => setNewTitle(e.target.value)}
              className="w-full mb-3 px-3 py-2 border rounded"
            />
            <textarea
              placeholder="Description (optional)"
              value={newDesc}
              onChange={(e) => setNewDesc(e.target.value)}
              className="w-full mb-4 px-3 py-2 border rounded"
            />

            <div className="flex justify-end space-x-3">
              <button
                onClick={handleCreate}
                className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
              >
                Create
              </button>
              <button
                onClick={() => setShowModal(false)}
                className="bg-gray-300 px-4 py-2 rounded hover:bg-gray-400"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );

}


export default TodoLists;