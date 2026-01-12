<template>
    <div class="flex flex-row items-center gap-2 justify-end">
        Sort by:
        <select v-model="sortColumn" class="select">
            <option value="id">ID</option>
            <option value="key">Key</option>
        </select>
        <button class="btn btn-ghost" @click="sortAscending = !sortAscending" :aria-label="sortAscending ? 'Ascending' : 'Descending'">
            <ArrowDownWideNarrow v-if="sortAscending" :size="18" />
            <ArrowUpNarrowWide v-else :size="18" />
        </button>
    </div>

    <button class="btn" :disabled="showSecret" @click="clickShowSecret">Show Secrets</button>
    <progress class="progress progress-primary w-56" :value="showSecretPercentRemaining" max="100" />

    <div class="overflow-x-auto rounded-box border border-base-content/5 bg-base-200 mt-4">
        <table class="table">
            <thead>
                <tr>
                    <th></th>
                    <th>Key</th>
                    <th>Value</th>
                    <th></th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="secret in tableSecrets" :key="secret.id">
                    <th>{{ secret.id }}</th>
                    <td>{{ secret.key }}</td>
                    <td><SecretDisplay :secret="secret.value" :showSecret /></td>
                    <td>
                        <button class="btn btn-outline btn-error" @click="deleteSecret(secret.id)">
                            <Trash :size="18" />
                        </button>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import SecretDisplay from "./SecretDisplay.vue";
import { Trash, ArrowDownWideNarrow, ArrowUpNarrowWide } from "lucide-vue-next";
import { toast } from "vue3-toastify";
import { SERVER_URL } from "@/main";

const showSecretTimeRemaining = ref(null);
let showSecretInterval: ReturnType<typeof setInterval>;
const showSecretPercentRemaining = computed(() => {
    if (showSecretTimeRemaining.value === null) {
        return 0;
    }
    return Math.floor((showSecretTimeRemaining.value / 10) * 100);
});
const showSecret = computed(() => showSecretPercentRemaining.value > 0);
function clickShowSecret() {
    // Start the timer to show the secret for 10 seconds
    showSecretTimeRemaining.value = 10;
    if (showSecretInterval) {
        clearInterval(showSecretInterval);
    }
    showSecretInterval = setInterval(() => {
        if (showSecretTimeRemaining.value !== null) {
            showSecretTimeRemaining.value -= 1;
            if (showSecretTimeRemaining.value <= 0) {
                showSecretTimeRemaining.value = 0;
                clearInterval(showSecretInterval!);
                showSecretInterval = null;
            }
        }
    }, 1000);
}

onMounted(() => fetchSecrets());
const secrets = ref([]);
async function fetchSecrets() {
    try {
        const response = await fetch(`${SERVER_URL}/secrets`);
        const data = await response.json();
        secrets.value = data.secrets;
    } catch (error) {
        console.error("Error fetching secrets:", error);
    }
}

const sortColumn = ref("id");
const sortAscending = ref(true);
const filterText = ref("");

const tableSecrets = computed(() => {
    const sortedSecrets = secrets.value.sort((a, b) => {
        if (sortColumn.value === "id") {
            return sortAscending.value ? a.id - b.id : b.id - a.id;
        } else if (sortColumn.value === "key") {
            return sortAscending.value ? a.key.localeCompare(b.key) : b.key.localeCompare(a.key);
        }
        return 0;
    });
    const filteredSecrets = sortedSecrets.filter((secret) => secret.key && secret.key.toLowerCase().includes(filterText.value.toLowerCase()));
    return filteredSecrets;
});

async function deleteSecret(id: number) {
    try {
        await fetch(`${SERVER_URL}/secrets/${id}`, { method: "DELETE" });
        await fetchSecrets();
    } catch (error) {
        console.error("Error deleting secret:", error);
        toast("Error deleting secret", { type: "error" });
    }
}
</script>
