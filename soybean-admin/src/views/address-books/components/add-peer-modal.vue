<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import { useNaiveForm } from '@/hooks/common/form';

defineOptions({
  name: 'AddPeerModal'
});

interface Props {
  addressBookId: number;
  addressBookName: string;
}

const props = defineProps<Props>();

interface Emits {
  (e: 'success'): void;
}

const emit = defineEmits<Emits>();

const visible = defineModel<boolean>('visible', { required: true });

const { formRef, validate } = useNaiveForm();

const title = computed(() => `Adicionar Peer - ${props.addressBookName}`);

const formModel = ref({
  rustdesk_id: '',
  alias: '',
  password: '',
  hostname: '',
  username: '',
  platform: 'Windows',
  tags: [] as string[]
});

const platformOptions = [
  { label: 'Windows', value: 'Windows' },
  { label: 'Linux', value: 'Linux' },
  { label: 'macOS', value: 'Mac OS' },
  { label: 'Android', value: 'Android' }
];

const tagOptions = computed(() => {
  return availableTags.value.map(tag => ({
    label: tag,
    value: tag
  }));
});

const loading = ref(false);
const searching = ref(false);
const availableTags = ref<string[]>([]);

watch(visible, async (val) => {
  if (val) {
    // Reset form when opening
    formModel.value = {
      rustdesk_id: '',
      alias: '',
      password: '',
      hostname: '',
      username: '',
      platform: 'Windows',
      tags: []
    };
    
    // Fetch available tags for this address book
    await fetchAvailableTags();
  }
});

// Fetch available tags for suggestions
async function fetchAvailableTags() {
  try {
    const { request } = await import('@/service/request');
    
    // Get tags directly from the address book
    const response = await request({
      url: `/address-books/${props.addressBookId}/tags`,
      method: 'get'
    });

    // Response.data should be array of tag names
    if (response.data && Array.isArray(response.data)) {
      availableTags.value = response.data;
    } else {
      availableTags.value = [];
    }
    
    console.log('Tags disponíveis:', availableTags.value);
  } catch (error) {
    console.error('Erro ao buscar tags disponíveis:', error);
    availableTags.value = [];
  }
}


// Auto-fill device info when RustDesk ID is entered
async function handleRustdeskIdBlur() {
  const trimmedId = formModel.value.rustdesk_id?.trim();
  if (!trimmedId) {
    return;
  }

  // Update the field with trimmed value
  formModel.value.rustdesk_id = trimmedId;

  // Only search if other fields are empty (don't override manual input)
  if (formModel.value.hostname !== '' || formModel.value.alias !== '') {
    return;
  }

  searching.value = true;
  try {
    const { request } = await import('@/service/request');
    
    // Search for device by RustDesk ID
    const response = await request({
      url: '/devices/list',
      method: 'get',
      params: {
        rustdesk_id: trimmedId,
        current: 1,
        size: 1
      }
    });

    const data = response.data as any;
    if (data && data.records && data.records.length > 0) {
      const device = data.records[0];
      
      // Auto-fill fields if they're empty
      if (!formModel.value.hostname) {
        formModel.value.hostname = device.hostname || '';
      }
      if (!formModel.value.username) {
        formModel.value.username = device.username || '';
      }
      // Alias = RustDesk ID to show correctly in RustDesk client
      if (!formModel.value.alias) {
        formModel.value.alias = device.rustdesk_id || device.hostname || '';
      }
      if (!formModel.value.platform || formModel.value.platform === 'Windows') {
        formModel.value.platform = device.os || 'Windows';
      }

      window.$message?.success('Informações do dispositivo preenchidas automaticamente!');
    } else {
      window.$message?.info('Dispositivo não encontrado. Preencha manualmente.');
    }
  } catch (error) {
    console.error('Erro ao buscar dispositivo:', error);
    // Não mostra erro, apenas deixa o usuário preencher manualmente
  } finally {
    searching.value = false;
  }
}

async function handleSubmit() {
  await validate();
  
  loading.value = true;
  try {
    // Call API to create peer using the request utility
    const { request } = await import('@/service/request');
    
    await request({
      url: `/address-books/${props.addressBookId}/peers`,
      method: 'post',
      data: {
        ab_id: props.addressBookId,
        rustdesk_id: formModel.value.rustdesk_id,
        alias: formModel.value.alias,
        password: formModel.value.password,
        hostname: formModel.value.hostname,
        username: formModel.value.username,
        platform: formModel.value.platform,
        tags: JSON.stringify(formModel.value.tags)
      }
    });

    window.$message?.success('Peer adicionado com sucesso!');
    emit('success');
    visible.value = false;
  } catch (error) {
    console.error('Erro ao adicionar peer:', error);
    window.$message?.error('Erro ao adicionar peer');
  } finally {
    loading.value = false;
  }
}

function handleClose() {
  visible.value = false;
}
</script>

<template>
  <NModal v-model:show="visible" :title="title" preset="card" class="w-700px">
    <NForm ref="formRef" :model="formModel" label-placement="left" :label-width="120">
      <NGrid :cols="24" :x-gap="18">
        <NFormItemGi :span="24" label="ID RustDesk" path="rustdesk_id">
          <NInput 
            v-model:value="formModel.rustdesk_id" 
            placeholder="Digite o ID do RustDesk e pressione Tab"
            :loading="searching"
            :input-props="{ autocomplete: 'off' }"
            @blur="handleRustdeskIdBlur"
          >
            <template #suffix>
              <NTooltip v-if="!searching">
                <template #trigger>
                  <icon-mdi-information-outline class="text-icon" />
                </template>
                Auto-preenche os dados se o dispositivo estiver cadastrado
              </NTooltip>
              <NSpin v-else size="small" />
            </template>
          </NInput>
        </NFormItemGi>
        
        <NFormItemGi :span="24" label="Alias (Apelido)" path="alias">
          <NInput 
            v-model:value="formModel.alias" 
            placeholder="Preenchido automaticamente (editável)"
            :input-props="{ autocomplete: 'off' }"
          />
        </NFormItemGi>
        
        <NFormItemGi :span="24" label="Senha" path="password">
          <NInput 
            v-model:value="formModel.password" 
            type="password"
            show-password-on="click"
            placeholder="Senha de acesso permanente (opcional)"
            :input-props="{ autocomplete: 'new-password' }"
          />
        </NFormItemGi>
        
        <NFormItemGi :span="12" label="Hostname" path="hostname">
          <NInput 
            v-model:value="formModel.hostname" 
            placeholder="Auto-preenchido (editável)"
            :input-props="{ autocomplete: 'off' }"
          />
        </NFormItemGi>
        
        <NFormItemGi :span="12" label="Username" path="username">
          <NInput 
            v-model:value="formModel.username" 
            placeholder="Auto-preenchido (editável)"
            :input-props="{ autocomplete: 'off' }"
          />
        </NFormItemGi>
        
        <NFormItemGi :span="12" label="Plataforma" path="platform">
          <NSelect 
            v-model:value="formModel.platform" 
            :options="platformOptions"
          />
        </NFormItemGi>
        
        <NFormItemGi :span="24" label="Tags" path="tags">
          <NSelect 
            v-model:value="formModel.tags" 
            :options="tagOptions"
            multiple
            filterable
            tag
            placeholder="Selecione ou crie novas tags"
            :input-props="{ autocomplete: 'off' }"
          />
          <template v-if="availableTags.length > 0">
            <NText depth="3" style="font-size: 12px; margin-top: 4px;margin-left: 4px; display: block;">
              Tags disponíveis: {{ availableTags.join(', ') }}
            </NText>
          </template>
        </NFormItemGi>
        
        <NFormItemGi :span="24">
          <NAlert type="info" size="small">
            <template #icon>
              <icon-mdi-lightbulb-on-outline />
            </template>
             <strong>Dica:</strong> Digite o ID RustDesk e pressione Tab. Os campos serão preenchidos automaticamente se o dispositivo estiver cadastrado.
          </NAlert>
        </NFormItemGi>
      </NGrid>
    </NForm>

    <template #footer>
      <NSpace justify="end">
        <NButton @click="handleClose">Cancelar</NButton>
        <NButton type="primary" :loading="loading" @click="handleSubmit">
          Adicionar
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>

<style scoped></style>
