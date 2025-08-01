const API_BASE = 'http://localhost:8000';

export async function createList(title, description = '') {
  const res = await fetch(`${API_BASE}/api/lists`, {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      title,
      description,
    }),
  });
  if (!res.ok) {
    throw new Error('Failed to create list');
  }
  const data = await res.json();
  return data.id;
}

export async function deleteListById(id) {
  const response = await fetch(`${API_BASE}/api/lists/${id}`, {
    method: "DELETE",
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error("Failed to delete list");
  }

  return await response.json(); 
}


export async function getLists() {
  const res = await fetch(`${API_BASE}/api/lists`, {
    method: 'GET',
    credentials: 'include',
  });
  if (!res.ok) return null;
  return res.json();
}

export async function getListById(id) {
  const res = await fetch(`${API_BASE}/api/lists/${id}`, {
    method: 'GET',
    credentials: 'include',
  });

  if (!res.ok) {
    console.error('Failed to fetch list by id:', res.status);
    return null;
  }

  try {
    const json = await res.json();
    return json?.data;
  } catch (err) {
    console.error('Error parsing JSON in getListById:', err);
    return null;
  }
}


export async function signIn(username, password) {
  const res = await fetch(`${API_BASE}/auth/sign-in`, {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  });
  return res.ok;
}

export async function signUp(name, username, password) {
  const res = await fetch(`${API_BASE}/auth/sign-up`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({name, username, password }),
  });
  return res.ok;
}

export async function updateList(id, title, description = '') {
  const res = await fetch(`${API_BASE}/api/lists/${id}`, {
    method: 'PUT',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ title, description }),
  });

  if (!res.ok) {
    throw new Error('Failed to update');
  }

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