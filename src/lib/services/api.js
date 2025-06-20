import { currentUser } from '$lib/stores/auth';

let bearerToken = null;
currentUser.subscribe(user => {
  bearerToken = user ? user.token : null;
});

const BASE_URL = '/api/v1'; // Assuming your API is prefixed

async function request(method, path, data) {
  const opts = { method, headers: {} };

  if (data) {
    opts.headers['Content-Type'] = 'application/json';
    opts.body = JSON.stringify(data);
  }

  if (bearerToken) {
    opts.headers['Authorization'] = `Bearer ${bearerToken}`;
  }

  try {
    const response = await fetch(`${BASE_URL}${path}`, opts);
    if (response.ok) {
      if (response.status === 204) return null; // No content
      return await response.json();
    } else {
      const errorBody = await response.text(); // Try to get error body
      console.error(`API Error: ${response.status} ${response.statusText}`, errorBody);
      // You might want to throw a custom error object here
      throw new Error(`API request failed: ${response.status} ${response.statusText} - ${errorBody}`);
    }
  } catch (err) {
    console.error(`Network or other error: ${path}`, err);
    throw err; // Re-throw the error to be caught by the caller
  }
}

// --- User/Profile ---
// This one is an assumption, not explicitly in backend steps but needed for profile pages
export const getUserProfile = (userID) => {
  // Assuming a backend route like GET /api/v1/users/{userID}/profile or similar
  // For now, let's use a placeholder or assume one of the existing routes can provide this
  // If backend has GET /api/v1/users/{userID} that returns full profile:
  // return request('GET', `/users/${userID}`);
  // This is a guess, the backend routes list didn't explicitly have a "get user profile" endpoint.
  // Let's assume for now that GET /users/{userID}/followers also returns basic user info or we adapt
  // For now, this is a placeholder function
  console.warn("getUserProfile: Using placeholder data. Implement with actual backend endpoint.");
  return Promise.resolve({ id: userID, name: `User ${userID}`, is_private: Math.random() > 0.5, followers_count: 0, following_count: 0 });
};


// --- Follow ---
export const followUser = (userID) => {
  return request('POST', `/users/${userID}/follow`);
};

export const unfollowUser = (userID) => {
  return request('DELETE', `/users/${userID}/unfollow`);
};

export const getFollowers = (userID) => {
  return request('GET', `/users/${userID}/followers`);
};

export const getFollowing = (userID) => {
  return request('GET', `/users/${userID}/following`);
};

// --- Follow Requests ---
export const getPendingFollowRequests = () => {
  return request('GET', '/follow_requests/pending');
};

export const acceptFollowRequest = (requesterID) => {
  // Path from backend: /api/v1/users/{requesterID}/accept_follow
  return request('POST', `/users/${requesterID}/accept_follow`);
};

export const declineFollowRequest = (requesterID) => {
  // Path from backend: /api/v1/users/{requesterID}/decline_follow
  return request('POST', `/users/${requesterID}/decline_follow`); // or DELETE if backend supports it
};

// Helper to check initial follow status - this is more complex
// It ideally needs a dedicated backend endpoint: GET /users/{profileUserID}/relationship
// For now, we can try to infer it by fetching lists, but it's inefficient.
// This function will be called from the profile page.
export const getFollowRelationship = async (profileUserID, loggedInUserID) => {
  if (!loggedInUserID || !profileUserID || loggedInUserID === profileUserID) {
    return null; // No relationship to check or it's own profile
  }
  try {
    // 1. Check if loggedInUser is following profileUser (status 'accepted')
    const followingList = await getFollowing(loggedInUserID); // List of users loggedInUser follows
    const isFollowing = followingList.find(user => user.id === profileUserID);
    if (isFollowing) {
      return { status: 'accepted' }; // Assumes GetFollowing only returns accepted ones
    }

    // 2. If not 'accepted', check if a 'pending' request exists from loggedInUser to profileUser
    // This is tricky. The backend's Follow() method creates 'pending' or 'accepted'.
    // The GET /follow_requests/pending is for requests *received by* loggedInUser.
    // We don't have an endpoint to see requests *sent by* loggedInUser that are pending.
    // This is a limitation of the current backend API design for this frontend feature.
    // For now, if not 'accepted', we can't easily determine 'pending_sent_by_me'.
    // The "Request Sent" button state will be hard to set accurately without this.
    // Let's assume the Follow() API call on a private user will result in a "pending" state
    // and the component will optimistically update to "Request Sent".
    // A more robust solution would be an API endpoint: GET /api/v1/users/{targetUserID}/relationship_status
    // that returns { status: 'none' | 'pending_them_to_me' | 'pending_me_to_them' | 'accepted' }

    // For this subtask, we will rely on the follow/unfollow buttons and optimistic updates.
    // The profile page will show "Follow" or "Request Follow". If clicked, it becomes "Unfollow" or "Request Sent".
    // "Request Sent" state will be mostly optimistic or based on a local component state after clicking "Request Follow".
    return { status: 'none' }; // Default to 'none' if not directly following.
  } catch (error) {
    console.error("Error determining follow relationship:", error);
    return { status: 'none' }; // Or handle error appropriately
  }
};
