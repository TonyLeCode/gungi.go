<script lang="ts">
</script>

<input class="replay-bar" type="range" name="replay" min="0" max="10" />
<div class="controls">
	<button class="button-primary">&lt;</button>
	<div>
		{5}/{1012}
	</div>
	<button class="button-primary"> &gt;</button>
	<button class="button-primary">stop</button>
	<button class="button-primary">autoplay</button>
	<input type="number" name="autoplay" value={50} />
</div>
<ol class="history-list">
	<li><button>sdf</button></li>
	<li><button>sdf</button></li>
	<li><button>sdf</button></li>
	<li><button>sdf</button></li>
	<li><button>sdf</button></li>
	<li><button>sdf</button></li>
	<li><button>sdf</button></li>
	<li><button>sdf</button></li>
	<li><button>sdf</button></li>
	<li><button>sdf</button></li>
	<li><button>sdf</button></li>
	<li><button>sdf</button></li>
	<li><button>sdf</button></li>
	<li><button>sdf</button></li>
</ol>

<style lang="scss">
	.replay-bar {
		width: 100%;
    margin-bottom: 1rem;
		// height: 4rem;
	}
	.controls {
		display: flex;
    background-color: rgb(var(--bg-2));
		box-shadow: 0px 5px 15px rgba(0, 0, 0, 0.07);
    margin-bottom: 1rem;
    padding: 1rem;
	}
	.history-list {
		background-color: rgb(var(--bg-2));
		box-shadow: 0px 5px 15px rgba(0, 0, 0, 0.07);
		padding: 1rem;
    display: flex;
    flex-direction: column;
    gap: .25rem;
    border-radius: 4px;
    height: 30rem;
    overflow: auto;
	}
	li {
    // background-color: rgb(230, 230, 230);
	}
  li button{
    width: 100%;
    text-align: left;
    padding: .25rem .5rem;
    border-radius: 4px;
    &:hover{
			background-color: rgb(var(--primary));
			color: white;
		}
  }
</style>
