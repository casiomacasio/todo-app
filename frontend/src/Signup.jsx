import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { signUp } from './api';

function Signup() {
  const [name, setName] = useState('');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleSignup = async (e) => {
    e.preventDefault();
    setError('');

    try {
      const ok = await signUp(name, username, password);
      if (ok) {
        navigate('/login');
      } else {
        setError('Signup failed');
      }
    } catch (err) {
      setError('Something went wrong');
    }
  };

  return (
    <div className="flex flex-col items-center p-6 gap-3">
      <h2 className="text-2xl mb-2">Sign Up</h2>
      {error && <p className="text-red-500">{error}</p>}
      <form onSubmit={handleSignup} className="flex flex-col gap-3 w-64">
        <input
          type="text"
          placeholder="Name"
          className="p-2 border rounded"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
        />
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
          className="bg-green-600 text-white p-2 rounded"
        >
          Sign Up
        </button>
      </form>
      <button
        onClick={() => navigate('/login')}
        className="bg-blue-600 text-white p-2 rounded w-64"
      >
        Login
      </button>
    </div>
  );
}

export default Signup;