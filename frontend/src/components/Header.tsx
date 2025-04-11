'use client';

import React from 'react';
import { useRouter } from 'next/navigation';
import { buttonStyle } from './AuthenticationStyle';

const headerContainerStyle: React.CSSProperties = {
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

export default function Header() {
  const router = useRouter();
  return (
    <header style={headerContainerStyle}>
      <button style={{ ...buttonStyle, margin: 0 }} onClick={() => router.back()}>
        Back
      </button>
      <button style={{ ...buttonStyle, margin: 0 }} onClick={() => router.push('/')}>
        Home
      </button>
    </header>
  );
}
