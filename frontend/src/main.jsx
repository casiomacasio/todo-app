import React, { useState, useEffect } from 'react';
import Login from './Login';
import TodoLists from './TodoLists';
import { checkAuth } from './api';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    checkAuth().then(authenticated => {
      setIsAuthenticated(authenticated);
    });
  }, []);

  return isAuthenticated ? (
    <TodoLists onLogout={() => setIsAuthenticated(false)} />
  ) : (
    <Login onLogin={() => setIsAuthenticated(true)} />
  );
}

export default App;
