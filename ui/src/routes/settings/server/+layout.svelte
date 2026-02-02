<script lang="ts">
    import {page} from '$app/state';

    let {children} = $props();

    const prefix = '/settings/server';

    const tabs = [
        {label: 'General', href: 'general', desc: 'Configure basic server settings'},
        {label: 'Users', href: 'users', desc: 'Manage user accounts'},
        {label: 'Metadata', href: 'metadata', desc: 'Where to get game information'},
        {label: 'Indexer', href: 'indexer', desc: 'Where to find and download games'},
        {label: 'Download', href: 'downloads', desc: 'How downloads are handled'},
    ];

    let currentPath = $derived(page.url.pathname);

    function withPrefix(ref: string) {
        return `${prefix}/${ref}`
    }

    let activeTab = $derived(tabs.find(tab => currentPath.startsWith(withPrefix(tab.href))));
</script>

<div class="flex flex-col h-full bg-background">
    <header class="px-8 pt-8 pb-0">
        <div class="flex items-baseline gap-4 mb-6">
            <h1 class="text-3xl font-bold tracking-tight text-foreground">
                Server Settings
            </h1>

            {#if activeTab}
                <span class="text-muted/30 text-2xl font-light">/</span>
                <p class="text-lg font-medium text-muted transition-all duration-300">
                    {activeTab.label}
                </p>
                <span class="text-muted/30 text-2xl font-light">/</span>
                <p class="text-lg font-medium text-muted transition-all duration-300">
                    {activeTab.desc}
                </p>
            {/if}
        </div>

        <nav class="flex gap-2 border-b border-border pb-2 relative">
            {#each tabs as tab}
                {@const active = currentPath.startsWith(withPrefix(tab.href))}
                <a
                        href={withPrefix(tab.href)}
                        class="group relative px-4 py-2 rounded-lg text-sm font-medium transition-all duration-300
                    {active
                        ? 'bg-frost-500/10 text-frost-400'
                        : 'text-muted hover:text-foreground hover:bg-panel'}"
                >
                    {tab.label}

                    {#if active}
                        <div class="absolute -bottom-2.25 left-0 right-0 h-0.5 bg-frost-500 rounded-full"></div>
                    {/if}
                </a>
            {/each}
        </nav>
    </header>

    <main class="flex-1 overflow-y-auto p-8">
        {@render children()}
    </main>
</div>
