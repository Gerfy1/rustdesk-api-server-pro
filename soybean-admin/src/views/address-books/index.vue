<script setup lang="tsx">
import { ref } from 'vue';
import { NButton, NTag, NPopconfirm, NSpace } from 'naive-ui';
import { fetchAddressBooksList, deleteAddressBook } from '@/service/api/address-books';
import { $t } from '@/locales';
import { useAppStore } from '@/store/modules/app';
import { useTable } from '@/hooks/common/table';
import TableHeader from './components/table-header.vue';
import SearchForm from './components/search.vue';
import AddressBookModal from './components/address-book-modal.vue';
import PeersDrawer from './components/peers-drawer.vue';

const appStore = useAppStore();

// Modal state
const modalVisible = ref(false);
const modalMode = ref<'add' | 'edit'>('add');
const currentRecord = ref<Api.AddressBooks.AddressBook | null>(null);

// Peers drawer state
const peersDrawerVisible = ref(false);
const currentAddressBookId = ref<number>(0);
const currentAddressBookName = ref<string>('');

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
  apiFn: fetchAddressBooksList,
  showTotal: true,
  apiParams: {
    current: 1,
    size: 10,
    name: null,
    owner: null
  },
  columns: () => [
    {
      key: 'id',
      title: 'ID',
      align: 'center',
      width: 60
    },
    {
      key: 'name',
      title: 'Nome',
      align: 'center',
      width: 200
    },
    {
      key: 'owner',
      title: 'Proprietário',
      align: 'center',
      width: 150
    },
    {
      key: 'peer_count',
      title: 'Peers',
      align: 'center',
      width: 100,
      render: (row: any) => {
        return (
          <NTag type="info" size="small">
            {row.peer_count} peers
          </NTag>
        );
      }
    },
    {
      key: 'shared',
      title: 'Compartilhado',
      align: 'center',
      width: 120,
      render: (row: any) => {
        return row.shared ? (
          <NTag type="success" size="small">✓ Sim</NTag>
        ) : (
          <NTag type="default" size="small">✗ Não</NTag>
        );
      }
    },
    {
      key: 'rule',
      title: 'Permissão',
      align: 'center',
      width: 120,
      render: (row: any) => {
        const ruleMap: Record<number, { text: string; type: any }> = {
          1: { text: 'Leitura', type: 'default' },
          2: { text: 'Ler/Escrever', type: 'warning' },
          3: { text: 'Total', type: 'success' }
        };
        const rule = ruleMap[row.rule] || { text: 'Desconhecido', type: 'error' };
        return <NTag type={rule.type} size="small">{rule.text}</NTag>;
      }
    },
    {
      key: 'note',
      title: 'Nota',
      align: 'center',
      width: 200,
      ellipsis: {
        tooltip: true
      }
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
    },
    {
      key: 'actions',
      title: 'Ações',
      align: 'center',
      width: 280,
      fixed: 'right',
      render: (row: any) => {
        return (
          <NSpace justify="center">
            <NButton size="small" type="info" onClick={() => viewPeers(row)}>
              Ver Peers
            </NButton>
            <NButton size="small" type="success" onClick={() => handleEdit(row)}>
              Editar
            </NButton>
            <NPopconfirm onPositiveClick={() => handleDelete(row.id)}>
              {{
                default: () => 'Tem certeza que deseja deletar este Address Book?',
                trigger: () => <NButton size="small" type="error">Deletar</NButton>
              }}
            </NPopconfirm>
          </NSpace>
        );
      }
    }
  ]
});

function handleAdd() {
  modalMode.value = 'add';
  currentRecord.value = null;
  modalVisible.value = true;
}

function handleEdit(record: Api.AddressBooks.AddressBook) {
  modalMode.value = 'edit';
  currentRecord.value = record;
  modalVisible.value = true;
}

async function handleDelete(id: number) {
  try {
    await deleteAddressBook(id);
    window.$message?.success('Address Book deletado com sucesso');
    getData();
  } catch (error) {
    window.$message?.error('Erro ao deletar Address Book');
  }
}

function viewPeers(record: Api.AddressBooks.AddressBook) {
  currentAddressBookId.value = record.id!;
  currentAddressBookName.value = record.name;
  peersDrawerVisible.value = true;
}

function handleModalSuccess() {
  modalVisible.value = false;
  getData();
}
</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <SearchForm v-model:model="searchParams" @reset="resetSearchParams" @search="getDataByPage" />

    <NCard title="Address Books" :bordered="false" size="small" class="sm:flex-1-hidden card-wrapper">
      <template #header-extra>
        <NSpace>
          <NButton type="primary" @click="handleAdd">
            <template #icon>
              <icon-ic-round-plus class="text-icon" />
            </template>
            Novo Address Book
          </NButton>
          <TableHeader v-model:columns="columnChecks" :loading="loading" @refresh="getData" />
        </NSpace>
      </template>
      <NDataTable
        :columns="columns"
        :data="data"
        size="small"
        :flex-height="!appStore.isMobile"
        :scroll-x="1200"
        :loading="loading"
        remote
        :row-key="(row: any) => row.id"
        :pagination="mobilePagination"
        class="sm:h-full"
      />
    </NCard>

    <!-- Modal for Add/Edit -->
    <AddressBookModal
      v-model:visible="modalVisible"
      :mode="modalMode"
      :record="currentRecord"
      @success="handleModalSuccess"
    />

    <!-- Drawer for viewing peers -->
    <PeersDrawer
      v-model:visible="peersDrawerVisible"
      :address-book-id="currentAddressBookId"
      :address-book-name="currentAddressBookName"
    />
  </div>
</template>

<style scoped></style>
