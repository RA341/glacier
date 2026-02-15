<script lang="ts">

    let {game = $bindable()} = $props();
</script>

<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <div class="bg-surface border border-border rounded-3xl p-6 space-y-4">
        <h3 class="text-sm font-bold text-foreground">Indexer Info</h3>
        <div class="space-y-4">
            <div class="p-4 bg-panel border border-border rounded-2xl">
                <p class="text-[9px] font-bold text-muted uppercase">Primary Source</p>
                <p class="text-sm font-bold text-frost-400">TODO: Indexer Selector</p>
            </div>
        </div>
    </div>
</div>