<script lang="ts">
  export default {
    name: 'AttackSource',
  }
</script>

<script setup lang="ts">
  const router = useRouter()
  const props = defineProps<{
    data: {
      name: string
      value: number
    }[]
  }>()

  const sourceList = ref<
    {
      name: string
      value: number
    }[]
  >([])

  const goToAlert = (ip: string) => {
    const resolveRouter = router.resolve({
      path: '/alerts/alarm-attack',
      query: { attackIp: ip },
    })
    window.open(resolveRouter.href, '_blank')
  }

  watchEffect(() => {
    sourceList.value = props.data
  })
  const arrList = computed(() => {
    const middleIndex = 5
    const leftArr = sourceList.value.slice().splice(0, middleIndex)
    const rightArr = sourceList.value.slice().splice(5)
    return [leftArr, rightArr]
  })
</script>

<template>
  <div class="attack-source-container">
    <ul class="sources-content">
      <li
        v-for="(item, index) in arrList[1]"
        :key="index"
        class="source right"
        :style="{
          transform: `translateX(${130 + (index > 2 ? 5 - index : index + 1) * 40}px) translateY(${index * 66 - 22}px)`,
        }"
        @click="() => goToAlert(item.name)"
      >
        {{ item.name }}
      </li>
      <li
        v-for="(item, index) in arrList[0]"
        :key="index"
        class="source left"
        :style="{
          transform: `translateX(-${130 + (index > 2 ? 5 - index : index + 1) * 40}px) translateY(${
            index * 66 - 22
          }px)`,
        }"
        @click="() => goToAlert(item.name)"
      >
        {{ item.name }}
      </li>
    </ul>
  </div>
</template>

<style scoped lang="scss">
  .attack-source-container {
    display: grid;
    .sources-content {
      width: 250px;
      height: 250px;
      place-self: center;
      background: url('@/assets/dashboard_images/attack-source.png') no-repeat center;
      background-size: 215px 215px;
      position: relative;
      margin: 14px 0;
      .source {
        width: 230px;
        height: 44px;
        padding: 0 10px;
        line-height: 44px;
        font-size: 24px;
        color: #91b1e9;
        position: absolute;
        transform: translateX(180px) translateY(-10px);
        overflow-x: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
        cursor: pointer;
        &::after {
          content: ' ';
          display: block;
          position: absolute;
          inset: 0;
          background: url('@/assets/dashboard_images/source-item.png') no-repeat center left;
        }
        &::before {
          content: '•';
          margin-right: 5px;
        }
        &.left {
          left: 0;
          text-align: right;
          padding-right: 16px;
          &::after {
            transform: rotateY(179deg);
          }
        }
        &.right {
          right: 0;
          text-indent: 5px;
        }

        &:nth-child(2) {
          transform: translateX(180px) translateY(50px);
        }

        &:nth-child(3) {
          transform: translateX(210px) translateY(110px);
        }

        &:nth-child(4) {
          transform: translateX(240px) translateY(170px);
        }

        &:nth-child(5) {
          transform: translateX(210px) translateY(230px);
        }

        &:nth-child(6) {
          transform: translateX(-180px) translateY(230px);
        }

        &:nth-child(7) {
          transform: translateX(-210px) translateY(170px);
        }

        &:nth-child(8) {
          transform: translateX(-240px) translateY(110px);
        }

        &:nth-child(9) {
          transform: translateX(-210px) translateY(50px);
        }

        &:nth-child(10) {
          transform: translateX(-180px) translateY(-10px);
        }
      }
    }
  }
</style>
