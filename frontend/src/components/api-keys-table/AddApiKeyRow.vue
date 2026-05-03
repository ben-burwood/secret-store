<template>
    <tr class="bg-success/10">
        <td>
            <input v-model="newName" class="input input-bordered w-full" placeholder="Name" @keyup.enter="addKey" />
        </td>
        <td colspan="2"></td>
        <td>
            <select multiple v-model="newScopes" class="select w-full" :size="Math.min(Math.max(secrets.length, 3), 6)">
                <option v-for="s in secrets" :key="s.id" :value="s.id">{{ s.key }}</option>
            </select>
        </td>
        <td>
            <ConfirmReject @confirm="addKey" @reject="cancel" :canConfirm="newName.trim() !== ''" />
        </td>
    </tr>
</template>

<script setup lang="ts">
import { ref } from "vue";
import ConfirmReject from "@/components/ConfirmReject.vue";
import { toast } from "vue3-toastify";
import { backendFetch } from "@/main";

defineProps<{
    secrets: { id: number; key: string }[];
}>();

const emit = defineEmits(["refresh", "cancel"]);

const newName = ref("");
const newScopes = ref<number[]>([]);

const cancel = () => {
    newName.value = "";
    newScopes.value = [];
    emit("cancel");
};

async function addKey() {
    const name = newName.value.trim();
    if (!name) return;
    try {
        const res = await backendFetch(`/api/keys/new`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ name, secret_ids: newScopes.value }),
        });
        if (!res.ok) {
            const body = await res.json().catch(() => null);
            toast(body?.error ?? "Error adding key", { type: "error" });
            return;
        }
        emit("refresh");
        cancel();
    } catch (error) {
        console.error("Error adding key:", error);
        toast("Error adding key", { type: "error" });
        cancel();
    }
}
</script>
