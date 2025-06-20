<script>
  import { onMount } from 'svelte';
  import { currentUser, initializeAuth, logout } from '$lib/stores/auth';
  import { goto } from '$app/navigation';

  onMount(() => {
    initializeAuth();
  });

  function handleLogout() {
    logout();
    goto('/'); // Redirect to home page after logout
  }

  // Reactive declaration to get current user ID
  let currentUserID = null;
  currentUser.subscribe(user => {
    if (user) {
      currentUserID = user.id;
    } else {
      currentUserID = null;
    }
  });
</script>

<header>
  <nav>
    <a href="/">Home</a>
    {#if $currentUser}
      <a href="/profile/{currentUserID}">My Profile</a>
      <a href="/follow-requests">Follow Requests</a>
      <button on:click={handleLogout}>Logout</button>
      <span>(Logged in as {$currentUser.id})</span>
    {:else}
      <a href="/login">Login</a> <!-- Assuming a /login route will be created -->
      <a href="/register">Register</a> <!-- Assuming a /register route will be created -->
    {/if}
  </nav>
</header>

<main>
  <slot></slot>
</main>

<style>
  :global(body) {
    font-family: Arial, sans-serif;
    margin: 0;
    padding: 0;
    background-color: #f4f4f4;
    color: #333;
  }
  header {
    background-color: #333;
    color: white;
    padding: 1em;
    text-align: center;
  }
  nav a, nav button {
    color: white;
    margin: 0 10px;
    text-decoration: none;
    background: none;
    border: none;
    cursor: pointer;
    font-size: 1em;
  }
  nav a:hover, nav button:hover {
    text-decoration: underline;
  }
  main {
    padding: 1em;
  }
</style>
