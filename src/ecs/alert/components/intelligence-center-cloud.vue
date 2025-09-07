<script lang="ts">
  export default {
    name: 'IntelligenceCenterCloud',
  }
</script>

<script setup lang="ts">
  import { getInfoCloudPageApi } from '@/api-ecs/alert'
  import { Search } from '@element-plus/icons-vue'
  const search = reactive({
    searchStr: '',
    pageNum: 1,
    pageSize: 10,
  })
  const show = ref(false)
  const infoData = ref()
  const total = ref(100)
  const loading = ref(true)
  const getInfoCloudPageHandle = async () => {
    loading.value = true
    show.value = !show.value || true
    const { data } = await getInfoCloudPageApi(search)
    infoData.value = data && data.records
    total.value = data?.total || 0
    loading.value = false
  }
  const getConfidence = (code: string) => {
    return ['低', '中', '高'].findIndex((i) => code === i) + 1
  }
  onMounted(() => {
    particlesJS('particles', {
      particles: {
        number: {
          value: 200,
          density: {
            enable: true,
            value_area: 750,
          },
        },
        color: {
          value: '#A495FF',
        },
        shape: {
          type: 'circle',
          stroke: {
            width: 0,
            color: 'red',
          },
          polygon: {
            nb_sides: 5,
          },
        },
        opacity: {
          value: 0.8,
          random: true,
          anim: {
            enable: true,
            speed: 1,
            opacity_min: 0.1,
            sync: false,
          },
        },
        size: {
          value: 5,
          random: true,
          anim: {
            enable: true,
            speed: 1,
            size_min: 0.1,
            sync: false,
          },
        },
        line_linked: {
          enable: true,
          distance: 150,
          color: '#A495FF',
          opacity: 0.4,
          width: 1,
        },
        move: {
          enable: true,
          speed: 0.8,
          direction: 'none',
          random: false,
          straight: false,
          out_mode: 'out',
          bounce: false,
          attract: {
            enable: false,
            rotateX: 600,
            rotateY: 1200,
          },
        },
      },
      interactivity: {
        detect_on: 'canvas',
        events: {
          onhover: {
            enable: true,
            mode: 'grab',
          },
          resize: true,
        },
        modes: {
          grab: {
            distance: 140,
            line_linked: {
              opacity: 1,
            },
          },
        },
      },
      retina_detect: true,
    })
  })
  onUnmounted(() => {
    const pJSDom = (window as any).pJSDom
    pJSDom[0].pJS.fn.vendors.destroypJS()
    ;(window as any).pJSDom = []
  })
</script>

<template>
  <div class="info-cloud-container">
    <div id="particles"></div>
    <div class="search-box" :class="{ move: show }">
      <h1>威胁情报检索中心</h1>
      <div style="position: relative">
        <el-input
          v-model.trim="search.searchStr"
          clearable
          placeholder="请输入有效域名、IP、URL、文件MD5等IOC进行查询"
          :prefix-icon="Search"
          size="large"
          style="width: 100%; border-radius: 8px; height: 46px"
        />
        <el-button
          :auto-insert-space="false"
          style="width: 98px; height: 42px; position: absolute; border-radius: 8px; top: 2px; right: 2px"
          type="primary"
          @click="
            () => {
              search.pageNum = 1
              getInfoCloudPageHandle()
            }
          "
        >
          检索
        </el-button>
      </div>
    </div>
    <Transition name="slide-fade">
      <div v-if="show" v-loading="loading" class="divvv" element-loading-background="rgba(16, 13, 42, 0.8)">
        <el-table :data="infoData" stripe style="width: 100%">
          <el-table-column label="威胁名称" prop="ruleDesc" show-overflow-tooltip width="240" />
          <el-table-column label="数据来源" prop="iocSource" width="120" />
          <el-table-column label="威胁等级" prop="hazardLevelStr" width="100">
            <template #default="{ row }">
              <span v-if="row.hazardLevelStr" :class="['attackResult', `attackResult-${row.hazardLevel}`]">
                <el-icon><WarnTriangleFilled /></el-icon>
                {{ row.hazardLevelStr }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="置信度" prop="confidence" width="80">
            <template #default="{ row }">
              <el-rate
                :colors="['#67C23A', '#67C23A', '#67C23A']"
                disabled
                disabled-void-color="#C7C6D4"
                :max="3"
                :model-value="getConfidence(row.confidence)"
              />
            </template>
          </el-table-column>
          <el-table-column label="常用攻击手段" prop="attackMethod" show-overflow-tooltip width="120" />
          <el-table-column label="报告时间" prop="reportDate" show-overflow-tooltip width="120" />
          <el-table-column label="威胁描述" prop="frontDescription" show-overflow-tooltip />
          <template #empty>
            <div class="empty">
              <img alt="暂无数据" :src="require('@/assets/alert_images/empty.svg')" width="130" />
              暂未查询到关联的恶意情报
            </div>
          </template>
        </el-table>
        <el-pagination
          v-model:current-page="search.pageNum"
          v-model:page-size="search.pageSize"
          background
          class="cloud_pagination"
          layout="total, sizes, prev, pager, next, jumper"
          :page-sizes="[10, 20, 30]"
          :total="total"
          @current-change="getInfoCloudPageHandle"
          @size-change="getInfoCloudPageHandle"
        />
      </div>
    </Transition>
  </div>
</template>

<style scoped lang="scss">
  .info-cloud-container {
    position: relative;
    overflow: hidden;
    height: 100%;
    #particles {
      position: absolute;
      height: 42vh;
      z-index: 1;
      inset: 0;
    }
    .search-box {
      width: 50vw;
      margin: 230px auto 0;
      z-index: 2;
      position: relative;
      transition: all 0.3s ease-out;
      &.move {
        margin-top: 30px;
      }
      h1 {
        font-weight: 600;
        font-size: 64px;
        color: #ffffff;
        line-height: 90px;
        text-align: center;
        margin-bottom: 50px;
      }
      :deep() {
        .el-input__wrapper {
          border-radius: 8px;
          padding-right: 110px;
        }
      }
    }
    .divvv {
      width: 75vw;
      min-width: 1000px;
      min-height: 300px;
      margin: 80px auto 0;
      background: rgba(16, 13, 42, 0.8);
      border-radius: 8px;
      border: 1px solid #332d56;
      z-index: 1;
      position: relative;
      overflow: hidden;
      .attackResult {
        font-size: 14px;
        width: 58px;
        height: 24px;
        border-radius: 4px;
        display: block;
        line-height: 22px;
        text-align: center;
        .el-icon {
          margin-right: -4px;
          vertical-align: -3px;
          font-size: 17px;
        }
        &-critical {
          background: #d60705;
        }
        &-high {
          background: #ff5553;
        }
        &-medium {
          background: #ff8538;
        }
        &-low {
          background: #ffc643;
        }
        &-safety {
          background: #32b45f;
        }
      }
      .empty {
        display: flex;
        flex-direction: column;
        align-items: center;
        color: #bebcd0;
      }
    }
    :deep() {
      .el-table {
        background-color: transparent;
        --el-table-border-color: #332d56;
        --el-table-tr-bg-color: transparent;
        --el-fill-color-lighter: transparent;
        thead {
          background: linear-gradient(180deg, #2d2370 0%, #160d3d 100%) !important;
          tr,
          th {
            background-color: transparent !important;
          }
        }
        .el-scrollbar__view:has(.el-table__empty-block) {
          height: 100%;
        }
        .el-table__cell:first-child {
          text-indent: 10px;
        }
        .el-table__row .el-table__cell .cell {
          color: #ffffff !important;
        }
        .el-table__body tr:hover > td.el-table__cell {
          background-color: transparent !important;
        }
        .el-scrollbar {
          height: calc(100vh - 540px);
          .el-scrollbar__thumb {
            background-color: #eaeaeb !important;
          }
        }
      }
      .cloud_pagination {
        margin: 20px;
        .el-input__wrapper {
          background-color: #1c173d;
          --el-input-border-color: #2b2554;
          --el-input-hover-border-color: #2b2554;
          --el-select-border-color-hover: #2b2554;
        }
        .el-pager .number,
        .el-pager .more,
        .btn-prev,
        .btn-next {
          background: #3a316f;
          color: #ffffff;
        }
        .number.is-active {
          background: #6954f0;
        }
      }
    }
    .slide-fade-enter-active {
      animation: bounce-in 0.5s ease;
    }
    .slide-fade-leave-active {
      animation: bounce-in 0.5s reverse ease;
    }
    @keyframes bounce-in {
      0% {
        transform: translateY(400px);
      }
      50% {
        opacity: 0.5;
      }
      100% {
        transform: translateY(0);
        opacity: 1;
      }
    }
  }
</style>
