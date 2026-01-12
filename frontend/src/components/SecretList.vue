<template>
    <div class="flex flex-row justify-between">
        <TimerButton @enabled="showSecrets = true" @disabled="showSecrets = false">Show Secrets</TimerButton>
        <button class="btn btn-primary" @click="showAddRow = true" :disabled="showAddRow">+ Add Secret</button>
    </div>

    <div class="overflow-x-auto rounded-box border border-base-content/5 bg-base-200 my-2">
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
                <AddSecretRow v-if="showAddRow" @refresh="fetchSecrets" @cancel="showAddRow = false" />
                <SecretRow v-for="secret in tableSecrets" :key="secret.id" :secret="secret" :showSecret="showSecrets" @refresh="fetchSecrets" />
            </tbody>
        </table>
    </div>

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
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import SecretDisplay from "./SecretDisplay.vue";
import { Trash, ArrowDownWideNarrow, ArrowUpNarrowWide } from "lucide-vue-next";
import ConfirmReject from "./secrets-table/ConfirmReject.vue";
import SecretRow from "./secrets-table/SecretRow.vue";
import AddSecretRow from "./secrets-table/AddSecretRow.vue";
import { toast } from "vue3-toastify";
import { SERVER_URL } from "@/main";
import TimerButton from "./secrets-table/TimerButton.vue";

const showAddRow = ref(false);

const showSecrets = ref(false);

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
onMounted(fetchSecrets);

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
</script>
