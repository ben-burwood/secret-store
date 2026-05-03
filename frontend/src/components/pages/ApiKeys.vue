<template>
    <div class="flex flex-row justify-between">
        <h2 class="text-xl font-semibold">API Keys</h2>
        <button class="btn btn-primary" @click="showAddRow = true" :disabled="showAddRow">+ Add API Key</button>
    </div>

    <div class="overflow-x-auto rounded-box border border-base-content/5 bg-base-200 my-2">
        <table class="table">
            <thead>
                <tr>
                    <th>Name</th>
                    <th>Key</th>
                    <th>Created</th>
                    <th>Scopes</th>
                    <th></th>
                </tr>
            </thead>
            <tbody>
                <AddApiKeyRow v-if="showAddRow" :secrets="secrets" @refresh="refresh" @cancel="showAddRow = false" />
                <ApiKeyRow v-for="apiKey in apiKeys" :key="apiKey.id" :apiKey="apiKey" :secrets="secrets" @refresh="refresh" />
                <tr v-if="!loading && apiKeys.length === 0 && !showAddRow">
                    <td colspan="5" class="text-center opacity-60">No API keys yet</td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { toast } from "vue3-toastify";
import { backendFetch } from "@/main";
import ApiKeyRow from "@/components/api-keys-table/ApiKeyRow.vue";
import AddApiKeyRow from "@/components/api-keys-table/AddApiKeyRow.vue";

type ApiKey = { id: number; name: string; key: string; created_at: string | null; secret_ids: number[] };
type Secret = { id: number; key: string };

const apiKeys = ref<ApiKey[]>([]);
const secrets = ref<Secret[]>([]);
const loading = ref(true);
const showAddRow = ref(false);

async function fetchApiKeys() {
    const res = await backendFetch(`/api/keys`);
    if (!res.ok) return;
    const data = await res.json();
    apiKeys.value = data.api_keys ?? [];
}

async function fetchSecrets() {
    const res = await backendFetch(`/secrets`);
    if (!res.ok) return;
    const data = await res.json();
    secrets.value = data.secrets ?? [];
}

async function refresh() {
    loading.value = true;
    try {
        await fetchApiKeys();
    } catch (error) {
        console.error("Error fetching API keys:", error);
        toast("Failed to fetch API keys", { type: "error" });
    } finally {
        loading.value = false;
    }
}

onMounted(async () => {
    await Promise.all([fetchApiKeys(), fetchSecrets()]);
    loading.value = false;
});
</script>
