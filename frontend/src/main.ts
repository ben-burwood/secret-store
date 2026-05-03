import { createApp, ref } from 'vue';
import App from './App.vue';
import './style.css';
import "vue3-toastify/dist/index.css";
import Vue3Toastify, { type ToastContainerOptions } from 'vue3-toastify';

export const SERVER_URL = `${window.location.protocol}//${window.location.hostname}:${window.location.port}/web`;

export const authed = ref<boolean | null>(null);

export async function backendFetch(path: string, init: RequestInit = {}): Promise<Response> {
  const res = await fetch(`${SERVER_URL}${path}`, { ...init, credentials: 'include' });
  if (res.status === 401 && authed.value !== false) {
    authed.value = false;
  }
  return res;
}

createApp(App).use(
  Vue3Toastify, { theme: "colored", position: "bottom-right" } as ToastContainerOptions,
).mount('#app')
