<script setup lang="ts">
import { ref, onMounted } from 'vue';

defineOptions({
  name: 'AuditStats'
});

interface Stats {
  top_devices: Array<{ rustdesk_id: string; count: number }>;
  top_users: Array<{ user_id: number; username: string; count: number }>;
  avg_duration: number;
  daily_stats: Array<{ date: string; count: number }>;
  total_connections: number;
}

const loading = ref(false);
const stats = ref<Stats>({
  top_devices: [],
  top_users: [],
  avg_duration: 0,
  daily_stats: [],
  total_connections: 0
});

async function fetchStats() {
  loading.value = true;
  try {
    const { request } = await import('@/service/request');
    const response = await request({
      url: '/audit/stats',
      method: 'get'
    });

    if (response.data) {
      stats.value = response.data;
    }
  } catch (error) {
    console.error('Erro ao buscar estatísticas:', error);
    window.$message?.error('Erro ao carregar estatísticas de auditoria');
  } finally {
    loading.value = false;
  }
}

// Format duration from seconds to human readable
function formatDuration(seconds: number): string {
  if (seconds < 60) return `${seconds}s`;
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}min`;
  const hours = Math.floor(minutes / 60);
  const remainingMinutes = minutes % 60;
  return `${hours}h ${remainingMinutes}min`;
}

onMounted(() => {
  fetchStats();
});
</script>

<template>
  <NGrid :cols="24" :x-gap="12" :y-gap="12" responsive="screen">
    <!-- Top 10 Dispositivos Mais Acessados -->
    <NGridItem :span="24" :md="12" :lg="6">
      <NCard title=" Top 10 Dispositivos" :bordered="false" size="small" class="stats-card">
        <NSpin :show="loading">
          <div class="stats-content">
            <NList v-if="stats.top_devices.length > 0" size="small">
              <NListItem v-for="(device, index) in stats.top_devices.slice(0, 5)" :key="device.rustdesk_id">
                <template #prefix>
                  <NTag :type="index < 3 ? 'warning' : 'default'" size="small">
                    {{ index + 1 }}
                  </NTag>
                </template>
                <NSpace justify="space-between" style="width: 100%">
                  <NText strong style="font-family: monospace; font-size: 12px">{{ device.rustdesk_id }}</NText>
                  <NBadge :value="device.count" type="info" />
                </NSpace>
              </NListItem>
            </NList>
            <NEmpty v-else description="Nenhum dado" size="small" />
          </div>
        </NSpin>
      </NCard>
    </NGridItem>

    <!-- Top 10 Usuários Mais Ativos -->
    <NGridItem :span="24" :md="12" :lg="6">
      <NCard title="Top 10 Usuários" :bordered="false" size="small" class="stats-card">
        <NSpin :show="loading">
          <div class="stats-content">
            <NList v-if="stats.top_users.length > 0" size="small">
              <NListItem v-for="(user, index) in stats.top_users.slice(0, 5)" :key="user.user_id">
                <template #prefix>
                  <NTag :type="index < 3 ? 'success' : 'default'" size="small">
                    {{ index + 1 }}
                  </NTag>
                </template>
                <NSpace justify="space-between" style="width: 100%">
                  <NText strong style="font-size: 12px">{{ user.username }}</NText>
                  <NBadge :value="user.count" type="success" />
                </NSpace>
              </NListItem>
            </NList>
            <NEmpty v-else description="Nenhum dado" size="small" />
          </div>
        </NSpin>
      </NCard>
    </NGridItem>

    <!-- Tempo Médio de Sessão -->
    <NGridItem :span="24" :md="12" :lg="6">
      <NCard title="Tempo Médio" :bordered="false" size="small" class="stats-card">
        <NSpin :show="loading">
          <NSpace vertical align="center" style="width: 100%; padding: 10px 0">
            <NStatistic>
              <template #label>
                <NText depth="3" style="font-size: 11px">Duração Média</NText>
              </template>
              <template #default>
                <NText style="font-size: 24px; font-weight: 600; color: var(--primary-color)">
                  {{ formatDuration(stats.avg_duration) }}
                </NText>
              </template>
            </NStatistic>
            <NText depth="3" style="font-size: 11px">
              {{ stats.total_connections }} conexões
            </NText>
          </NSpace>
        </NSpin>
      </NCard>
    </NGridItem>

    <!-- Conexões nos Últimos 7 Dias -->
    <NGridItem :span="24" :md="12" :lg="6">
      <NCard title="Últimos 7 Dias" :bordered="false" size="small" class="stats-card">
        <NSpin :show="loading">
          <div class="stats-content">
            <NList v-if="stats.daily_stats.length > 0" size="small">
              <NListItem v-for="day in stats.daily_stats.slice(0, 5)" :key="day.date">
                <NSpace justify="space-between" style="width: 100%">
                  <NText depth="2" style="font-size: 12px">{{ day.date }}</NText>
                  <NTag type="success" size="small">{{ day.count }}</NTag>
                </NSpace>
              </NListItem>
            </NList>
            <NEmpty v-else description="Nenhum dado" size="small" />
          </div>
        </NSpin>
      </NCard>
    </NGridItem>
  </NGrid>
</template>

<style scoped>
:deep(.n-statistic) {
  text-align: center;
}

.stats-card {
  height: 100%;
  max-height: 280px;
}

.stats-content {
  max-height: 220px;
  overflow-y: auto;
  overflow-x: hidden;
}

/* Scrollbar customization */
.stats-content::-webkit-scrollbar {
  width: 6px;
}

.stats-content::-webkit-scrollbar-track {
  background: transparent;
}

.stats-content::-webkit-scrollbar-thumb {
  background: var(--n-border-color);
  border-radius: 3px;
}

.stats-content::-webkit-scrollbar-thumb:hover {
  background: var(--n-text-color-disabled);
}
</style>
