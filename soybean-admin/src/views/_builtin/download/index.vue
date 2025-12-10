<script setup lang="ts">
import { computed, ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { getPaletteColorByNumber, mixColor } from '@sa/color';
import { useThemeStore } from '@/store/modules/theme';
import { useAppStore } from '@/store/modules/app';

const themeStore = useThemeStore();
const appStore = useAppStore();
const router = useRouter();

// URL base do backend (em produção será a mesma origem, em dev usa proxy)
const getApiBaseUrl = () => {
  // Em produção, usar a mesma origem
  if (import.meta.env.PROD) {
    return '';
  }
  // Em desenvolvimento, usar o proxy configurado ou URL direta
  // O backend está em localhost:12345
  return 'http://localhost:12345';
};

const apiBaseUrl = getApiBaseUrl();

const bgThemeColor = computed(() =>
  themeStore.darkMode ? getPaletteColorByNumber(themeStore.themeColor, 600) : themeStore.themeColor
);

const bgColor = computed(() => {
  const COLOR_WHITE = '#ffffff';
  const ratio = themeStore.darkMode ? 0.5 : 0.2;
  return mixColor(COLOR_WHITE, themeStore.themeColor, ratio);
});

interface DownloadItem {
  name: string;
  platform: string;
  icon: string;
  url: string;
  size: string;
  external: boolean;
}

const loading = ref(true);
const downloads = ref<DownloadItem[]>([]);

// Ícones e labels por plataforma
const platformInfo: Record<string, { icon: string; label: string }> = {
  windows: { icon: 'logos-microsoft-windows', label: 'Windows' },
  macos: { icon: 'logos-apple', label: 'macOS' },
  linux: { icon: 'logos-linux-tux', label: 'Linux' }
};

// Formatar tamanho do arquivo
function formatSize(bytes: number): string {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
}

// Carregar lista de instaladores da API
async function loadInstallers() {
  loading.value = true;
  try {
    const response = await fetch(`${apiBaseUrl}/api/download/list`);
    const data = await response.json();
    
    if (data.installers && data.installers.length > 0) {
      downloads.value = data.installers.map((item: any) => ({
        name: item.name,
        platform: platformInfo[item.platform]?.label || item.platform,
        icon: platformInfo[item.platform]?.icon || 'mdi-download',
        url: item.external ? item.url : `${apiBaseUrl}${item.url}`,
        size: item.size > 0 ? formatSize(item.size) : '',
        external: item.external || false
      }));
    } else {
      // Fallback para links padrão se não houver instaladores
      downloads.value = [
        {
          name: 'Windows Installer',
          platform: 'Windows',
          icon: 'logos-microsoft-windows',
          url: `${apiBaseUrl}/api/download/windows`,
          size: '',
          external: false
        },
        {
          name: 'macOS Installer',
          platform: 'macOS',
          icon: 'logos-apple',
          url: `${apiBaseUrl}/api/download/macos`,
          size: '',
          external: false
        },
        {
          name: 'Linux Installer',
          platform: 'Linux',
          icon: 'logos-linux-tux',
          url: `${apiBaseUrl}/api/download/linux`,
          size: '',
          external: false
        }
      ];
    }
  } catch (error) {
    console.error('Error loading installers:', error);
    // Em caso de erro, mostrar links padrão
    downloads.value = [
      {
        name: 'Windows Installer',
        platform: 'Windows',
        icon: 'logos-microsoft-windows',
        url: `${apiBaseUrl}/api/download/windows`,
        size: '',
        external: false
      },
      {
        name: 'macOS Installer',
        platform: 'macOS',
        icon: 'logos-apple',
        url: `${apiBaseUrl}/api/download/macos`,
        size: '',
        external: false
      },
      {
        name: 'Linux Installer',
        platform: 'Linux',
        icon: 'logos-linux-tux',
        url: `${apiBaseUrl}/api/download/linux`,
        size: '',
        external: false
      }
    ];
  } finally {
    loading.value = false;
  }
}

const handleDownload = (item: DownloadItem) => {
  // Criar um link temporário para download direto (sem popup)
  const link = document.createElement('a');
  link.href = item.url;
  link.download = item.name;
  link.style.display = 'none';
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
};

onMounted(() => {
  loadInstallers();
});
</script>

<template>
  <div class="relative size-full flex-center overflow-hidden" :style="{ backgroundColor: bgColor }">
    <WaveBg :theme-color="bgThemeColor" />
    <NCard :bordered="false" class="relative z-4 w-auto rd-12px">
      <div class="w-500px lt-sm:w-350px">
        <header class="flex-y-center justify-between">
          <SystemLogo class="text-64px text-primary lt-sm:text-48px" />
          <h3 class="text-28px text-primary font-500 lt-sm:text-22px">MT Remoto Download</h3>
          <div class="i-flex-col">
            <ThemeSchemaSwitch
              :theme-schema="themeStore.themeScheme"
              :show-tooltip="false"
              class="text-20px lt-sm:text-18px"
              @switch="themeStore.toggleThemeScheme"
            />
            <LangSwitch
              :lang="appStore.locale"
              :lang-options="appStore.localeOptions"
              :show-tooltip="false"
              @change-lang="appStore.changeLocale"
            />
          </div>
        </header>

        <main class="pt-24px">
          <h3 class="text-18px text-primary font-medium mb-16px">
            Escolha sua plataforma para baixar o cliente MT Remoto
          </h3>

          <NSpace vertical :size="16">
            <NCard
              v-for="item in downloads"
              :key="item.platform"
              hoverable
              class="cursor-pointer transition-all hover:shadow-lg"
              @click="handleDownload(item)"
            >
              <div class="flex items-center gap-16px">
                <div class="text-48px">
                  <span :class="item.icon"></span>
                </div>
                <div class="flex-1">
                  <div class="text-16px font-medium">{{ item.name }}</div>
                  <div class="text-14px text-gray-500">{{ item.platform }}</div>
                  <div v-if="item.size" class="text-12px text-gray-400">{{ item.size }}</div>
                </div>
                <NButton type="primary" size="large">
                  <template #icon>
                    <icon-mdi-download />
                  </template>
                  Download
                </NButton>
              </div>
            </NCard>
          </NSpace>

          <NDivider />

          <div class="text-center text-14px text-gray-500">
            <p>Após o download, execute o instalador e siga as instruções.</p>
            <p class="mt-8px">
              Precisa de ajuda?
              <NButton text type="primary" @click="router.push('/login')">
                Faça login
              </NButton>
              para acessar o suporte.
            </p>
          </div>
        </main>
      </div>
    </NCard>
  </div>
</template>

<style scoped>
.transition-all {
  transition: all 0.3s ease;
}
</style>
