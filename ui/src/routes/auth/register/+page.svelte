<script lang="ts">
    import {fade, fly} from 'svelte/transition';
    import {KeyIcon, LoaderIcon, TriangleAlert, UserIcon} from '@lucide/svelte';
    import {glacierPubCli} from "$lib/api/api";
    import {AuthService} from "$lib/gen/auth/v1/auth_pb";
    import {createRPCRunner} from "$lib/api/svelte-api.svelte";
    import {goto} from "$app/navigation";

    let username = $state("");
    let password = $state("");
    let passwordVerify = $state("");

    const authSrv = glacierPubCli(AuthService)
    const regRpc = createRPCRunner(() => authSrv.register({
        username: username,
        password: password,
        passwordVerify: passwordVerify,
    }))


    async function handleRegister(e: Event) {
        e.preventDefault();
        await regRpc.runner()
        if (!regRpc.error) {
            await goto("login")
        }
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
                <UserIcon size={18} class="absolute left-4 top-1/2 -translate-y-1/2 text-muted/50"/>
                <input
                        type="text"
                        bind:value={username}
                        placeholder="John117"
                        required
                        class="w-full bg-panel border border-border rounded-2xl py-3.5 pl-12 pr-4 outline-none focus:border-frost-500 transition-all text-sm"
                />
            </div>
        </div>

        <div class="space-y-2">
            <label class="text-[10px] font-bold text-muted uppercase tracking-widest ml-1">Password</label>
            <div class="relative">
                <KeyIcon size={18} class="absolute left-4 top-1/2 -translate-y-1/2 text-muted/50"/>
                <input
                        type="password"
                        bind:value={password}
                        placeholder="••••••••"
                        required
                        class="w-full bg-panel border border-border rounded-2xl py-3.5 pl-12 pr-4 outline-none focus:border-frost-500 transition-all text-sm"
                />
            </div>
        </div>
        <div class="space-y-2">
            <label class="text-[10px] font-bold text-muted uppercase tracking-widest ml-1">Verify Password</label>
            <div class="relative">
                <KeyIcon size={18} class="absolute left-4 top-1/2 -translate-y-1/2 text-muted/50"/>
                <input
                        type="password"
                        bind:value={passwordVerify}
                        placeholder="••••••••"
                        required
                        class="w-full bg-panel border border-border rounded-2xl py-3.5 pl-12 pr-4 outline-none focus:border-frost-500 transition-all text-sm"
                />
            </div>
        </div>

        <button
                type="submit"
                disabled={regRpc.loading}
                class="w-full py-4 bg-frost-500 text-background font-bold rounded-2xl hover:bg-frost-400 active:scale-[0.98] transition-all flex items-center justify-center gap-2 shadow-lg shadow-frost-500/20"
        >
            {#if regRpc.loading}
                <LoaderIcon size={20} class="animate-spin"/>
            {:else}
                Create Account
            {/if}
        </button>
    </form>
</div>


{#if regRpc.error}
    <div class="fixed inset-0 z-110 flex items-center justify-center p-6" transition:fade={{ duration: 100 }}>
        <!-- Darker backdrop for the error alert -->
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div class="absolute inset-0 bg-black/40 backdrop-blur-sm" onclick={regRpc.clear}></div>

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
                {regRpc.error}
            </p>

            <div class="flex justify-end pt-2">
                <button
                        onclick={regRpc.clear}
                        class="px-5 py-2 bg-panel border border-border rounded-xl text-sm font-bold hover:text-foreground transition-all"
                >
                    Dismiss
                </button>
            </div>
        </div>
    </div>
{/if}