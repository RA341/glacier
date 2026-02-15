<script lang="ts">
    import {fade} from "svelte/transition";
    import FileManager from "./FileManager.svelte";

    let {game = $bindable()} = $props();
</script>

<div class="bg-surface border border-border rounded-3xl overflow-hidden h-150 shadow-sm">
    <FileManager game={game.ID}/>
</div>

