<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <NCard title="DocHelp" :bordered="false" size="small" class="sm:flex-1-hidden card-wrapper">
      <template #header-extra>
        <NButton type="primary" @click="openCreateTicket">
          <template #icon>
            <icon-mdi-alert-circle-outline />
          </template>
          {{ $t('page.dochelp.reportProblem') }}
        </NButton>
      </template>

      <NTabs v-model:value="activeTab" type="line" animated>
        <!-- Knowledge Base Tab -->
        <NTabPane name="kb" :tab="$t('route.dochelp_kb')">
          <div class="flex-col gap-16px">
            <!-- Search and Filters -->
            <NSpace>
              <NInput
                v-model:value="kbSearchQuery"
                :placeholder="$t('page.dochelp.searchPlaceholder')"
                clearable
                @input="handleKBSearch"
              >
                <template #prefix>
                  <icon-mdi-magnify />
                </template>
              </NInput>

              <NSelect
                v-model:value="selectedCategory"
                :options="categoryOptions"
                :placeholder="$t('page.dochelp.selectCategory')"
                clearable
                style="width: 200px"
                @update:value="handleKBSearch"
              />

              <NCheckbox v-model:checked="showPinnedOnly" @update:checked="handleKBSearch">
                {{ $t('page.dochelp.pinnedOnly') }}
              </NCheckbox>
            </NSpace>

            <!-- Articles List -->
            <NDataTable
              :columns="kbColumns"
              :data="kbData"
              :loading="kbLoading"
              :pagination="kbPagination"
              :scroll-x="1100"
              @update:page="handleKBPageChange"
            />
          </div>
        </NTabPane>

        <!-- Admin Tab (Categories & Articles Management) -->
        <NTabPane v-if="canManageTickets" name="admin" tab="Admin">
          <NCard :bordered="false" class="h-full">
            <NScrollbar style="max-height: calc(100vh - 280px)">
              <NSpace vertical size="large">
                <!-- Categories Section -->
                <div>
                  <NSpace justify="space-between" style="margin-bottom: 16px">
                    <h3>Categorias</h3>
                    <NButton type="primary" @click="openCategoryModal">
                      <template #icon>
                        <icon-mdi-plus />
                      </template>
                      Nova Categoria
                    </NButton>
                  </NSpace>
                  
                  <NDataTable
                    :columns="categoryAdminColumns"
                    :data="categories"
                    :pagination="false"
                    size="small"
                  />
                </div>

                <NDivider />

                <!-- Articles Section -->
                <div>
                  <NSpace justify="space-between" style="margin-bottom: 16px">
                    <h3>Artigos</h3>
                    <NButton type="primary" @click="openArticleModal">
                      <template #icon>
                        <icon-mdi-plus />
                      </template>
                      Novo Artigo
                    </NButton>
                  </NSpace>
                  
                  <NDataTable
                    :columns="articleAdminColumns"
                    :data="kbData"
                    :loading="kbLoading"
                    :pagination="kbPagination"
                    size="small"
                  />
                </div>
              </NSpace>
            </NScrollbar>
          </NCard>
        </NTabPane>

        <!-- My Tickets Tab -->
        <NTabPane name="mytickets" :tab="$t('route.dochelp_tickets')">
          <div class="flex-col gap-16px">
            <NSpace>
              <NSelect
                v-model:value="myTicketStatusFilter"
                :options="statusOptions"
                :placeholder="$t('page.dochelp.filterByStatus')"
                clearable
                style="width: 200px"
                @update:value="handleMyTicketsSearch"
              />

              <NSelect
                v-model:value="myTicketPriorityFilter"
                :options="priorityOptions"
                :placeholder="$t('page.dochelp.filterByPriority')"
                clearable
                style="width: 200px"
                @update:value="handleMyTicketsSearch"
              />
            </NSpace>

            <NDataTable
              :columns="ticketColumns"
              :data="myTicketsData"
              :loading="myTicketsLoading"
              :pagination="myTicketsPagination"
              :scroll-x="1300"
              @update:page="handleMyTicketsPageChange"
            />
          </div>
        </NTabPane>

        <!-- All Tickets Tab (Support Staff Only) -->
        <NTabPane v-if="canManageTickets" name="alltickets" :tab="$t('route.dochelp_alltickets')">
          <div class="flex-col gap-16px">
            <NSpace>
              <NSelect
                v-model:value="allTicketStatusFilter"
                :options="statusOptions"
                :placeholder="$t('page.dochelp.filterByStatus')"
                clearable
                style="width: 200px"
                @update:value="handleAllTicketsSearch"
              />

              <NSelect
                v-model:value="allTicketPriorityFilter"
                :options="priorityOptions"
                :placeholder="$t('page.dochelp.filterByPriority')"
                clearable
                style="width: 200px"
                @update:value="handleAllTicketsSearch"
              />
            </NSpace>

            <NDataTable
              :columns="ticketColumns"
              :data="allTicketsData"
              :loading="allTicketsLoading"
              :pagination="allTicketsPagination"
              :scroll-x="1300"
              @update:page="handleAllTicketsPageChange"
            />
          </div>
        </NTabPane>
      </NTabs>
    </NCard>

    <!-- Create Ticket Modal -->
    <NModal v-model:show="showCreateTicketModal" preset="card" title="Report a Problem" style="width: 700px">
      <NForm ref="ticketFormRef" :model="ticketForm" :rules="ticketFormRules" label-placement="left" label-width="120">
        <NFormItem :label="$t('page.dochelp.ticketTitle')" path="title">
          <NInput v-model:value="ticketForm.title" :placeholder="$t('page.dochelp.ticketTitlePlaceholder')" />
        </NFormItem>

        <NFormItem :label="$t('page.dochelp.description')" path="description">
          <NInput
            v-model:value="ticketForm.description"
            type="textarea"
            :placeholder="$t('page.dochelp.descriptionPlaceholder')"
            :rows="6"
          />
        </NFormItem>

        <NFormItem :label="$t('page.dochelp.priority')" path="priority">
          <NSelect v-model:value="ticketForm.priority" :options="priorityOptions" />
        </NFormItem>

        <NFormItem :label="$t('page.dochelp.category')">
          <NSelect v-model:value="ticketForm.category_id" :options="categoryOptions" clearable />
        </NFormItem>

        <NFormItem label="Anexos">
          <NUpload
            v-model:file-list="ticketAttachments"
            :max="5"
            :custom-request="handleUploadFile"
            list-type="text"
          >
            <NButton>📎 Anexar Arquivos (máx. 10MB)</NButton>
          </NUpload>
        </NFormItem>
      </NForm>

      <template #footer>
        <NSpace justify="end">
          <NButton @click="showCreateTicketModal = false">{{ $t('common.cancel') }}</NButton>
          <NButton type="primary" :loading="submitting" @click="handleCreateTicket">{{
            $t('common.confirm')
          }}</NButton>
        </NSpace>
      </template>
    </NModal>

    <!-- Category Modal -->
    <NModal v-model:show="showCategoryModal" preset="card" :title="categoryModalTitle" style="width: 500px">
      <NForm ref="categoryFormRef" :model="categoryForm" label-placement="left" label-width="100">
        <NFormItem label="Nome" path="name">
          <NInput v-model:value="categoryForm.name" placeholder="Nome da categoria" />
        </NFormItem>

        <NFormItem label="Ícone" path="icon">
          <NInput v-model:value="categoryForm.icon" placeholder="Ex: mdi:help-circle" />
        </NFormItem>

        <NFormItem label="Ordem" path="order">
          <NInputNumber v-model:value="categoryForm.order" :min="0" style="width: 100%" />
        </NFormItem>
      </NForm>

      <template #footer>
        <NSpace justify="end">
          <NButton @click="showCategoryModal = false">{{ $t('common.cancel') }}</NButton>
          <NButton type="primary" :loading="submitting" @click="handleSaveCategory">
            {{ $t('common.confirm') }}
          </NButton>
        </NSpace>
      </template>
    </NModal>

    <!-- Article Modal -->
    <NModal v-model:show="showArticleModal" preset="card" :title="articleModalTitle" style="width: 900px">
      <NForm ref="articleFormRef" :model="articleForm" label-placement="left" label-width="120">
        <NFormItem label="Título" path="title">
          <NInput v-model:value="articleForm.title" placeholder="Título do artigo" />
        </NFormItem>

        <NFormItem label="Categoria" path="category_id">
          <NSelect v-model:value="articleForm.category_id" :options="categoryOptions" />
        </NFormItem>

        <NFormItem label="Conteúdo" path="content">
          <NInput
            v-model:value="articleForm.content"
            type="textarea"
            placeholder="Conteúdo do artigo (suporta Markdown)"
            :rows="15"
          />
        </NFormItem>

        <NFormItem label="Tags" path="tags">
          <NInput v-model:value="articleForm.tags" placeholder="Separadas por vírgula: tag1, tag2" />
        </NFormItem>

        <NFormItem label="Fixar no topo">
          <NCheckbox v-model:checked="articleForm.is_pinned">Destacar este artigo</NCheckbox>
        </NFormItem>
      </NForm>

      <template #footer>
        <NSpace justify="end">
          <NButton @click="showArticleModal = false">{{ $t('common.cancel') }}</NButton>
          <NButton type="primary" :loading="submitting" @click="handleSaveArticle">
            {{ $t('common.confirm') }}
          </NButton>
        </NSpace>
      </template>
    </NModal>

    <!-- Ticket Details Modal -->
    <NModal 
      v-model:show="showTicketModal" 
      preset="card" 
      :title="ticketModalTitle"
      style="width: 900px"
      :segmented="{ content: true }"
    >
      <NScrollbar style="max-height: 600px">
        <NSpace vertical size="large">
          <!-- Ticket Info -->
          <div v-if="currentTicket">
            <NDescriptions bordered :column="2" size="small">
              <NDescriptionsItem label="ID">
                <NTag type="info">#{{ currentTicket.id }}</NTag>
              </NDescriptionsItem>
              
              <NDescriptionsItem label="Status">
                <NTag :type="getStatusType(currentTicket.status)">
                  {{ getStatusLabel(currentTicket.status) }}
                </NTag>
              </NDescriptionsItem>
              
              <NDescriptionsItem label="Prioridade">
                <NTag :type="getPriorityType(currentTicket.priority)">
                  {{ getPriorityLabel(currentTicket.priority) }}
                </NTag>
              </NDescriptionsItem>
              
              <NDescriptionsItem label="Criado por">
                {{ currentTicket.creator_name || 'Desconhecido' }}
              </NDescriptionsItem>
              
              <NDescriptionsItem label="Criado em" :span="2">
                {{ formatDate(currentTicket.created_at) }}
              </NDescriptionsItem>
              
              <NDescriptionsItem label="Título" :span="2">
                <strong>{{ currentTicket.title }}</strong>
              </NDescriptionsItem>
              
              <NDescriptionsItem label="Descrição" :span="2">
                <div style="white-space: pre-wrap">{{ currentTicket.description }}</div>
              </NDescriptionsItem>
            </NDescriptions>

            <!-- Change Status (SUPPORT_N2+ only) -->
            <div v-if="canManageTickets" style="margin-top: 16px">
              <NSpace>
                <NSelect
                  v-model:value="ticketStatusUpdate"
                  :options="statusOptions"
                  placeholder="Alterar status"
                  style="width: 200px"
                />
                <NButton type="primary" :loading="updatingTicket" @click="handleUpdateTicketStatus">
                  Atualizar Status
                </NButton>
              </NSpace>
            </div>
          </div>

          <NDivider />

          <!-- Comments Section -->
          <div>
            <NSpace justify="space-between" style="margin-bottom: 16px">
              <h3>Comentários ({{ ticketComments.length }})</h3>
            </NSpace>

            <!-- Comments List -->
            <NSpace vertical size="large" style="width: 100%">
              <NCard
                v-for="comment in ticketComments"
                :key="comment.id"
                size="small"
                :bordered="true"
              >
                <template #header>
                  <NSpace align="center">
                    <NAvatar round size="small">
                      {{ comment.username?.charAt(0).toUpperCase() }}
                    </NAvatar>
                    <span><strong>{{ comment.username }}</strong></span>
                    <NTag v-if="comment.is_internal" type="warning" size="small">Interno</NTag>
                  </NSpace>
                </template>
                
                <div style="white-space: pre-wrap">{{ comment.comment }}</div>
                
                <!-- Attachments -->
                <div v-if="comment.attachments" style="margin-top: 12px">
                  <NDivider style="margin: 8px 0">Anexos</NDivider>
                  <NSpace>
                    <a
                      v-for="(url, idx) in JSON.parse(comment.attachments)"
                      :key="idx"
                      :href="url"
                      target="_blank"
                      style="color: #2080f0"
                    >
                      📎 {{ url.split('/').pop() }}
                    </a>
                  </NSpace>
                </div>
                
                <template #footer>
                  <NText depth="3" style="font-size: 12px">
                    {{ formatDate(comment.created_at) }}
                  </NText>
                </template>
              </NCard>

              <NEmpty v-if="ticketComments.length === 0" description="Nenhum comentário ainda" />
            </NSpace>

            <!-- Add Comment Form -->
            <NCard title="Adicionar Comentário" size="small" style="margin-top: 16px">
              <NForm>
                <NFormItem>
                  <NInput
                    v-model:value="newComment"
                    type="textarea"
                    placeholder="Escreva seu comentário..."
                    :rows="4"
                  />
                </NFormItem>
                
                <NFormItem>
                  <NUpload
                    v-model:file-list="commentAttachments"
                    :max="5"
                    :custom-request="handleUploadFile"
                    list-type="text"
                  >
                    <NButton size="small">📎 Anexar Arquivos</NButton>
                  </NUpload>
                </NFormItem>
                
                <NFormItem v-if="canManageTickets">
                  <NCheckbox v-model:checked="commentIsInternal">
                    Comentário interno (visível apenas para staff)
                  </NCheckbox>
                </NFormItem>
              </NForm>
              
              <template #footer>
                <NSpace justify="end">
                  <NButton type="primary" :loading="addingComment" @click="handleAddComment">
                    <template #icon>
                      <icon-mdi-send />
                    </template>
                    Enviar Comentário
                  </NButton>
                </NSpace>
              </template>
            </NCard>
          </div>
        </NSpace>
      </NScrollbar>
    </NModal>
  </div>
</template>

<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue';
import type { DataTableColumns, FormInst, FormRules, UploadCustomRequestOptions, UploadFileInfo } from 'naive-ui';
import { NBadge, NButton, NSpace, NTag, NCheckbox, NInputNumber, NPopconfirm, NDivider, NScrollbar, NDescriptions, NDescriptionsItem, NAvatar, NEmpty, NText, NUpload } from 'naive-ui';
import { fetchKBArticles, fetchKBCategories, fetchTickets, createTicket, uploadDocHelpFile } from '@/service/api';
import { usePermission } from '@/hooks/business/auth';
import { $t } from '@/locales';

// Helper para tradução sem erro de tipo
const t = (key: string) => ($t as any)(key);

const { hasPermission, USER_ROLES } = usePermission();
const canManageTickets = computed(() => hasPermission(USER_ROLES.SUPPORT_N2));

// Active Tab
const activeTab = ref('kb');

// ===== KNOWLEDGE BASE =====
const kbLoading = ref(false);
const kbData = ref<Api.DocHelp.Article[]>([]);
const kbSearchQuery = ref('');
const selectedCategory = ref<number | null>(null);
const showPinnedOnly = ref(false);
const kbPagination = reactive({
  page: 1,
  pageSize: 20,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100]
});

const kbColumns = computed<DataTableColumns<Api.DocHelp.Article>>(() => [
  {
    title: t('page.dochelp.title'),
    key: 'title',
    width: 300,
    ellipsis: { tooltip: true },
    render: row => {
      return row.is_pinned
        ? h('div', { class: 'flex items-center gap-2' }, [
            h(NTag, { type: 'warning', size: 'small' }, { default: () => 'Pinned' }),
            row.title
          ])
        : row.title;
    }
  },
  {
    title: t('page.dochelp.category'),
    key: 'category_name',
    width: 150
  },
  {
    title: t('page.dochelp.views'),
    key: 'views',
    width: 100,
    render: row => {
      return h(NBadge, { value: row.views, color: '#2080f0' });
    }
  },
  {
    title: t('page.dochelp.author'),
    key: 'author_name',
    width: 120
  },
  {
    title: t('page.dochelp.createdAt'),
    key: 'created_at',
    width: 180
  },
  {
    title: $t('common.action'),
    key: 'actions',
    width: 120,
    fixed: 'right',
    render: row => {
      return h(
        NButton,
        {
          size: 'small',
          type: 'primary',
          onClick: () => handleViewArticle(row.id!)
        },
        { default: () => t('common.view') }
      );
    }
  }
]);

async function loadKBArticles() {
  kbLoading.value = true;
  try {
    const params: Api.DocHelp.ArticleParams = {
      current: kbPagination.page,
      size: kbPagination.pageSize
    };

    if (kbSearchQuery.value) {
      params.search = kbSearchQuery.value;
    }

    if (selectedCategory.value) {
      params.category_id = selectedCategory.value;
    }

    if (showPinnedOnly.value) {
      params.pinned = 'true';
    }

    const { data } = await fetchKBArticles(params);
    kbData.value = (data as any).records;
    kbPagination.itemCount = (data as any).total;
  } catch (error) {
    window.$message?.error('Failed to load articles');
  } finally {
    kbLoading.value = false;
  }
}

function handleKBSearch() {
  kbPagination.page = 1;
  loadKBArticles();
}

function handleKBPageChange(page: number) {
  kbPagination.page = page;
  loadKBArticles();
}

function handleViewArticle(id: number) {
  // TODO: Navigate to article detail page or open modal
  console.log('View article:', id);
}

// ===== TICKETS =====
const myTicketsLoading = ref(false);
const allTicketsLoading = ref(false);
const myTicketsData = ref<Api.DocHelp.Ticket[]>([]);
const allTicketsData = ref<Api.DocHelp.Ticket[]>([]);
const myTicketStatusFilter = ref<number | null>(null);
const myTicketPriorityFilter = ref<number | null>(null);
const allTicketStatusFilter = ref<number | null>(null);
const allTicketPriorityFilter = ref<number | null>(null);

const myTicketsPagination = reactive({
  page: 1,
  pageSize: 20,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100]
});

const allTicketsPagination = reactive({
  page: 1,
  pageSize: 20,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100]
});

const statusOptions = computed(() => [
  { label: t('page.dochelp.statusOpen'), value: 1 },
  { label: t('page.dochelp.statusInAnalysis'), value: 2 },
  { label: t('page.dochelp.statusResolved'), value: 3 }
]);

const priorityOptions = computed(() => [
  { label: t('page.dochelp.priorityLow'), value: 1 },
  { label: t('page.dochelp.priorityMedium'), value: 2 },
  { label: t('page.dochelp.priorityHigh'), value: 3 },
  { label: t('page.dochelp.priorityCritical'), value: 4 }
]);

const ticketColumns = computed<DataTableColumns<Api.DocHelp.Ticket>>(() => [
  {
    title: 'ID',
    key: 'id',
    width: 80
  },
  {
    title: t('page.dochelp.ticketTitle'),
    key: 'title',
    width: 250,
    ellipsis: { tooltip: true }
  },
  {
    title: t('page.dochelp.status'),
    key: 'status',
    width: 130,
    render: row => {
      const colors: Record<number, 'info' | 'warning' | 'success'> = { 1: 'info', 2: 'warning', 3: 'success' };
      const labels = {
        1: t('page.dochelp.statusOpen'),
        2: t('page.dochelp.statusInAnalysis'),
        3: t('page.dochelp.statusResolved')
      };
      return h(NTag, { type: colors[row.status] || 'info' }, { default: () => labels[row.status as keyof typeof labels] });
    }
  },
  {
    title: t('page.dochelp.priority'),
    key: 'priority',
    width: 120,
    render: row => {
      console.log('[Priority Debug] Row:', row, 'Priority value:', row.priority, 'Type:', typeof row.priority);
      const colors: Record<number, 'default' | 'info' | 'warning' | 'error'> = { 1: 'default', 2: 'info', 3: 'warning', 4: 'error' };
      const labels = {
        1: t('page.dochelp.priorityLow'),
        2: t('page.dochelp.priorityMedium'),
        3: t('page.dochelp.priorityHigh'),
        4: t('page.dochelp.priorityCritical')
      };
      // Se priority for 0 ou inválido, não retorna nada (mostra blank)
      if (!row.priority || row.priority < 1 || row.priority > 4) {
        return h('span', { style: { color: '#ff0000' } }, `[Invalid: ${row.priority}]`);
      }
      return h(NTag, { type: colors[row.priority] || 'default' }, { default: () => labels[row.priority as keyof typeof labels] });
    }
  },
  {
    title: t('page.dochelp.creator'),
    key: 'creator_name',
    width: 120
  },
  {
    title: t('page.dochelp.comments'),
    key: 'comment_count',
    width: 100,
    render: row => {
      return h(NBadge, { value: row.comment_count, color: '#18a058' });
    }
  },
  {
    title: t('page.dochelp.createdAt'),
    key: 'created_at',
    width: 180
  },
  {
    title: $t('common.action'),
    key: 'actions',
    width: 120,
    fixed: 'right',
    render: row => {
      return h(
        NButton,
        {
          size: 'small',
          type: 'primary',
          onClick: () => handleViewTicket(row.id!)
        },
        { default: () => t('common.view') }
      );
    }
  }
]);

async function loadMyTickets() {
  myTicketsLoading.value = true;
  try {
    const params: Api.DocHelp.TicketParams = {
      current: myTicketsPagination.page,
      size: myTicketsPagination.pageSize
    };

    if (myTicketStatusFilter.value) {
      params.status = myTicketStatusFilter.value;
    }

    if (myTicketPriorityFilter.value) {
      params.priority = myTicketPriorityFilter.value;
    }

    const { data } = await fetchTickets(params);
    myTicketsData.value = (data as any).records;
    myTicketsPagination.itemCount = (data as any).total;
  } catch (error) {
    window.$message?.error('Failed to load tickets');
  } finally {
    myTicketsLoading.value = false;
  }
}

async function loadAllTickets() {
  if (!canManageTickets.value) return;

  allTicketsLoading.value = true;
  try {
    const params: Api.DocHelp.TicketParams = {
      current: allTicketsPagination.page,
      size: allTicketsPagination.pageSize
    };

    if (allTicketStatusFilter.value) {
      params.status = allTicketStatusFilter.value;
    }

    if (allTicketPriorityFilter.value) {
      params.priority = allTicketPriorityFilter.value;
    }

    const { data } = await fetchTickets(params);
    allTicketsData.value = (data as any).records;
    allTicketsPagination.itemCount = (data as any).total;
  } catch (error) {
    window.$message?.error('Failed to load tickets');
  } finally {
    allTicketsLoading.value = false;
  }
}

function handleMyTicketsSearch() {
  myTicketsPagination.page = 1;
  loadMyTickets();
}

function handleAllTicketsSearch() {
  allTicketsPagination.page = 1;
  loadAllTickets();
}

function handleMyTicketsPageChange(page: number) {
  myTicketsPagination.page = page;
  loadMyTickets();
}

function handleAllTicketsPageChange(page: number) {
  allTicketsPagination.page = page;
  loadAllTickets();
}

// ===== TICKET DETAILS & COMMENTS =====
const showTicketModal = ref(false);
const currentTicket = ref<Api.DocHelp.Ticket | null>(null);
const ticketComments = ref<Api.DocHelp.Comment[]>([]);
const newComment = ref('');
const commentIsInternal = ref(false);
const addingComment = ref(false);
const updatingTicket = ref(false);
const ticketStatusUpdate = ref<number | null>(null);

const ticketModalTitle = computed(() => {
  return currentTicket.value ? `Ticket #${currentTicket.value.id} - ${currentTicket.value.title}` : 'Ticket';
});

function getStatusType(status: number): 'info' | 'warning' | 'success' {
  const types: Record<number, 'info' | 'warning' | 'success'> = { 1: 'info', 2: 'warning', 3: 'success' };
  return types[status] || 'info';
}

function getStatusLabel(status: number): string {
  return statusOptions.value.find(opt => opt.value === status)?.label || 'Desconhecido';
}

function getPriorityType(priority: number): 'default' | 'info' | 'warning' | 'error' {
  const types: Record<number, 'default' | 'info' | 'warning' | 'error'> = { 1: 'default', 2: 'info', 3: 'warning', 4: 'error' };
  return types[priority] || 'default';
}

function getPriorityLabel(priority: number): string {
  return priorityOptions.value.find(opt => opt.value === priority)?.label || 'Desconhecido';
}

function formatDate(dateString: string | undefined): string {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleString('pt-BR');
}

async function handleViewTicket(id: number) {
  showTicketModal.value = true;
  
  try {
    // Load ticket details
    const { fetchTicket, fetchTicketComments } = await import('@/service/api');
    const { data: ticketData } = await fetchTicket(id);
    currentTicket.value = ticketData as any;
    ticketStatusUpdate.value = (ticketData as any)?.status;
    
    // Load comments
    const { data: commentsData } = await fetchTicketComments(id);
    ticketComments.value = (commentsData as any)?.comments || [];
  } catch (error) {
    console.error('Failed to load ticket:', error);
    window.$message?.error('Erro ao carregar ticket');
  }
}

async function handleAddComment() {
  if (!newComment.value.trim() || !currentTicket.value) {
    window.$message?.warning('Digite um comentário');
    return;
  }

  addingComment.value = true;
  try {
    const { addTicketComment, fetchTicketComments } = await import('@/service/api');
    
    // Include uploaded files URLs as JSON
    const attachments = uploadedFiles.value.length > 0 ? JSON.stringify(uploadedFiles.value) : undefined;
    
    await addTicketComment(currentTicket.value.id!, {
      comment: newComment.value,
      is_internal: commentIsInternal.value,
      attachments: attachments as any
    });
    
    window.$message?.success('Comentário adicionado');
    
    // Reload only comments without closing modal or reloading entire ticket
    const { data: commentsData } = await fetchTicketComments(currentTicket.value.id!);
    ticketComments.value = (commentsData as any)?.comments || [];
    
    // Clear form
    newComment.value = '';
    commentIsInternal.value = false;
    commentAttachments.value = [];
    uploadedFiles.value = [];
  } catch (error) {
    console.error('Failed to add comment:', error);
    window.$message?.error('Erro ao adicionar comentário');
  } finally {
    addingComment.value = false;
  }
}

async function handleUpdateTicketStatus() {
  if (!currentTicket.value || ticketStatusUpdate.value === null) return;

  updatingTicket.value = true;
  try {
    const { updateTicket } = await import('@/service/api');
    await updateTicket(currentTicket.value.id!, {
      status: ticketStatusUpdate.value
    });
    
    window.$message?.success('Status atualizado');
    currentTicket.value.status = ticketStatusUpdate.value;
    
    // Reload ticket lists (commented to prevent page refresh interruption)
    // Uncomment if you want to auto-reload lists after status update
    // loadMyTickets();
    // if (canManageTickets.value) {
    //   loadAllTickets();
    // }
  } catch (error) {
    console.error('Failed to update ticket:', error);
    window.$message?.error('Erro ao atualizar status');
  } finally {
    updatingTicket.value = false;
  }
}

// ===== CREATE TICKET =====
const showCreateTicketModal = ref(false);
const ticketFormRef = ref<FormInst | null>(null);
const submitting = ref(false);
const ticketForm = reactive<Api.DocHelp.TicketCreate>({
  title: '',
  description: '',
  priority: 2
});

// File attachments
const ticketAttachments = ref<UploadFileInfo[]>([]);
const commentAttachments = ref<UploadFileInfo[]>([]);
const uploadedFiles = ref<string[]>([]); // URLs dos arquivos enviados

const ticketFormRules: FormRules = {
  title: [{ required: true, message: 'Title is required', trigger: 'blur' }],
  description: [{ required: true, message: 'Description is required', trigger: 'blur' }],
  priority: [{ required: true, type: 'number', message: 'Priority is required', trigger: 'change' }]
};

function openCreateTicket() {
  showCreateTicketModal.value = true;
}

// Custom file upload handler
async function handleUploadFile({ file, onFinish, onError }: UploadCustomRequestOptions) {
  try {
    const { data } = await uploadDocHelpFile(file.file as File);
    if (data) {
      uploadedFiles.value.push(data.url);
      window.$message?.success(`Arquivo ${file.name} enviado com sucesso`);
    }
    onFinish();
  } catch (error) {
    console.error('Upload error:', error);
    window.$message?.error(`Erro ao enviar ${file.name}`);
    onError();
  }
}

async function handleCreateTicket() {
  await ticketFormRef.value?.validate(async errors => {
    if (!errors) {
      submitting.value = true;
      try {
        await createTicket(ticketForm);
        window.$message?.success('Ticket created successfully');
        showCreateTicketModal.value = false;
        // Reset form
        ticketForm.title = '';
        ticketForm.description = '';
        ticketForm.priority = 2;
        ticketForm.category_id = undefined;
        // Reload tickets
        loadMyTickets();
      } catch (error) {
        window.$message?.error('Failed to create ticket');
      } finally {
        submitting.value = false;
      }
    }
  });
}

// ===== CATEGORIES =====
const categories = ref<Api.DocHelp.Category[]>([]);
const categoryOptions = computed(() => {
  return categories.value.map(cat => ({ label: cat.name, value: cat.id! }));
});

async function loadCategories() {
  try {
    const response = await fetchKBCategories();
    console.log('Categories response:', response);
    
    if (response && response.data) {
      categories.value = (response.data as any).categories || [];
      console.log('Loaded categories:', categories.value);
    } else {
      console.error('Invalid response format:', response);
      categories.value = [];
    }
  } catch (error) {
    console.error('Failed to load categories:', error);
    window.$message?.error('Erro ao carregar categorias');
    categories.value = [];
  }
}

// ===== CATEGORY MANAGEMENT =====
const showCategoryModal = ref(false);
const categoryFormRef = ref<FormInst | null>(null);
const editingCategoryId = ref<number | null>(null);

const categoryForm = reactive({
  name: '',
  icon: 'mdi:help-circle',
  order: 0
});

const categoryModalTitle = computed(() => {
  return editingCategoryId.value ? 'Editar Categoria' : 'Nova Categoria';
});

const categoryAdminColumns = computed<DataTableColumns<Api.DocHelp.Category>>(() => [
  {
    title: 'ID',
    key: 'id',
    width: 80
  },
  {
    title: 'Nome',
    key: 'name',
    width: 200
  },
  {
    title: 'Ícone',
    key: 'icon',
    width: 150
  },
  {
    title: 'Ordem',
    key: 'order',
    width: 100
  },
  {
    title: 'Ações',
    key: 'actions',
    width: 180,
    render: row => {
      // Only show actions if ID exists
      if (!row.id) {
        return h('span', { style: 'color: #999' }, 'ID inválido');
      }
      
      return h(NSpace, {}, {
        default: () => [
          h(NButton, {
            size: 'small',
            onClick: () => openEditCategory(row)
          }, { default: () => 'Editar' }),
          h(NButton, {
            size: 'small',
            type: 'error',
            onClick: () => handleDeleteCategory(row.id!)
          }, { default: () => 'Excluir' })
        ]
      });
    }
  }
]);

function openCategoryModal() {
  editingCategoryId.value = null;
  categoryForm.name = '';
  categoryForm.icon = 'mdi:help-circle';
  categoryForm.order = 0;
  showCategoryModal.value = true;
}

function openEditCategory(category: Api.DocHelp.Category) {
  editingCategoryId.value = category.id!;
  categoryForm.name = category.name;
  categoryForm.icon = category.icon || 'mdi:help-circle';
  categoryForm.order = category.order || 0;
  showCategoryModal.value = true;
}

async function handleSaveCategory() {
  submitting.value = true;
  try {
    const { createKBCategory, updateKBCategory } = await import('@/service/api');
    
    if (editingCategoryId.value) {
      await updateKBCategory(editingCategoryId.value, categoryForm);
      window.$message?.success('Categoria atualizada com sucesso');
    } else {
      await createKBCategory(categoryForm);
      window.$message?.success('Categoria criada com sucesso');
    }
    
    showCategoryModal.value = false;
    await loadCategories();
  } catch (error) {
    window.$message?.error('Erro ao salvar categoria');
  } finally {
    submitting.value = false;
  }
}

async function handleDeleteCategory(id: number) {
  // Validate ID
  if (!id || id === undefined) {
    window.$message?.error('ID da categoria inválido');
    return;
  }
  
  if (!confirm('Tem certeza que deseja excluir esta categoria?')) return;
  
  try {
    const { deleteKBCategory } = await import('@/service/api');
    await deleteKBCategory(id);
    window.$message?.success('Categoria excluída com sucesso');
    await loadCategories();
  } catch (error: any) {
    console.error('Error deleting category:', error);
    window.$message?.error(error?.message || 'Erro ao excluir categoria');
  }
}

// ===== ARTICLE MANAGEMENT =====
const showArticleModal = ref(false);
const articleFormRef = ref<FormInst | null>(null);
const editingArticleId = ref<number | null>(null);

const articleForm = reactive({
  title: '',
  category_id: null as number | null,
  content: '',
  tags: '',
  is_pinned: false
});

const articleModalTitle = computed(() => {
  return editingArticleId.value ? 'Editar Artigo' : 'Novo Artigo';
});

const articleAdminColumns = computed<DataTableColumns<Api.DocHelp.Article>>(() => [
  {
    title: 'ID',
    key: 'id',
    width: 80
  },
  {
    title: 'Título',
    key: 'title',
    width: 300,
    ellipsis: { tooltip: true }
  },
  {
    title: 'Categoria',
    key: 'category_name',
    width: 150
  },
  {
    title: 'Visualizações',
    key: 'views',
    width: 120
  },
  {
    title: 'Fixado',
    key: 'is_pinned',
    width: 100,
    render: row => {
      return row.is_pinned ? h(NTag, { type: 'warning' }, { default: () => 'Sim' }) : 'Não';
    }
  },
  {
    title: 'Ações',
    key: 'actions',
    width: 180,
    render: row => {
      return h(NSpace, {}, {
        default: () => [
          h(NButton, {
            size: 'small',
            onClick: () => openEditArticle(row)
          }, { default: () => 'Editar' }),
          h(NButton, {
            size: 'small',
            type: 'error',
            onClick: () => handleDeleteArticle(row.id!)
          }, { default: () => 'Excluir' })
        ]
      });
    }
  }
]);

function openArticleModal() {
  editingArticleId.value = null;
  articleForm.title = '';
  articleForm.category_id = null;
  articleForm.content = '';
  articleForm.tags = '';
  articleForm.is_pinned = false;
  showArticleModal.value = true;
}

function openEditArticle(article: Api.DocHelp.Article) {
  editingArticleId.value = article.id!;
  articleForm.title = article.title;
  articleForm.category_id = article.category_id;
  articleForm.content = article.content;
  articleForm.tags = article.tags || '';
  articleForm.is_pinned = article.is_pinned || false;
  showArticleModal.value = true;
}

async function handleSaveArticle() {
  submitting.value = true;
  try {
    const { createKBArticle, updateKBArticle } = await import('@/service/api');
    
    if (editingArticleId.value) {
      await updateKBArticle(editingArticleId.value, articleForm as any);
      window.$message?.success('Artigo atualizado com sucesso');
    } else {
      await createKBArticle(articleForm as any);
      window.$message?.success('Artigo criado com sucesso');
    }
    
    showArticleModal.value = false;
    await loadKBArticles();
  } catch (error) {
    window.$message?.error('Erro ao salvar artigo');
  } finally {
    submitting.value = false;
  }
}

async function handleDeleteArticle(id: number) {
  if (!confirm('Tem certeza que deseja excluir este artigo?')) return;
  
  try {
    const { deleteKBArticle } = await import('@/service/api');
    await deleteKBArticle(id);
    window.$message?.success('Artigo excluído com sucesso');
    await loadKBArticles();
  } catch (error) {
    window.$message?.error('Erro ao excluir artigo');
  }
}

// ===== LIFECYCLE =====
onMounted(() => {
  loadCategories();
  loadKBArticles();
  loadMyTickets();
  if (canManageTickets.value) {
    loadAllTickets();
  }
});
</script>

<style scoped></style>
