export const API_BASE = 'http://localhost:8080/api';

let accessToken: string | null =
  typeof window !== 'undefined'
    ? localStorage.getItem('accessToken')
    : null;

export function setAccessToken(token: string) {
  accessToken = token;
  if (typeof window !== 'undefined') {
    localStorage.setItem('accessToken', token);
  }
}

function getAccessToken() {
  return accessToken;
}

async function tryRefreshToken(): Promise<boolean> {
  try {
    const res = await fetch(`${API_BASE}/user/refresh`, {
      method: 'GET',
      credentials: 'include',
    });
    if (!res.ok) {
      console.warn('[tryRefreshToken] Refresh failed');
      return false;
    }
    const data = await res.json();
    console.log('[tryRefreshToken] New access token:', data.access_token);
    setAccessToken(data.access_token);
    return true;
  } catch (err) {
    console.error('[tryRefreshToken] Exception:', err);
    return false;
  }
}

export async function apiRequest(endpoint: string, options: RequestInit = {}): Promise<any> {
  const token = getAccessToken();

  const mergedOptions: RequestInit = {
    ...options,
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
      Authorization: token ? `Bearer ${token}` : '',
    },
  };

  let response = await fetch(endpoint, mergedOptions);
  if (response.status === 401) {
    console.warn('[apiRequest] Token possibly expired, attempting refresh...');
    const success = await tryRefreshToken();
    if (!success) {
      throw new Error('Unauthorized - Refresh token invalid or expired');
    }
    const newToken = getAccessToken();
    mergedOptions.headers = {
      ...mergedOptions.headers,
      Authorization: newToken ? `Bearer ${newToken}` : '',
    };
    response = await fetch(endpoint, mergedOptions);
  }

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.error || `Request failed with status ${response.status}`);
  }

  return response.json();
}
