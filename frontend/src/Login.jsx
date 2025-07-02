import React, { useState } from 'react';
import { signIn, signUp } from './api';

function Login({ onLogin }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [status, setStatus] = useState('');

  const handleSignIn = async () => {
    const ok = await signIn(username, password);
    if (ok) onLogin();
    else setStatus('Login failed');
  };

  const handleSignUp = async () => {
    const ok = await signUp(username, password);
    setStatus(ok ? 'Signup successful! Now sign in.' : 'Signup failed');
  };

  return (
    <div>
      <h1>Login / Signup</h1>
      <input value={username} onChange={e => setUsername(e.target.value)} placeholder="Username" />
      <input type="password" value={password} onChange={e => setPassword(e.target.value)} placeholder="Password" />
      <br />
      <button onClick={handleSignIn}>Sign In</button>
      <button onClick={handleSignUp}>Sign Up</button>
      <p>{status}</p>
    </div>
  );
}

export default Login;
