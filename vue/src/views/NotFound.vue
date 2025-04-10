<template>
  <div class="app-container">
    <h1>404 - 页面未找到</h1>
    <p>您访问的页面不存在。</p>
    <el-button type="primary" @click="backToIndex">返回首页</el-button>
  </div>
</template>

<script>
export default {
  methods: {
    backToIndex() {
      this.$router.push('/')
    }
  }
}
</script>

<style scoped>
.app-container {
  padding: 20px;
  text-align: center;
  max-width: 100%;
  box-sizing: border-box;
}

@media (max-width: 768px) {
  .app-container {
    padding: 10px;
  }

  h1 {
    font-size: 24px;
  }

  p {
    font-size: 14px;
  }

  .el-button {
    width: 100%;
    margin-top: 10px;
  }
}

@media (max-width: 480px) {
  .app-container {
    padding: 5px;
  }

  h1 {
    font-size: 20px;
  }

  .el-button {
    font-size: 12px;
  }
}
</style>