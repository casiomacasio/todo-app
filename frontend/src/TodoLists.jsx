import React, { useEffect, useState } from 'react';
import { getLists, logout } from './api';

function TodoLists({ onLogout }) {
  const [lists, setLists] = useState([]);

  useEffect(() => {
    getLists().then(data => {
      setLists(data || []);
    });
  }, []);

  const handleLogout = async () => {
    await logout();
    onLogout();
  };

  return (
    <div>
      <h1>Your Todo Lists</h1>
      <button onClick={handleLogout}>Logout</button>
      <ul>
        {lists.map(list => (
          <li key={list.id}>{list.title}</li>
        ))}
      </ul>
    </div>
  );
}

export default TodoLists;
