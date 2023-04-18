<script lang='ts'>
  import BoardSimple from '$lib/components/BoardSimple.svelte';

  export let data;
  $: console.log(data.data)
</script>

<svelte:head>
	<title>Gungi.go</title>
</svelte:head>


<h1>Gungi.go</h1>
<div class='gameList'>
  {#each data.data as game}
    <BoardSimple gameData={game} />
  {/each}
  <!-- <BoardSimple />
  <BoardSimple />
  <BoardSimple />
  <BoardSimple /> -->
</div>

<style>
  .gameList{
    display:grid;
    margin: 0 auto;
    gap: 1rem;
    grid-template-columns: repeat(auto-fit, minmax(15rem, 1fr));
    max-width: 70rem;
    padding: 1rem;
  }
</style>
