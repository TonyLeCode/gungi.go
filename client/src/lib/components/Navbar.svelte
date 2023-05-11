<script lang="ts">
	import type { Session } from '@supabase/supabase-js';

	export let session: Session | null;
	$: session = session;
	// $: console.log("session!!", session)
</script>

<nav class="navbar">
	<a href="/">Gungi.go</a>
	<ul class="nav-list">
		{#if session}
			<a href="/home">home</a>
			<a href="/play">play</a>
			<!-- <a href="/games">games</a> -->
			<!-- TODO learning, puzzles, resources, library -->
		{/if}
		<a href="/rules">rules</a>
	</ul>
	<ul class="nav-account">
		{#if session}
			<!-- TODO settings, friends, notifications  -->
			<a href="/logout">logout</a>
		{:else}
			<a class="account" href="/login">login</a>
			<a href="/register">register</a>
		{/if}
	</ul>
</nav>

<style>
	.navbar {
		display: flex;
		/* justify-content: space-between; */
		/* background-color:red; */
		/* max-width: 1000px; */
		padding: 0.5rem 12rem;
	}
	.nav-list {
		display: flex;
		gap: 1rem;
		margin-left: 2rem;
	}
	.nav-account{
		display: flex;
		gap: 1rem;
		margin-left: auto;
	}
</style>
