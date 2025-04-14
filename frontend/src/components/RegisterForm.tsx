'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { InputField } from './InputField';
import FormLayout from './FormLayout';
import { buttonStyle, formStyle, linkStyle } from './AuthenticationStyle';

// Import our new registerUser function
import { registerUser } from '../utils/api';

export default function RegisterForm() {
  const router = useRouter();

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Basic front-end validation
    if (!username || !password || !confirmPassword) {
      setError('All fields are required');
      return;
    }
    if (password !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    // Clear old errors
    setError('');

    try {
      // Call our registerUser function
      const result = await registerUser(username, password);
      // result could look like: { email: "nasssuh@example.com", user_id: "..." }
      console.log('[RegisterForm] Registration successful:', result);

      // Redirect to login or wherever you want
      router.push('/login');
    } catch (err: any) {
      // Avoid console.error to stop Next.js dev overlay
      console.warn('[RegisterForm] Registration error:', err.message);
      setError(err.message || 'Registration failed');
    }
  };

  return (
    <FormLayout title="Register">
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
        <InputField
          id="confirmPassword"
          label="Confirm Password"
          type="password"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
        />

        {error && <div style={{ color: 'red', textAlign: 'center' }}>{error}</div>}

        <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
          <button type="submit" style={buttonStyle}>
            Register
          </button>
          <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
            <div style={linkStyle} onClick={() => router.push('/login')}>
              Already registered? Login here
            </div>
          </div>
        </div>
      </form>
    </FormLayout>
  );
}
