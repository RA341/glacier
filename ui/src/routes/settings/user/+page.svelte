<script lang="ts">
    import {
        Edit3Icon,
        FingerprintPattern,
        KeyIcon,
        LoaderIcon,
        LockIcon,
        SaveIcon,
        ShieldCheckIcon,
        UserIcon
    } from "@lucide/svelte";
    import {callRPC, glacierCli} from "$lib/api/api";
    import {UserService} from "$lib/gen/user/v1/user_pb";
    import {createRPCRunner} from "$lib/api/svelte-api.svelte";
    import {onMount} from "svelte";
    import {fade} from "svelte/transition";
    import {getSnackbarCtx} from "$lib/components/snackbar/snackbar-provider.svelte";

    const userSrv = glacierCli(UserService);

    let selfRpc = createRPCRunner(() => userSrv.self({}));

    let newUsername = $state("");
    let newPassword = $state("");
    let confirmPassword = $state("");

    let isUpdatingUsername = $state(false);
    let isUpdatingPassword = $state(false);

    $effect(() => {
        if (selfRpc.value?.user && !newUsername) {
            newUsername = selfRpc.value.user.Username;
            console.log($inspect(selfRpc.value.user));
        }
    });

    async function loadUser() {
        await selfRpc.runner();
    }

    const sm = getSnackbarCtx()

    async function upUser(username?: string, pass?: string) {
        if (!selfRpc.value?.user?.id) {
            sm.push("Cannot update, empty user info", 'warn')
            return
        }

        const {err} = await callRPC(() => userSrv.edit({
            user: {
                id: selfRpc.value!.user!.id,
                Role: selfRpc.value!.user!.Role,
                Username: username ?? "",
                Password: pass ?? "",
            }
        }))

        if (err) {
            sm.push(`error occurred while updating: ${err}`, 'error')
        }
    }

    onMount(async () => {
        await loadUser()
    });

    async function handleUpdateUsername() {
        if (!newUsername || newUsername === selfRpc.value?.user?.Username) return;

        isUpdatingUsername = true;
        await upUser(newUsername);
        isUpdatingUsername = false;
        sm.push("Username updated successfully", "success");

        await loadUser();
    }

    async function handleUpdatePassword() {
        if (!newPassword || newPassword !== confirmPassword) {
            sm.push("Passwords do not match", "warn");
            return;
        }

        isUpdatingPassword = true;
        await upUser("", newPassword);
        isUpdatingPassword = false;

        newPassword = ""
        confirmPassword = ""
        sm.push("Password updated successfully", "success");

        await loadUser();
    }
</script>

<div class="p-10 mx-auto space-y-8" in:fade>
    <!-- Header -->
    <header class="flex items-center gap-4 px-2">
        <div class="p-3 bg-panel border border-border rounded-2xl text-frost-400">
            <FingerprintPattern size={28}/>
        </div>
        <div>
            <h1 class="text-3xl font-bold tracking-tight text-foreground">My Profile</h1>
            <p class="text-sm text-muted">Manage your account</p>
        </div>
    </header>

    {#if selfRpc.loading && !selfRpc.value}
        <div class="flex flex-col items-center justify-center h-64 text-muted gap-4">
            <LoaderIcon class="animate-spin text-frost-500" size={32}/>
            <p class="animate-pulse text-sm">Retrieving account details...</p>
        </div>
    {:else if selfRpc.value}
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">

            <!-- LEFT: Profile Preview -->
            <div class="md:col-span-1">
                <div class="bg-surface border border-border rounded-4xl p-8 flex flex-col items-center text-center gap-4 shadow-xl sticky top-8">
                    <div class="w-24 h-24 bg-panel border-4 border-surface rounded-3xl flex items-center justify-center text-frost-400 shadow-2xl">
                        {#if selfRpc.value.user?.Role?.toLowerCase() === 'admin'}
                            <ShieldCheckIcon size={48} strokeWidth={1.5}/>
                        {:else}
                            <UserIcon size={48} strokeWidth={1.5}/>
                        {/if}
                    </div>

                    <div>
                        <h2 class="text-2xl font-black text-foreground truncate max-w-full px-2">
                            {selfRpc.value.user?.Username}
                        </h2>
                        <span class="inline-block mt-1 px-3 py-1 rounded-full bg-frost-500/10 border border-frost-500/20 text-[10px] font-bold text-frost-400 uppercase tracking-widest">
                            {selfRpc.value.user?.Role || 'Unknown'}
                        </span>
                    </div>
                </div>
            </div>

            <!-- RIGHT: Settings Forms -->
            <div class="md:col-span-2 space-y-6">

                <!-- Section 1: General Settings (Username) -->
                <div class="bg-surface border border-border rounded-4xl p-8 shadow-xl space-y-6">
                    <div class="flex items-center gap-3">
                        <div class="p-2 bg-panel rounded-xl text-muted">
                            <Edit3Icon size={20}/>
                        </div>
                        <h3 class="text-lg font-bold">General Settings</h3>
                    </div>

                    <div class="space-y-4">
                        <div class="space-y-2">
                            <label class="text-[10px] font-bold text-muted uppercase tracking-widest ml-1">Display
                                Username</label>
                            <div class="relative">
                                <UserIcon size={18} class="absolute left-4 top-1/2 -translate-y-1/2 text-muted/40"/>
                                <input
                                        type="text"
                                        bind:value={newUsername}
                                        class="w-full bg-panel border border-border rounded-2xl py-4 pl-12 pr-4 outline-none focus:border-frost-500 transition-all text-sm font-bold"
                                />
                            </div>
                        </div>

                        <div class="flex justify-end">
                            <button
                                    onclick={handleUpdateUsername}
                                    disabled={isUpdatingUsername || newUsername === selfRpc.value?.user?.Username}
                                    class="px-8 py-3 bg-panel border border-border text-foreground rounded-xl text-xs font-bold hover:border-frost-500 transition-all flex items-center gap-2 disabled:opacity-30"
                            >
                                {#if isUpdatingUsername}
                                    <LoaderIcon size={14} class="animate-spin"/>
                                {:else}
                                    <SaveIcon size={14}/>
                                {/if}
                                Save Username
                            </button>
                        </div>
                    </div>
                </div>

                <!-- Section 2: Security Settings (Password) -->
                <div class="bg-surface border border-border rounded-4xl p-8 shadow-xl space-y-6">
                    <div class="flex items-center gap-3">
                        <div class="p-2 bg-panel rounded-xl text-muted">
                            <LockIcon size={20}/>
                        </div>
                        <h3 class="text-lg font-bold">Security</h3>
                    </div>

                    <div class="space-y-5">
                        <div class="space-y-2">
                            <label class="text-[10px] font-bold text-muted uppercase tracking-widest ml-1">New
                                Password</label>
                            <div class="relative">
                                <KeyIcon size={18} class="absolute left-4 top-1/2 -translate-y-1/2 text-muted/40"/>
                                <input
                                        type="password"
                                        bind:value={newPassword}
                                        placeholder="••••••••"
                                        class="w-full bg-panel border border-border rounded-2xl py-4 pl-12 pr-4 outline-none focus:border-frost-500 transition-all text-sm"
                                />
                            </div>
                        </div>

                        <div class="space-y-2">
                            <label class="text-[10px] font-bold text-muted uppercase tracking-widest ml-1">Confirm
                                Password</label>
                            <div class="relative">
                                <LockIcon size={18} class="absolute left-4 top-1/2 -translate-y-1/2 text-muted/40"/>
                                <input
                                        type="password"
                                        bind:value={confirmPassword}
                                        placeholder="••••••••"
                                        class="w-full bg-panel border border-border rounded-2xl py-4 pl-12 pr-4 outline-none focus:border-frost-500 transition-all text-sm"
                                />
                            </div>
                        </div>

                        <button
                                onclick={handleUpdatePassword}
                                disabled={isUpdatingPassword || !newPassword}
                                class="w-full py-4 bg-frost-500 text-background font-bold rounded-2xl hover:bg-frost-400 transition-all flex items-center justify-center gap-2 shadow-lg shadow-frost-500/20 disabled:opacity-30"
                        >
                            {#if isUpdatingPassword}
                                <LoaderIcon size={20} class="animate-spin"/>
                            {:else}
                                <SaveIcon size={18}/>
                                Update Password
                            {/if}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    {/if}
</div>