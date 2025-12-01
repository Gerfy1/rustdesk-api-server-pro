<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import { createAddressBook, updateAddressBook } from '@/service/api/address-books';
import { useNaiveForm } from '@/hooks/common/form';

defineOptions({
  name: 'AddressBookModal'
});

interface Props {
  mode: 'add' | 'edit';
  record: Api.AddressBooks.AddressBook | null;
}

const props = defineProps<Props>();

interface Emits {
  (e: 'success'): void;
}

const emit = defineEmits<Emits>();

const visible = defineModel<boolean>('visible', { required: true });

const { formRef, validate } = useNaiveForm();

const title = computed(() => (props.mode === 'add' ? 'Adicionar Address Book' : 'Editar Address Book'));

const formModel = ref({
  user_id: 1,
  name: '',
  note: '',
  rule: 3,
  max_peer: 0,
  shared: false
});

const ruleOptions = [
  { label: 'Leitura', value: 1 },
  { label: 'Ler/Escrever', value: 2 },
  { label: 'Controle Total', value: 3 }
];

watch(visible, (val) => {
  if (val && props.mode === 'edit' && props.record) {
    formModel.value = {
      user_id: props.record.user_id,
      name: props.record.name,
      note: props.record.note || '',
      rule: props.record.rule,
      max_peer: props.record.max_peer,
      shared: props.record.shared
    };
  } else if (val && props.mode === 'add') {
    formModel.value = {
      user_id: 1,
      name: '',
      note: '',
      rule: 3,
      max_peer: 0,
      shared: false
    };
  }
});

async function handleSubmit() {
  await validate();
  
  try {
    if (props.mode === 'add') {
      await createAddressBook(formModel.value);
      window.$message?.success('Address Book criado com sucesso');
    } else if (props.record?.id) {
      await updateAddressBook(props.record.id, formModel.value);
      window.$message?.success('Address Book atualizado com sucesso');
    }
    emit('success');
  } catch (error) {
    window.$message?.error('Erro ao salvar Address Book');
  }
}

function handleClose() {
  visible.value = false;
}
</script>

<template>
  <NModal v-model:show="visible" :title="title" preset="card" class="w-700px">
    <NForm ref="formRef" :model="formModel" label-placement="left" :label-width="120">
      <NFormItem label="Nome" path="name" required>
        <NInput v-model:value="formModel.name" placeholder="Nome do Address Book" />
      </NFormItem>
      
      <NFormItem label="Nota" path="note">
        <NInput
          v-model:value="formModel.note"
          type="textarea"
          placeholder="Adicione uma nota (opcional)"
          :rows="3"
        />
      </NFormItem>

      <NFormItem label="Permissão" path="rule">
        <NSelect v-model:value="formModel.rule" :options="ruleOptions" />
      </NFormItem>

      <NFormItem label="Max Peers" path="max_peer">
        <NInputNumber
          v-model:value="formModel.max_peer"
          :min="0"
          placeholder="0 = ilimitado"
          class="w-full"
        />
      </NFormItem>

      <NFormItem label="Compartilhado" path="shared">
        <NSwitch v-model:value="formModel.shared">
          <template #checked>Sim</template>
          <template #unchecked>Não</template>
        </NSwitch>
      </NFormItem>
    </NForm>

    <template #footer>
      <NSpace justify="end">
        <NButton @click="handleClose">Cancelar</NButton>
        <NButton type="primary" @click="handleSubmit">Salvar</NButton>
      </NSpace>
    </template>
  </NModal>
</template>

<style scoped></style>
