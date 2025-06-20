<script>
  import { onMount } from 'svelte';
  import { currentUser } from '$lib/stores/auth';
  import * as api from '$lib/services/api';

  let loggedInUserID = null;
  currentUser.subscribe(user => {
    loggedInUserID = user ? user.id : null;
  });

  let pendingRequests = [];
  let isLoading = true;
  let error = null;

  onMount(async () => {
    if (!loggedInUserID) {
      // Optional: redirect to login or show message if not logged in
      // For now, just won't load data. User should be redirected by layout or similar.
      isLoading = false;
      error = "You need to be logged in to see follow requests.";
      return;
    }
    await loadPendingRequests();
  });

  async function loadPendingRequests() {
    isLoading = true;
    error = null;
    try {
      pendingRequests = await api.getPendingFollowRequests();
    } catch (err) {
      console.error("Error loading pending follow requests:", err);
      error = err.message;
    } finally {
      isLoading = false;
    }
  }

  async function handleAcceptRequest(requesterID) {
    try {
      await api.acceptFollowRequest(requesterID);
      // Remove from list optimistically
      pendingRequests = pendingRequests.filter(req => req.id !== requesterID);
    } catch (err) {
      alert(`Error accepting request: ${err.message}`);
      // Optionally reload list: await loadPendingRequests();
    }
  }

  async function handleDeclineRequest(requesterID) {
    try {
      await api.declineFollowRequest(requesterID);
      // Remove from list optimistically
      pendingRequests = pendingRequests.filter(req => req.id !== requesterID);
    } catch (err) {
      alert(`Error declining request: ${err.message}`);
      // Optionally reload list: await loadPendingRequests();
    }
  }
</script>

<div>
  <h2>Pending Follow Requests</h2>

  {#if isLoading}
    <p>Loading requests...</p>
  {:else if error}
    <p style="color: red;">Error: {error}</p>
  {:else if pendingRequests.length === 0}
    <p>No pending follow requests.</p>
  {:else}
    <ul>
      {#each pendingRequests as request (request.id)}
        <li>
          <span>{request.first_name || request.nickname || request.id}</span>
          <div>
            <button on:click={() => handleAcceptRequest(request.id)}>Accept</button>
            <button on:click={() => handleDeclineRequest(request.id)}>Decline</button>
          </div>
        </li>
      {/each}
    </ul>
  {/if}
</div>

<style>
  ul {
    list-style-type: none;
    padding: 0;
  }
  li {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px;
    border: 1px solid #eee;
    margin-bottom: 5px;
    background-color: #fff;
  }
  button {
    margin-left: 10px;
    padding: 5px 10px;
    cursor: pointer;
  }
</style>
