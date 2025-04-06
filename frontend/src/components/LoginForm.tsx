'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { InputField } from './InputField';
import FormLayout from './FormLayout';
import { buttonStyle, formStyle, linkStyle } from './AuthenticationStyle';
import { loginUser } from '../utils/api';

export default function LoginForm() {
  const router = useRouter();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [accessToken, setAccessToken] = useState(''); // Holds the access token in memory

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!username || !password) {
      setError('Username and password are required');
      return;
    }

    try {
      console.log('[LoginForm] Sending login request for:', username);
      
      // Call the login API utility
      const data = await loginUser(username, password);
      
      console.log('[LoginForm] Login successful:', data);
      
      // Store the access token in component state (optional in memory)
      setAccessToken(data.access_token);

      // Also store the token in local storage
      localStorage.setItem('accessToken', data.access_token);

      // Navigate to the home page (or wherever you like)
      router.push('/');
    } catch (err: any) {
      console.error('[LoginForm] Login error:', err.message);
      setError(err.message || 'Login failed');
    }
  };

  return (
    <FormLayout title="Login">
      <form onSubmit={handleSubmit} style={formStyle}>
        <InputField
          id="username"
          label="Username"
          type="text"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <InputField
          id="password"
          label="Password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />

        {error && <div style={{ color: 'red', textAlign: 'center' }}>{error}</div>}

        <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
          <button type="submit" style={buttonStyle}>
            Login
          </button>
          <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
            <div style={linkStyle} onClick={() => router.push('/register')}>
              Register new account
            </div>
          </div>
        </div>
      </form>
    </FormLayout>
  );
}
