import { writable } from 'svelte/store';

// Initial state: no user logged in
const initialUser = null;
/*
User object structure could be:
{
  id: "string",
  email: "string",
  // ... other user props
  token: "jwt_token_string"
}
*/

export const currentUser = writable(initialUser);

// Function to simulate login
export function login(userData, token) {
  currentUser.set({ ...userData, token });
  // In a real app, token might be stored in localStorage
  localStorage.setItem('authToken', token);
  localStorage.setItem('userID', userData.id);
}

// Function to simulate logout
export function logout() {
  currentUser.set(null);
  localStorage.removeItem('authToken');
  localStorage.removeItem('userID');
}

// Function to check auth state on app load
export function initializeAuth() {
  const token = localStorage.getItem('authToken');
  const userID = localStorage.getItem('userID');
  if (token && userID) {
    // In a real app, you might want to verify the token with the backend here
    // For now, just restore basic user info.
    // You'd typically fetch user details from the backend using the token.
    currentUser.set({ id: userID, token: token });
  }
}
