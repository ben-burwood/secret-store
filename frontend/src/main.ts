import { createApp } from 'vue';
import App from './App.vue';
import './style.css';
import "vue3-toastify/dist/index.css";
import Vue3Toastify, { type ToastContainerOptions } from 'vue3-toastify';

export const SERVER_URL = `${window.location.protocol}//${window.location.hostname}:${window.location.port}/web`;

createApp(App).use(
  Vue3Toastify, { theme: "colored", position: "bottom-right" } as ToastContainerOptions,
).mount('#app')
