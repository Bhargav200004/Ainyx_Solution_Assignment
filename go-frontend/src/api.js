const BASE = '/api';

async function request(path, options = {}) {
  const res = await fetch(`${BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  });

  if (res.status === 204) return null;

  const data = await res.json();

  if (!res.ok) {
    const msg = data.error || 'Something went wrong';
    const err = new Error(msg);
    err.details = data.details;
    throw err;
  }

  return data;
}

export function listUsers(page = 1, limit = 10) {
  return request(`/users?page=${page}&limit=${limit}`);
}

export function getUserById(id) {
  return request(`/users/${id}`);
}

export function createUser(name, dob) {
  return request('/users', {
    method: 'POST',
    body: JSON.stringify({ name, dob }),
  });
}

export function updateUser(id, name, dob) {
  return request(`/users/${id}`, {
    method: 'PUT',
    body: JSON.stringify({ name, dob }),
  });
}

export function deleteUser(id) {
  return request(`/users/${id}`, { method: 'DELETE' });
}
