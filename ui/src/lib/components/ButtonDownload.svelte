<script lang="ts">
    import {
        ChevronDownIcon,
        Download,
        DownloadIcon,
        FileSearchIcon,
        FolderOpenIcon,
        Hammer,
        LoaderIcon,
        PauseIcon,
        PlayIcon,
        RefreshCwIcon,
        Replace,
        Trash2Icon,
        TriangleAlert,
        XIcon
    } from '@lucide/svelte';
    import {fade, fly} from 'svelte/transition';
    import {onMount} from "svelte";
    import {createRPCRunner} from "$lib/api/rpc.svelte.js";
    import {callRPC, Frost, frostCli} from "$lib/api/api";
    import {FrostLibraryService, GamePlaySchema} from "$lib/gen/frost_library/v1/frost_library_pb";
    import {getSnackbarCtx} from "$lib/components/snackbar/snackbar-provider.svelte";
    import {trimPrefix} from "$lib/api/strings";
    import {handleProcessStatus} from "$lib/api/websockets";
    import {create} from "@bufbuild/protobuf";
    import {getDialogCtx} from "$lib/components/dialog/dialog.svelte";

    let {gameId}: { gameId: bigint } = $props();

    const llService = frostCli(FrostLibraryService);
    const sm = getSnackbarCtx();

    const localGameRpc = createRPCRunner(() => llService.getByGameId({id: gameId, localDownload: true}));
    let isInstalled = $derived(!!localGameRpc.value?.download);

    let localGame = $derived(localGameRpc.value?.download);
    let gameDownloadState = $derived(localGame?.Download?.Status);

    let isComplete = $derived(gameDownloadState && gameDownloadState === "StatusComplete")

    let isSending = $state(false);
    let isMenuOpen = $state(false);
    let errorMessage = $state<string | null>(null);

    onMount(() => {
        localGameRpc.runner();
    });

    $effect(() => {
        // if the game is in the library but not yet complete
        if (isInstalled && !isComplete) {
            const interval = setInterval(() => {
                localGameRpc.runner();
            }, 3000);

            return () => clearInterval(interval);
        }
    });

    let isRunning = $state(false);
    let cleanupWs: (() => void) | null = null; // To store the unsubscribe/close function

    function watchProcess() {
        const exe = localGame?.Play?.LaunchExe ?? localGame?.game?.file?.Exe;
        if (!localGame?.ID || cleanupWs) {
            console.log("not watching process because fields are empty", localGame?.game?.file?.Exe, localGame?.Play?.LaunchExe, !localGame?.ID)
            return;
        }

        const proc = `${Frost.base}/launcher/running/${localGame.ID}`;

        cleanupWs = handleProcessStatus({
            url: proc,
            onDone: () => {
                isRunning = false;
                cleanupWs = null;
                launchFinalExePicker({override: false})
            },
            onConnect: () => {
                isRunning = true;
            }
        });
    }

    let launchPicker = $state(false)
    let customExePath = $state("");
    $effect(() => {
        customExePath = localGame?.Play?.LaunchExe ?? ""
    })

    let isInstaller = $derived(localGame?.game?.Source?.GameType === "Installer")
    let hasLaunchExe = $derived(localGame?.Play?.LaunchExe)

    function launchFinalExePicker({override = false}: { override?: boolean }) {
        if (!override && (!isInstaller || hasLaunchExe)) return;
        launchPicker = true
    }

    async function pickFinalExe() {
        const {err} = await callRPC(() => llService.edit({
                LocalGame: {
                    ID: localGame!.ID,
                    Play: create(GamePlaySchema, {
                        LaunchExe: customExePath
                    })
                }
            })
        )
        if (err) sm.push(`Could not save exe ${err}`, 'error');

        await localGameRpc.runner()
        launchPicker = false
    }

    async function triggerFilePicker() {
        const {val, err} = await callRPC(() => llService.launchFilePicker({
            baseDir: localGame?.Play?.LaunchExe ?? localGame?.Download?.DownloadPath ?? ""
        }))

        if (err) sm.push(`Could not launch FilePicker: ${err}`, 'error');
        if (val?.path === "") sm.push("empty file path", 'warn');
        else customExePath = val!.path
    }

    $effect(() => {
        if (isComplete) {
            watchProcess();
        }
        return () => {
            cleanupWs?.()
        };
    });

    async function handlePrimaryAction() {
        if (isRunning) {
            sm.push("Game is already running", "info");
            return;
        }

        if (isInstalled) {
            const playExe = localGame?.Play?.LaunchExe
            if (!isInstaller || playExe) {
                sm.push("launching game", "info")
            } else if (localGame?.game?.file?.Exe) {
                sm.push("launching installer", "info")
            }

            await launchGame();
            return;
        }

        await download();
    }

    async function launchGame() {
        const {err} = await callRPC(() => llService.launch({Id: gameId}))
        if (err) {
            sm.push(`Error launching game ${err}`, 'error');
        }

        watchProcess();
    }

    async function download() {
        if (isSending) return;
        isSending = true;
        const {err} = await callRPC(() => llService.download({gameId, downloadFolder: "./downloads"}));
        if (err) errorMessage = String(err);
        else sm.push("Download started", "success");
        isSending = false;

        await localGameRpc.runner();
    }

    async function recheck() {
        isMenuOpen = false;

        sm.push("Recheck initiated", "info");
        const {err} = await callRPC(() => llService.download({gameId: gameId, recheck: true}));
        if (err) sm.push(String(err), "error");

        isMenuOpen = false;

        await localGameRpc.runner();
    }

    const dm = getDialogCtx()

    async function forceDownload() {
        const ok = await dm.confirm(
            "Delete and Redownload",
            "This will delete existing files and download all files again, consider verifying files if you are not sure",
            "warn"
        )
        if (!ok) return

        isMenuOpen = false;
        const {err} = await callRPC(() => llService.download({
            gameId: gameId,
            force: true,
            recheck: true
        }));

        if (err) sm.push(String(err), "error");
        else sm.push("Force redownload started", "info");

        await localGameRpc.runner();
    }

    async function deleteLocalGame() {
        const gameID = localGameRpc.value?.download?.GameID;
        if (!gameID) {
            return
        }

        if (localGame?.game?.Source?.GameType == "Installer") {
            const ok = await dm.confirm(
                "Easy there",
                "This will delete the installer files, the installed game files will not be affected. Continue?",
                "error"
            )
            if (!ok) return;
        } else {
            const ok = await dm.confirm(
                "WHOA, YOU SURE ABOUT THAT",
                "This will ALL game files, you will not be able to play again without redownloading, Are you sure?",
                "error"
            )
            if (!ok) return;
        }

        isMenuOpen = false;
        const {err} = await callRPC(() => llService.delete({
            id: gameID,
        }));

        if (err) sm.push(String(err), "error");
        else sm.push("Game files deleted", "info");

        await localGameRpc.runner();
    }
</script>

<div class="relative inline-flex items-stretch gap-1">
    <!-- Main Action Button -->
    <div class="relative inline-flex items-stretch gap-1">
        <button
                onclick={handlePrimaryAction}
                disabled={isSending || localGameRpc.loading}
                class="px-8 py-2 bg-frost-500 text-background text-sm font-bold hover:bg-frost-400 disabled:opacity-50 transition-all flex items-center gap-2 shadow-lg shadow-frost-500/20 active:scale-[0.98]
        {isInstalled ? 'rounded-l-xl' : 'rounded-xl'}"
        >
            {#if (localGameRpc.loading && !localGameRpc.value?.download) || isSending}
                <LoaderIcon size={16} class="animate-spin"/>
                <span>Checking files...</span>
            {:else if gameDownloadState && !isComplete}
                <span>{trimPrefix(gameDownloadState, "Status")}</span>
            {:else if isInstalled}
                {#if hasLaunchExe}
                    {#if isRunning}
                        <PauseIcon size={16} fill="currentColor"/>
                        <span>Playing</span>
                    {:else}
                        <PlayIcon size={16} fill="currentColor"/>
                        <span>Play</span>
                    {/if}
                {:else if isInstaller}
                    {#if isRunning}
                        <Download size={16} fill="currentColor"/>
                        <span>Installing</span>
                    {:else}
                        <Download size={16} fill="currentColor"/>
                        <span>Install</span>
                    {/if}
                {/if}

            {:else}
                <DownloadIcon size={16}/>
                <span>Install</span>
            {/if}
        </button>

        <!-- Dropdown Toggle - Only visible if installed -->
        {#if isInstalled}
            <button
                    onclick={() => isMenuOpen = !isMenuOpen}
                    class="px-3 bg-frost-500 text-background rounded-r-xl border-l border-background/10 hover:bg-frost-400 transition-all shadow-lg shadow-frost-500/20 active:scale-[0.98]"
            >
                <ChevronDownIcon size={16} class="transition-transform {isMenuOpen ? 'rotate-180' : ''}"/>
            </button>
        {/if}

        {#if isMenuOpen}
            <!-- svelte-ignore a11y_click_events_have_key_events -->
            <!-- svelte-ignore a11y_no_static_element_interactions -->
            <div class="fixed inset-0 z-110" onclick={() => isMenuOpen = false}></div>

            <div
                    class="absolute right-0 top-full mt-2 w-48 bg-surface border border-border rounded-xl shadow-2xl overflow-hidden z-120"
                    transition:fly={{ y: 5, duration: 150 }}
            >
                <button onclick={(e) => {
                    e.preventDefault()
                    e.stopPropagation()
                    launchFinalExePicker({override: true})
                }}
                        class="w-full flex items-center gap-2 px-4 py-3 text-xs font-bold text-indigo-500 hover:text-frost-400 hover:bg-panel transition-all">
                    <Replace size={14}/>
                    Set Game Exe
                </button>
                {#if isInstaller}
                    <!--                    // todo-->
                    <!--                    <button onclick={(e) => {-->
                    <!--                        e.preventDefault()-->
                    <!--                        e.stopPropagation()-->
                    <!--                        launchFinalExePicker({override: true})-->
                    <!--                    }}-->
                    <!--                            class="w-full flex items-center gap-2 px-4 py-3 text-xs font-bold text-green-700 hover:text-frost-400 hover:bg-panel transition-all">-->
                    <!--                        <RefreshCwIcon size={14}/>-->
                    <!--                        Reinstall-->
                    <!--                    </button>-->
                {/if}
                <button onclick={recheck}
                        class="w-full flex items-center gap-2 px-4 py-3 text-xs font-bold text-green-700 hover:text-frost-400 hover:bg-panel transition-all">
                    <RefreshCwIcon size={14}/>
                    Recheck Files
                </button>
                <button onclick={forceDownload}
                        class="w-full flex items-center gap-2 px-4 py-3 text-xs font-bold text-yellow-400 hover:bg-red-500/10 transition-all">
                    <Hammer size={14}/>
                    Force Download
                </button>

                <button onclick={deleteLocalGame}
                        class="w-full flex items-center gap-2 px-4 py-3 text-xs font-bold text-red-400 hover:bg-red-500/10 transition-all">
                    <Trash2Icon size={14}/>
                    Delete
                </button>
            </div>
        {/if}
    </div>
</div>


<!-- Custom EXE Picker Dialog -->
{#if launchPicker}
    <div class="fixed inset-0 z-150 flex items-center justify-center p-6" transition:fade={{ duration: 150 }}>
        <div class="absolute inset-0 bg-background/80 backdrop-blur-md" onclick={() => launchPicker = false}></div>
        <div class="relative w-full max-w-md bg-surface border border-border rounded-2xl shadow-2xl p-6 flex flex-col gap-4"
             transition:fly={{ y: 20, duration: 300 }}>

            <div class="flex items-center justify-between">
                <div class="flex items-center gap-3 text-frost-400">
                    <div class="p-2 bg-frost-500/10 rounded-lg">
                        <FileSearchIcon size={20}/>
                    </div>
                    <h3 class="font-bold uppercase tracking-tight text-sm">Set Game Executable</h3>
                </div>
                <button onclick={() => launchPicker = false} class="text-muted hover:text-foreground">
                    <XIcon size={20}/>
                </button>
            </div>

            <p class="text-xs text-muted leading-relaxed">
                If the game installed to a custom location or requires a specific .exe to run, enter the relative or
                absolute path below.
            </p>

            <div class="space-y-2">
                <label for="exe-path" class="text-[10px] font-bold uppercase text-muted px-1">Executable Path</label>
                <div class="relative group">
                    <input
                            id="exe-path"
                            type="text"
                            bind:value={customExePath}
                            placeholder="e.g. Binaries/Win64/Game.exe"
                            class="w-full bg-panel border border-border rounded-xl pl-4 pr-12 py-2 text-sm focus:outline-none focus:border-frost-500 transition-colors"
                    />
                    <!-- Picker Button -->
                    <button
                            onclick={triggerFilePicker}
                            title="Browse files"
                            class="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 rounded-lg hover:bg-white/5 text-muted hover:text-frost-400 transition-colors"
                    >
                        <FolderOpenIcon size={18}/>
                    </button>
                </div>
            </div>

            <div class="flex justify-end gap-3 mt-2">
                <button onclick={() => launchPicker = false}
                        class="px-4 py-2 text-xs font-bold text-muted hover:text-foreground transition-all">
                    Cancel
                </button>
                <button onclick={pickFinalExe}
                        class="px-6 py-2 bg-frost-500 text-background rounded-xl text-xs font-bold hover:bg-frost-400 transition-all shadow-lg shadow-frost-500/20">
                    Save Changes
                </button>
            </div>
        </div>
    </div>
{/if}

<!-- Error Dialog -->
{#if errorMessage}
    <div class="fixed inset-0 z-150 flex items-center justify-center p-6" transition:fade={{ duration: 150 }}>
        <div class="absolute inset-0 bg-background/80 backdrop-blur-md" onclick={() => errorMessage = null}></div>
        <div class="relative w-full max-w-sm bg-surface border border-red-500/30 rounded-2xl shadow-2xl p-6 flex flex-col gap-4"
             transition:fly={{ y: 20, duration: 300 }}>
            <div class="flex items-center justify-between">
                <div class="flex items-center gap-3 text-red-400">
                    <div class="p-2 bg-red-500/10 rounded-lg">
                        <TriangleAlert size={20}/>
                    </div>
                    <h3 class="font-bold">Operation Failed</h3>
                </div>
                <button onclick={() => errorMessage = null} class="text-muted hover:text-foreground">
                    <XIcon size={20}/>
                </button>
            </div>
            <p class="text-sm text-muted leading-relaxed">{errorMessage}</p>
            <div class="flex justify-end mt-2">
                <button onclick={() => errorMessage = null}
                        class="px-6 py-2 bg-panel border border-border rounded-xl text-sm font-bold hover:text-foreground transition-all">
                    Dismiss
                </button>
            </div>
        </div>
    </div>
{/if}

<style>
    /* Prevent text selection and ensure button stays together */
    .inline-flex {
        white-space: nowrap;
    }
</style>
