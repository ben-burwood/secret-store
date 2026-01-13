<template>
    <Layout v-model:selectedPage="selectedPage">
        <component :is="currentComponent" />
    </Layout>
</template>

<script setup lang="ts">
import Layout from "@/components/Layout.vue";
import Generate from "@/components/pages/Generate.vue";
import Home from "@/components/pages/Home.vue";
import Authentication from "@/components/pages/Authentication.vue";
import Settings from "@/components/pages/Settings.vue";
import { ref, computed } from "vue";
import { Page } from "@/types/page";

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
</script>
