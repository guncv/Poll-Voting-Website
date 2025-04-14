'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { createQuestion } from '../../utils/api';
import {
  formStyle,
  inputStyle,
  buttonStyle,
  titleStyle,
} from '../../components/AuthenticationStyle';
import Header from '../../components/Header';

export default function CreateQuestionPage() {
  const router = useRouter();

  // State for form fields
  const [text, setText] = useState('');
  const [firstChoice, setFirstChoice] = useState('');
  const [secondChoice, setSecondChoice] = useState('');
  const [milestones, setMilestones] = useState('');
  const [followUps, setFollowUps] = useState('');
  const [groupID, setGroupID] = useState('');
  const [error, setError] = useState('');
  const [successMessage, setSuccessMessage] = useState('');

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError('');
    setSuccessMessage('');

    if (!text || !firstChoice || !secondChoice) {
      setError('Please fill in all required fields (Question Text, First Choice, Second Choice).');
      return;
    }

    const payload = {
      text,
      first_choice: firstChoice,
      second_choice: secondChoice,
      milestones,
      follow_ups: followUps,
      group_id: groupID,
    };

    try {
      const response = await createQuestion(payload);
      setSuccessMessage(response.message || 'Question created successfully');

      // Reset fields 
      setText('');
      setFirstChoice('');
      setSecondChoice('');
      setMilestones('');
      setFollowUps('');
      setGroupID('');
    } catch (err: any) {
      setError(err.message || 'Failed to create question.');
    }
  }

  return (
    <div
      style={{
        padding: '2rem',
        paddingTop: 'calc(120px + 2rem)',  
        maxWidth: '800px',
        margin: '0 auto',
        textAlign: 'center',
      }}
    >
      <Header />
      <h1 style={titleStyle}>Create Question</h1>
      <form onSubmit={handleSubmit} style={formStyle}>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '1rem' }}>
          <label style={{ width: '150px', textAlign: 'left', marginRight: '1rem' }}>Question Text:</label>
          <input
            type="text"
            value={text}
            onChange={(e) => setText(e.target.value)}
            style={inputStyle}
            placeholder="Enter your question"
            required
          />
        </div>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '1rem' }}>
          <label style={{ width: '150px', textAlign: 'left', marginRight: '1rem' }}>First Choice:</label>
          <input
            type="text"
            value={firstChoice}
            onChange={(e) => setFirstChoice(e.target.value)}
            style={inputStyle}
            placeholder="Enter first choice"
            required
          />
        </div>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '1rem' }}>
          <label style={{ width: '150px', textAlign: 'left', marginRight: '1rem' }}>Second Choice:</label>
          <input
            type="text"
            value={secondChoice}
            onChange={(e) => setSecondChoice(e.target.value)}
            style={inputStyle}
            placeholder="Enter second choice"
            required
          />
        </div>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '1rem' }}>
          <label style={{ width: '150px', textAlign: 'left', marginRight: '1rem' }}>Milestones:</label>
          <input
            type="text"
            value={milestones}
            onChange={(e) => setMilestones(e.target.value)}
            style={inputStyle}
            placeholder="e.g., 100:m1,200:m2"
          />
        </div>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '1rem' }}>
          <label style={{ width: '150px', textAlign: 'left', marginRight: '1rem' }}>Follow-ups:</label>
          <input
            type="text"
            value={followUps}
            onChange={(e) => setFollowUps(e.target.value)}
            style={inputStyle}
            placeholder="Optional"
          />
        </div>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '1rem' }}>
          <label style={{ width: '150px', textAlign: 'left', marginRight: '1rem' }}>Group ID:</label>
          <input
            type="text"
            value={groupID}
            onChange={(e) => setGroupID(e.target.value)}
            style={inputStyle}
            placeholder="Optional"
          />
        </div>
        {error && (
          <p style={{ color: 'red', marginTop: '1rem', textAlign: 'left' }}>{error}</p>
        )}
        {successMessage && (
          <p style={{ color: 'green', marginTop: '1rem', textAlign: 'left' }}>
            {successMessage}
          </p>
        )}
        <button type="submit" style={buttonStyle}>
          Create Question
        </button>
      </form>
    </div>
  );
}
