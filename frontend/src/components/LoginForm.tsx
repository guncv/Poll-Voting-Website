'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { InputField } from './InputField';
import FormLayout from './FormLayout';
import { buttonStyle, formStyle, linkStyle } from './AuthenticationStyle';

// We only need loginUser from real calls
import { loginUser } from '../utils/api';

export default function LoginForm() {
  const router = useRouter();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    console.log('NEXT_PUBLIC_API_PATH:', process.env.NEXT_PUBLIC_API_PATH);
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!username || !password) {
      setError('Username and password are required');
      return;
    }

    try {
      console.log('[LoginForm] Sending login request for:', username);

      const data = await loginUser(username, password);
      console.log('[LoginForm] Login successful:', data);

      // data.access_token is automatically stored in local storage
      // by loginUser calling setAccessToken. So you can just route to /
      router.push('/');
    } catch (err: any) {
      console.warn('[LoginForm] Login error:', err.message);
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
