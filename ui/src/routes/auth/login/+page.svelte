<script lang="ts">
    import {TriangleAlert, ArrowRightIcon, KeyIcon, LoaderIcon, User} from '@lucide/svelte';
    import {fade, fly} from 'svelte/transition';
    import {glacierPubCli} from "$lib/api/api";
    import {AuthService} from "$lib/gen/auth/v1/auth_pb";
    import {createRPCRunner} from "$lib/api/svelte-api.svelte";
    import {goto} from "$app/navigation";

    let username = $state("");
    let password = $state("");

    const authSrv = glacierPubCli(AuthService)
    const loginRpc = createRPCRunner(() => authSrv.login({
        username: username,
        password: password,
    }))

    async function handleLogin(e: Event) {
        e.preventDefault();
        await loginRpc.runner()

        if (!loginRpc.error) {
            await goto("/library", {replaceState: true})
        }
    }
</script>

<div class="space-y-6" in:fade={{ duration: 300 }}>
    <div class="space-y-1">
        <h2 class="text-xl font-bold text-foreground">Welcome back</h2>
        <p class="text-sm text-muted">Please enter your credentials to continue.</p>
    </div>

    <form onsubmit={handleLogin} class="space-y-4">
        <div class="space-y-2">
            <label for="email" class="text-[10px] font-bold text-muted uppercase tracking-widest ml-1">
                Username
            </label>
            <div class="relative">
                <User size={18} class="absolute left-4 top-1/2 -translate-y-1/2 text-muted/50"/>
                <input
                        type="text"
                        id="username"
                        bind:value={username}
                        placeholder="Username"
                        required
                        class="w-full bg-panel border border-border rounded-2xl py-3.5 pl-12 pr-4 outline-none focus:border-frost-500 transition-all text-sm"
                />
            </div>
        </div>

        <div class="space-y-2">
            <div class="flex justify-between items-center px-1">
                <label for="password"
                       class="text-[10px] font-bold text-muted uppercase tracking-widest">Password</label>
                <a href="/auth/forgot"
                   class="text-[10px] font-bold text-frost-500 hover:text-frost-400 uppercase tracking-widest">Forgot?</a>
            </div>
            <div class="relative">
                <KeyIcon size={18} class="absolute left-4 top-1/2 -translate-y-1/2 text-muted/50"/>
                <input
                        type="password"
                        id="password"
                        bind:value={password}
                        placeholder="••••••••"
                        required
                        class="w-full bg-panel border border-border rounded-2xl py-3.5 pl-12 pr-4 outline-none focus:border-frost-500 transition-all text-sm"
                />
            </div>
        </div>

        <button
                type="submit"
                disabled={loginRpc.loading}
                class="w-full py-4 bg-frost-500 text-background font-bold rounded-2xl hover:bg-frost-400 active:scale-[0.98] transition-all flex items-center justify-center gap-2 shadow-lg shadow-frost-500/20 disabled:opacity-50"
        >
            {#if loginRpc.loading}
                <LoaderIcon size={20} class="animate-spin"/>
            {:else}
                Sign In
                <ArrowRightIcon size={18}/>
            {/if}
        </button>

    </form>


    <button
            type="submit"
            onclick={() => goto("/api/server/public/auth/oidc")}
            class="w-full py-4 bg-frost-500 text-background font-bold rounded-2xl hover:bg-frost-400 active:scale-[0.98] transition-all flex items-center justify-center gap-2 shadow-lg shadow-frost-500/20 disabled:opacity-50"
    >
        {#if loginRpc.loading}
            <LoaderIcon size={20} class="animate-spin"/>
        {:else}
            Use OIDC
            <ArrowRightIcon size={18}/>
        {/if}
    </button>
</div>

{#if loginRpc.error}
    <div class="fixed inset-0 z-110 flex items-center justify-center p-6" transition:fade={{ duration: 100 }}>
        <!-- Darker backdrop for the error alert -->
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div class="absolute inset-0 bg-black/40 backdrop-blur-sm" onclick={loginRpc.clear}></div>

        <div
                class="relative w-full max-w-sm bg-surface border border-red-500/30 rounded-2xl shadow-2xl p-6 flex flex-col gap-4"
                transition:fly={{ y: 10, duration: 200 }}
        >
            <div class="flex items-center gap-3 text-red-400">
                <div class="p-2 bg-red-500/10 rounded-lg">
                    <TriangleAlert size={20}/>
                </div>
                <h3 class="font-bold">Action Failed</h3>
            </div>

            <p class="text-sm text-muted leading-relaxed">
                {loginRpc.error}
            </p>

            <div class="flex justify-end pt-2">
                <button
                        onclick={loginRpc.clear}
                        class="px-5 py-2 bg-panel border border-border rounded-xl text-sm font-bold hover:text-foreground transition-all"
                >
                    Dismiss
                </button>
            </div>
        </div>
    </div>
{/if}