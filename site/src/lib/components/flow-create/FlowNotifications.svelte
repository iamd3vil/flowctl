<script lang="ts">
    import MultiReceiverSelector from "$lib/components/shared/MultiReceiverSelector.svelte";

    let {
        notifications = $bindable(),
        addNotification,
        availableMessengers,
        messengerConfigs,
        disabled = false,
    }: {
        notifications: any[];
        addNotification: () => void;
        availableMessengers: string[];
        messengerConfigs: Record<string, any>;
        disabled?: boolean;
    } = $props();

    function removeNotification(index: number) {
        notifications.splice(index, 1);
    }

    const eventOptions = [
        { value: "on_success", label: "On Success" },
        { value: "on_failure", label: "On Failure" },
        { value: "on_waiting", label: "On Waiting" },
        { value: "on_cancelled", label: "On Cancelled" },
    ];

    function onChannelChange(notification: any) {
        // Initialize config keys from the schema so bindings work immediately
        const schema = messengerConfigs[notification.channel];
        if (!schema?.properties) {
            notification.config = {};
            return;
        }
        const config: Record<string, any> = {};
        for (const [key, property] of Object.entries(schema.properties) as [string, any][]) {
            // Preserve existing value if present, otherwise set a typed default
            if (notification.config && notification.config[key] !== undefined) {
                config[key] = notification.config[key];
            } else if (property.type === "array" || property.widget === "userselector") {
                config[key] = [];
            } else {
                config[key] = "";
            }
        }
        notification.config = config;
    }

    function toggleEvent(notification: any, eventValue: string) {
        if (!notification.events) {
            notification.events = [];
        }
        const index = notification.events.indexOf(eventValue);
        if (index > -1) {
            notification.events.splice(index, 1);
        } else {
            notification.events.push(eventValue);
        }
    }
</script>

<!-- Flow Notifications Section -->
<div>
    <div class="flex items-center justify-between mb-6">
        <div>
            <h3 class="text-base font-medium text-foreground">
                Flow Notifications
            </h3>
            <p class="mt-1 text-sm text-muted-foreground">
                Configure notifications for flow execution events
            </p>
        </div>
        {#if !disabled}
            <button
                onclick={addNotification}
                class="px-4 py-2 text-sm font-medium bg-primary-500 text-white rounded-md hover:bg-primary-600 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 cursor-pointer"
            >
                + Add Notification
            </button>
        {/if}
    </div>

    <div class="space-y-4">
        {#each notifications as notification, index (index)}
            <div class="border border-border rounded-lg p-4 relative">
                {#if !disabled}
                    <button
                        onclick={() => removeNotification(index)}
                        class="absolute top-4 right-4 text-muted-foreground hover:text-danger-600 cursor-pointer"
                    >
                        <svg
                            class="w-5 h-5"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M6 18L18 6M6 6l12 12"
                            />
                        </svg>
                    </button>
                {/if}

                <div class="space-y-4">
                    <!-- Channel -->
                    <div>
                        <label
                            class="block text-sm font-medium text-foreground mb-1"
                            >Channel *</label
                        >
                        <select
                            bind:value={notification.channel}
                            onchange={() => onChannelChange(notification)}
                            required
                            class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm"
                        >
                            <option value="">
                                {availableMessengers.length === 0
                                    ? "No messengers available"
                                    : "Select channel"}
                            </option>
                            {#each availableMessengers as messenger}
                                <option value={messenger}>
                                    {messenger.charAt(0).toUpperCase() +
                                        messenger.slice(1)}
                                </option>
                            {/each}
                        </select>
                    </div>

                    <!-- Events -->
                    <div>
                        <label
                            class="block text-sm font-medium text-foreground mb-2"
                            >Events *</label
                        >
                        <div class="space-y-2">
                            {#each eventOptions as event}
                                <label class="flex items-center cursor-pointer">
                                    <input
                                        type="checkbox"
                                        checked={notification.events?.includes(
                                            event.value,
                                        )}
                                        onchange={() =>
                                            toggleEvent(
                                                notification,
                                                event.value,
                                            )}
                                        class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-input rounded"
                                    />
                                    <span class="ml-2 text-sm text-foreground"
                                        >{event.label}</span
                                    >
                                </label>
                            {/each}
                        </div>
                    </div>

                    <!-- Dynamic messenger config fields -->
                    {#if notification.channel && messengerConfigs[notification.channel]?.properties}
                        {@const schema =
                            messengerConfigs[notification.channel]}
                        {#each Object.entries(schema.properties) as [key, property]}
                            {@const isRequired =
                                schema.required?.includes(key)}
                            {@const label = property.title || key}
                            {@const description = property.description || ""}

                            <div>
                                {#if property.widget === "userselector"}
                                    <label
                                        class="block text-sm font-medium text-foreground mb-2"
                                    >
                                        {label}
                                        {#if isRequired}<span
                                                class="text-red-500"
                                                >*</span
                                            >{/if}
                                    </label>
                                    <MultiReceiverSelector
                                        bind:selectedReceivers={notification
                                            .config[key]}
                                    />
                                {:else}
                                    <label
                                        class="block text-sm font-medium text-foreground mb-1"
                                    >
                                        {label}
                                        {#if isRequired}<span
                                                class="text-red-500"
                                                >*</span
                                            >{/if}
                                    </label>
                                    <input
                                        type="text"
                                        bind:value={notification.config[key]}
                                        placeholder={property.placeholder ||
                                            ""}
                                        class="w-full px-3 py-2 text-foreground bg-card border border-input rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm"
                                    />
                                {/if}
                                {#if description}
                                    <p
                                        class="mt-1 text-xs text-muted-foreground"
                                    >
                                        {description}
                                    </p>
                                {/if}
                            </div>
                        {/each}
                    {/if}
                </div>
            </div>
        {/each}

        {#if notifications.length === 0}
            <div class="text-center py-8 text-muted-foreground">
                <svg
                    class="mx-auto h-12 w-12 text-muted-foreground mb-3"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
                    />
                </svg>
                <p>No notifications configured yet</p>
                {#if !disabled}
                    <button
                        onclick={addNotification}
                        class="mt-2 text-sm text-primary-600 hover:text-primary-700 font-medium cursor-pointer"
                    >
                        Add your first notification
                    </button>
                {/if}
            </div>
        {/if}
    </div>
</div>
