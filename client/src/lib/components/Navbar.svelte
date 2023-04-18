<script lang='ts'>
  import type { Session } from "@supabase/supabase-js";

	export let session: Session | null;
	$: session = session
	// $: console.log("session!!", session)
</script>

<nav class="navbar">
	<a href="/">Gungi.go</a>
	<ul class="nav-list">
		<a href="/play">play</a>
		<a href="/games">games</a>
		{#if session}
			<a href="/logout">logout</a>
		{:else}
			<a href="/login">login</a>
			<a href="/register">register</a>
		{/if}
	</ul>
</nav>

<style>
	.navbar {
		display: flex;
		justify-content: space-between;
		/* background-color:red; */
		/* max-width: 1000px; */
		padding: 0.5rem 12rem;
	}
	.nav-list {
		display: flex;
		gap: 1rem;
	}
</style>
