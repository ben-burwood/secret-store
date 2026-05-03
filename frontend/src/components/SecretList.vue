<template>
    <div class="flex flex-row justify-between">
        <TimerButton @enabled="showSecrets = true" @disabled="showSecrets = false">Show Secrets</TimerButton>
        <button class="btn btn-primary" @click="showAddRow = true" :disabled="showAddRow">+ Add Secret</button>
    </div>

    <div class="overflow-x-auto rounded-box border border-base-content/5 bg-base-200 my-2">
        <table class="table">
            <thead>
                <tr>
                    <th>Key</th>
                    <th>Value</th>
                    <th>Tag</th>
                    <th></th>
                </tr>
            </thead>
            <tbody>
                <AddSecretRow v-if="showAddRow" @refresh="emit('refresh')" @cancel="showAddRow = false" />
                <SecretRow v-for="secret in tableSecrets" :key="secret.id" :secret="secret" :showSecret="showSecrets" @refresh="emit('refresh')" />
            </tbody>
        </table>
    </div>

    <datalist id="existing-tags">
        <option v-for="tag in uniqueTags" :key="tag" :value="tag.toUpperCase()" />
    </datalist>

    <div v-if="props.secrets.length > 1" class="flex flex-row items-center gap-2 justify-end">
        <template v-if="uniqueTags.length > 0">
            Tag:
            <select v-model="tagFilter" class="select">
                <option value="">All</option>
                <option v-for="tag in uniqueTags" :key="tag" :value="tag">{{ tag.toUpperCase() }}</option>
            </select>
        </template>
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
import { ref, computed } from "vue";
import { ArrowDownWideNarrow, ArrowUpNarrowWide } from "lucide-vue-next";
import SecretRow from "./secrets-table/SecretRow.vue";
import AddSecretRow from "./secrets-table/AddSecretRow.vue";
import TimerButton from "./secrets-table/TimerButton.vue";

type Secret = { id: number; key: string; value: string; tag?: string | null };

const showAddRow = ref(false);
const showSecrets = ref(false);

const props = defineProps<{ secrets: Secret[] }>();

const emit = defineEmits(["refresh"]);

const sortColumn = ref("id");
const sortAscending = ref(true);
const filterText = ref("");
const tagFilter = ref("");

const uniqueTags = computed(() => {
    const tags = new Set<string>();
    for (const secret of props.secrets) {
        if (secret.tag) tags.add(secret.tag);
    }
    return Array.from(tags).sort((a, b) => a.localeCompare(b));
});

const tableSecrets = computed(() => {
    const sortedSecrets = [...props.secrets].sort((a, b) => {
        if (sortColumn.value === "id") {
            return sortAscending.value ? a.id - b.id : b.id - a.id;
        } else if (sortColumn.value === "key") {
            return sortAscending.value ? a.key.localeCompare(b.key) : b.key.localeCompare(a.key);
        }
        return 0;
    });
    const needle = filterText.value.toLowerCase();
    return sortedSecrets.filter((secret) => {
        if (!secret.key || !secret.key.toLowerCase().includes(needle)) return false;
        if (tagFilter.value && secret.tag !== tagFilter.value) return false;
        return true;
    });
});
</script>
