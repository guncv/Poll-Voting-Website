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

interface MilestoneForm {
  popUpScore: string;  // e.g. "100"
  text: string;
  firstChoice: string;
  secondChoice: string;
}

export default function CreateQuestionPage() {
  const router = useRouter();

  // ------------- Main Question Fields -------------
  const [text, setText] = useState('');
  const [firstChoice, setFirstChoice] = useState('');
  const [secondChoice, setSecondChoice] = useState('');

  // ------------- Milestones Toggle -------------
  const [hasMilestones, setHasMilestones] = useState(false);
  const [milestones, setMilestones] = useState<MilestoneForm[]>([]);

  // ------------- Messages -------------
  const [error, setError] = useState('');
  const [successMessage, setSuccessMessage] = useState('');

  // Add a new blank milestone form
  function addMilestone() {
    setMilestones((prev) => [
      ...prev,
      { popUpScore: '', text: '', firstChoice: '', secondChoice: '' },
    ]);
  }

  // Remove a milestone
  function removeMilestone(index: number) {
    setMilestones((prev) => prev.filter((_, i) => i !== index));
  }

  // Update a milestone field
  function updateMilestoneField(index: number, field: keyof MilestoneForm, value: string) {
    setMilestones((prev) =>
      prev.map((m, i) => (i === index ? { ...m, [field]: value } : m))
    );
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError('');
    setSuccessMessage('');

    // Validate the main question
    if (!text || !firstChoice || !secondChoice) {
      setError('Please fill in the main question text, first choice, and second choice.');
      return;
    }

    // If user checked "hasMilestones" but gave no milestones
    if (hasMilestones && milestones.length === 0) {
      setError('You indicated you want milestone questions, but none are added.');
      return;
    }

    try {
      let finalMilestonesStr = ''; 
      let mainFollowUp = '';    
      let groupId = '';         // We'll only create a group_id if we have milestones

      // ---------------- CREATE MILESTONES (IF ANY) ----------------
      if (hasMilestones && milestones.length > 0) {
        groupId = crypto.randomUUID(); // generate a group_id since we do have milestones

        // Sort ascending by popUpScore
        const sorted = [...milestones].sort((a, b) => {
          const sA = parseInt(a.popUpScore, 10) || 0;
          const sB = parseInt(b.popUpScore, 10) || 0;
          return sA - sB;
        });

        const milestoneIDsAsc: string[] = [];
        let nextId = '';

        // Create them in descending order
        for (let i = sorted.length - 1; i >= 0; i--) {
          const m = sorted[i];
          if (!m.popUpScore || !m.text || !m.firstChoice || !m.secondChoice) {
            throw new Error('One of your milestone questions is missing required fields.');
          }

          const milestonePayload: any = {
            text: m.text,
            first_choice: m.firstChoice,
            second_choice: m.secondChoice,
            milestones: '',
            follow_ups: nextId,
          };

          // If we have a non-empty groupId, attach it
          milestonePayload.group_id = groupId;

          const milestoneRes = await createQuestion(milestonePayload);
          const milestoneId = milestoneRes.question_id;
          if (!milestoneId) {
            throw new Error('Milestone creation did not return question_id.');
          }

          nextId = milestoneId;
          milestoneIDsAsc.unshift(milestoneId);
        }

        // The earliest milestone is milestoneIDsAsc[0]
        mainFollowUp = milestoneIDsAsc[0];

        // Build "score:questionId" parts
        const parts: string[] = [];
        for (let i = 0; i < sorted.length; i++) {
          const score = sorted[i].popUpScore.trim();
          const mid = milestoneIDsAsc[i];
          parts.push(`${score}:${mid}`);
        }
        finalMilestonesStr = parts.join(',');
      }

      // ---------------- CREATE MAIN QUESTION ----------------
      // We'll build the main question payload dynamically
      const mainPayload: any = {
        text,
        first_choice: firstChoice,
        second_choice: secondChoice,
        milestones: finalMilestonesStr,
        follow_ups: mainFollowUp,
      };

      // Only attach group_id if we actually have one
      if (groupId) {
        mainPayload.group_id = groupId;
      }

      const mainRes = await createQuestion(mainPayload);
      if (!mainRes.question_id) {
        throw new Error('Main question creation did not return question_id.');
      }

      // ---------------- DONE ----------------
      if (hasMilestones && milestones.length > 0) {
        setSuccessMessage(
          `Main question + ${milestones.length} milestone(s) created successfully!`
        );
      } else {
        setSuccessMessage('Main question (no milestones) created successfully!');
      }

      resetForm();
    } catch (err: any) {
      console.error(err);
      setError(err.message || 'Failed to create question(s).');
    }
  }

  function resetForm() {
    setText('');
    setFirstChoice('');
    setSecondChoice('');
    setHasMilestones(false);
    setMilestones([]);
  }

  return (
    <div
      style={{
        padding: '2rem',
        paddingTop: 'calc(120px + 2rem)',
        maxWidth: '900px',
        margin: '0 auto',
        textAlign: 'center',
      }}
    >
      <Header />
      <h1 style={titleStyle}>Create Question + Chained Milestones</h1>

      {error && <p style={{ color: 'red', marginBottom: '1rem' }}>{error}</p>}
      {successMessage && <p style={{ color: 'green', marginBottom: '1rem' }}>{successMessage}</p>}

      <form onSubmit={handleSubmit} style={formStyle}>
        {/* MAIN QUESTION FIELDS */}
        <h2 style={{ textAlign: 'left' }}>Main Question</h2>
        <div style={{ display: 'flex', marginBottom: '1rem' }}>
          <label style={{ width: '150px', textAlign: 'left', marginRight: '1rem' }}>
            Question Text:
          </label>
          <input
            type="text"
            value={text}
            onChange={(e) => setText(e.target.value)}
            style={inputStyle}
            required
          />
        </div>

        <div style={{ display: 'flex', marginBottom: '1rem' }}>
          <label style={{ width: '150px', textAlign: 'left', marginRight: '1rem' }}>
            First Choice:
          </label>
          <input
            type="text"
            value={firstChoice}
            onChange={(e) => setFirstChoice(e.target.value)}
            style={inputStyle}
            required
          />
        </div>

        <div style={{ display: 'flex', marginBottom: '1rem' }}>
          <label style={{ width: '150px', textAlign: 'left', marginRight: '1rem' }}>
            Second Choice:
          </label>
          <input
            type="text"
            value={secondChoice}
            onChange={(e) => setSecondChoice(e.target.value)}
            style={inputStyle}
            required
          />
        </div>

        {/* TOGGLE FOR MILESTONES */}
        <div style={{ textAlign: 'left', margin: '1rem 0' }}>
          <label>
            <input
              type="checkbox"
              checked={hasMilestones}
              onChange={(e) => setHasMilestones(e.target.checked)}
            />
            &nbsp;Include milestone questions?
          </label>
        </div>

        {/* MILESTONE DETAILS IF HAS MILESTONES */}
        {hasMilestones && (
          <div
            style={{
              border: '1px solid #ccc',
              padding: '1rem',
              marginBottom: '1rem',
              textAlign: 'left',
            }}
          >
            <h3>Milestone Questions</h3>
            <p style={{ fontStyle: 'italic', marginBottom: '1rem' }}>
              Each milestone needs a numeric pop-up score, plus text & choices.
              The code automatically chains each milestone to the next.
            </p>

            {milestones.map((m, idx) => (
              <div
                key={idx}
                style={{ marginBottom: '1rem', padding: '0.5rem', border: '1px solid #ddd' }}
              >
                <div style={{ display: 'flex', marginBottom: '0.5rem' }}>
                  <label style={{ width: '110px', marginRight: '0.5rem' }}>
                    Pop-Up Score:
                  </label>
                  <input
                    type="text"
                    value={m.popUpScore}
                    onChange={(e) => updateMilestoneField(idx, 'popUpScore', e.target.value)}
                    style={inputStyle}
                    placeholder="e.g. 100, 150"
                  />
                </div>

                <div style={{ display: 'flex', marginBottom: '0.5rem' }}>
                  <label style={{ width: '110px', marginRight: '0.5rem' }}>
                    Text:
                  </label>
                  <input
                    type="text"
                    value={m.text}
                    onChange={(e) => updateMilestoneField(idx, 'text', e.target.value)}
                    style={inputStyle}
                  />
                </div>

                <div style={{ display: 'flex', marginBottom: '0.5rem' }}>
                  <label style={{ width: '110px', marginRight: '0.5rem' }}>
                    1st Choice:
                  </label>
                  <input
                    type="text"
                    value={m.firstChoice}
                    onChange={(e) => updateMilestoneField(idx, 'firstChoice', e.target.value)}
                    style={inputStyle}
                  />
                </div>

                <div style={{ display: 'flex', marginBottom: '0.5rem' }}>
                  <label style={{ width: '110px', marginRight: '0.5rem' }}>
                    2nd Choice:
                  </label>
                  <input
                    type="text"
                    value={m.secondChoice}
                    onChange={(e) => updateMilestoneField(idx, 'secondChoice', e.target.value)}
                    style={inputStyle}
                  />
                </div>

                <button
                  type="button"
                  style={{ ...buttonStyle, backgroundColor: '#f44336' }}
                  onClick={() => removeMilestone(idx)}
                >
                  Remove
                </button>
              </div>
            ))}

            <button type="button" style={buttonStyle} onClick={addMilestone}>
              + Add Another Milestone
            </button>
          </div>
        )}

        {/* SUBMIT */}
        <button type="submit" style={buttonStyle}>
          Create Question
        </button>
      </form>
    </div>
  );
}
