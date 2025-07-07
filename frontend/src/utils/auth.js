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

// 解析JWT获取用户信息
export function getUserId() {
  const token = localStorage.getItem('token');
  if (!token) return null;
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload.user_id;
  } catch {
    return null;
  }
} 