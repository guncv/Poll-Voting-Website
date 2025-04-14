'use client';

import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import {
  fetchCacheToday,
  fetchCacheQuestionByID,
  voteOnQuestion,
  getProfile,
} from '../../utils/api';
import { buttonStyle } from '../../components/AuthenticationStyle';
import styles from '../../components/VotingPage.module.css';

interface CacheQuestion {
  question_id: string;
  user_id: string;
  text: string;
  first_choice: string;
  second_choice: string;
  total_participants: number;
  first_choice_count: number;
  second_choice_count: number;
  milestones: string;
  follow_ups: string;
  group_id: string;
}

export default function VotingPage() {
  const router = useRouter();

  const [tab, setTab] = useState<'random' | 'single' | 'all'>('random');

  const [userID, setUserID] = useState('');
  const [userEmail, setUserEmail] = useState('');
  const [profileError, setProfileError] = useState('');

  const [randomQ, setRandomQ] = useState<CacheQuestion | null>(null);
  const [randomError, setRandomError] = useState('');
  const [hasVotedRandom, setHasVotedRandom] = useState(false);
  const [randomSuccess, setRandomSuccess] = useState('');

  const [singleID, setSingleID] = useState('');
  const [singleQ, setSingleQ] = useState<CacheQuestion | null>(null);
  const [singleError, setSingleError] = useState('');
  const [hasVotedSingle, setHasVotedSingle] = useState(false);
  const [singleSuccess, setSingleSuccess] = useState('');

  const [allQuestions, setAllQuestions] = useState<CacheQuestion[]>([]);
  const [allError, setAllError] = useState('');

  useEffect(() => {
    async function loadProfile() {
      const token = localStorage.getItem('accessToken');
      if (!token) {
        router.push('/login');
        return;
      }
      try {
        const userData = await getProfile();
        setUserID(userData.user_id);
        setUserEmail(userData.email);
      } catch (err: any) {
        setProfileError(err.message || 'Cannot load user profile.');
        router.push('/login');
      }
    }
    loadProfile();
  }, [router]);

  const headerStyle: React.CSSProperties = {
    position: 'fixed',
    top: 0,
    left: 0,
    right: 0,
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: '1rem',
    backgroundColor: '#FF98A0', 
    borderBottom: '1px solid #ccc',
    zIndex: 1000,
  };

  // =============== RANDOM TAB ===============
  async function handleFetchRandom() {
    setRandomError('');
    setRandomQ(null);
    setHasVotedRandom(false);
    setRandomSuccess('');

    if (!userID) {
      setRandomError('No user ID found. Please log in first.');
      return;
    }
    try {
      const data = await fetchCacheToday(); 
      if (!data.questions?.length) {
        throw new Error('No questions found for today.');
      }
      const randIdx = Math.floor(Math.random() * data.questions.length);
      setRandomQ(data.questions[randIdx]);
    } catch (err: any) {
      setRandomError(err.message || 'Error fetching random question.');
    }
  }

  async function handleVoteRandom(isFirst: boolean) {
    if (!randomQ) return;
    setRandomSuccess('');
    try {
      const response = await voteOnQuestion({
        question_id: randomQ.question_id,
        is_first_choice: isFirst,
        user_id: userID,
      });
      if (response.already_voted) {
        setRandomError('You have already voted on this question!');
      } else {
        setHasVotedRandom(true);
        setRandomSuccess('Vote success!');
        const updated = await fetchCacheQuestionByID(randomQ.question_id);
        setRandomQ(updated);
      }
    } catch (err: any) {
      setRandomError(err.message || 'Error voting on random question.');
    }
  }

  function renderRandomChoiceLabel(isFirst: boolean) {
    if (!randomQ) return '';
    const label = isFirst ? randomQ.first_choice : randomQ.second_choice;
    const count = isFirst ? randomQ.first_choice_count : randomQ.second_choice_count;
    if (!hasVotedRandom) {
      return label;
    }
    return `${label} (${count})`;
  }

  // =============== SINGLE TAB ===============
  async function handleFetchSingle() {
    setSingleError('');
    setSingleQ(null);
    setHasVotedSingle(false);
    setSingleSuccess('');

    if (!singleID) {
      setSingleError('Please enter a question ID.');
      return;
    }
    if (!userID) {
      setSingleError('No user ID found. Please log in first.');
      return;
    }
    try {
      const question = await fetchCacheQuestionByID(singleID);
      setSingleQ(question);
    } catch (err: any) {
      setSingleError(err.message || 'Error fetching single question.');
    }
  }

  async function handleVoteSingle(isFirst: boolean) {
    if (!singleQ) return;
    setSingleSuccess('');
    try {
      const response = await voteOnQuestion({
        question_id: singleQ.question_id,
        is_first_choice: isFirst,
        user_id: userID,
      });
      if (response.already_voted) {
        setSingleError('You have already voted on this question!');
      } else {
        setHasVotedSingle(true);
        setSingleSuccess('Vote success!');
        const updated = await fetchCacheQuestionByID(singleQ.question_id);
        setSingleQ(updated);
      }
    } catch (err: any) {
      setSingleError(err.message || 'Error voting on single question.');
    }
  }

  function renderSingleChoiceLabel(isFirst: boolean) {
    if (!singleQ) return '';
    const label = isFirst ? singleQ.first_choice : singleQ.second_choice;
    const count = isFirst ? singleQ.first_choice_count : singleQ.second_choice_count;
    if (!hasVotedSingle) {
      return label;
    }
    return `${label} (${count})`;
  }

  // =============== ALL TAB ===============
  async function handleFetchAll() {
    setAllError('');
    setAllQuestions([]);

    if (!userID) {
      setAllError('No user ID found. Please log in first.');
      return;
    }
    try {
      const data = await fetchCacheToday();
      if (!data.questions) {
        throw new Error('Invalid response. Expected { "questions": [...] }.');
      }
      setAllQuestions(data.questions);
    } catch (err: any) {
      setAllError(err.message || 'Error fetching all questions.');
    }
  }

  async function handleSelectQuestion(qid: string) {
    setTab('single');
    setSingleID(qid);
    setSingleError('');
    setSingleQ(null);
    setHasVotedSingle(false);
    setSingleSuccess('');
    try {
      const question = await fetchCacheQuestionByID(qid);
      setSingleQ(question);
    } catch (err: any) {
      setSingleError(err.message || 'Error loading selected question.');
    }
  }

  return (
    <div className={styles.container} style={{ paddingTop: '100px' }}>
      <header style={headerStyle}>
        <button style={buttonStyle} onClick={() => router.back()}>
          Back
        </button>
        <button style={buttonStyle} onClick={() => router.push('/')}>
          Home
        </button>
      </header>

      <h1 style={{ marginBottom: '1rem' }}>Voting Page</h1>
      <p>User Email: {userEmail || '(loading...)'}</p>
      {profileError && <p className={styles.error}>{profileError}</p>}

      <div className={styles.tabBar}>
        <button style={buttonStyle} onClick={() => setTab('random')}>
          Random
        </button>
        <button style={buttonStyle} onClick={() => setTab('single')}>
          Single
        </button>
        <button style={buttonStyle} onClick={() => setTab('all')}>
          All
        </button>
      </div>

      {tab === 'random' && (
        <div className={styles.tabContainer}>
          <h2 style={{ textAlign: 'center' }}>Random Question</h2>
          <button style={buttonStyle} onClick={handleFetchRandom}>
            Fetch Random
          </button>
          {randomError && <p className={styles.error}>{randomError}</p>}
          {randomSuccess && <p style={{ color: 'green', marginTop: '1rem' }}>{randomSuccess}</p>}
          {randomQ && (
            <div style={{ marginTop: '1rem' }}>
              <p><strong>Question:</strong> {randomQ.text}</p>
              <p><strong>Participants:</strong> {randomQ.total_participants}</p>
              <p>
                <strong>First Choice:</strong> {renderRandomChoiceLabel(true)}
              </p>
              <p>
                <strong>Second Choice:</strong> {renderRandomChoiceLabel(false)}
              </p>
              <div
                style={{
                  marginTop: '1rem',
                  display: 'flex',
                  gap: '1rem',
                  justifyContent: 'center',
                }}
              >
                <button style={buttonStyle} onClick={() => handleVoteRandom(true)}>
                  Vote {randomQ.first_choice}
                </button>
                <button style={buttonStyle} onClick={() => handleVoteRandom(false)}>
                  Vote {randomQ.second_choice}
                </button>
              </div>
            </div>
          )}
        </div>
      )}

      {tab === 'single' && (
        <div className={styles.tabContainer}>
          <h2 style={{ textAlign: 'center' }}>Single Question</h2>
          <div
            style={{
              margin: '1rem auto',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between',
              maxWidth: '500px',
            }}
          >
            <label style={{ marginRight: '1rem', fontWeight: '500' }}>
              Question ID:
            </label>
            <input
              type="text"
              placeholder="Enter question ID"
              value={singleID}
              onChange={(e) => setSingleID(e.target.value)}
              style={{ width: '250px', textAlign: 'left', padding: '5px' }}
            />
          </div>
          <button style={buttonStyle} onClick={handleFetchSingle}>
            Fetch Single
          </button>
          {singleError && <p className={styles.error}>{singleError}</p>}
          {singleSuccess && <p style={{ color: 'green', marginTop: '1rem' }}>{singleSuccess}</p>}
          {singleQ && (
            <div style={{ marginTop: '1rem' }}>
              <p><strong>Question:</strong> {singleQ.text}</p>
              <p><strong>Participants:</strong> {singleQ.total_participants}</p>
              <p>
                <strong>First Choice:</strong> {renderSingleChoiceLabel(true)}
              </p>
              <p>
                <strong>Second Choice:</strong> {renderSingleChoiceLabel(false)}
              </p>
              <div
                style={{
                  marginTop: '1rem',
                  display: 'flex',
                  gap: '1rem',
                  justifyContent: 'center',
                }}
              >
                <button style={buttonStyle} onClick={() => handleVoteSingle(true)}>
                  Vote {singleQ.first_choice}
                </button>
                <button style={buttonStyle} onClick={() => handleVoteSingle(false)}>
                  Vote {singleQ.second_choice}
                </button>
              </div>
            </div>
          )}
        </div>
      )}

      {tab === 'all' && (
        <div className={styles.tabContainer}>
          <h2 style={{ textAlign: 'center' }}>All Questions (Today)</h2>
          <button style={buttonStyle} onClick={handleFetchAll}>
            Fetch All
          </button>
          {allError && <p className={styles.error}>{allError}</p>}
          {allQuestions.length > 0 && (
            <div style={{ marginTop: '1rem' }}>
              {allQuestions.map((q) => (
                <div key={q.question_id} className={styles.card}>
                  <p><strong>Question:</strong> {q.text}</p>
                  <p><strong>Participants:</strong> {q.total_participants}</p>
                  <p><strong>First Choice:</strong> {q.first_choice}</p>
                  <p><strong>Second Choice:</strong> {q.second_choice}</p>
                  <button
                    style={{ ...buttonStyle, marginTop: '0.5rem' }}
                    onClick={() => handleSelectQuestion(q.question_id)}
                  >
                    View &amp; Vote
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
