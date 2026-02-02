<script lang="ts">
    import {
        CircleAlert,
        EllipsisVertical,
        LoaderIcon,
        PlusIcon,
        RefreshCcw,
        SearchIcon,
        ShieldCheckIcon,
        Trash2Icon,
        UserCogIcon,
        UserIcon,
        UserMinusIcon
    } from "@lucide/svelte";
    import {callRPC, glacierCli} from "$lib/api/api";
    import {type User, UserService} from "$lib/gen/user/v1/user_pb";
    import {createRPCRunner} from "$lib/api/svelte-api.svelte";
    import {onMount} from "svelte";
    import {fade} from "svelte/transition";
    import UserDialog from "./UserDialog.svelte";
    import {getSnackbarCtx} from "$lib/components/snackbar/snackbar-provider.svelte";

    const userSrv = glacierCli(UserService);
    let query = $state("");

    let listRpc = createRPCRunner(() => userSrv.list({query: query}));

    function refreshList() {
        listRpc.runner();
    }

    onMount(() => {
        refreshList()
    });

    function getRoleTheme(role: string = "") {
        const r = role.toLowerCase();
        if (r === 'admin') return 'text-frost-400 bg-frost-400/10 border-frost-400/20';
        return 'text-muted bg-panel border-border';
    }

    let isOpen = $state(false)

    let user = $state<User | undefined>()

    function editUser(input: User) {
        isOpen = true
        user = input
    }

    let sm = getSnackbarCtx()

    async function deleteUser(input: User) {
        const {err} = await callRPC(() => userSrv.delete({id: input.id}))
        sm.push(`error deleting user: ${err}`, 'error')
        refreshList()
    }

    function onClose() {
        isOpen = false;
        refreshList()
    }
</script>

<div class="space-y-6 max-w-1l mx-auto">
    <!-- HEADER & SEARCH -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 px-2">
        <div class="flex items-center gap-3">
            <button
                    onclick={() => {
                        isOpen = true
                        user = undefined
                    }}
                    class="flex items-center gap-2 px-6 py-2.5 bg-frost-500 text-background rounded-xl text-sm font-bold hover:bg-frost-400 transition-all shadow-lg shadow-frost-500/20 active:scale-95"
            >
                <PlusIcon size={18}/>
                Add
            </button>

            <div class="relative group">
                <SearchIcon size={16}
                            class="absolute left-4 top-1/2 -translate-y-1/2 text-muted group-focus-within:text-frost-400 transition-colors"/>
                <input
                        type="text"
                        bind:value={query}
                        onkeydown={(e) => e.key === 'Enter' && refreshList()}
                        placeholder="Filter by username..."
                        class="bg-panel border border-border rounded-xl py-2 pl-11 pr-4 outline-none focus:border-frost-500 transition-all text-sm w-full sm:w-64"
                />
            </div>
            <button
                    onclick={refreshList}
                    class="p-2 bg-panel border border-border rounded-xl text-muted hover:text-frost-400 transition-all"
            >
                <RefreshCcw size={18} class={listRpc.loading ? 'animate-spin' : ''}/>
            </button>
        </div>
    </header>

    <!-- LIST AREA -->
    <main class="min-h-100">
        {#if listRpc.loading && !listRpc.value}
            <div class="flex flex-col items-center justify-center h-64 text-muted gap-4">
                <LoaderIcon class="animate-spin text-frost-500" size={32}/>
                <p class="animate-pulse text-sm font-medium">Fetching users...</p>
            </div>
        {:else if listRpc.error}
            <div class="flex flex-col items-center justify-center h-64 text-red-400 gap-3 bg-red-500/5 border border-red-500/10 rounded-3xl">
                <CircleAlert size={32}/>
                <p class="text-sm font-medium">{listRpc.error}</p>
            </div>
        {:else if !listRpc.value?.users || listRpc.value.users.length === 0}
            <div class="flex flex-col items-center justify-center h-64 border-2 border-dashed border-border rounded-3xl text-muted/30">
                <UserMinusIcon size={48} strokeWidth={1} class="mb-2"/>
                <p class="text-sm font-medium">No users found</p>
            </div>
        {:else}
            <!-- THE LIST STACK -->
            <div class="flex flex-col gap-2">
                {#each listRpc.value.users as user (user.Username)}
                    <div role="button"
                         tabindex="0"
                         onclick={() => editUser(user)}
                         onkeydown={(e) => {
                                if (e.key === 'Enter' || e.key === ' ') {
                                    e.preventDefault();
                                    editUser(user);
                                }
                            }}
                         class="group flex items-center justify-between p-3 pl-4 bg-panel/20 border border-border rounded-xl hover:bg-panel/40 hover:border-frost-500/30 transition-all cursor-pointer"
                         in:fade={{ duration: 200 }}
                    >
                        <div class="flex items-center gap-4 min-w-0">
                            <!-- Role-based Icon -->
                            <div class="w-10 h-10 shrink-0 bg-surface border border-border rounded-lg flex items-center justify-center text-muted group-hover:text-frost-400 transition-colors">
                                {#if user.Role?.toLowerCase() === 'admin'}
                                    <ShieldCheckIcon size={20}/>
                                {:else}
                                    <UserIcon size={20}/>
                                {/if}
                            </div>

                            <div class="flex flex-col min-w-0">
                                <span class="font-bold text-foreground text-sm truncate">{user.Username}</span>
                            </div>
                        </div>

                        <div class="flex items-center gap-6">
                            <!-- Role Badge -->
                            <span class="hidden sm:inline-block text-[10px] font-bold px-2 py-0.5 rounded border uppercase tracking-widest {getRoleTheme(user.Role)}">
                                {user.Role || 'User'}
                            </span>

                            <!-- Actions -->
                            <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-all translate-x-1 group-hover:translate-x-0 pr-1">
                                <button
                                        onclick={(e) => {
                                            e.preventDefault();
                                            e.stopPropagation();
                                            editUser(user)
                                        }}
                                        class="p-2 text-muted hover:text-frost-400 hover:bg-frost-500/10 rounded-lg transition-all"
                                        title="Edit Permissions">
                                    <UserCogIcon size={18}/>
                                </button>
                                <button
                                        onclick={(e) => {
                                            e.preventDefault();
                                            e.stopPropagation();
                                            deleteUser(user)
                                        }}
                                        class="p-2 text-muted hover:text-red-400 hover:bg-red-500/10 rounded-lg transition-all"
                                        title="Remove User">
                                    <Trash2Icon size={18}/>
                                </button>
                            </div>

                            <!-- Small indicator for mobile or non-hover -->
                            <div class="group-hover:hidden text-muted/20">
                                <EllipsisVertical size={18}/>
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    </main>
</div>

<UserDialog open={isOpen} editUser={user} onClose={onClose}/>

<style>
    /* Prevent text selection on names */
    span {
        user-select: none;
    }
</style>
