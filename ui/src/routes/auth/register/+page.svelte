<script lang="ts">
    import { MailIcon, KeyIcon, UserIcon, ShieldCheckIcon, LoaderIcon } from '@lucide/svelte';
    import { fade } from 'svelte/transition';

    let username = $state("");
    let email = $state("");
    let password = $state("");
    let isLoading = $state(false);

    async function handleRegister(e: Event) {
        e.preventDefault();
        isLoading = true;
        // Logic here
        setTimeout(() => isLoading = false, 1500);
    }
</script>

<div class="space-y-6" in:fade={{ duration: 300 }}>
    <div class="space-y-1">
        <h2 class="text-xl font-bold text-foreground">Create account</h2>
<!--        <p class="text-sm text-muted">Join this server.</p>-->
    </div>

    <form onsubmit={handleRegister} class="space-y-4">
        <div class="space-y-2">
            <label class="text-[10px] font-bold text-muted uppercase tracking-widest ml-1">Username</label>
            <div class="relative">
                <UserIcon size={18} class="absolute left-4 top-1/2 -translate-y-1/2 text-muted/50" />
                <input
                        type="text"
                        bind:value={username}
                        placeholder="johndoe"
                        required
                        class="w-full bg-panel border border-border rounded-2xl py-3.5 pl-12 pr-4 outline-none focus:border-frost-500 transition-all text-sm"
                />
            </div>
        </div>

        <div class="space-y-2">
            <label class="text-[10px] font-bold text-muted uppercase tracking-widest ml-1">Email Address</label>
            <div class="relative">
                <MailIcon size={18} class="absolute left-4 top-1/2 -translate-y-1/2 text-muted/50" />
                <input
                        type="email"
                        bind:value={email}
                        placeholder="name@example.com"
                        required
                        class="w-full bg-panel border border-border rounded-2xl py-3.5 pl-12 pr-4 outline-none focus:border-frost-500 transition-all text-sm"
                />
            </div>
        </div>

        <div class="space-y-2">
            <label class="text-[10px] font-bold text-muted uppercase tracking-widest ml-1">Password</label>
            <div class="relative">
                <KeyIcon size={18} class="absolute left-4 top-1/2 -translate-y-1/2 text-muted/50" />
                <input
                        type="password"
                        bind:value={password}
                        placeholder="••••••••"
                        required
                        class="w-full bg-panel border border-border rounded-2xl py-3.5 pl-12 pr-4 outline-none focus:border-frost-500 transition-all text-sm"
                />
            </div>
        </div>

        <button
                type="submit"
                disabled={isLoading}
                class="w-full py-4 bg-frost-500 text-background font-bold rounded-2xl hover:bg-frost-400 active:scale-[0.98] transition-all flex items-center justify-center gap-2 shadow-lg shadow-frost-500/20"
        >
            {#if isLoading}
                <LoaderIcon size={20} class="animate-spin" />
            {:else}
                Create Account
            {/if}
        </button>
    </form>
</div>
