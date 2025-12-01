<script setup lang="ts">
import { onMounted, ref, onUnmounted } from 'vue';
import { $t } from '@/locales';
import { fetchDevicesOnlineList } from '@/service/api';

defineOptions({
  name: 'CreativityBanner'
});

interface Stats {
  online: number;
  offline: number;
  total: number;
  connections: number;
}

const stats = ref<Stats>({
  online: 0,
  offline: 0,
  total: 0,
  connections: 0
});

const loading = ref(false);
let intervalId: number | null = null;

async function fetchStats() {
  loading.value = true;
  try {
    const { data, error } = await fetchDevicesOnlineList();
    if (!error && data) {
      const onlineDevices = data.filter((d: any) => d.is_online);
      stats.value.online = onlineDevices.length;
      stats.value.offline = data.length - onlineDevices.length;
      stats.value.total = data.length;
      // Count total connections from all devices
      stats.value.connections = data.reduce((sum: number, d: any) => sum + (d.conns || 0), 0);
    }
  } catch (error) {
    console.error('Failed to fetch device stats:', error);
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  fetchStats();
  // Auto-refresh every 30 seconds
  intervalId = window.setInterval(fetchStats, 30000);
});

onUnmounted(() => {
  if (intervalId) {
    clearInterval(intervalId);
  }
});
</script>

<template>
  <NCard :title="$t('page.home.realtimeStats')" :bordered="false" size="small" class="h-full card-wrapper">
    <NSpin :show="loading">
      <div class="h-full">
        <NGrid :cols="2" :x-gap="12" :y-gap="12">
          <NGi>
            <div class="stat-item online">
              <div class="stat-icon">
                <icon-mdi-circle class="text-success" />
              </div>
              <div class="stat-content">
                <div class="stat-label">{{ $t('page.home.devicesOnline') }}</div>
                <div class="stat-value text-success">{{ stats.online }}</div>
              </div>
            </div>
          </NGi>
          <NGi>
            <div class="stat-item offline">
              <div class="stat-icon">
                <icon-mdi-circle class="text-error" />
              </div>
              <div class="stat-content">
                <div class="stat-label">{{ $t('page.home.devicesOffline') }}</div>
                <div class="stat-value text-error">{{ stats.offline }}</div>
              </div>
            </div>
          </NGi>
          <NGi>
            <div class="stat-item total">
              <div class="stat-icon">
                <icon-mdi-devices class="text-primary" />
              </div>
              <div class="stat-content">
                <div class="stat-label">{{ $t('page.home.totalDevices') }}</div>
                <div class="stat-value text-primary">{{ stats.total }}</div>
              </div>
            </div>
          </NGi>
          <NGi>
            <div class="stat-item connections">
              <div class="stat-icon">
                <icon-mdi-connection class="text-warning" />
              </div>
              <div class="stat-content">
                <div class="stat-label">{{ $t('page.home.activeConnections') }}</div>
                <div class="stat-value text-warning">{{ stats.connections }}</div>
              </div>
            </div>
          </NGi>
        </NGrid>
      </div>
    </NSpin>
  </NCard>
</template>

<style scoped>
.stat-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background-color: var(--n-card-color);
  border-radius: 8px;
  transition: all 0.3s ease;
  border: 1px solid var(--n-border-color);
}

.stat-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-icon {
  font-size: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-label {
  font-size: 12px;
  color: var(--n-text-color-disabled);
  font-weight: 500;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  line-height: 1;
}

.text-success {
  color: #2ecc71;
}

.text-error {
  color: #e74c3c;
}

.text-primary {
  color: #95a5a6;
}

.text-warning {
  color: #3498db;
}
</style>
