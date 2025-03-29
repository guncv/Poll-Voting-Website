'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { InputField } from './InputField';
import FormLayout from './FormLayout';
import { buttonStyle, formStyle, linkStyle } from './AuthenticationStyle';

export default function LoginForm() {
  const router = useRouter();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!username || !password) {
      setError('Username and password are required');
      return;
    }
    router.push('/');
  };

  return (
    <FormLayout title="Login">
      <form onSubmit={handleSubmit} style={formStyle}>

        <InputField id="username" label="Username" type="text" value={username} onChange={e => setUsername(e.target.value)} />
        <InputField id="password" label="Password" type="password" value={password} onChange={e => setPassword(e.target.value)} />

        {error && <div style={{ color: 'red', textAlign: 'center' }}>{error}</div>}

        <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
          <button type="submit" style={buttonStyle}>Login</button>
          <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
            <div style={linkStyle} onClick={() => router.push('/register')}>Register new account</div>
          </div>
        </div>
        
      </form>
    </FormLayout>
  );
}
