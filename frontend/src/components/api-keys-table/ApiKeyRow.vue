<template>
    <tr :class="{ 'bg-success/10': isEditing }">
        <td>{{ apiKey.name }}</td>
        <td>
            <SecretDisplay :secret="apiKey.key" :showSecret="false" />
        </td>
        <td>{{ created }}</td>
        <td>
            <template v-if="isEditing">
                <select multiple v-model="editScopes" class="select w-full" :size="Math.min(Math.max(secrets.length, 3), 6)">
                    <option v-for="s in secrets" :key="s.id" :value="s.id">{{ s.key }}</option>
                </select>
            </template>
            <template v-else>
                <div class="flex flex-wrap gap-1">
                    <span v-if="apiKey.secret_ids.length === 0" class="badge badge-outline">All secrets</span>
                    <span v-for="s in scopedSecrets" :key="s.id" class="badge badge-outline">{{ s.label }}</span>
                </div>
            </template>
        </td>
        <td>
            <div v-if="!isEditing" class="flex flex-row gap-2">
                <button class="btn btn-outline btn-primary btn-square" @click="startEdit" aria-label="Edit scopes">
                    <Pencil :size="18" />
                </button>
                <button class="btn btn-outline btn-primary btn-square" @click="regenerate" aria-label="Regenerate key">
                    <RefreshCw :size="18" />
                </button>
                <button class="btn btn-outline btn-error btn-square" @click="remove" aria-label="Delete key">
                    <Trash :size="18" />
                </button>
            </div>
            <ConfirmReject v-else @confirm="saveScopes" @reject="cancelEdit" :canConfirm="true" />
        </td>
    </tr>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { RefreshCw, Trash, Pencil } from "lucide-vue-next";
import SecretDisplay from "@/components/SecretDisplay.vue";
import ConfirmReject from "@/components/ConfirmReject.vue";
import { toast } from "vue3-toastify";
import { backendFetch } from "@/main";

const props = defineProps<{
    apiKey: {
        id: number;
        name: string;
        key: string;
        created_at: string | null;
        secret_ids: number[];
    };
    secrets: { id: number; key: string }[];
}>();

const emit = defineEmits(["refresh"]);

const isEditing = ref(false);
const editScopes = ref<number[]>([]);

const created = computed(() => (props.apiKey.created_at ? new Date(props.apiKey.created_at).toLocaleString() : "—"));

const scopedSecrets = computed(() => {
    const byId = new Map(props.secrets.map((s) => [s.id, s.key]));
    return props.apiKey.secret_ids.map((id) => ({ id, label: byId.get(id) ?? "Unknown" }));
});

function startEdit() {
    editScopes.value = [...props.apiKey.secret_ids];
    isEditing.value = true;
}

function cancelEdit() {
    isEditing.value = false;
}

async function saveScopes() {
    try {
        const res = await backendFetch(`/api/keys/${props.apiKey.id}/scopes`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ secret_ids: editScopes.value }),
        });
        if (!res.ok) {
            const body = await res.json().catch(() => null);
            toast(body?.error ?? "Error saving scopes", { type: "error" });
            return;
        }
        emit("refresh");
        isEditing.value = false;
    } catch (error) {
        console.error("Error saving scopes:", error);
        toast("Error saving scopes", { type: "error" });
    }
}

async function regenerate() {
    try {
        const res = await backendFetch(`/api/keys/${props.apiKey.id}/regenerate`, { method: "POST" });
        if (!res.ok) {
            const body = await res.json().catch(() => null);
            toast(body?.error ?? "Error regenerating key", { type: "error" });
            return;
        }
        emit("refresh");
    } catch (error) {
        console.error("Error regenerating key:", error);
        toast("Error regenerating key", { type: "error" });
    }
}

async function remove() {
    try {
        await backendFetch(`/api/keys/${props.apiKey.id}`, { method: "DELETE" });
        emit("refresh");
    } catch (error) {
        console.error("Error deleting key:", error);
        toast("Error deleting key", { type: "error" });
    }
}
</script>
