'use client';

import { useEffect, useState } from 'react';
import { fetchPoll, sendVote } from '../utils/api';
import { PollData } from '../types/pollData';

const Voting = () => {
  const [poll, setPoll] = useState<PollData | null>(null);
  const [isVoting, setIsVoting] = useState(false);
  const [voted, setVoted] = useState(false);
  const [selectedChoice, setSelectedChoice] = useState<string | null>(null);
  console.log("selectedChoice", selectedChoice);

  useEffect(() => {
    loadPoll();
  }, []);

  const loadPoll = async () => {
    const data = await fetchPoll();
    setPoll(data);
    setVoted(false);
    setSelectedChoice(null);
  };

  const handleVote = async (choice: string) => {
    if (!poll || isVoting) return;
    setIsVoting(true);
    await sendVote(poll.id, choice);
    setVoted(true);
    setIsVoting(false);
    setTimeout(() => {
      loadPoll();
    }, 2000);
  };

  if (!poll) return <div>Loading...</div>;

  return (
    <main style={{ textAlign: 'center', marginTop: '3rem' }}>
      <h1>{poll.question}</h1>
      {voted && <p style={{ fontSize: '30px', marginTop: '2rem'}}>Thanks for voting! 🎉</p>}
      {poll.choices.map((choice, index) => (
        <button
          key={choice}
          onClick={() => {
            setSelectedChoice(choice);
            handleVote(choice);
          }}
          disabled={voted}
          style={{margin: '3rem',
            backgroundColor: selectedChoice === choice ? '#FF98A0' :  selectedChoice === null ? index === 0 ? 'green' : 'red' : 'white', 
            borderRadius: '10px' , 
            border: 'none', 
            color: 'black', 
            padding: '1rem 2rem', 
            fontSize: '30px', 
            fontFamily: 'Poppins, sans-serif' ,
            cursor: 'pointer'
        }}
        >
          {choice}
        </button>
      ))}
      
    </main>
  );
};

export default Voting;
