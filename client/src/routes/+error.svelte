<script>
	import { page } from '$app/stores';
  console.log($page.error)
</script>

<main>
	<h1>{$page.status == 404 ? 'Page Not Found' : 'Server Error'}</h1>
	{#if $page.status == 404}
		<p>Looks like you're lost. Go back <a href='/'>home</a></p>
	{:else}
		<p>Looks like something went wrong in the server. Try again later.</p>
	{/if}
</main>

<style>
	main {
		max-width: 70rem;
		margin: auto;
		text-align: center;
		font-size: 1.5rem;
		min-height: calc(100vh - 3rem);
		display: flex;
		justify-content: center;
		flex-direction: column;
	}
	h1 {
		font-size: 2rem;
		font-weight: 600;
	}
	p {
		margin-bottom: 12rem;
	}
  a{
    color: rgb(255, 94, 0);
		padding: 0;
  }
</style>
