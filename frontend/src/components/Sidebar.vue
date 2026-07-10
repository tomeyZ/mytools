<template>
  <div class="sidebar-container">
    <!-- 顶部功能列表 -->
    <div class="sidebar-header">
      <h2 class="sidebar-title">功能列表</h2>
      <ul class="menu-list">
        <li
            v-for="item in menuItems"
            :key="item.id"
            :class="{ active: activeItem === item.id }"
            @click="switchTool(item.id)"
        >
          {{ item.name }}
        </li>
      </ul>
    </div>

    <!-- 底部版本信息 -->
    <div class="sidebar-footer">
      <p class="version-info">版本：{{ currentVersion }}</p>
    </div>

  </div>
</template>
<script>
export default {
  data() {
    return {
      menuItems: [
        { id: 'time', name: '🕒 时区转换' },
        { id: 'md5', name: '🔐 MD5加密' },
        { id: 'json', name: '{.} JSON美化' },
        { id: 'ip', name: '🌐 IP地址查询' },
        { id: 'aes', name: '🔐 AES加解密' },
        { id: 'rsa', name: '🔑 RSA加解密' }
      ],
      activeItem: 'time',
      currentVersion: ''
    }
  },
  mounted() {
    this.loadVersion()
  },
  methods: {
    switchTool(itemId) {
      this.activeItem = itemId
      this.$emit('tool-change', itemId)
    },
    async loadVersion() {
      try {
        this.currentVersion = await window.go.handler.VersionHandler.GetCurrentVersion()
      } catch (error) {
        console.error('获取版本号失败:', error)
        this.currentVersion = '1.0.0'
      }
    }
  }
}
</script>

<style scoped>
.sidebar-container {
  width: 220px;
  background: #f8f9fa;
  height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  box-shadow: 0 0 10px rgba(0,0,0,0.05);
}

.sidebar-header {
  padding: 20px;
  flex: 1;
  overflow-y: auto;
}

.sidebar-title {
  font-size: 23px;
  font-weight: bold;
  color: #333;
  margin-bottom: 14px;
}

.menu-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.menu-list li {
  padding: 12px 16px;
  color: #555;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
}

.menu-list li:hover {
  background: #e9ecef;
  color: #333;
}

.menu-list li.active {
  background: #42b983;
  color: #fff;
}

.sidebar-footer {
  padding: 10px 15px 15px 15px;
  background: #fff;
  border-top: 1px solid #eee;
  text-align: center;
  flex-shrink: 0;
}

.version-info {
  font-size: 14px;
  color: #999;
}
</style>