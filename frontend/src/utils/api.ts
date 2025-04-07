// src/utils/api.ts

import { PollData } from '../types/pollData';
import { apiRequest, setAccessToken } from './apiClient';

const API_BASE = 'http://localhost:8080/api';

/* ===================== Mock Data ===================== */
const MOCK_QUESTIONS = [
  {
    question_id: 'mock-q1',
    archive_date: '2025-03-20',
    question_text: 'Which color do you prefer?',
    first_choice: 'Red',
    second_choice: 'Blue',
    total_participants: 100,
    first_choice_count: 70,
    second_choice_count: 30,
    created_by: '00000000-0000-0000-0000-000000000000',
    created_at: '2025-03-20T10:00:00Z',
  },
  {
    question_id: 'mock-q2',
    archive_date: '2025-03-21',
    question_text: 'Cats or Dogs?',
    first_choice: 'Cats',
    second_choice: 'Dogs',
    total_participants: 50,
    first_choice_count: 20,
    second_choice_count: 30,
    created_by: '00000000-0000-0000-0000-000000000000',
    created_at: '2025-03-21T12:00:00Z',
  },
  {
    question_id: 'mock-q3',
    archive_date: '2025-03-22',
    question_text: 'Best programming language?',
    first_choice: 'Go',
    second_choice: 'Rust',
    total_participants: 80,
    first_choice_count: 40,
    second_choice_count: 40,
    created_by: '00000000-0000-0000-0000-000000000000',
    created_at: '2025-03-22T15:30:00Z',
  },
];

// Mock single-question fetch by ID
export async function fetchQuestionByIDMock(id: string): Promise<any> {
  const found = MOCK_QUESTIONS.find((q) => q.question_id === id);
  if (!found) {
    throw new Error(`Mock: Question with ID ${id} not found`);
  }
  return found;
}

// Mock all-questions fetch
export async function fetchAllQuestionsMock(): Promise<any[]> {
  return MOCK_QUESTIONS;
}

// Mock poll data
const MOCK_POLLS: PollData[] = [
  { id: 'poll1', question: 'Which do you prefer for frontend development?', choices: ['React', 'Vue'] },
  { id: 'poll2', question: 'Which programming language do you enjoy most?', choices: ['TypeScript', 'JavaScript'] },
  { id: 'poll3', question: 'Favorite state management tool?', choices: ['Redux', 'Zustand'] },
  { id: 'poll4', question: 'Best CSS solution for large apps?', choices: ['Tailwind CSS', 'Styled Components'] },
  { id: 'poll5', question: 'Preferred backend stack?', choices: ['Node.js', 'Django'] },
];

// Mock poll fetching
export async function fetchPoll(): Promise<PollData> {
  const random = MOCK_POLLS[Math.floor(Math.random() * MOCK_POLLS.length)];
  console.log('[MOCK] Fetching poll...');
  return new Promise((resolve) => {
    setTimeout(() => resolve(random), 500);
  });
}

// Mock sending a vote
export async function sendVote(pollId: string, choice: string) {
  console.log(`[MOCK] Sending vote for ${choice} in poll ${pollId}`);
  return new Promise((resolve) => {
    setTimeout(() => resolve({ success: true }), 500);
  });
}

/* ===================== Real API Calls ===================== */

// Real call for last archived question
export async function fetchLastArchivedQuestion(): Promise<any> {
  // We rely on apiRequest, which handles the token in headers & refresh
  return apiRequest(`${API_BASE}/question/last`);
}

// Real login endpoint
export async function loginUser(email: string, password: string): Promise<{ access_token: string }> {
  const res = await fetch(`${API_BASE}/user/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
    credentials: 'include', // so refresh cookie is stored
  });

  const data = await res.json();
  if (!res.ok) {
    throw new Error(data.error || 'Login failed');
  }

  // Store access token in local state
  setAccessToken(data.access_token);
  return data;
}

// Real get profile
export async function getProfile(): Promise<any> {
  return apiRequest(`${API_BASE}/user/profile`, { method: 'GET' });
}

// Real logout
export async function logoutUser(): Promise<void> {
  // This might not need apiRequest, because we don't need the token for /logout
  const res = await fetch(`${API_BASE}/user/logout`, {
    method: 'GET',
    credentials: 'include',
  });
  if (!res.ok) {
    const errorData = await res.json().catch(() => ({}));
    throw new Error(errorData.error || 'Logout failed');
  }
  setAccessToken(''); // Clear from local storage
}

// Real register
export async function registerUser(email: string, password: string): Promise<any> {
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
