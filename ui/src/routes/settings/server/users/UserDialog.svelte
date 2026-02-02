<script lang="ts">
    import {
        XIcon, UserIcon, LockIcon, ShieldIcon,
        SaveIcon, UserPlusIcon, UserCogIcon, ChevronDownIcon
    } from "@lucide/svelte";
    import {fade, fly} from 'svelte/transition';
    import {type User, UserService} from "$lib/gen/user/v1/user_pb";
    import {glacierCli} from "$lib/api/api";
    import {createRPCRunner} from "$lib/api/svelte-api.svelte.ts";
    import {callRPC} from "$lib/api/api.ts";
    import {getSnackbarCtx} from "$lib/components/snackbar/snackbar-provider.svelte";

    interface Props {
        open: boolean;
        editUser?: User;
        onClose?: () => void;
    }

    let {onClose, editUser, open = $bindable()}: Props = $props();

    let username = $state("");
    let password = $state("");
    let role = $state("User");


    $effect(() => {
        if (open) {
            rolesRpc.runner()
            username = "";
            password = "";
            role = "User";
        }

        if (editUser) {
            username = editUser.Username;
            password = "";
            role = editUser.Role || "User";
        }
    });

    function handleClose() {
        open = false;
        if (onClose) onClose();
    }

    const userSrv = glacierCli(UserService);

    let rolesRpc = createRPCRunner(() => userSrv.listRoles({}))
    const snackbarManager = getSnackbarCtx();

    async function handleSave() {
        const user = {
            id: BigInt(editUser?.id ?? 0),
            Username: username,
            Password: password,
            Role: role,
        }

        if (editUser) {
            const {err} = await callRPC(() => userSrv.edit({user: user}))
            if (err) {
                snackbarManager.push(`Error creating ${err}`, 'error')
            }
        } else {
            const {err} = await callRPC(() => userSrv.new({user: user}))
            if (err) {
                snackbarManager.push(`Error editing ${err}`, 'error')
            }
        }

        handleClose();
    }
</script>

{#if open}
    <div class="fixed inset-0 z-100 flex items-center justify-center p-6" transition:fade={{ duration: 150 }}>
        <!-- Backdrop -->
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div class="absolute inset-0 bg-background/80 backdrop-blur-md" onclick={handleClose}></div>

        <!-- Dialog Content -->
        <div
                class="relative w-full max-w-md bg-surface border border-border rounded-4xl shadow-2xl overflow-hidden flex flex-col"
                transition:fly={{ y: 20, duration: 300 }}
        >
            <!-- Header -->
            <div class="flex items-center justify-between p-6 border-b border-border bg-panel/30">
                <div class="flex items-center gap-3">
                    <div class="p-2.5 bg-frost-500/10 text-frost-400 rounded-xl">
                        {#if editUser}
                            <UserCogIcon size={20}/>
                        {:else}
                            <UserPlusIcon size={20}/>
                        {/if}
                    </div>
                    <h2 class="text-xl font-bold text-foreground">
                        {editUser ? 'Edit User' : 'Create User'}
                    </h2>
                </div>
                <button onclick={handleClose} class="text-muted hover:text-foreground transition-colors p-1">
                    <XIcon size={20}/>
                </button>
            </div>

            <!-- Form Body -->
            <div class="p-8 space-y-6">
                <!-- Username -->
                <div class="space-y-2">
                    <label class="text-[10px] font-bold text-muted uppercase tracking-widest ml-1">Username</label>
                    <div class="relative">
                        <UserIcon size={18} class="absolute left-4 top-1/2 -translate-y-1/2 text-muted/50"/>
                        <input
                                type="text"
                                bind:value={username}
                                placeholder="johndoe"
                                class="w-full bg-panel border border-border rounded-2xl py-3.5 pl-12 pr-4 outline-none focus:border-frost-500 transition-all text-sm"
                        />
                    </div>
                </div>

                <!-- Password -->
                <div class="space-y-2">
                    <label class="text-[10px] font-bold text-muted uppercase tracking-widest ml-1">
                        {'Password'}
                    </label>
                    <div class="relative">
                        <LockIcon size={18} class="absolute left-4 top-1/2 -translate-y-1/2 text-muted/50"/>
                        <input
                                type="password"
                                bind:value={password}
                                placeholder="••••••••"
                                class="w-full bg-panel border border-border rounded-2xl py-3.5 pl-12 pr-4 outline-none focus:border-frost-500 transition-all text-sm"
                        />
                    </div>
                    {#if editUser}
                        <p class="text-[10px] text-muted/50 px-1 italic">Leave blank to keep existing password</p>
                    {/if}
                </div>

                <!-- Role -->
                <div class="space-y-2">
                    <label class="text-[10px] font-bold text-muted uppercase tracking-widest ml-1">System Role</label>
                    <div class="relative">
                        <ShieldIcon size={18} class="absolute left-4 top-1/2 -translate-y-1/2 text-muted/50"/>
                        <select
                                bind:value={role}
                                class="w-full bg-panel border border-border rounded-2xl py-3.5 pl-12 pr-10 outline-none focus:border-frost-500 transition-all text-sm appearance-none cursor-pointer font-bold"
                        >
                            {#each rolesRpc.value?.roles as r}
                                <option value={r.Name}>{r.Name}</option>
                            {/each}
                        </select>
                        <div class="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-muted/50">
                            <ChevronDownIcon size={16}/>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Footer -->
            <div class="p-6 border-t border-border bg-panel/30 flex justify-end gap-3">
                <button
                        onclick={handleClose}
                        class="px-6 py-2.5 rounded-xl border border-border hover:bg-panel transition-all text-sm font-bold text-muted hover:text-foreground"
                >
                    Cancel
                </button>
                <button
                        onclick={handleSave}
                        class="px-8 py-2.5 rounded-xl bg-frost-500 text-background font-bold text-sm hover:bg-frost-400 transition-all active:scale-95 shadow-lg shadow-frost-500/20 flex items-center gap-2"
                >
                    <SaveIcon size={16}/>
                    {editUser ? 'Save Changes' : 'Create User'}
                </button>
            </div>
        </div>
    </div>
{/if}

<style>
    select {
        background-image: none;
    }
</style>
