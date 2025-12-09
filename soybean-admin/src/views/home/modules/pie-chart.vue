<script setup lang="ts">
import { ref, nextTick } from 'vue';
import { $t } from '@/locales';
import { useEcharts } from '@/hooks/common/echarts';
import { fetchPieCharts } from '@/service/api/home';

defineOptions({
  name: 'PieChart'
});

interface PieDataItem {
  name: string;
  value: number;
}

const showModal = ref(false);
const chartData = ref<PieDataItem[]>([]);

const { domRef, updateOptions } = useEcharts(() => ({
  title: {
    text: $t('page.home.operatingSystem'),
    left: 'center'
  },
  tooltip: {
    trigger: 'item',
    formatter: '{b}: {c} ({d}%)'
  },
  legend: {
    type: 'scroll',
    bottom: '0%',
    left: 'center',
    itemStyle: {
      borderWidth: 0
    }
  },
  series: [
    {
      name: $t('page.home.operatingSystem'),
      type: 'pie',
      radius: ['40%', '70%'],
      center: ['50%', '45%'],
      avoidLabelOverlap: true,
      itemStyle: {
        borderRadius: 10,
        borderColor: '#fff',
        borderWidth: 1
      },
      label: {
        show: false
      },
      emphasis: {
        label: {
          show: true,
          fontSize: '12',
          fontWeight: 'bold'
        }
      },
      data: [] as PieDataItem[]
    }
  ]
}));

// Configuração do gráfico expandido (modal)
const { domRef: modalDomRef, updateOptions: updateModalOptions } = useEcharts(() => ({
  title: {
    text: $t('page.home.operatingSystem'),
    left: 'center',
    textStyle: {
      fontSize: 18
    }
  },
  tooltip: {
    trigger: 'item',
    formatter: '{b}: {c} ({d}%)'
  },
  legend: {
    type: 'scroll',
    orient: 'vertical',
    right: '5%',
    top: 'middle',
    itemStyle: {
      borderWidth: 0
    },
    textStyle: {
      fontSize: 14
    }
  },
  series: [
    {
      name: $t('page.home.operatingSystem'),
      type: 'pie',
      radius: ['35%', '65%'],
      center: ['40%', '50%'],
      avoidLabelOverlap: true,
      itemStyle: {
        borderRadius: 10,
        borderColor: '#fff',
        borderWidth: 2
      },
      label: {
        show: true,
        fontSize: 12,
        formatter: '{b}: {d}%'
      },
      labelLine: {
        show: true
      },
      emphasis: {
        label: {
          show: true,
          fontSize: '16',
          fontWeight: 'bold'
        }
      },
      data: [] as PieDataItem[]
    }
  ]
}));

async function fetchChartsData() {
  const pie = await fetchPieCharts();
  
  const processedData: PieDataItem[] = (pie.data as PieDataItem[])?.map(item => ({
    ...item,
    name: item.name || 'Unknown'
  })) || [];
  
  chartData.value = processedData;
  
  updateOptions(opt => {
    opt.series = [
      {
        name: $t('page.home.operatingSystem'),
        type: 'pie',
        radius: ['40%', '70%'],
        center: ['50%', '45%'],
        avoidLabelOverlap: true,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 1
        },
        label: {
          show: false
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '12',
            fontWeight: 'bold'
          }
        },
        data: processedData
      }
    ];
    return opt;
  });
}

function openExpandedView() {
  showModal.value = true;
  nextTick(() => {
    updateModalOptions(opt => {
      opt.series = [
        {
          name: $t('page.home.operatingSystem'),
          type: 'pie',
          radius: ['35%', '65%'],
          center: ['40%', '50%'],
          avoidLabelOverlap: true,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: {
            show: true,
            fontSize: 12,
            formatter: '{b}: {d}%'
          },
          labelLine: {
            show: true
          },
          emphasis: {
            label: {
              show: true,
              fontSize: '16',
              fontWeight: 'bold'
            }
          },
          data: chartData.value
        }
      ];
      return opt;
    });
  });
}

async function init() {
  fetchChartsData();
}

// init
init();
</script>

<template>
  <NCard :bordered="false" class="card-wrapper">
    <template #header-extra>
      <NButton quaternary circle size="small" @click="openExpandedView">
        <template #icon>
          <icon-mdi-fullscreen />
        </template>
      </NButton>
    </template>
    <div ref="domRef" class="h-360px overflow-hidden"></div>
  </NCard>

  <!-- Modal Expandido -->
  <NModal v-model:show="showModal" preset="card" :title="$t('page.home.operatingSystem')" style="width: 90%; max-width: 900px;">
    <div ref="modalDomRef" class="h-500px"></div>
  </NModal>
</template>

<style scoped></style>
