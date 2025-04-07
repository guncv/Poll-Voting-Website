'use client';

import React, { useState } from 'react';

// Import the mock calls for single/multiple
import {
  fetchQuestionByIDMock,
  fetchAllQuestionsMock,
  // Import the real call for last archived
  fetchLastArchivedQuestion,
} from '../../utils/api';

interface Question {
  question_id: string;
  archive_date: string;
  question_text: string;
  first_choice: string;
  second_choice: string;
  total_participants: number;
  first_choice_count: number;
  second_choice_count: number;
  created_by: string;
  created_at: string;
}

export default function QuestionPage() {
  // We store an access token for the "popular" tab
  const [accessToken, setAccessToken] = useState('');

  // Tab control
  const [tab, setTab] = useState<'single' | 'multiple' | 'continuous' | 'popular'>('single');

  // SINGLE
  const [singleID, setSingleID] = useState('');
  const [singleQuestion, setSingleQuestion] = useState<Question | null>(null);
  const [singleError, setSingleError] = useState('');

  // MULTIPLE
  const [allQuestions, setAllQuestions] = useState<Question[]>([]);
  const [multiError, setMultiError] = useState('');

  // CONTINUOUS
  const [continuousQuestions, setContinuousQuestions] = useState<Question[]>([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [contError, setContError] = useState('');

  // POPULAR
  const [popularQuestion, setPopularQuestion] = useState<Question | null>(null);
  const [popError, setPopError] = useState('');

  // SINGLE (mock)
  async function handleFetchSingle() {
    setSingleQuestion(null);
    setSingleError('');

    if (!singleID) {
      setSingleError('Please enter a question ID');
      return;
    }

    try {
      const q = await fetchQuestionByIDMock(singleID);
      setSingleQuestion(q);
    } catch (err: any) {
      setSingleError(err.message || 'Error fetching single question');
    }
  }

  // MULTIPLE (mock)
  async function handleFetchAll() {
    setAllQuestions([]);
    setMultiError('');

    try {
      let questions = await fetchAllQuestionsMock();
      // Sort ascending by date
      questions.sort((a, b) =>
        new Date(a.archive_date).getTime() - new Date(b.archive_date).getTime()
      );
      setAllQuestions(questions);
    } catch (err: any) {
      setMultiError(err.message || 'Error fetching all questions');
    }
  }

  // CONTINUOUS (mock)
  async function handleFetchForContinuous() {
    setContinuousQuestions([]);
    setCurrentIndex(0);
    setContError('');

    try {
      let questions = await fetchAllQuestionsMock();
      questions.sort((a, b) =>
        new Date(a.archive_date).getTime() - new Date(b.archive_date).getTime()
      );
      setContinuousQuestions(questions);
      setCurrentIndex(0);
    } catch (err: any) {
      setContError(err.message || 'Error fetching questions for continuous mode');
    }
  }

  function handleAnswerQuestion() {
    setCurrentIndex((prev) => prev + 1);
  }

  // POPULAR (real)
  async function handleFetchPopular() {
    setPopularQuestion(null);
    setPopError('');

    if (!accessToken) {
      setPopError('No access token provided.');
      return;
    }

    try {
      const q = await fetchLastArchivedQuestion(accessToken);
      setPopularQuestion(q);
    } catch (err: any) {
      setPopError(err.message || 'Error fetching most popular question');
    }
  }

  // current question in continuous mode
  const currentQuestion = continuousQuestions[currentIndex];

  return (
    <div style={{ padding: '1rem' }}>
      <h1>Question Demo Page</h1>

      <div style={{ marginBottom: '1rem' }}>
        <label>Access Token (for "Popular" tab): </label>
        <input
          type="text"
          placeholder="Paste your JWT"
          value={accessToken}
          onChange={(e) => setAccessToken(e.target.value)}
          style={{ width: '300px' }}
        />
      </div>

      <div style={{ marginBottom: '1rem' }}>
        <button onClick={() => setTab('single')}>Single</button>
        <button onClick={() => setTab('multiple')}>Multiple</button>
        <button onClick={() => setTab('continuous')}>Continuous</button>
        <button onClick={() => setTab('popular')}>Popular (Yesterday)</button>
      </div>

      {/* SINGLE MODE (Mock) */}
      {tab === 'single' && (
        <div style={{ border: '1px solid #ccc', padding: '1rem' }}>
          <h2>Single Question (Mock)</h2>
          <p>Enter a question ID (e.g., "mock-q1", "mock-q2", or "mock-q3"):</p>
          <input
            type="text"
            value={singleID}
            onChange={(e) => setSingleID(e.target.value)}
          />
          <button onClick={handleFetchSingle}>Fetch</button>

          {singleError && <p style={{ color: 'red' }}>{singleError}</p>}

          {singleQuestion && (
            <div style={{ marginTop: '1rem' }}>
              <p><strong>ID:</strong> {singleQuestion.question_id}</p>
              <p><strong>Date:</strong> {singleQuestion.archive_date}</p>
              <p><strong>Text:</strong> {singleQuestion.question_text}</p>
              <p><strong>1st Choice:</strong> {singleQuestion.first_choice}</p>
              <p><strong>2nd Choice:</strong> {singleQuestion.second_choice}</p>
              <p><strong>Participants:</strong> {singleQuestion.total_participants}</p>
              <p><strong>1st Choice Count:</strong> {singleQuestion.first_choice_count}</p>
              <p><strong>2nd Choice Count:</strong> {singleQuestion.second_choice_count}</p>
              <p><strong>Created By:</strong> {singleQuestion.created_by}</p>
              <p><strong>Created At:</strong> {singleQuestion.created_at}</p>
            </div>
          )}
        </div>
      )}

      {/* MULTIPLE MODE (Mock) */}
      {tab === 'multiple' && (
        <div style={{ border: '1px solid #ccc', padding: '1rem' }}>
          <h2>All Questions (Mock, Sorted by Date)</h2>
          <button onClick={handleFetchAll}>Fetch All</button>
          {multiError && <p style={{ color: 'red' }}>{multiError}</p>}

          {allQuestions.length > 0 && (
            <ul style={{ marginTop: '1rem' }}>
              {allQuestions.map((q) => (
                <li key={q.question_id}>
                  <strong>ID:</strong> {q.question_id} &nbsp;
                  <strong>Date:</strong> {q.archive_date} &nbsp;
                  <strong>Text:</strong> {q.question_text}
                </li>
              ))}
            </ul>
          )}
        </div>
      )}

      {/* CONTINUOUS MODE (Mock) */}
      {tab === 'continuous' && (
        <div style={{ border: '1px solid #ccc', padding: '1rem' }}>
          <h2>Continuous (Mock)</h2>
          <button onClick={handleFetchForContinuous}>Start</button>
          {contError && <p style={{ color: 'red' }}>{contError}</p>}

          {continuousQuestions.length > 0 && currentIndex < continuousQuestions.length && (
            <div style={{ marginTop: '1rem' }}>
              <p><strong>ID:</strong> {currentQuestion.question_id}</p>
              <p><strong>Date:</strong> {currentQuestion.archive_date}</p>
              <p><strong>Text:</strong> {currentQuestion.question_text}</p>
              <button onClick={handleAnswerQuestion}>Answer & Next</button>
            </div>
          )}
          {continuousQuestions.length > 0 && currentIndex >= continuousQuestions.length && (
            <p>All questions answered!</p>
          )}
        </div>
      )}

      {/* POPULAR (REAL) */}
      {tab === 'popular' && (
        <div style={{ border: '1px solid #ccc', padding: '1rem' }}>
          <h2>Most Popular (Yesterday) [REAL]</h2>
          <button onClick={handleFetchPopular}>Get Last Archived</button>
          {popError && <p style={{ color: 'red' }}>{popError}</p>}

          {popularQuestion && (
            <div style={{ marginTop: '1rem' }}>
              <p><strong>ID:</strong> {popularQuestion.question_id}</p>
              <p><strong>Date:</strong> {popularQuestion.archive_date}</p>
              <p><strong>Text:</strong> {popularQuestion.question_text}</p>
              <p><strong>1st Choice:</strong> {popularQuestion.first_choice}</p>
              <p><strong>2nd Choice:</strong> {popularQuestion.second_choice}</p>
              <p><strong>Participants:</strong> {popularQuestion.total_participants}</p>
              <p><strong>1st Choice Count:</strong> {popularQuestion.first_choice_count}</p>
              <p><strong>2nd Choice Count:</strong> {popularQuestion.second_choice_count}</p>
              <p><strong>Created By:</strong> {popularQuestion.created_by}</p>
              <p><strong>Created At:</strong> {popularQuestion.created_at}</p>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
