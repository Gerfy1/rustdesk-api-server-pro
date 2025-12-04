<script setup lang="tsx">
import { onMounted, onUnmounted, ref } from 'vue';
import { NTag, NBadge } from 'naive-ui';
import { fetchDevicesList } from '@/service/api/devices';
import { $t } from '@/locales';
import { useAppStore } from '@/store/modules/app';
import { useTable } from '@/hooks/common/table';
import { AuditTypeOptions } from '@/constants/business';
import TableHeader from './components/table-header.vue';
import AuditBaseLogsSearch from './components/search.vue';

const appStore = useAppStore();

const auditTypes: any = {
  '0': '',
  '1': 'error',
  '2': 'success'
};

// Auto-refresh configuration
let autoRefreshTimer: NodeJS.Timeout | null = null;
const autoRefreshInterval = 30000; // 30 seconds

// Helper function to format relative time
const formatRelativeTime = (dateString: string) => {
  if (!dateString) return 'N/A';
  const date = new Date(dateString);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffSec = Math.floor(diffMs / 1000);
  const diffMin = Math.floor(diffSec / 60);
  const diffHour = Math.floor(diffMin / 60);
  const diffDay = Math.floor(diffHour / 24);

  if (diffSec < 60) return 'Agora';
  if (diffMin < 60) return `há ${diffMin} min`;
  if (diffHour < 24) return `há ${diffHour}h`;
  return `há ${diffDay}d`;
};

// Format memory: convert MB to GB if needed, or keep original format
const formatMemory = (memory: string) => {
  if (!memory) return '-';
  
  // If it's already in GB format (e.g., "5.83 GB"), extract and fix
  const gbMatch = memory.match(/(\d+\.?\d*)\s*GB/i);
  if (gbMatch) {
    const gb = parseFloat(gbMatch[1]);
    // Round to whole number if close to it
    const rounded = Math.round(gb);
    if (Math.abs(gb - rounded) < 0.1) {
      return `${rounded} GB`;
    }
    return `${gb.toFixed(1)} GB`;
  }
  
  // If it's in MB format (e.g., "6144 MB")
  const mbMatch = memory.match(/(\d+)\s*MB/i);
  if (mbMatch) {
    const mb = parseInt(mbMatch[1]);
    const gb = mb / 1024;
    // If it's a round number of GB, show as integer
    if (gb >= 1 && gb % 1 === 0) {
      return `${gb} GB`;
    }
    // If close to a round number, round it
    const rounded = Math.round(gb);
    if (Math.abs(gb - rounded) < 0.1) {
      return `${rounded} GB`;
    }
    // Otherwise show 1 decimal place
    return `${gb.toFixed(1)} GB`;
  }
  
  // Return as-is if format not recognized
  return memory;
};

const {
  columns,
  columnChecks,
  data,
  getData,
  getDataByPage,
  loading,
  mobilePagination,
  searchParams,
  resetSearchParams
} = useTable({
  apiFn: fetchDevicesList,
  showTotal: true,
  apiParams: {
    current: 1,
    size: 10,
    hostname: null,
    username: null,
    rustdesk_id: null,
    status: null
  },
  columns: () => [
    {
      key: 'id',
      title: 'ID',
      align: 'center',
      width: 60
    },
    {
      key: 'is_online',
      title: 'Status',
      align: 'center',
      width: 100,
      render: (row: any) => {
        if (row.is_online) {
          return <NTag type="success" size="small">🟢 Online</NTag>;
        }
        return <NTag type="error" size="small">🔴 Offline</NTag>;
      }
    },
    {
      key: 'rustdesk_id',
      title: $t('dataMap.device.rustdesk_id'),
      align: 'center',
      width: 120
    },
    {
      key: 'hostname',
      title: $t('dataMap.device.hostname'),
      align: 'center',
      width: 150
    },
    {
      key: 'username',
      title: $t('dataMap.device.username'),
      align: 'center',
      width: 120
    },
    {
      key: 'ip_address',
      title: 'IP',
      align: 'center',
      width: 140,
      render: (row: any) => row.ip_address || 'N/A'
    },
    {
      key: 'conns',
      title: 'Conexões Ativas',
      align: 'center',
      width: 110,
      render: (row: any) => {
        if (row.conns > 0) {
          return <NBadge value={row.conns} type="success" />;
        }
        return <span style="color: #999">0</span>;
      }
    },
    {
      key: 'total_accesses',
      title: 'Total Acessos',
      align: 'center',
      width: 110,
      render: (row: any) => {
        if (row.total_accesses > 0) {
          return <NBadge value={row.total_accesses} color="#3b82f6" />;
        }
        return <span style="color: #999">0</span>;
      }
    },
    {
      key: 'last_connection_at',
      title: 'Última Conexão',
      align: 'center',
      width: 140,
      render: (row: any) => {
        if (!row.last_connection_at) {
          return <span style="color: #999">Nunca</span>;
        }
        const relativeTime = formatRelativeTime(row.last_connection_at);
        return <span title={row.last_connection_at}>{relativeTime}</span>;
      }
    },
    {
      key: 'last_seen_at',
      title: 'Último Heartbeat',
      align: 'center',
      width: 140,
      render: (row: any) => {
        const relativeTime = formatRelativeTime(row.last_seen_at);
        return <span title={row.last_seen_at} style="color: #999; font-size: 12px">{relativeTime}</span>;
      }
    },
    {
      key: 'version',
      title: $t('dataMap.device.version'),
      align: 'center',
      width: 100
    },
    {
      key: 'os',
      title: $t('dataMap.device.os'),
      align: 'center',
      width: 120
    },
    {
      key: 'memory',
      title: 'RAM Total',
      align: 'center',
      width: 100,
      render: (row: any) => formatMemory(row.memory)
    },
    {
      key: 'created_at',
      title: 'Criado em',
      align: 'center',
      width: 100,
      render: (row: any) => {
        if (!row.created_at) return '-';
        const [date, time] = row.created_at.split(' ');
        return (
          <div style="line-height: 1.4">
            <div style="font-size: 12px">{date}</div>
            <div style="font-size: 11px; color: #999">{time}</div>
          </div>
        );
      }
    }
  ]
});

// Setup auto-refresh
onMounted(() => {
  autoRefreshTimer = setInterval(() => {
    if (!loading.value) {
      getData();
    }
  }, autoRefreshInterval);
});

onUnmounted(() => {
  if (autoRefreshTimer) {
    clearInterval(autoRefreshTimer);
  }
});

</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <AuditBaseLogsSearch v-model:model="searchParams" @reset="resetSearchParams" @search="getDataByPage" />

    <NCard :title="$t('route.devices')" :bordered="false" size="small" class="sm:flex-1-hidden card-wrapper">
      <template #header-extra>
        <NSpace>
          <NTag type="info" size="small">Auto-refresh: 30s</NTag>
          <TableHeader v-model:columns="columnChecks" :loading="loading" @refresh="getData" />
        </NSpace>
      </template>
      <NDataTable
        :columns="columns"
        :data="data"
        size="small"
        :flex-height="!appStore.isMobile"
        :scroll-x="1540"
        :loading="loading"
        remote
        :row-key="(row: any) => row.id"
        :pagination="mobilePagination"
        class="sm:h-full"
      />
    </NCard>
  </div>
</template>

<style scoped></style>
