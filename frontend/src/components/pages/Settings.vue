<template>
    <BackupRestore />
    <Import class="mt-5" />

    <div class="card bg-base-200 shadow-md max-w-md mx-auto mt-5">
        <div class="card-body">
            <h2 class="card-title">Session</h2>
            <button class="btn btn-error" @click="logout" :disabled="loggingOut">
                {{ loggingOut ? "Logging out..." : "Logout" }}
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import BackupRestore from "@/components/BackupRestore.vue";
import Import from "@/components/Import.vue";
import { backendFetch, authed } from "@/main";
import { toast } from "vue3-toastify";

const loggingOut = ref(false);

async function logout() {
    loggingOut.value = true;
    try {
        await backendFetch("/auth/logout", { method: "POST" });
    } catch (error) {
        console.error("Error logging out:", error);
        toast("Error logging out", { type: "error" });
    } finally {
        loggingOut.value = false;
        authed.value = false;
    }
}
</script>
