const API_BASE = '';

export async function signIn(username, password) {
  const res = await fetch(`${API_BASE}/auth/sign-in`, {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  });
  return res.ok;
}

export async function signUp(username, password) {
  const res = await fetch(`${API_BASE}/auth/sign-up`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  });
  return res.ok;
}

export async function getLists() {
  const res = await fetch(`${API_BASE}/api/lists`, {
    method: 'GET',
    credentials: 'include',
  });
  if (!res.ok) return null;
  return res.json();
}

export async function logout() {
  await fetch(`${API_BASE}/auth/logout`, {
    method: 'POST',
    credentials: 'include',
  });
}

export async function checkAuth() {
  const res = await fetch(`${API_BASE}/auth/refresh`, {
    method: 'POST',
    credentials: 'include',
  });
  return res.ok;
}

let refreshInterval;

export function initAutoRefresh(intervalMs = 4 * 60 * 1000) {
  if (refreshInterval) clearInterval(refreshInterval);
  refreshInterval = setInterval(async () => {
    await checkAuth();
  }, intervalMs);
}