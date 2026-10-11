<script lang="ts">
	import { EyeIcon, EyeOffIcon, InfoIcon, LockIcon } from '@lucide/svelte';
	import CopyButton from '#lib/components/CopyButton.svelte';
	import Self from './ConfigField.svelte';
	import { splitCamelCase } from '#lib/api/strings.js';
	import FolderPicker from '#lib/components/FolderPicker.svelte';

	let { schema, level = 0 } = $props();

	let visibleSecrets = $state<Record<string, boolean>>({});

	function toggleSecret(key: string) {
		visibleSecrets[key] = !visibleSecrets[key];
	}
</script>

<div
	class="space-y-6"
	class:ml-6={level > 0}
	class:border-l={level > 0}
	class:border-border={level > 0}
	class:pl-6={level > 0}
>
	{#each Object.entries(schema) as [key, item]}
		<div class="flex flex-col gap-2">
			{#if item.IsStruct}
				<!-- SECTION HEADER -->
				<div class="pt-4 first:pt-0">
					<h3
						class="flex items-center gap-2 text-sm font-black tracking-widest text-frost-400 uppercase"
					>
						{key}
						{#if item.Help}
							<div class="group relative">
								<InfoIcon size={12} class="cursor-help text-muted/40" />
								<div
									class="pointer-events-none absolute top-0 left-full z-50 ml-2 w-48 rounded border border-border bg-panel p-2 text-[10px] font-medium tracking-normal text-white normal-case opacity-0 transition-opacity group-hover:opacity-100"
								>
									{item.Help}
								</div>
							</div>
						{/if}
					</h3>
					<div class="mt-2 h-px w-full bg-border opacity-50"></div>
				</div>

				{#if item.Nested}
					<Self schema={item.Nested} level={level + 1} />
				{/if}
			{:else}
				<!-- LEAF INPUTS -->
				<div
					class="group flex flex-col gap-1.5 rounded-xl p-2 transition-colors hover:bg-panel/20"
					class:opacity-60={item.EnvSet}
				>
					<div class="flex items-center justify-between">
						<div class="flex items-center gap-2">
							<label class="text-[11px] font-bold tracking-wider text-muted uppercase"
								>{splitCamelCase(key)}</label
							>

							{#if item.EnvSet}
								<div class="group/env relative">
									<LockIcon size={12} class="text-amber-500/60" />
									<!-- ENV TOOLTIP -->
									<div
										class="pointer-events-none absolute top-1/2 left-full z-[60] ml-2 w-64 -translate-y-1/2 rounded-lg border border-amber-500/30 bg-surface p-3 text-[10px] font-medium tracking-normal text-muted normal-case opacity-0 shadow-xl transition-opacity group-hover/env:opacity-100"
									>
										<p class="mb-1 font-bold tracking-tighter text-amber-400 uppercase">
											Locked by Environment
										</p>
										This value is set via<code class="rounded bg-panel px-1 text-foreground"
											>{item.Env}</code
										>. Changes must be made via the env var or by removing the env then editing.
									</div>
								</div>
							{/if}
						</div>
						{#if item.Default}
							{#if item.Default}
								<CopyButton label="Default" value={item.Default} />
							{/if}
						{/if}

						{#if item.Env}
							{#if item.Default}
								<CopyButton label="Env" value={item.Env} />
							{/if}
						{/if}
					</div>

					<div class="relative flex items-center gap-2">
						{#if item.IsFolder}
							<FolderPicker disabled={item.EnvSet} bind:value={item.Value} />
						{:else if item.FieldType === 'bool'}
							<label
								class="relative inline-flex items-center"
								class:cursor-not-allowed={item.EnvSet}
								class:cursor-pointer={!item.EnvSet}
							>
								<input
									type="checkbox"
									bind:checked={item.Value}
									disabled={item.EnvSet}
									class="peer sr-only"
								/>
								<div
									class="peer h-6 w-11 rounded-full border border-border bg-panel opacity-100 peer-checked:bg-frost-500 peer-focus:outline-none peer-disabled:opacity-40 after:absolute after:start-[2px] after:top-[2px] after:h-5 after:w-5 after:rounded-full after:border after:border-border after:bg-muted after:transition-all after:content-[''] peer-checked:after:translate-x-full peer-checked:after:border-white peer-checked:after:bg-background rtl:peer-checked:after:-translate-x-full"
								></div>
							</label>
						{:else if item.FieldType === 'int'}
							<input
								type="number"
								bind:value={item.Value}
								disabled={item.EnvSet}
								class="w-full rounded-lg border border-border bg-panel px-3 py-2 text-sm transition-colors outline-none focus:border-frost-500 disabled:cursor-not-allowed"
							/>
						{:else if item.IsSecret}
							<div class="relative flex-1">
								<input
									type={visibleSecrets[key] ? 'text' : 'password'}
									bind:value={item.Value}
									disabled={item.EnvSet}
									class="w-full rounded-lg border border-border bg-panel py-2 pr-10 pl-3 text-sm transition-colors outline-none focus:border-frost-500 disabled:cursor-not-allowed"
								/>
								<button
									type="button"
									onclick={() => toggleSecret(key)}
									class="absolute top-1/2 right-2 -translate-y-1/2 text-muted hover:text-frost-400"
								>
									{#if visibleSecrets[key]}
										<EyeOffIcon size={16} />
									{:else}
										<EyeIcon size={16} />
									{/if}
								</button>
							</div>
						{:else if Array.isArray(item.Value)}
							<input
								type="text"
								value={item.Value.join(', ')}
								disabled={item.EnvSet}
								oninput={(e) =>
									(item.Value = e.currentTarget.value.split(',').map((s) => s.trim()))}
								class="w-full rounded-lg border border-border bg-panel px-3 py-2 font-mono text-sm transition-colors outline-none focus:border-frost-500 disabled:cursor-not-allowed"
							/>
						{:else}
							<input
								type="text"
								bind:value={item.Value}
								disabled={item.EnvSet}
								class="w-full rounded-lg border border-border bg-panel px-3 py-2 text-sm transition-colors outline-none focus:border-frost-500 disabled:cursor-not-allowed"
							/>
						{/if}
					</div>

					{#if item.Help}
						<p class="text-[10px] italic">{item.Help}</p>
					{/if}
				</div>
			{/if}
		</div>
	{/each}
</div>
