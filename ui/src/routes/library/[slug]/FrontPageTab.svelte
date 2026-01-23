<script lang="ts">
    import {fade} from "svelte/transition";
    import type {Game} from "$lib/gen/library/v1/library_pb";

    let {game = $bindable(null)}: { game: Game | null } = $props();
</script>

<div class="space-y-6" in:fade>
    <div class="p-6 bg-surface border border-border rounded-3xl space-y-3">
        <h3 class="text-[10px] font-bold text-muted uppercase tracking-[0.2em]">Summary</h3>
        <p class="text-lg font-medium leading-relaxed">
            {game?.Meta?.Summary || "No summary"}
        </p>
    </div>

    <div class="p-8 bg-surface border border-border rounded-3xl space-y-4">
        <h3 class="text-[10px] font-bold text-muted uppercase tracking-[0.2em]">Description</h3>
        <p class="text-muted leading-relaxed text-sm whitespace-pre-line">
            {game?.Meta?.Description || "No description found"}
        </p>
    </div>
</div>
