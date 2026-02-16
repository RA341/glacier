<script lang="ts">
    import {page} from '$app/state';
    import {goto} from '$app/navigation';
    import {fade, fly} from 'svelte/transition';
    import {
        CheckCircle2Icon,
        CircleAlert,
        Globe,
        GlobeIcon,
        LoaderIcon,
        LockIcon,
        PencilIcon,
        RefreshCwIcon,
        UserPlusIcon
    } from '@lucide/svelte';
    import {appName, callRPC, frostCli, isFrost} from "$lib/api/api";
    import {onMount} from 'svelte';
    import {ConfigService} from "$lib/gen/config/v1/config_pb";
    import {getSnackbarCtx} from "$lib/components/snackbar/snackbar-provider.svelte";
    import {createRPCRunner} from "$lib/api/rpc.svelte";
    import {getDialogCtx} from "$lib/components/dialog/dialog.svelte";

    let {children} = $props();

    const sm = getSnackbarCtx();
    const confSrv = frostCli(ConfigService);

    let activeTab = $derived(page.url.pathname.includes('register') ? 'register' : 'login');

    function navigate(path: string) {
        goto(`/auth/${path}`, {replaceState: true});
    }

    let glacierUrlInput = $state("");
    let isSettingUrl = $state(false);
    let forceEdit = $state(false);
    let glacierUrlCheckRpc = createRPCRunner(() => confSrv.getField({}));
    let isGlacierReady = $derived(!isFrost || (!!glacierUrlCheckRpc?.value?.Value && !glacierUrlCheckRpc?.value?.Error));
    let showEditor = $derived(isFrost && (!isGlacierReady || forceEdit));

    onMount(() => {
        if (isFrost) {
            glacierUrlCheckRpc.runner();
        }
    });

    const dm = getDialogCtx()

    async function handleSetUrl() {
        if (!glacierUrlInput) return;
        isSettingUrl = true;
        forceEdit = false;

        console.log("Setting Glacier URL to:", glacierUrlInput);
        const {err} = await callRPC(() => confSrv.setField({Value: glacierUrlInput}));
        if (err) {
            await dm.alert("Could not connect to glacier", err, 'error')
        } else {
            await glacierUrlCheckRpc.runner();
        }

        isSettingUrl = false;
    }


</script>

<div class="min-h-screen w-full bg-background flex items-center justify-center p-6 relative overflow-hidden">
    <div class="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-frost-500/10 blur-[120px] rounded-full"></div>
    <div class="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-frost-900/20 blur-[120px] rounded-full"></div>

    <div class="w-full max-w-110 z-10" transition:fly={{ y: 20, duration: 600 }}>
        <div class="flex flex-col items-center gap-3 mb-8">
            <img src="/favicon.svg" alt="logo" height="60" width="60" class="drop-shadow-2xl"/>
            <h1 class="text-3xl font-black tracking-tighter text-foreground uppercase">{appName}</h1>
            <p class="text-center text-xs font-bold uppercase tracking-[0.2em] text-frost-50/50">
                <span class="text-frost-500 font-black">COLD</span> Storage for your games
            </p>
        </div>

        {#if showEditor}
            {#if isFrost && glacierUrlCheckRpc.loading && !glacierUrlCheckRpc.value}
                <div class="flex flex-col items-center justify-center py-20 text-muted gap-4" in:fade>
                    <LoaderIcon class="animate-spin text-frost-500" size={32}/>
                    <p class="text-[10px] font-bold uppercase tracking-widest animate-pulse">Checking Instance
                        Status...</p>
                </div>

            {:else}
                <!-- GLACIER URL ENTRY VIEW -->
                <div class="bg-surface border border-border rounded-4xl shadow-2xl p-8 space-y-4" in:fly={{ y: 10 }}>
                    <div class="space-y-2">
                        <h2 class="text-xl font-bold text-foreground">Connect Glacier</h2>
                        <p class="text-sm text-muted">Enter url to a Glacier instance</p>
                    </div>

                    {#if glacierUrlCheckRpc.error}
                        <div class="flex items-start gap-3 p-4 rounded-2xl bg-red-500/5 border border-red-500/20 text-red-400">
                            <CircleAlert size={18} class="shrink-0"/>
                            <div class="text-xs">
                                <p class="font-bold">Connection Failed</p>
                                <p class="opacity-80">{glacierUrlCheckRpc.error}</p>
                            </div>
                        </div>
                    {/if}

                    <div class="space-y-4">
                        <div class="space-y-2">
                            <div class="relative">
                                <GlobeIcon size={18} class="absolute left-4 top-1/2 -translate-y-1/2 text-muted/50"/>
                                <input type="url"
                                       bind:value={glacierUrlInput}
                                       placeholder="glacier.yourdomain.com"
                                       class="w-full bg-panel border border-border rounded-2xl py-3.5 pl-12 pr-4 outline-none focus:border-frost-500 transition-all text-sm"
                                       onkeydown={(e) => {
                                            if (e.key === 'Enter' && glacierUrlInput && !isSettingUrl) {
                                                handleSetUrl();
                                            }
                                        }}
                                />
                            </div>
                        </div>

                        <button
                                onclick={handleSetUrl}
                                disabled={isSettingUrl || !glacierUrlInput}
                                class="w-full py-4 bg-frost-500 text-background font-bold rounded-2xl hover:bg-frost-400 active:scale-[0.98] transition-all flex items-center justify-center gap-2 shadow-lg shadow-frost-500/20 disabled:opacity-50"
                        >
                            {#if isSettingUrl}
                                <LoaderIcon size={20} class="animate-spin"/>
                                Connecting...
                            {:else}
                                <Globe size={18}/>
                                Connect
                            {/if}
                        </button>
                    </div>
                </div>
            {/if}
        {:else}
            {#if isFrost && glacierUrlCheckRpc.value}
                <div class="flex items-center justify-center mb-6" in:fade>
                    <div class="flex items-center gap-3 px-4 py-2 bg-frost-500/5 border border-frost-500/10 rounded-full shadow-sm">
                        <CheckCircle2Icon size={14} class="text-green-500"/>
                        <span class="text-[10px] font-mono text-frost-400/80 truncate max-w-50">
                            {glacierUrlCheckRpc.value.Value}
                        </span>
                        <div class="w-px h-3 bg-border mx-1"></div>
                        <button
                                onclick={() => {
                                glacierUrlInput = glacierUrlCheckRpc?.value?.Value ?? "";
                                forceEdit = true;
                            }}
                                class="text-[10px] font-black uppercase tracking-widest text-muted hover:text-frost-400 transition-colors flex items-center gap-1"
                        >
                            <PencilIcon size={10}/>
                            Edit
                        </button>
                    </div>
                </div>
            {/if}

            <div class="bg-surface border border-border rounded-4xl shadow-2xl overflow-hidden flex flex-col" in:fade>
                <div class="p-2 bg-panel/30 border-b border-border flex gap-1">
                    <button
                            onclick={() => navigate('login')}
                            class="flex-1 flex items-center justify-center gap-2 py-3 rounded-2xl text-sm font-bold transition-all
                        {activeTab === 'login' ? 'bg-surface border border-border text-frost-400 shadow-sm' : 'text-muted hover:text-foreground hover:bg-panel'}"
                    >
                        <LockIcon size={16}/>
                        Login
                    </button>
                    <button onclick={() => navigate('register')}
                            class="flex-1 flex items-center justify-center gap-2 py-3 rounded-2xl text-sm font-bold transition-all
                        {activeTab === 'register' ? 'bg-surface border border-border text-frost-400 shadow-sm' : 'text-muted hover:text-foreground hover:bg-panel'}"
                    >
                        <UserPlusIcon size={16}/>
                        Register
                    </button>
                </div>
                <main class="p-8">
                    {@render children()}
                </main>
            </div>

            <button
                    onclick={() => glacierUrlCheckRpc.runner()}
                    class="mt-6 mx-auto flex items-center gap-2 text-[10px] font-bold text-muted/40 uppercase tracking-widest hover:text-frost-500 transition-colors"
            >
                <RefreshCwIcon size={12}/>
                Refresh Instance Check
            </button>
        {/if}
    </div>
</div>

<style>
    :global(body) {
        background-color: #0a0a0c;
    }
</style>
