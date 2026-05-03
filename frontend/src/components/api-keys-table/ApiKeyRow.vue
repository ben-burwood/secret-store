<template>
    <tr>
        <td>{{ apiKey.name }}</td>
        <td>
            <SecretDisplay :secret="apiKey.key" :showSecret="false" />
        </td>
        <td>{{ created }}</td>
        <td>
            <div class="flex flex-row gap-2">
                <button class="btn btn-outline btn-primary btn-square" @click="regenerate" aria-label="Regenerate key">
                    <RefreshCw :size="18" />
                </button>
                <button class="btn btn-outline btn-error btn-square" @click="remove" aria-label="Delete key">
                    <Trash :size="18" />
                </button>
            </div>
        </td>
    </tr>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { RefreshCw, Trash } from "lucide-vue-next";
import SecretDisplay from "@/components/SecretDisplay.vue";
import { toast } from "vue3-toastify";
import { backendFetch } from "@/main";

const props = defineProps<{
    apiKey: {
        id: number;
        name: string;
        key: string;
        created_at: string | null;
    };
}>();

const emit = defineEmits(["refresh"]);

const created = computed(() => (props.apiKey.created_at ? new Date(props.apiKey.created_at).toLocaleString() : "—"));

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
