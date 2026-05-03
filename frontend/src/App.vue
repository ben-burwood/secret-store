<template>
    <div v-if="authed === null" class="min-h-screen flex items-center justify-center bg-base-300">
        <span class="loading loading-dots loading-lg"></span>
    </div>
    <Login v-else-if="!authed" @authenticated="authed = true" />
    <Layout v-else v-model:selectedPage="selectedPage">
        <component :is="currentComponent" />
    </Layout>
</template>

<script setup lang="ts">
import Layout from "@/components/Layout.vue";
import Generate from "@/components/pages/Generate.vue";
import Home from "@/components/pages/Home.vue";
import Authentication from "@/components/pages/Authentication.vue";
import Settings from "@/components/pages/Settings.vue";
import Login from "@/components/pages/Login.vue";
import { ref, computed, onMounted } from "vue";
import { Page } from "@/types/page";
import { authed, backendFetch } from "@/main";

const selectedPage = ref(Page.HOME);
const currentComponent = computed(() => {
    switch (selectedPage.value) {
        case Page.HOME:
            return Home;
        case Page.GENERATE:
            return Generate;
        case Page.AUTHENTICATION:
            return Authentication;
        case Page.SETTINGS:
            return Settings;
        default:
            return Home;
    }
});

onMounted(async () => {
    try {
        const res = await backendFetch("/auth/status");
        authed.value = res.ok;
    } catch {
        authed.value = false;
    }
});
</script>
