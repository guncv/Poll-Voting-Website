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
