'use client';

import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { getProfile, logoutUser } from '../utils/api';

export default function HomePage() {
  const router = useRouter();
  const [profile, setProfile] = useState<any>(null);
  const [error, setError] = useState('');

  // Click handler for "Check Profile"
  async function handleCheckProfile() {
    setError('');
    setProfile(null);

    const token = localStorage.getItem('accessToken');
    if (!token) {
      setError('No access token found. Please login first.');
      return;
    }

    try {
      const userData = await getProfile(token);
      setProfile(userData);
      console.log('[HomePage] Profile data:', userData);
    } catch (err: any) {
      console.error('[HomePage] Error fetching profile:', err.message);
      setError(err.message || 'Failed to fetch profile');
    }
  }

  // Click handler for "Logout"
  async function handleLogout() {
    setError('');
    try {
      // 1) Call the server to clear the refresh token cookie
      await logoutUser();
      
      // 2) Remove the access token from local storage or in-memory
      localStorage.removeItem('accessToken'); 
      
      // 3) (Optional) Redirect or update UI to reflect logged-out state
      router.push('/login');
      
      console.log('Logout successful');
    } catch (err: any) {
      setError(err.message || 'Logout failed');
    }
  }

  return (
    <div>
      <h1>Home</h1>
      {/* Existing Login/Register buttons */}
      <button onClick={() => router.push('/login')}>Login</button>
      <button onClick={() => router.push('/register')}>Register</button>

      {/* New button: Check Profile */}
      <button onClick={handleCheckProfile}>Check Profile</button>

      {/* New button: Logout */}
      <button onClick={handleLogout}>Logout</button>

      {/* Display errors if any */}
      {error && <p style={{ color: 'red' }}>{error}</p>}

      {/* Show user profile data if available */}
      {profile && (
        <div style={{ marginTop: '1rem' }}>
          <h2>User Profile</h2>
          <pre>{JSON.stringify(profile, null, 2)}</pre>
        </div>
      )}
    </div>
  );
}
