<script lang="ts">
    import {page} from '$app/state';

    let {children} = $props();

    const tabs = [
        {label: 'General', href: '/settings/general'},
        {label: 'Profile', href: '/settings/profile'},
    ];

    let currentPath = $derived(page.url.pathname);
</script>

<div class="flex flex-col h-full bg-background">
    <!-- Header Section -->
    <header class="px-8 pt-8 pb-0">
        <h1 class="text-3xl font-bold tracking-tight text-foreground mb-6">Settings</h1>

        <!-- Tab List Container -->
        <!-- Added 'pb-2' to give space between the pills and the bottom border -->
        <nav class="flex gap-2 border-b border-border pb-2 relative">
            {#each tabs as tab}
                {@const active = currentPath.startsWith(tab.href)}
                <a
                        href={tab.href}
                        class="group relative px-4 py-2 rounded-lg text-sm font-medium transition-all duration-300
                    {active
                        ? 'bg-frost-500/10 text-frost-400'
                        : 'text-muted hover:text-foreground hover:bg-panel'}"
                >
                    {tab.label}

                    {#if active}
                        <!-- Optional: Small high-intensity underline indicator -->
                        <div
                                class="absolute -bottom-2.25 left-0 right-0 bg-frost-500 rounded-full"
                        ></div>
                    {/if}
                </a>
            {/each}
        </nav>
    </header>

    <!-- Content Area -->
    <main class="flex-1 overflow-y-auto p-8">
        {@render children()}
    </main>
</div>

<style>
    nav {
        scrollbar-width: none;
        -ms-overflow-style: none;
    }

    nav::-webkit-scrollbar {
        display: none;
    }
</style>
