<template>
  <div class="time-converter-container">
    <div class="time-converter-header">
      <h2>时区转换工具</h2>
      <div class="timezone-selector">
        <label for="timezone">选择时区：</label>
        <select id="timezone" class="timezone-dropdown" v-model="timezone">
          <option
              v-for="tz in timezoneOptions"
              :key="tz.value"
              :value="tz.value"
          >
            {{ tz.text }}
          </option>
        </select>
      </div>
    </div>

    <div class="current-time-bar">
      <div class="current-time-item">
        <span class="time-label">当前时间：</span>
        <span class="time-value">{{ currentTime }}</span>
      </div>
      <div class="current-time-item">
        <span class="time-label">当前时间戳：</span>
        <span class="time-value">{{ currentTimestamp }}</span>
      </div>
      <!-- 添加暂停/开始按钮 -->
      <div class="time-control">
        <button class="time-btn" @click="toggleTimeUpdate">
          {{ isPaused ? '开始' : '暂停' }}
        </button>
      </div>
    </div>

    <div class="conversion-section">
      <div class="conversion-card">
        <h2>时间戳 → 日期</h2>
        <div class="input-group">
          <input type="text" v-model="timestampInput" class="timestamp-input" placeholder="输入时间戳" @input="clearResultIfEmpty('timestamp')" >
          <button class="convert-btn" @click="convertToDate">转换</button>
        </div>
        <div class="result-box">
          <span class="result-label">结果：</span>
          <span class="result-value">{{ dateResult }}</span>
        </div>
      </div>

      <div class="conversion-card">
        <h2>日期 → 时间戳</h2>
        <div class="input-group">
          <input type="text" v-model="dateInput" class="datetime-input" placeholder="YYYY-MM-DD HH:mm:ss" @input="clearResultIfEmpty('date')">
          <button class="convert-btn" @click="convertToTimestamp">转换</button>
        </div>
        <div class="result-box">
          <span class="result-label">结果：</span>
          <span class="result-value">{{ timestampResult }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      timezone: "Asia/Shanghai",
      timezoneOptions: [
        { text: "中国标准时间 (UTC+8)", value: "Asia/Shanghai" },
        { text: "菲律宾标准时间 (UTC+8)", value: "Asia/Manila" },
        { text: "印度标准时间 (UTC+5:30)", value: "Asia/Kolkata" },
        { text: "印尼西部时间 (UTC+7)", value: "Asia/Jakarta" },
        { text: "越南标准时间 (UTC+7)", value: "Asia/Ho_Chi_Minh" },
        { text: "墨西哥中部时间 (UTC-6)", value: "America/Mexico_City" },
        { text: "美国东部时间 (UTC-5)", value: "America/New_York" },
        { text: "美国太平洋时间 (UTC-8)", value: "America/Los_Angeles" },
      ],
      currentTime: "",
      currentTimestamp: "",
      timestampInput: "",
      dateInput: "",
      dateResult: "",
      timestampResult: "",
      timer: null,
      isPaused: false
    }
  },
  mounted() {
    this.updateCurrentTime();
    this.timer = setInterval(this.updateCurrentTime, 1000);
  },
  beforeDestroy() {
    if (this.timer) {
      clearInterval(this.timer);
    }
  },
  methods: {
    clearTimer() {
      if (this.timer) {
        clearInterval(this.timer);
        this.timer = null;
      }
    },
    // 切换时间更新状态
    toggleTimeUpdate() {
      this.isPaused = !this.isPaused;

      if (this.isPaused) {
        this.clearTimer(); // 暂停时清除定时器
      } else {
        // 恢复时重新启动定时器
        this.timer = setInterval(this.updateCurrentTime, 1000);
      }
    },
    updateCurrentTime() {
      const now = new Date();
      this.currentTimestamp = Math.floor(now.getTime() / 1000);

      // Format the date according to the selected timezone
      const options = {
        timeZone: this.timezone,
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false
      };
      const formatter = new Intl.DateTimeFormat('zh-CN', options);
      const parts = formatter.formatToParts(now);
      const year = parts.find(p => p.type === 'year').value;
      const month = parts.find(p => p.type === 'month').value;
      const day = parts.find(p => p.type === 'day').value;
      const hour = parts.find(p => p.type === 'hour').value.padStart(2, '0');
      const minute = parts.find(p => p.type === 'minute').value.padStart(2, '0');
      const second = parts.find(p => p.type === 'second').value.padStart(2, '0');
      this.currentTime = `${year}-${month}-${day} ${hour}:${minute}:${second}`;
    },
    async convertToDate() {
      if (!this.timestampInput) {
        this.dateResult = "";
        return;
      }
      this.dateResult = await window.go.handler.TimeHandler.TimestampToDate(
          parseInt(this.timestampInput),
          this.timezone
      );
    },
    async convertToTimestamp() {
      if (!this.dateInput) {
        this.timestampResult = "";
        return;
      }
      this.timestampResult = await window.go.handler.TimeHandler.DateToTimestamp(
          this.dateInput,
          this.timezone
      );
    },
    // 新增方法：检查输入是否为空并清空对应结果
    clearResultIfEmpty(type) {
      if (type === 'timestamp' && !this.timestampInput.trim()) {
        this.dateResult = '';
      } else if (type === 'date' && !this.dateInput.trim()) {
        this.timestampResult = '';
      }
    }
  },
  watch: {
    timezone() {
      this.updateCurrentTime();
    }
  }
}
</script>

<style scoped>
:global(html), :global(body) {
  overflow-y: hidden;
  height: 100%;
  margin: 0;
  padding: 0;
}
.time-converter-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  background-color: white;
  border-radius: 10px;
  font-family: "思源宋体", "Noto Serif SC", serif;
}

.time-converter-header {
  margin-bottom: 30px;
  text-align: center;
}

.time-converter-header h1 {
  color: #2c3e50;
  margin-bottom: 15px;
}

.timezone-selector {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

.timezone-dropdown {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 5px;
  background-color: white;
  font-size: 14px;
}

.conversion-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.conversion-card {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.conversion-card h2 {
  color: #3498db;
  font-size: 18px;
  margin-bottom: 15px;
}

.input-group {
  display: flex;
  gap: 10px;
  margin-bottom: 15px;
}

.input-group input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 5px;
  font-size: 14px;
}

.convert-btn {
  padding: 8px 16px;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.convert-btn:hover {
  background-color: #2980b9;
}

.result-box {
  padding: 12px;
  background-color: #f8f9fa;
  border-radius: 5px;
  border-left: 4px solid #3498db;
  text-align: left;
}

.result-label {
  font-weight: bold;
  color: #555;
  margin-right: 8px;
}

.result-value {
  color: #2c3e50;
  //margin-left: 5px;
}

.current-time-bar {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 30px;
  margin-bottom: 20px;
  padding: 12px;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.current-time-item {
  display: flex;
  align-items: center;
}

.time-label {
  font-weight: bold;
  color: #3498db;
  margin-right: 8px;
}

.time-value {
  font-family: 'Courier New', monospace;
  color: #2c3e50;
  font-size: 15px;
}
.time-control {
  display: flex;
  align-items: center;
}
.time-btn {
  padding: 6px 12px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.2s;
}
</style>