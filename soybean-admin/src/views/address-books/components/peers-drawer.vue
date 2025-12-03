<script setup lang="tsx">
import { ref, watch } from 'vue';
import { NTag, NBadge, NDataTable, NButton, NSpace, NPopconfirm } from 'naive-ui';
import { fetchAddressBookPeers, importDevicesAsPeers } from '@/service/api/address-books';
import AddPeerModal from './add-peer-modal.vue';

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

// Add Peer Modal state
const addPeerModalVisible = ref(false);

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
    width: 140,
    render: (row: Api.AddressBooks.Peer) => {
      // Show rustdesk_id (the actual unique identifier)
      return <span style="font-weight: 500">{row.rustdesk_id}</span>;
    }
  },
  {
    key: 'alias',
    title: 'Alias / Nome',
    align: 'center' as const,
    width: 140,
    render: (row: Api.AddressBooks.Peer) => {
      return <span style="color: #666">{row.alias || '-'}</span>;
    }
  },
  {
    key: 'tags',
    title: 'Tags',
    align: 'center' as const,
    width: 180,
    render: (row: Api.AddressBooks.Peer) => {
      if (!row.tags || row.tags.length === 0) {
        return <span style="color: #999; font-size: 12px">Sem tags</span>;
      }
      try {
        const tags = typeof row.tags === 'string' ? JSON.parse(row.tags) : row.tags;
        return (
          <NSpace size="small" justify="center" style="width: 100%">
            {tags.map((tag: string, index: number) => (
              <NTag key={index} type="info" size="small">{tag}</NTag>
            ))}
          </NSpace>
        );
      } catch (e) {
        return <span style="color: #999; font-size: 12px">-</span>;
      }
    }
  },
  {
    key: 'hostname',
    title: 'Hostname',
    align: 'center' as const,
    width: 150,
    render: (row: Api.AddressBooks.Peer) => {
      return <span style="font-size: 13px; color: #666">{row.hostname || '-'}</span>;
    }
  },
  {
    key: 'username',
    title: 'Username',
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
  },
  {
    key: 'actions',
    title: 'Ações',
    align: 'center' as const,
    width: 100,
    render: (row: Api.AddressBooks.Peer) => {
      return (
        <NPopconfirm onPositiveClick={() => handleDeletePeer(row.id)}>
          {{
            default: () => 'Tem certeza que deseja deletar este peer?',
            trigger: () => (
              <NButton type="error" size="small" quaternary>
                {{
                  icon: () => <icon-mdi-delete class="text-icon" />,
                  default: () => 'Deletar'
                }}
              </NButton>
            )
          }}
        </NPopconfirm>
      );
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

async function handleDeletePeer(peerId: number) {
  loading.value = true;
  try {
    const { request } = await import('@/service/request');
    
    await request({
      url: `/address-books/${props.addressBookId}/peers/${peerId}`,
      method: 'delete'
    });

    window.$message?.success('Peer deletado com sucesso!');
    await loadPeers(); // Reload peers list
  } catch (error) {
    console.error('Erro ao deletar peer:', error);
    window.$message?.error('Erro ao deletar peer');
  } finally {
    loading.value = false;
  }
}

function handleAddPeer() {
  addPeerModalVisible.value = true;
}

async function handleAddPeerSuccess() {
  await loadPeers(); // Reload peers after adding
}

function handleClose() {
  visible.value = false;
}
</script>

<template>
  <NDrawer v-model:show="visible" :width="1000" placement="right">
    <NDrawerContent :title="`Peers em: ${addressBookName}`" closable>
      <div class="flex-col gap-16px">
        <!-- Barra de estatísticas e ações principais -->
        <NCard :bordered="false" size="small" class="shadow-sm">
          <NSpace justify="space-between" align="center">
            <NText strong>
              Total de Peers: {{ peers.length }}
            </NText>
            
            <!-- Botões de ação em destaque -->
            <NSpace>
              <NButton type="success" size="medium" strong @click="handleAddPeer">
                <template #icon>
                  <icon-mdi-account-plus class="text-icon" />
                </template>
                Adicionar Peer Manualmente
              </NButton>
              
              <NPopconfirm @positive-click="handleImportDevices">
                <template #trigger>
                  <NButton type="primary" size="medium" :loading="loading">
                    <template #icon>
                      <icon-mdi-cloud-download class="text-icon" />
                    </template>
                    Importar Devices Online
                  </NButton>
                </template>
                Importar todos os dispositivos online como peers? Dispositivos já existentes serão ignorados.
              </NPopconfirm>
            </NSpace>
          </NSpace>
        </NCard>

        <!-- Tabela de peers -->
        <NDataTable
          :columns="columns"
          :data="peers"
          :loading="loading"
          size="small"
          :scroll-x="1200"
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

  <!-- Add Peer Modal -->
  <AddPeerModal
    v-model:visible="addPeerModalVisible"
    :address-book-id="addressBookId"
    :address-book-name="addressBookName"
    @success="handleAddPeerSuccess"
  />
</template>

<style scoped></style>
