import { apiRequest, setAccessToken, API_BASE } from './apiClient';

export async function fetchLastArchivedQuestion(): Promise<any> {
  return apiRequest(`${API_BASE}/question/last`, { method: 'GET' });
}

export async function loginUser(
  email: string,
  password: string
): Promise<{ access_token: string }> {
  const res = await fetch(`${API_BASE}/user/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
    credentials: 'include', 
  });

  const data = await res.json();
  if (!res.ok) {
    throw new Error(data.error || 'Login failed');
  }

  setAccessToken(data.access_token);
  return data;
}

export async function getProfile(): Promise<any> {
  return apiRequest(`${API_BASE}/user/profile`, { method: 'GET' });
}

export async function logoutUser(): Promise<void> {
  const res = await fetch(`${API_BASE}/user/logout`, {
    method: 'GET',
    credentials: 'include',
  });
  if (!res.ok) {
    const errorData = await res.json().catch(() => ({}));
    throw new Error(errorData.error || 'Logout failed');
  }
  setAccessToken('');
}

export async function registerUser(
  email: string,
  password: string
): Promise<any> {
  const res = await fetch(`${API_BASE}/user/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
    credentials: 'include',
  });

  const data = await res.json();
  if (!res.ok) {
    throw new Error(data.error || 'Registration failed');
  }
  return data;
}

export async function fetchCacheToday(): Promise<any> {
  return apiRequest(`${API_BASE}/question/cache/today`, { method: 'GET' });
}

export async function fetchCacheQuestionByID(id: string): Promise<any> {
  return apiRequest(`${API_BASE}/question/cache/${id}`, { method: 'GET' });
}

export async function voteOnQuestion(payload: {
  question_id: string;
  is_first_choice: boolean;
  user_id: string;
}): Promise<{
  question_id: string;
  total_participants: number;
  first_choice_count: number;
  second_choice_count: number;
  newly_revealed_ids: string[];
  already_voted: boolean;
}> {
  return apiRequest(`${API_BASE}/question/vote`, {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function createQuestion(payload: {
  text: string;
  first_choice: string;
  second_choice: string;
  milestones?: string;
  follow_ups?: string;
  group_id?: string;
}): Promise<any> {
  return apiRequest(`${API_BASE}/question/cache`, {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}