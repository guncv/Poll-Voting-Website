// src/utils/apiClient.ts

// We'll keep the access token in a module-level variable
let accessToken: string | null = typeof window !== 'undefined'
  ? localStorage.getItem('accessToken')
  : null;

/**
 * setAccessToken:
 *  - updates our module-level variable
 *  - persists to localStorage (so page reload doesn't lose the token)
 */
export function setAccessToken(token: string) {
  accessToken = token;
  if (typeof window !== 'undefined') {
    localStorage.setItem('accessToken', token);
  }
}

/**
 * getAccessToken:
 *  - returns our in-memory token
 */
function getAccessToken(): string | null {
  return accessToken;
}

/**
 * A function to call your /refresh endpoint if access token is expired.
 * The refresh token is in an HttpOnly cookie, so just call /refresh
 * with credentials:'include'. If it returns a new access token, store it.
 */
async function tryRefreshToken(): Promise<boolean> {
  try {
    const res = await fetch('http://localhost:8080/api/user/refresh', {
      method: 'GET', // or POST if your refresh endpoint is POST
      credentials: 'include', // ensures the refresh token cookie is sent
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
    console.error('[tryRefreshToken] Exception during refresh:', err);
    return false;
  }
}

/**
 * apiRequest: a custom fetch wrapper that auto-includes credentials,
 * attaches the current access token, and on 401 attempts to refresh once.
 */
export async function apiRequest(
  endpoint: string,
  options: RequestInit = {}
): Promise<any> {
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

  // First request attempt
  let response = await fetch(endpoint, mergedOptions);

  // If 401, try refreshing
  if (response.status === 401) {
    console.warn('[apiRequest] Token possibly expired, attempting refresh...');
    const success = await tryRefreshToken();
    if (!success) {
      throw new Error('Unauthorized - Refresh token invalid or expired');
    }
    // Refresh succeeded, retry the original request with the new token
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
