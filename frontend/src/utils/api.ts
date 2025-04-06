import axios from 'axios';
import { PollData } from '../types/pollData';

const API_BASE = 'http://localhost:8000/api';

// export async function fetchPoll(): Promise<PollData> {
//   const { data } = await axios.get<PollData>(`${API_BASE}/poll/random`);
//   return data;
// }

// export async function sendVote(pollId: string, choice: string) {
//   const { data } = await axios.post(`${API_BASE}/poll/vote`, {
//     pollId,
//     choice,
//   });
//   return data;
// }

const MOCK_POLLS: PollData[] = [
  {
    id: 'poll1',
    question: 'Which do you prefer for frontend development?',
    choices: ['React', 'Vue'],
  },
  {
    id: 'poll2',
    question: 'Which programming language do you enjoy most?',
    choices: ['TypeScript', 'JavaScript'],
  },
  {
    id: 'poll3',
    question: 'Favorite state management tool?',
    choices: ['Redux', 'Zustand'],
  },
  {
    id: 'poll4',
    question: 'Best CSS solution for large apps?',
    choices: ['Tailwind CSS', 'Styled Components'],
  },
  {
    id: 'poll5',
    question: 'Preferred backend stack?',
    choices: ['Node.js', 'Django'],
  },
];

export async function fetchPoll(): Promise<PollData> {
  const random = MOCK_POLLS[Math.floor(Math.random() * MOCK_POLLS.length)];
  console.log('[MOCK] Fetching poll...');
  return new Promise((resolve) => {
    setTimeout(() => resolve(random), 500);
  });
}

export async function sendVote(pollId: string, choice: string) {
  console.log(`[MOCK] Sending vote for ${choice} in poll ${pollId}`);
  return new Promise((resolve) => {
    setTimeout(() => resolve({ success: true }), 500);
  });
}

export async function loginUser(email: string, password: string): Promise<{ access_token: string }> {
  const response = await fetch('http://localhost:8080/api/user/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
    credentials: 'include', // Ensures cookies (for the refresh token) are sent/received.
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || 'Login failed');
  }

  return data; // Expects an object containing { access_token: string }
}

export async function getProfile(accessToken: string): Promise<any> {
  const response = await fetch('http://localhost:8080/api/user/profile', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${accessToken}`, 
      'Content-Type': 'application/json'
    },
    credentials: 'include',
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.error || 'Failed to fetch profile');
  }

  return response.json();
}

export async function logoutUser(): Promise<void> {
  const response = await fetch('http://localhost:8080/api/user/logout', {
    method: 'GET', // or 'POST' if your endpoint is POST
    credentials: 'include', // ensures the refresh token cookie is sent
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.error || 'Logout failed');
  }
}