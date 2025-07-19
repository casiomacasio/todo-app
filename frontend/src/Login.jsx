import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { signIn } from './api';

function Login({onLogin}) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleLogin = async (e) => {
    e.preventDefault();
    setError('');

    const success = await signIn(username, password);
    if (success) {
      onLogin(); 
      navigate('/todos');
    } else {
      setError('Login failed. Please check your credentials.');
    }
  };

  return (
    <div className="flex flex-col items-center p-6 gap-3">
      <h2 className="text-2xl mb-2">Login</h2>
      {error && <p className="text-red-500">{error}</p>}
      <form onSubmit={handleLogin} className="flex flex-col gap-3 w-64">
        <input
          type="text"
          placeholder="Username"
          className="p-2 border rounded"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
        />
        <input
          type="password"
          placeholder="Password"
          className="p-2 border rounded"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <button
          type="submit"
          className="bg-blue-600 text-white p-2 rounded"
        >
          Login
        </button>
      </form>
      <button
        onClick={() => navigate('/signup')}
        className="bg-green-600 text-white p-2 rounded w-64"
      >
        Sign Up
      </button>
    </div>
  );
}

export default Login;