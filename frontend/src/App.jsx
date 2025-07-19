import { useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { checkAuth } from './api';
import Login from './Login';
import Signup from './Signup';
import Todos from './TodoLists';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(null);

  useEffect(() => {
    (async () => {
      try {
        const ok = await checkAuth();
        setIsAuthenticated(ok);
      } catch {
        setIsAuthenticated(false);
      }
    })();
  }, []);

  if (isAuthenticated === null) return <p>Loading...</p>;

  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login onLogin={() => setIsAuthenticated(true)} />} />
        <Route path="/signup" element={<Signup />} />
        <Route
          path="/todos"
          element={isAuthenticated ? <Todos onLogout={() => setIsAuthenticated(false)}/> : <Navigate to="/login" />}
        />
        <Route
          path="/"
          element={<Navigate to={isAuthenticated ? '/todos' : '/login'} />}
        />
      </Routes>
    </Router>
  );
}

export default App;
