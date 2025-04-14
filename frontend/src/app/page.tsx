'use client';

import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { getProfile, logoutUser, fetchLastArchivedQuestion } from '../utils/api';
import { buttonStyle } from '../components/AuthenticationStyle';

export default function HomePage() {
  const router = useRouter();
  const [profile, setProfile] = useState<any>(null);
  const [lastQuestion, setLastQuestion] = useState<any>(null);
  const [error, setError] = useState('');

  const containerStyle: React.CSSProperties = {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    height: '100vh',
    textAlign: 'center',
  };

  // Style for the last archived question card
  const questionCardStyle: React.CSSProperties = {
    border: '1px solid #ccc',
    padding: '1rem',
    margin: '1rem 0',
    borderRadius: '8px',
    textAlign: 'left',
    width: '80%',
    maxWidth: '500px',
  };

  useEffect(() => {
    async function loadProfile() {
      const token = localStorage.getItem('accessToken');
      if (!token) {
        router.push('/login');
        return;
      }
      try {
        const userData = await getProfile();
        setProfile(userData);
      } catch (err: any) {
        console.error('[HomePage] Error fetching profile:', err.message);
        router.push('/login');
      }
    }

    async function loadLastQuestion() {
      try {
        const questionData = await fetchLastArchivedQuestion();
        setLastQuestion(questionData);
      } catch (err: any) {
        console.error('[HomePage] Error fetching last archived question:', err.message);
      }
    }

    loadProfile();
    loadLastQuestion();
  }, [router]);

  async function handleLogout() {
    setError('');
    try {
      await logoutUser();
      localStorage.removeItem('accessToken');
      router.push('/login');
    } catch (err: any) {
      console.error('[HomePage] Logout error:', err.message);
      setError(err.message || 'Logout failed');
    }
  }

  if (error) {
    return <p style={{ color: 'red' }}>{error}</p>;
  }

  if (!profile) {
    return <p>Loading your profile...</p>;
  }

  return (
    <div style={containerStyle}>
      <h1>Welcome to the Home Page</h1>
      <p>
        Hello, <strong>{profile.email ?? 'User'}</strong>!
      </p>

      {lastQuestion && (
        <div style={questionCardStyle}>
          <h2>Popular Question</h2>
          <p>
            <strong>Question From Date:</strong>{' '}
            {new Date(lastQuestion.archive_date).toLocaleDateString()}
          </p>
          <p>
            <strong>Question:</strong> {lastQuestion.question_text}
          </p>
          <p>
            <strong>{lastQuestion.first_choice}:</strong> {lastQuestion.first_choice_count} votes
          </p>
          <p>
            <strong>{lastQuestion.second_choice}:</strong> {lastQuestion.second_choice_count} votes
          </p>
          {/* <p>
            <strong>Total Participants:</strong> {lastQuestion.total_participants}
          </p> */}
        </div>
      )}

      <div style={{ marginTop: '1rem', display: 'flex', gap: '1rem' }}>
        <button style={buttonStyle} onClick={() => router.push('/create')}>
          Create Question
        </button>
        <button style={buttonStyle} onClick={() => router.push('/voting')}>
          Go to Voting Page
        </button>
        <button style={buttonStyle} onClick={handleLogout}>
          Logout
        </button>
      </div>
    </div>
  );
}
