<script lang="ts">
    import {
        CircleAlert,
        GlobeIcon,
        PowerIcon,
        RefreshCwIcon,
        LoaderIcon,
        SaveIcon,
        SettingsIcon
    } from '@lucide/svelte';
    import {fade, fly} from 'svelte/transition';

    import {callRPC, Frost, frostCli} from "$lib/api/api";
    import {ConfigService} from "$lib/gen/config/v1/config_pb";
    import {createRPCRunner} from "$lib/api/rpc.svelte";
    import {onMount} from "svelte";
    import {getSnackbarCtx} from "$lib/components/snackbar/snackbar-provider.svelte";
    import ConfigField from "$lib/components/ConfigField.svelte";
    import {extractValues} from "$lib/components/util.svelte";
    import {getRandomIntInclusive} from "$lib/api/byte-math.ts";
    import {getDialogCtx} from "$lib/components/dialog/dialog.svelte.ts";

    const sm = getSnackbarCtx()
    const libConf = frostCli(ConfigService)
    let configSchema = $state({});

    let isSaving = $state(false);

    async function handleSave() {
        isSaving = true

        const payload = extractValues(configSchema);

        const {err} = await callRPC(() => libConf.set({
            configSchema: JSON.stringify(payload)
        }))

        if (err) {
            sm.push(`Could not update config: ${err}`, 'error')
        }

        await configRpc.runner()

        isSaving = false
    }


    const configRpc = createRPCRunner(() => libConf.get({}))

    onMount(() => {
        configRpc.runner()
    })

    $effect(() => {
        if (configRpc.value) {
            configSchema = JSON.parse(configRpc.value.configSchema)
        }
    })

    let isRestarting = $state(false);
    let isShuttingDown = $state(false);
    let restartStatus = $state("Waiting for server...");

    async function handleRestart() {
        try {
            isRestarting = true;
            restartStatus = "Restarting client standby...";

            // Trigger restart
            await fetch(`${Frost.base}/restart`);

            // Give the server 2 seconds to actually shut down before we start pinging
            setTimeout(startPingLoop, 2000);
        } catch (e) {
            isRestarting = false;
            sm.push("Restart request failed", "error");
        }
    }

    function startPingLoop() {
        restartStatus = "Reconnecting to Glacier...";
        const messages = [
            "Test test",
            "Finding my marbles"
        ]

        const interval = setInterval(async () => {
            try {
                // We use a short timeout so the UI doesn't hang on a stuck socket
                const controller = new AbortController();
                const timeoutId = setTimeout(() => controller.abort(), 900);

                const res = await fetch(`${Frost.base}/hello`, {signal: controller.signal});

                if (res.ok) {
                    clearInterval(interval);
                    restartStatus = "Connection established!";

                    // Final delay for internal services to warm up
                    setTimeout(() => {
                        window.location.reload();
                    }, 1000);
                }
            } catch (e) {
                // Fetch error is expected while server is down
                restartStatus = messages[getRandomIntInclusive(0, messages.length)];
            }
        }, 1000);
    }

    const dm = getDialogCtx()

    async function handleShutdown() {
        if (!(await dm.confirm("You sure about that", "Just checkin", 'error'))) {
            return
        }

        try {
            isShuttingDown = true;
            await fetch(`${Frost.base}/shutdown`);
            sm.push("Shutdown command sent", "warn");
        } catch (e) {
            isShuttingDown = false;
            sm.push("Shutdown request failed", "error");
        }
    }
</script>

{#if configRpc.loading && !configRpc.value}
    <!-- THEMED LOADING STATE -->
    <div class="flex flex-col items-center justify-center h-[60vh] gap-4 text-muted">
        <LoaderIcon class="animate-spin text-frost-500" size={40}/>
        <p class="text-sm font-bold uppercase tracking-widest animate-pulse">Loading Configuration...</p>
    </div>

{:else if configRpc.error}
    <!-- THEMED ERROR STATE -->
    <div class="max-w-md mx-auto mt-20 p-6 bg-red-500/5 border border-red-500/20 rounded-3xl flex flex-col items-center text-center gap-4">
        <CircleAlert size={48} class="text-red-400"/>
        <div>
            <h3 class="font-bold text-foreground">Failed to load config</h3>
            <p class="text-sm text-muted mt-1">{configRpc.error}</p>
        </div>
        <button onclick={() => configRpc.runner()}
                class="px-6 py-2 bg-panel border border-border rounded-xl text-xs font-bold hover:bg-surface transition-all">
            Retry Connection
        </button>
    </div>

{:else}
    <div class="max-w-4xl mx-auto p-8 min-h-screen">
        <!-- STICKY HEADER -->
        <header class="sticky top-0 z-30 bg-background/80 backdrop-blur-md border-b border-border mb-10 -mx-8 px-8 py-6 flex justify-between items-center">
            <div class="flex items-center gap-4">
                <div class="p-3 bg-panel border border-border rounded-2xl text-frost-400">
                    <SettingsIcon size={24}/>
                </div>
                <div>
                    <h1 class="text-2xl font-black uppercase tracking-tight text-foreground leading-tight">Frost
                        Settings</h1>
                    <p class="text-xs text-muted font-medium uppercase tracking-wider">Configure the frost client</p>
                </div>
            </div>

            <div class="flex gap-2">
                <button
                        onclick={handleRestart}
                        disabled={isSaving || isRestarting}
                        class="flex items-center gap-2 px-5 py-2.5 bg-panel border border-border text-muted rounded-xl text-xs font-bold hover:text-frost-400 transition-all active:scale-95"
                >
                    <RefreshCwIcon size={16} class={isRestarting ? 'animate-spin' : ''}/>
                    Restart
                </button>

                <button
                        onclick={handleShutdown}
                        disabled={isSaving || isRestarting}
                        class="flex items-center gap-2 px-5 py-2.5 bg-panel border border-border text-muted rounded-xl text-xs font-bold hover:text-red-400 transition-all active:scale-95"
                >
                    <PowerIcon size={16}/>
                    Shutdown
                </button>

                <div class="w-px h-8 bg-border mx-2"></div>

                <button
                        onclick={handleSave}
                        disabled={isSaving || isRestarting}
                        class="flex items-center gap-2 px-8 py-2.5 bg-frost-500 text-background rounded-xl text-sm font-bold hover:bg-frost-400 transition-all shadow-lg shadow-frost-500/20 active:scale-95 disabled:opacity-50"
                >
                    {#if isSaving}
                        <LoaderIcon size={18} class="animate-spin"/>
                    {:else}
                        <SaveIcon size={18}/>
                        Save Changes
                    {/if}
                </button>
            </div>
        </header>

        <div class="space-y-12 pb-20">
            <ConfigField schema={configSchema}/>
        </div>
    </div>
{/if}

<!-- BLOCKING RESTART MODAL -->
{#if isRestarting}
    <div class="fixed inset-0 z-200 flex items-center justify-center p-6" transition:fade={{ duration: 200 }}>
        <div class="absolute inset-0 bg-background/90 backdrop-blur-xl"></div>

        <div
                class="relative w-full max-w-sm bg-surface border border-border rounded-[2.5rem] shadow-2xl p-10 flex flex-col items-center text-center gap-6"
                transition:fly={{ y: 20, duration: 400 }}
        >
            <div class="relative">
                <div class="p-5 bg-frost-500/10 rounded-full text-frost-500 animate-pulse">
                    <RefreshCwIcon size={48} strokeWidth={1.5} class="animate-spin-slow"/>
                </div>
                <div class="absolute -bottom-1 -right-1 p-2 bg-background border border-border rounded-full text-frost-400">
                    <GlobeIcon size={16}/>
                </div>
            </div>

            <div class="space-y-2">
                <h2 class="text-xl font-bold text-foreground">Restarting Server</h2>
                <p class="text-sm text-muted font-medium">{restartStatus}</p>
            </div>

            <!-- INDETERMINATE PROGRESS BAR -->
            <div class="w-full h-1.5 bg-panel rounded-full overflow-hidden border border-border">
                <div class="h-full bg-frost-500 w-1/3 rounded-full animate-progress-loop"></div>
            </div>

            <p class="text-[10px] text-muted uppercase tracking-[0.2em] font-bold opacity-50">
                So hows it going, hows the family
            </p>
        </div>
    </div>
{/if}

<style>
    :global(.animate-spin-slow) {
        animation: spin 3s linear infinite;
    }

    @keyframes progress-loop {
        0% {
            transform: translateX(-100%);
        }
        100% {
            transform: translateX(300%);
        }
    }

    .animate-progress-loop {
        animation: progress-loop 1.5s infinite linear;
    }
</style>