<script setup lang="ts">
import { computed } from 'vue';

defineOptions({
  name: 'ChangeLogs'
});

interface LogItem {
  id: number;
  content: string;
  version: string;
  time: string;
}

const logs = computed<LogItem[]>(() => [
  { id: 3, content: 'Sistema de status online/offline dos dispositivos com atualização automática', version: 'v1.3.0', time: '2025-12-01' },
  { id: 2, content: 'Gerenciamento completo de Address Books com visualização de peers', version: 'v1.3.0', time: '2025-12-01' },
  { id: 1, content: 'Sincronização automática do sistema operacional no dashboard', version: 'v1.3.0', time: '2025-12-01' }
]);
</script>

<template>
  <NCard :title="$t('page.home.changeLogs')" :bordered="false" size="small" segmented class="card-wrapper">
    <NList>
      <NListItem v-for="item in logs" :key="item.id">
        <template #prefix>
          <SoybeanAvatar class="size-48px!" />
        </template>
        <NThing :title="item.content" :title-extra="item.version" :description="item.time" />
      </NListItem>
    </NList>
  </NCard>
</template>

<style scoped></style>
