<template>
  <div class="update-dialog" v-if="showDialog">
    <div class="dialog-content">
      <div class="dialog-header">
        <div class="update-icon">
          <svg viewBox="0 0 24 24" width="24" height="24">
<!--            <path fill="currentColor" d="M17.65,6.35C16.2,4.9 14.21,4 12,4A8,8 0 0,0 4,12A8,8 0 0,0 12,20C15.73,20 18.84,17.45 19.73,14H17.65C16.83,16.33 14.61,18 12,18A6,6 0 0,1 6,12A6,6 0 0,1 12,6C13.66,6 15.14,6.69 16.22,7.78L13,11H20V4L17.65,6.35Z"/>-->
          </svg>
        </div>
        <h3>发现新版本</h3>
        <button class="close-btn" @click="closeDialog">
          <svg viewBox="0 0 24 24" width="18" height="18">
            <path fill="currentColor" d="M19,6.41L17.59,5L12,10.59L6.41,5L5,6.41L10.59,12L5,17.59L6.41,19L12,13.41L17.59,19L19,17.59L13.41,12L19,6.41Z"/>
          </svg>
        </button>
      </div>
      <div class="dialog-body">
        <div class="version-info">
          <div class="version-badge">NEW</div>
          <p>
            最新版本: <span class="version-text">{{ latestVersion }}</span>
            <span class="version-arrow">→</span>
            当前版本: <span class="version-text current">{{ currentVersion }}</span>
          </p>
        </div>

        <div class="changelog" v-if="hasUpdate && changelog">
          <h4>更新内容 <span class="changelog-badge">{{ changelog.length }}项更新</span></h4>
          <ul class="compact-list">
            <li v-for="(item, index) in changelog" :key="index">
              <span class="bullet">•</span> {{ item }}
            </li>
          </ul>
        </div>
      </div>
      <div class="dialog-footer" v-if="hasUpdate">
        <button class="later-btn" @click="laterUpdate">
          稍后提醒我
        </button>
        <button class="download-btn" @click="downloadUpdate">
          <svg viewBox="0 0 24 24" width="18" height="18">
            <path fill="currentColor" d="M5,20H19V18H5M19,9H15V3H9V9H5L12,16L19,9Z"/>
          </svg>
          立即更新
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    showDialog: Boolean,
    currentVersion: String,
    latestVersion: String,
    changelog: Array
  },
  data() {
    return {
      hasUpdate: true
    }
  },
  methods: {
    closeDialog() {
      this.$emit('close')
    },
    laterUpdate() {
      this.closeDialog()
    },
    downloadUpdate() {
      this.closeDialog()
      window.open('https://your-download-url.com', '_blank')
    }
  }
}
</script>

<style scoped>
.update-dialog {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(0,0,0,0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 9999;
  backdrop-filter: blur(3px);
}

.dialog-content {
  width: 460px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.15);
  overflow: hidden;
  transform: translateY(-20px);
  opacity: 0;
  animation: fadeIn 0.3s ease-out forwards;
}

@keyframes fadeIn {
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.dialog-header {
  display: flex;
  align-items: center;
  padding: 20px 24px;
  background: linear-gradient(135deg, #4acf8b, #359e6b);
  color: white;
  position: relative;
}

.dialog-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  flex: 1;
}

.dialog-body {
  padding: 24px;
  color: #333;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid #f0f0f0;
}

.dialog-footer button {
  padding: 10px 20px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 14px;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 6px;
}

.update-icon {
  margin-right: 12px;
  display: flex;
  align-items: center;
}

.update-icon svg {
  width: 24px;
  height: 24px;
}

.close-btn {
  background: rgba(255,255,255,0.2);
  border: none;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
  color: white;
}

.close-btn:hover {
  background: rgba(255,255,255,0.3);
}

.version-info {
  background: #f8f9ff;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 20px;
  position: relative;
  overflow: hidden;
}

.version-badge {
  position: absolute;
  top: 0;
  right: 0;
  background: #ff4757;
  color: white;
  font-size: 12px;
  font-weight: bold;
  padding: 4px 10px;
  border-radius: 0 8px 0 8px;
}

.version-arrow {
  margin: 0 15px 0 15px;
  //color: #7ab898;
  font-weight: bold;
}

.version-info p {
  margin: 0;
  font-size: 15px;
  display: flex;
  align-items: center;
}

.version-text {
  font-weight: 600;
  color: #359e6b;
}

.version-text.current {
  color: #666;
}

.compact-list {
  margin: 0;
  padding: 0;
  list-style-type: none;
}

.compact-list li {
  margin-bottom: 6px;
  padding-left: 0;
  display: flex;
  align-items: flex-start;
}

.changelog {
  margin-top: 12px;
}

.changelog h4 {
  margin: 0 0 8px 10px;
  font-size: 16px;
  display: flex;
  align-items: center;
}

.changelog-badge {
  background: #e8f5ee;
  color: #359e6b;
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 10px;
  margin-left: 8px;
}

.changelog ul {
  margin: 0;
  //padding: 0 20px 0 0;
}

.changelog li {
  margin-top: 10px;
  margin-bottom: 10px;
  font-size: 14px;
  line-height: 1.5;
  list-style-type: none;
  position: relative;
  padding-left: 14px;
}

.changelog li:last-child {
  margin-bottom: 0;
}

.bullet {
  margin-right: 6px;
  color: #359e6b;
  font-weight: bold;
}

.later-btn {
  background: #f5f5f5;
  color: #666;
}

.later-btn:hover {
  background: #eaeaea;
}

.download-btn {
  background: linear-gradient(135deg, #4acf8b, #359e6b);
  box-shadow: 0 2px 8px rgba(53, 158, 107, 0.3);
  color: white;
}

.download-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(53, 158, 107, 0.4);
}

.download-btn:active {
  transform: translateY(0);
}
</style>