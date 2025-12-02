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

const loading = ref(false);

watch(visible, (val) => {
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
  }
});

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
            placeholder="Digite o ID do RustDesk"
            :input-props="{ autocomplete: 'off' }"
          />
        </NFormItemGi>
        
        <NFormItemGi :span="24" label="Alias (Apelido)" path="alias">
          <NInput 
            v-model:value="formModel.alias" 
            placeholder="Apelido para identificação"
            :input-props="{ autocomplete: 'off' }"
          />
        </NFormItemGi>
        
        <NFormItemGi :span="24" label="Senha" path="password">
          <NInput 
            v-model:value="formModel.password" 
            type="password"
            show-password-on="click"
            placeholder="Senha de acesso (opcional)"
            :input-props="{ autocomplete: 'new-password' }"
          />
        </NFormItemGi>
        
        <NFormItemGi :span="12" label="Hostname" path="hostname">
          <NInput 
            v-model:value="formModel.hostname" 
            placeholder="Nome do host"
            :input-props="{ autocomplete: 'off' }"
          />
        </NFormItemGi>
        
        <NFormItemGi :span="12" label="Username" path="username">
          <NInput 
            v-model:value="formModel.username" 
            placeholder="Nome de usuário"
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
          <NDynamicTags v-model:value="formModel.tags" />
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
