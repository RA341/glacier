<script lang="ts">
    import {page} from '$app/state';
    import {DownloadIcon, LibraryIcon, LogOutIcon, MountainSnow, SearchIcon, Snowflake, UserIcon} from "@lucide/svelte";
    import {glacierPubCli, isFrost} from "$lib/api/api";
    import {AuthService} from "$lib/gen/auth/v1/auth_pb";
    import {goto} from "$app/navigation";
    import {getUserCtx} from "$lib/components/user/provider.svelte";

    const isActive = (href: string) => page.url.pathname.startsWith(href);

    const links = [
        {label: 'Library', href: '/library', icon: LibraryIcon},
        {label: 'Search', href: '/search', icon: SearchIcon},
        {label: 'Downloads', href: '/downloads', icon: DownloadIcon},
    ];

    const user = getUserCtx()

    const footerLinks = [
        {label: 'Profile', href: '/settings/user', icon: UserIcon}
    ];

    if (user.isOmni) {
        footerLinks.unshift({label: 'Server Settings', href: '/settings/glacier', icon: MountainSnow});
    }

    if (isFrost) {
        footerLinks.unshift({label: 'Frost Settings', href: '/settings/frost', icon: Snowflake});
    }

    const auth = glacierPubCli(AuthService)

    async function logout() {
        await auth.logout({})
        await goto("/auth/login", {replaceState: true})
    }
</script>

<aside class="flex flex-col w-16 h-full border-r border-border bg-surface items-center py-6">
    <!-- Logo Area -->
    <div class="mb-8 flex items-center justify-center">
        <img src={"/favicon.svg"} alt="Logo" class="w-8 h-8 rounded-lg shadow-frost"/>

    </div>

    <!-- Main Navigation -->
    <nav class="flex flex-col flex-1 gap-4">
        {#each links as link}
            {@const active = isActive(link.href)}
            <a
                    href={link.href}
                    title={link.label}
                    class="group relative flex items-center justify-center w-10 h-10 rounded-xl transition-all
            duration-300
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
                    class="group relative flex items-center justify-center w-10 h-10 rounded-xl transition-all
            duration-300
            {active
                ? 'bg-frost-500/10 text-frost-400 shadow-frost'
                : 'text-muted hover:text-frost-400 hover:bg-panel'}"
            >
                <link.icon size={22} strokeWidth={2}/>
            </a>
        {/each}

        <button
                onclick={logout}
                class="group relative flex items-center justify-center w-10 h-10 rounded-xl transition-all duration-300"
        >
            <LogOutIcon size={22} strokeWidth={2}/>
        </button>
    </div>
</aside>
