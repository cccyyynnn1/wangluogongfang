<script lang="ts">
  export default {
    name: 'SituationScreen', // 态势大屏
  }
</script>

<script setup lang="ts">
  import SituationList from './components/situation-screen/situation-list.vue'
  import SituationConfig from './components/situation-screen/situation-config.vue'
  import { getSystemConfigApi } from '~/src/api-ecs/system'
  // @ts-ignore
  import defaultImg from '@/assets/index_images/WechatIMG12.jpg'

  type objType = {
    [key: string]: any
  }

  const urlOBJ = {
    screen_word_map: '/dashboard-screen',
    screen_highway_liaoning: '/screen-highway/liaoning',
  }

  const obj: objType = {
    screen_highway_liaoning: 'screen_highway_liaoning',
    screen_word_map: 'screen_word_map',
  }

  const configData = reactive<objType>({
    screen_word_map: undefined,
    screen_highway_liaoning: undefined,
  })

  const chapter = ref('chapter_1')

  const backEvent = (val: string) => {
    chapter.value = val
  }

  const curItem = ref()

  const otherConfig = (val: any) => {
    curItem.value = val
  }

  const alldata = ref<any>([])

  const KEYS = 'screen_highway_liaoning,screen_word_map'

  const getData = async () => {
    const { data } = await getSystemConfigApi({ keys: KEYS })
    alldata.value = []
    for (const key in configData) {
      data[obj[key]].value = data[obj[key]].value && JSON.parse(data[obj[key]].value)
      configData[key] = data[obj[key]]
      // @ts-ignore
      alldata.value.push(data[obj[key]])
    }
  }

  const handleSkip = (item: any) => {
    // @ts-ignore
    if (urlOBJ[item.key]) {
      // @ts-ignore
      const url = `${window.location.origin}#${urlOBJ[item.key]}`
      window.open(url)
    } else {
      ElMessage.error('功能尚未开发!')
    }
  }

  const handleReflash = () => {
    getData()
  }

  onMounted(() => {
    getData()
  })
</script>

<template>
  <div class="Situation-container">
    <div v-if="chapter == 'chapter_1'" class="chapter_1">
      <div class="top">
        <div class="title">态势大屏</div>
        <el-button @click="chapter = 'chapter_2'">
          <el-icon><Setting /></el-icon>
        </el-button>
      </div>
      <div class="content">
        <el-row :gutter="15">
          <template v-for="item in alldata" :key="item.id">
            <el-col v-if="item.value.enable" class="warp" :span="6" style="margin-bottom: 15px">
              <el-image class="img" lazy :src="item.value.thumbnail || defaultImg" @click="handleSkip(item)" />
              <!-- <el-image v-else class="img" lazy :src="" /> -->
            </el-col>
          </template>
        </el-row>
      </div>
    </div>
    <SituationList
      v-else-if="chapter == 'chapter_2'"
      :list-date="alldata"
      @on-back-event="backEvent"
      @on-reflash="getData()"
      @on-to-other="otherConfig"
    />
    <SituationConfig
      v-else-if="chapter == 'chapter_3'"
      :remark="curItem"
      @on-back-event="backEvent"
      @on-reflash="handleReflash"
    />
  </div>
</template>

<style scoped lang="scss">
  .top {
    display: flex;
    align-items: center;
    justify-content: space-between;
    .title {
      height: 28px;
      font-size: 20px;
      font-weight: 500;
      color: #303133;
      line-height: 28px;
    }
  }

  .content {
    margin-top: 25px;
    height: calc(100vh - 260px);
    width: 100%;
    overflow-y: auto;
    // background-color: antiquewhite;
    &::-webkit-scrollbar {
      width: 0;
      height: 0;
    }

    .warp {
      position: relative;
      width: 100%;
      padding-top: 15%;
    }
    .img {
      position: absolute;
      top: 0;
      left: 0;
      width: 95%;
      height: 100%;
      &:hover {
        cursor: pointer;
      }
      img {
        width: 100%;
        height: 100%;
      }
    }
  }
</style>
