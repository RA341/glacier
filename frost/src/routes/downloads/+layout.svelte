<script lang="ts">
    import {DownloadIcon} from '@lucide/svelte';
    import {isFrost} from "$lib/api/api";
    import {page} from '$app/state';
    import {goto} from '$app/navigation';

    let {children} = $props();

    let activeTab = $derived.by(() => {
        const path = page.url.pathname;
        if (path.includes('/downloads/frost')) return 'frost';
        if (path.includes('/downloads/glacier')) return 'server';
        return '';
    });

    function handleTabChange(tab: string) {
        const path = tab === 'frost' ? '/downloads/frost' : '/downloads/glacier';
        goto(path, {replaceState: true, noScroll: true});
    }
</script>

<div class="max-w-7xl mx-auto p-6 space-y-8 bg-background text-foreground min-h-screen">
    <header class="flex items-center justify-between border-b border-border pb-6">
        <div class="flex items-center gap-4">
            <div class="p-3 bg-frost-500/10 text-frost-400 rounded-2xl border border-frost-500/20">
                <DownloadIcon size={28}/>
            </div>
            <div>
                <h1 class="text-3xl font-bold tracking-tight">Downloads</h1>
                <p class="text-muted text-sm">Monitor and manage active downloads.</p>
            </div>
        </div>

        <div class="flex gap-1 bg-panel p-1 rounded-xl border border-border">
            {#if isFrost}
                <button
                        onclick={() => handleTabChange('frost')}
                        class="px-6 py-2 rounded-lg text-sm font-bold transition-all {activeTab === 'frost' ? 'bg-surface shadow-sm text-frost-400' : 'text-muted hover:text-foreground'}"
                >
                    Frost
                </button>
            {/if}

            <button
                    onclick={() => handleTabChange('server')}
                    class="px-6 py-2 rounded-lg text-sm font-bold transition-all {activeTab === 'server' ? 'bg-surface shadow-sm text-frost-400' : 'text-muted hover:text-foreground'}"
            >
                Server
            </button>
        </div>
    </header>

    <main>
        {@render children()}
    </main>
</div>
