<script>
  import { currentUser } from '$lib/stores/auth';
</script>

<div>
  <h1>Welcome to Social Network</h1>
  {#if $currentUser}
    <p>You are logged in as {$currentUser.id}.</p>
    <p><a href="/profile/{$currentUser.id}">View Your Profile</a></p>
    <p><a href="/follow-requests">Check Follow Requests</a></p>
  {:else}
    <p>Please <a href="/login">login</a> or <a href="/register">register</a>.</p>
    <p>To test functionality without full login/register UI, you can use the auth.js functions `login({id: 'testUser', email:'test@example.com'}, 'fakeToken')` from the browser console after navigating to a page, then refresh.</p>
  {/if}
</div>

<style>
  div {
    text-align: center;
    padding: 2em;
  }
</style>
