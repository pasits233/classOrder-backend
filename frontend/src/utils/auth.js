export function getRole() {
  return localStorage.getItem('role');
}

export function setRole(role) {
  localStorage.setItem('role', role);
}

export function logout() {
  localStorage.removeItem('role');
  localStorage.removeItem('token');
} 