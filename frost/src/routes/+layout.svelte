<script lang="ts">
    import './layout.css';
    import {page} from '$app/state';
    import favicon from '$lib/assets/favicon.svg';
    import {LibraryIcon, SearchIcon, SettingsIcon, UserIcon} from "@lucide/svelte";

    let {children} = $props();

    const links = [
        {label: 'Library', href: '/library', icon: LibraryIcon},
        {label: 'Search', href: '/search', icon: SearchIcon}
    ];

    const footerLinks = [
        {label: 'Settings', href: '/settings/general', icon: SettingsIcon},
        {label: 'Profile', href: '/settings/profile', icon: UserIcon}
    ];

    const isActive = (href: string) => page.url.pathname.startsWith(href);
</script>

<svelte:head>
    <link rel="icon" href={favicon}/>
</svelte:head>

<div class="flex h-screen w-full overflow-hidden bg-background text-foreground">

    <!-- Sidebar -->
    <aside class="flex flex-col w-16 h-full border-r border-border bg-surface items-center py-6">

        <!-- Logo Area -->
        <div class="mb-8 flex items-center justify-center">
            <img src={favicon} alt="Logo" class="w-8 h-8 rounded-lg shadow-frost"/>
        </div>

        <!-- Main Navigation -->
        <nav class="flex flex-col flex-1 gap-4">
            {#each links as link}
                {@const active = isActive(link.href)}
                <a
                        href={link.href}
                        title={link.label}
                        class="group relative flex items-center justify-center w-10 h-10 rounded-xl transition-all duration-300
                            {active
                                ? 'bg-frost-500/10 text-frost-400 shadow-frost'
                                : 'text-muted hover:text-frost-400 hover:bg-panel'}"
                >

                    <link.icon
                            size={22}
                            strokeWidth={active ? 2.5 : 2}
                            class="transition-transform duration-200 {active ? 'scale-110' : 'group-active:scale-90'}"
                    />
                </a>
            {/each}
        </nav>

        <!-- Footer Navigation -->
        <div class="flex flex-col gap-4">
            {#each footerLinks as link}
                {@const active = isActive(link.href)}
                <a
                        href={link.href}
                        title={link.label}
                        class="group relative flex items-center justify-center w-10 h-10 rounded-xl transition-all duration-300
                            {active
                                ? 'bg-frost-500/10 text-frost-400 shadow-frost'
                                : 'text-muted hover:text-frost-400 hover:bg-panel'}"
                >
                    <link.icon size={22} strokeWidth={2}/>
                </a>
            {/each}
        </div>

    </aside>

    <!-- Main Content Area -->
    <main class="relative flex-1 overflow-y-auto">
        {@render children()}
    </main>
</div>

<style>
    a:hover::after {
        content: attr(title);
        position: absolute;
        left: 3.5rem;
        background: var(--color-panel);
        color: var(--color-foreground);
        padding: 4px 8px;
        border-radius: 4px;
        font-size: 0.75rem;
        white-space: nowrap;
        border: 1px solid var(--color-border);
        pointer-events: none;
        z-index: 50;
    }

    ::-webkit-scrollbar {
        width: 6px;
    }

    ::-webkit-scrollbar-track {
        background: transparent;
    }

    ::-webkit-scrollbar-thumb {
        background: var(--border);
        border-radius: 10px;
    }

    ::-webkit-scrollbar-thumb:hover {
        background: #444;
    }
</style>
