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

  // Reordered: 'all', 'random', 'single'
  const [tab, setTab] = useState<'all' | 'random' | 'single'>('all');

  const [userID, setUserID] = useState('');
  const [userEmail, setUserEmail] = useState('');
  const [profileError, setProfileError] = useState('');

  // ------------------ ALL TAB ------------------
  const [allQuestions, setAllQuestions] = useState<CacheQuestion[]>([]);
  const [allError, setAllError] = useState('');

  // ------------------ RANDOM TAB ------------------
  const [randomQ, setRandomQ] = useState<CacheQuestion | null>(null);
  const [randomError, setRandomError] = useState('');
  const [hasVotedRandom, setHasVotedRandom] = useState(false);
  const [randomSuccess, setRandomSuccess] = useState('');

  // ------------------ SINGLE TAB ------------------
  const [singleID, setSingleID] = useState('');
  const [singleQ, setSingleQ] = useState<CacheQuestion | null>(null);
  const [singleError, setSingleError] = useState('');
  const [hasVotedSingle, setHasVotedSingle] = useState(false);
  const [singleSuccess, setSingleSuccess] = useState('');

  // ------------------ MILESTONE POPUP ------------------
  const [popupQuestion, setPopupQuestion] = useState<CacheQuestion | null>(null);
  const [showPopup, setShowPopup] = useState(false);
  const [milestoneMsg, setMilestoneMsg] = useState('');

  // ------------------ LOAD PROFILE ------------------
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

  // ------------------ HELPER FUNCTIONS ------------------
  function isMainQuestion(q: CacheQuestion) {
    const hasM = !!q.milestones?.trim();
    const hasF = !!q.follow_ups?.trim();
    const hasG = !!q.group_id?.trim();

    // (A) Standalone => no milestones, no follow_ups, no group_id
    if (!hasM && !hasF && !hasG) return true;
    // (B) Has milestones => definitely main
    if (hasM) return true;

    // Otherwise, milestone question
    return false;
  }

  function parseMilestones(mStr: string) {
    if (!mStr.trim()) return [];
    return mStr.split(',').map((part) => {
      const [scoreStr, qid] = part.split(':');
      return {
        score: parseInt(scoreStr.trim(), 10) || 0,
        question_id: qid,
      };
    });
  }

  // Check if main question’s total_participants triggers a milestone
  async function checkAndShowMilestone(mainQ: CacheQuestion) {
    const arr = parseMilestones(mainQ.milestones);
    if (!arr.length) return;

    const total = mainQ.total_participants;
    let best: null | { score: number; question_id: string } = null;
    for (let m of arr) {
      if (m.score <= total) {
        if (!best || m.score > best.score) {
          best = m;
        }
      }
    }
    if (!best) return;

    try {
      const milestoneQ = await fetchCacheQuestionByID(best.question_id);
      setPopupQuestion(milestoneQ);
      setShowPopup(true);
      setMilestoneMsg(`You've reached the milestone at ${best.score} votes!`);
    } catch (err: any) {
      console.error('Could not fetch milestone question:', err);
    }
  }

  // ------------------ ALL TAB LOGIC ------------------
  async function handleFetchAll() {
    setAllError('');
    setAllQuestions([]);

    if (!userID) {
      return; // avoid "No user ID" error
    }
    try {
      const data = await fetchCacheToday();
      if (!data.questions) {
        throw new Error('Invalid response. Expected { "questions": [...] }.');
      }
      const mainOnly = data.questions.filter(isMainQuestion);
      setAllQuestions(mainOnly);
    } catch (err: any) {
      setAllError(err.message || 'Error fetching all questions.');
    }
  }

  // Clicking "View & Vote" from the All tab => sets tab=single & singleID
  function handleSelectQuestion(qid: string) {
    setTab('single');
    setSingleID(qid);
    setSingleError('');
    setSingleQ(null);
    setHasVotedSingle(false);
    setSingleSuccess('');
  }

  // ------------------ RANDOM TAB LOGIC ------------------
  async function handleFetchRandom() {
    setRandomError('');
    setRandomQ(null);
    setHasVotedRandom(false);
    setRandomSuccess('');

    if (!userID) return; // skip if user isn't loaded
    try {
      const data = await fetchCacheToday();
      if (!data.questions?.length) {
        throw new Error('No questions found for today.');
      }
      const mainQs = data.questions.filter(isMainQuestion);
      if (!mainQs.length) {
        throw new Error('No main questions found (all are milestones?).');
      }
      const randIdx = Math.floor(Math.random() * mainQs.length);
      setRandomQ(mainQs[randIdx]);
    } catch (err: any) {
      setRandomError(err.message || 'Error fetching random question.');
    }
  }

  async function handleVoteRandom(isFirst: boolean) {
    if (!randomQ) return;
    setRandomSuccess('');
    try {
      const resp = await voteOnQuestion({
        question_id: randomQ.question_id,
        is_first_choice: isFirst,
        user_id: userID,
      });
      if (resp.already_voted) {
        setRandomError('You have already voted on this question!');
      } else {
        setHasVotedRandom(true);
        setRandomSuccess('Vote success!');
        const updated = await fetchCacheQuestionByID(randomQ.question_id);
        setRandomQ(updated);

        // Check milestone
        await checkAndShowMilestone(updated);
      }
    } catch (err: any) {
      setRandomError(err.message || 'Error voting on random question.');
    }
  }

  function renderRandomChoiceLabel(isFirst: boolean) {
    if (!randomQ) return '';
    const label = isFirst ? randomQ.first_choice : randomQ.second_choice;
    const count = isFirst ? randomQ.first_choice_count : randomQ.second_choice_count;
    if (!hasVotedRandom) return label;
    return `${label} (${count})`;
  }

  // ------------------ SINGLE TAB LOGIC ------------------
  async function handleFetchSingle() {
    setSingleError('');
    setSingleQ(null);
    setHasVotedSingle(false);
    setSingleSuccess('');

    if (!singleID) {
      setSingleError('Please enter a question ID.');
      return;
    }
    if (!userID) return; // skip if user not loaded

    try {
      const q = await fetchCacheQuestionByID(singleID);
      if (!isMainQuestion(q)) {
        throw new Error(`Question ${singleID} is a milestone, not a main question.`);
      }
      setSingleQ(q);
    } catch (err: any) {
      setSingleError(err.message || 'Error fetching single question.');
    }
  }

  async function handleVoteSingle(isFirst: boolean) {
    if (!singleQ) return;
    setSingleSuccess('');
    try {
      const resp = await voteOnQuestion({
        question_id: singleQ.question_id,
        is_first_choice: isFirst,
        user_id: userID,
      });
      if (resp.already_voted) {
        setSingleError('You have already voted on this question!');
      } else {
        setHasVotedSingle(true);
        setSingleSuccess('Vote success!');
        const updated = await fetchCacheQuestionByID(singleQ.question_id);
        setSingleQ(updated);
        await checkAndShowMilestone(updated);
      }
    } catch (err: any) {
      setSingleError(err.message || 'Error voting on single question.');
    }
  }

  function renderSingleChoiceLabel(isFirst: boolean) {
    if (!singleQ) return '';
    const label = isFirst ? singleQ.first_choice : singleQ.second_choice;
    const count = isFirst ? singleQ.first_choice_count : singleQ.second_choice_count;
    if (!hasVotedSingle) return label;
    return `${label} (${count})`;
  }

  // ------------------ MILESTONE POPUP ------------------
  async function handleVoteMilestone(isFirst: boolean) {
    if (!popupQuestion) return;
    try {
      const resp = await voteOnQuestion({
        question_id: popupQuestion.question_id,
        is_first_choice: isFirst,
        user_id: userID,
      });
      if (resp.already_voted) {
        alert('You already voted on this milestone question!');
      } else {
        await fetchCacheQuestionByID(popupQuestion.question_id);
        alert('Milestone vote success!');
      }
    } catch (err: any) {
      alert(err.message || 'Error voting on milestone question.');
    }
    setPopupQuestion(null);
    setShowPopup(false);
    setMilestoneMsg('');
  }

  function handleCloseMilestone() {
    setPopupQuestion(null);
    setShowPopup(false);
    setMilestoneMsg('');
  }

  // ------------------ AUTO-FETCH ON TAB CHANGES ------------------
  // 1) If tab=all and userID loaded => fetchAll
  useEffect(() => {
    if (tab === 'all' && userID) {
      handleFetchAll();
    }
  }, [tab, userID]);

  // 2) If tab=random and userID loaded => fetchRandom
  useEffect(() => {
    if (tab === 'random' && userID) {
      handleFetchRandom();
    }
  }, [tab, userID]);

  // 3) If tab=single, userID, and singleID => auto-fetch single
  //    So if user just came from All tab with a question ID, we auto load it
  useEffect(() => {
    if (tab === 'single' && userID && singleID.trim()) {
      handleFetchSingle();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [tab, userID, singleID]);

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

      {milestoneMsg && (
        <p style={{ color: 'green', fontWeight: 'bold', marginTop: '1rem' }}>
          {milestoneMsg}
        </p>
      )}

      <div className={styles.tabBar}>
        <button style={buttonStyle} onClick={() => setTab('all')}>
          All
        </button>
        <button style={buttonStyle} onClick={() => setTab('random')}>
          Random
        </button>
        <button style={buttonStyle} onClick={() => setTab('single')}>
          Single
        </button>
      </div>

      {/* ALL TAB */}
      {tab === 'all' && (
        <div className={styles.tabContainer}>
          <h2 style={{ textAlign: 'center' }}>All Questions (Today)</h2>
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

      {/* RANDOM TAB */}
      {tab === 'random' && (
        <div className={styles.tabContainer}>
          <h2 style={{ textAlign: 'center' }}>Random Question</h2>
          {randomError && <p className={styles.error}>{randomError}</p>}
          {randomSuccess && <p style={{ color: 'green', marginTop: '1rem' }}>{randomSuccess}</p>}
          {randomQ && (
            <div style={{ marginTop: '1rem' }}>
              <p><strong>Question:</strong> {randomQ.text}</p>
              <p><strong>Participants:</strong> {randomQ.total_participants}</p>
              <p><strong>First Choice:</strong> {renderRandomChoiceLabel(true)}</p>
              <p><strong>Second Choice:</strong> {renderRandomChoiceLabel(false)}</p>
              <div style={{ marginTop: '1rem', display: 'flex', gap: '1rem', justifyContent: 'center' }}>
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

      {/* SINGLE TAB */}
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
              <p><strong>First Choice:</strong> {renderSingleChoiceLabel(true)}</p>
              <p><strong>Second Choice:</strong> {renderSingleChoiceLabel(false)}</p>
              <div style={{ marginTop: '1rem', display: 'flex', gap: '1rem', justifyContent: 'center' }}>
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

      {/* MILESTONE POPUP */}
      {showPopup && popupQuestion && (
        <div style={{
          position: 'fixed',
          top: 0, left: 0, right: 0, bottom: 0,
          backgroundColor: 'rgba(0,0,0,0.6)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          zIndex: 9999,
        }}>
          <div style={{
            backgroundColor: '#fff',
            padding: '2rem',
            borderRadius: '8px',
            width: '400px',
            maxWidth: '80%',
            textAlign: 'center',
          }}>
            <h3>Milestone Question</h3>
            <p><strong>Question:</strong> {popupQuestion.text}</p>
            <p><strong>Participants:</strong> {popupQuestion.total_participants}</p>
            <p>
              <strong>First Choice:</strong>
              {popupQuestion.first_choice} ({popupQuestion.first_choice_count})
            </p>
            <p>
              <strong>Second Choice:</strong>
              {popupQuestion.second_choice} ({popupQuestion.second_choice_count})
            </p>

            <div style={{ marginTop: '1rem', display: 'flex', gap: '1rem', justifyContent: 'center' }}>
              <button style={buttonStyle} onClick={() => handleVoteMilestone(true)}>
                Vote {popupQuestion.first_choice}
              </button>
              <button style={buttonStyle} onClick={() => handleVoteMilestone(false)}>
                Vote {popupQuestion.second_choice}
              </button>
            </div>

            <button
              type="button"
              style={{ ...buttonStyle, marginTop: '1rem' }}
              onClick={handleCloseMilestone}
            >
              Close
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
