<script>
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { currentUser } from '$lib/stores/auth';
  import * as api from '$lib/services/api';

  let profileUserID = $page.params.userID;
  let loggedInUserID = null;
  currentUser.subscribe(user => {
    loggedInUserID = user ? user.id : null;
  });

  let profileUser = null;
  let followers = [];
  let following = [];
  let followStatus = null; // 'accepted', 'pending_sent_by_me', 'none'
  let isLoading = true;
  let error = null;

  onMount(async () => {
    await loadProfileData();
  });

  async function loadProfileData() {
    isLoading = true;
    error = null;
    try {
      // Fetch profile user details
      // getUserProfile is a placeholder in api.js - needs actual backend endpoint
      profileUser = await api.getUserProfile(profileUserID);

      // Fetch followers and following lists (these should be only 'accepted' from backend)
      followers = await api.getFollowers(profileUserID);
      following = await api.getFollowing(profileUserID);

      // Determine follow status between loggedInUser and profileUser
      if (loggedInUserID && loggedInUserID !== profileUserID) {
        // Attempt to determine if already following
        const loggedInUserFollowingList = await api.getFollowing(loggedInUserID);
        const isActuallyFollowing = loggedInUserFollowingList.find(u => u.id === profileUserID);
        if (isActuallyFollowing) {
          followStatus = 'accepted';
        } else {
          // This is where it gets tricky without a dedicated endpoint for "pending_sent_by_me"
          // For now, if not 'accepted', we assume 'none'.
          // "Request Sent" state will be set optimistically after clicking "Request Follow"
          followStatus = 'none';
        }
      }

    } catch (err) {
      console.error("Error loading profile data:", err);
      error = err.message;
    } finally {
      isLoading = false;
    }
  }

  async function handleFollow() {
    if (!loggedInUserID) return alert('You must be logged in to follow users.');
    try {
      await api.followUser(profileUserID);
      // Optimistic update
      if (profileUser.is_private) {
        followStatus = 'pending_sent_by_me';
      } else {
        followStatus = 'accepted';
        // Refresh followers list for current profile if it's public and I just followed them
        // Or better, just add loggedInUser to the profileUser's followers list locally.
        // For simplicity, could just reload all data: await loadProfileData();
      }
      await loadProfileData(); // Reload to reflect changes
    } catch (err) {
      alert(`Error following user: ${err.message}`);
    }
  }

  async function handleUnfollow() {
    if (!loggedInUserID) return;
    try {
      await api.unfollowUser(profileUserID);
      followStatus = 'none';
      await loadProfileData(); // Reload to reflect changes
    } catch (err)
 {
      alert(`Error unfollowing user: ${err.message}`);
    }
  }

  // Reactive statements to derive button states
  $: isOwnProfile = loggedInUserID === profileUserID;
  $: canFollow = !isOwnProfile && followStatus === 'none' && !profileUser?.is_private;
  $: canRequestFollow = !isOwnProfile && followStatus === 'none' && profileUser?.is_private;
  $: canUnfollow = !isOwnProfile && followStatus === 'accepted';
  $: isRequestSent = !isOwnProfile && followStatus === 'pending_sent_by_me';

</script>

{#if isLoading}
  <p>Loading profile...</p>
{:else if error}
  <p style="color: red;">Error: {error}</p>
{:else if profileUser}
  <div>
    <h2>{profileUser.name || profileUser.id}'s Profile</h2>
    <p>ID: {profileUser.id}</p>
    <p>Is Private: {profileUser.is_private ? 'Yes' : 'No'}</p>

    {#if !isOwnProfile}
      <div>
        {#if canFollow}
          <button on:click={handleFollow}>Follow</button>
        {:else if canRequestFollow}
          <button on:click={handleFollow}>Request Follow</button>
        {:else if canUnfollow}
          <button on:click={handleUnfollow}>Unfollow</button>
        {:else if isRequestSent}
          <button disabled>Request Sent</button>
        {/if}
      </div>
    {/if}

    <hr />
    <h4>Followers ({followers.length})</h4>
    {#if followers.length > 0}
      <ul>
        {#each followers as follower}
          <li><a href="/profile/{follower.id}">{follower.first_name || follower.id}</a></li>
        {/each}
      </ul>
    {:else}
      <p>No followers yet.</p>
    {/if}

    <hr />
    <h4>Following ({following.length})</h4>
    {#if following.length > 0}
      <ul>
        {#each following as followedUser}
          <li><a href="/profile/{followedUser.id}">{followedUser.first_name || followedUser.id}</a></li>
        {/each}
      </ul>
    {:else}
      <p>Not following anyone yet.</p>
    {/if}
  </div>
{:else}
  <p>User not found.</p>
{/if}

<style>
  button {
    margin: 5px;
    padding: 8px 12px;
    cursor: pointer;
  }
  button[disabled] {
    cursor: not-allowed;
    opacity: 0.6;
  }
  ul {
    list-style-type: none;
    padding: 0;
  }
  li {
    margin-bottom: 5px;
  }
</style>
