import React, { useState, useEffect } from 'react';
import Login from './Login';
import TodoLists from './TodoLists';
import { checkAuth, initAutoRefresh } from './api';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    checkAuth().then(ok => {
      setIsAuthenticated(ok);
      if (ok) initAutoRefresh();
    });
  }, []);

  return isAuthenticated ? (
    <TodoLists onLogout={() => setIsAuthenticated(false)} />
  ) : (
    <Login onLogin={() => {
      setIsAuthenticated(true);
      initAutoRefresh();
    }} />
  );
}

export default App;
