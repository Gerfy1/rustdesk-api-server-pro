<script setup lang="tsx">
import { ref, watch } from 'vue';
import { NTag, NBadge, NDataTable, NButton, NSpace, NPopconfirm } from 'naive-ui';
import { fetchAddressBookPeers, importDevicesAsPeers } from '@/service/api/address-books';

defineOptions({
  name: 'PeersDrawer'
});

interface Props {
  addressBookId: number;
  addressBookName: string;
}

const props = defineProps<Props>();

const visible = defineModel<boolean>('visible', { required: true });

const loading = ref(false);
const peers = ref<Api.AddressBooks.Peer[]>([]);

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

const columns = [
  {
    key: 'is_online',
    title: 'Status',
    align: 'center' as const,
    width: 100,
    render: (row: Api.AddressBooks.Peer) => {
      if (row.is_online) {
        return <NTag type="success" size="small">🟢 Online</NTag>;
      }
      return <NTag type="error" size="small">🔴 Offline</NTag>;
    }
  },
  {
    key: 'rustdesk_id',
    title: 'RustDesk ID',
    align: 'center' as const,
    width: 120
  },
  {
    key: 'hostname',
    title: 'Hostname',
    align: 'center' as const,
    width: 150
  },
  {
    key: 'username',
    title: 'Username',
    align: 'center' as const,
    width: 120
  },
  {
    key: 'alias',
    title: 'Alias',
    align: 'center' as const,
    width: 120
  },
  {
    key: 'platform',
    title: 'Plataforma',
    align: 'center' as const,
    width: 120
  },
  {
    key: 'ip_address',
    title: 'IP',
    align: 'center' as const,
    width: 140
  },
  {
    key: 'last_seen_at',
    title: 'Última Conexão',
    align: 'center' as const,
    width: 140,
    render: (row: Api.AddressBooks.Peer) => {
      const relativeTime = formatRelativeTime(row.last_seen_at);
      return <span title={row.last_seen_at}>{relativeTime}</span>;
    }
  }
];

watch(visible, async (val) => {
  if (val && props.addressBookId) {
    await loadPeers();
  }
});

async function loadPeers() {
  loading.value = true;
  try {
    const response = await fetchAddressBookPeers(props.addressBookId);
    peers.value = response.data?.records || [];
  } catch (error) {
    window.$message?.error('Erro ao carregar peers');
  } finally {
    loading.value = false;
  }
}

async function handleImportDevices() {
  loading.value = true;
  try {
    const response = await importDevicesAsPeers(props.addressBookId);
    const data = response.data as any;
    window.$message?.success(
      `Importação concluída! ${data.imported} dispositivos importados, ${data.skipped} ignorados (já existentes)`
    );
    await loadPeers(); // Reload peers list
  } catch (error) {
    window.$message?.error('Erro ao importar dispositivos');
  } finally {
    loading.value = false;
  }
}

function handleClose() {
  visible.value = false;
}
</script>

<template>
  <NDrawer v-model:show="visible" :width="1000" placement="right">
    <NDrawerContent :title="`Peers em: ${addressBookName}`" closable>
      <div class="flex-col gap-16px">
        <NSpace justify="space-between" align="center">
          <NSpace>
            <NText>
              Total de peers: <NBadge :value="peers.length" type="info" />
            </NText>
            <NText v-if="peers.filter((p: Api.AddressBooks.Peer) => p.is_online).length > 0">
              Online: <NBadge :value="peers.filter((p: Api.AddressBooks.Peer) => p.is_online).length" type="success" />
            </NText>
          </NSpace>
          <NPopconfirm @positive-click="handleImportDevices">
            <template #trigger>
              <NButton type="primary" size="small" :loading="loading">
                <template #icon>
                  <icon-mdi-cloud-download class="text-icon" />
                </template>
                Importar Devices Online
              </NButton>
            </template>
            Importar todos os dispositivos online como peers? Dispositivos já existentes serão ignorados.
          </NPopconfirm>
        </NSpace>

        <NDataTable
          :columns="columns"
          :data="peers"
          :loading="loading"
          size="small"
          :scroll-x="1100"
          :max-height="600"
        />
      </div>

      <template #footer>
        <NSpace justify="end">
          <NButton @click="handleClose">Fechar</NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped></style>
