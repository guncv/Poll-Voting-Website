'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { InputField } from './InputField';
import FormLayout from './FormLayout';
import { buttonStyle, formStyle, linkStyle } from './AuthenticationStyle';

export default function RegisterForm() {
  const router = useRouter();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!username || !password || !confirmPassword) {
      setError('All fields are required');
      return;
    } if (password !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }
    router.push('/login');
  };

  return (
    <FormLayout title="Register">
      <form onSubmit={handleSubmit} style={formStyle}>

        <InputField id="username" label="Username" type="text" value={username} onChange={e => setUsername(e.target.value)} />
        <InputField id="password" label="Password" type="password" value={password} onChange={e => setPassword(e.target.value)} />
        <InputField id="confirmPassword" label="Confirm Password" type="password" value={confirmPassword} onChange={e => setConfirmPassword(e.target.value)} />

        {error && <div style={{ color: 'red', textAlign: 'center' }}>{error}</div>}

        <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
          <button type="submit" style={buttonStyle}>Register</button>
          <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
            <div style={linkStyle} onClick={() => router.push('/login')}>Already registered? Login here</div>
          </div>
        </div>

      </form>
    </FormLayout>
  );
}
