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
      <NCard :bordered="false" size="small" class="stats-card card-devices">
        <template #header>
          <div class="card-header">
            <icon-mdi-monitor class="card-icon" />
            <span class="card-title">Top Dispositivos</span>
          </div>
        </template>
        <NSpin :show="loading">
          <div class="stats-content">
            <NList v-if="stats.top_devices.length > 0" size="small">
              <NListItem v-for="(device, index) in stats.top_devices.slice(0, 5)" :key="device.rustdesk_id" class="list-item-hover">
                <template #prefix>
                  <div class="rank-number" :class="`rank-${index + 1}`">
                    {{ index + 1 }}
                  </div>
                </template>
                <NSpace justify="space-between" style="width: 100%">
                  <NText class="item-text">{{ device.rustdesk_id }}</NText>
                  <NBadge :value="device.count" :color="index < 3 ? '#f59e0b' : '#3b82f6'" />
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
      <NCard :bordered="false" size="small" class="stats-card card-users">
        <template #header>
          <div class="card-header">
            <icon-mdi-account-group class="card-icon" />
            <span class="card-title">Top Usuários</span>
          </div>
        </template>
        <NSpin :show="loading">
          <div class="stats-content">
            <NList v-if="stats.top_users.length > 0" size="small">
              <NListItem v-for="(user, index) in stats.top_users.slice(0, 5)" :key="user.user_id" class="list-item-hover">
                <template #prefix>
                  <div class="rank-number" :class="`rank-${index + 1}`">
                    {{ index + 1 }}
                  </div>
                </template>
                <NSpace justify="space-between" style="width: 100%">
                  <NText class="item-text">{{ user.username }}</NText>
                  <NBadge :value="user.count" :color="index < 3 ? '#10b981' : '#3b82f6'" />
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
      <NCard :bordered="false" size="small" class="stats-card card-duration">
        <template #header>
          <div class="card-header">
            <icon-mdi-clock-outline class="card-icon" />
            <span class="card-title">Tempo Médio</span>
          </div>
        </template>
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
      <NCard :bordered="false" size="small" class="stats-card card-daily">
        <template #header>
          <div class="card-header">
            <icon-mdi-chart-line class="card-icon" />
            <span class="card-title">Últimos 7 Dias</span>
          </div>
        </template>
        <NSpin :show="loading">
          <div class="stats-content">
            <NList v-if="stats.daily_stats.length > 0" size="small">
              <NListItem v-for="day in stats.daily_stats.slice(0, 7)" :key="day.date" class="list-item-hover">
                <NSpace justify="space-between" style="width: 100%" align="center">
                  <NText class="date-text">{{ day.date }}</NText>
                  <div class="daily-count">
                    <span class="count-value">{{ day.count }}</span>
                  </div>
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
/* Card headers with icons */
.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  font-size: 14px;
}

.card-icon {
  font-size: 18px;
  opacity: 0.8;
}

.card-title {
  font-weight: 500;
}

/* Card color accents */
.card-devices :deep(.n-card-header) {
  border-left: 3px solid #f59e0b;
}

.card-users :deep(.n-card-header) {
  border-left: 3px solid #10b981;
}

.card-duration :deep(.n-card-header) {
  border-left: 3px solid #3b82f6;
}

.card-daily :deep(.n-card-header) {
  border-left: 3px solid #8b5cf6;
}

/* Rank numbers */
.rank-number {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  transition: all 0.2s ease;
}

.rank-1 {
  background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
  color: white;
}

.rank-2 {
  background: linear-gradient(135deg, #94a3b8 0%, #64748b 100%);
  color: white;
}

.rank-3 {
  background: linear-gradient(135deg, #fb923c 0%, #f97316 100%);
  color: white;
}

.rank-4, .rank-5 {
  background: #f1f5f9;
  color: #64748b;
}

/* List items */
.list-item-hover {
  transition: all 0.15s ease;
  border-radius: 6px;
}

.list-item-hover:hover {
  background: rgba(59, 130, 246, 0.05);
}

.item-text {
  font-size: 13px;
  font-weight: 500;
  color: #374151;
}

/* Date text styling */
.date-text {
  font-size: 12px;
  color: #6b7280;
  font-weight: 500;
}

/* Daily count badge */
.daily-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 32px;
  height: 22px;
  padding: 0 8px;
  background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%);
  border-radius: 11px;
  font-size: 11px;
  font-weight: 600;
  color: white;
}

:deep(.n-statistic) {
  text-align: center;
}

.stats-card {
  height: 100%;
  max-height: 280px;
  transition: all 0.2s ease;
}

.stats-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.stats-content {
  max-height: 220px;
  overflow-y: auto;
  overflow-x: hidden;
}

/* Scrollbar customization */
.stats-content::-webkit-scrollbar {
  width: 5px;
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
