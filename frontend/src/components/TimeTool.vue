<template>
  <div class="time-converter-container">
    <!-- 头部选择时区区域 -->
    <div class="timezone-bar">
      <div class="timezone-selector">
        <span class="selector-label">
          <svg viewBox="0 0 24 24" width="16" height="16">
            <path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
          </svg>
          选择时区
        </span>
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
      <div class="timezone-code">
        <span class="code-label">当前时区</span>
        <span class="code-value">{{ timezone }}</span>
      </div>
    </div>

    <!-- 当前时间显示卡片 -->
    <div class="current-time-card">
      <div class="time-display-section">
        <div class="time-item">
          <div class="time-icon">
            <svg viewBox="0 0 24 24" width="20" height="20">
              <path fill="currentColor" d="M11.99 2C6.47 2 2 6.48 2 12s4.47 10 9.99 10C17.52 22 22 17.52 22 12S17.52 2 11.99 2zM12 20c-4.42 0-8-3.58-8-8s3.58-8 8-8 8 3.58 8 8-3.58 8-8 8zm.5-13H11v6l5.25 3.15.75-1.23-4.5-2.67z"/>
            </svg>
          </div>
          <div class="time-content">
            <span class="time-label">当前时间</span>
            <span class="time-value">{{ currentTime }}</span>
          </div>
        </div>
        <div class="time-divider"></div>
        <div class="time-item">
          <div class="time-icon timestamp-icon">
            <svg viewBox="0 0 24 24" width="20" height="20">
              <path fill="currentColor" d="M9 11H7v2h2v-2zm4 0h-2v2h2v-2zm4 0h-2v2h2v-2zm2-7h-1V2h-2v2H8V2H6v2H5c-1.11 0-1.99.9-1.99 2L3 20c0 1.1.89 2 2 2h14c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 16H5V9h14v11z"/>
            </svg>
          </div>
          <div class="time-content">
            <span class="time-label">当前时间戳</span>
            <span class="time-value mono">{{ currentTimestamp }}</span>
          </div>
        </div>
      </div>
      <button class="pause-btn" :class="{ 'paused': isPaused }" @click="toggleTimeUpdate">
        <svg v-if="!isPaused" viewBox="0 0 24 24" width="16" height="16">
          <path fill="currentColor" d="M6 19h4V5H6v14zm8-14v14h4V5h-4z"/>
        </svg>
        <svg v-else viewBox="0 0 24 24" width="16" height="16">
          <path fill="currentColor" d="M8 5v14l11-7z"/>
        </svg>
        {{ isPaused ? '继续' : '暂停' }}
      </button>
    </div>

    <!-- 转换区域 -->
    <div class="conversion-section">
      <!-- 时间戳转日期 -->
      <div class="conversion-card">
        <div class="card-header">
          <div class="card-icon">
            <svg viewBox="0 0 24 24" width="24" height="24">
              <path fill="currentColor" d="M11.99 2C6.47 2 2 6.48 2 12s4.47 10 9.99 10C17.52 22 22 17.52 22 12S17.52 2 11.99 2zM12 20c-4.42 0-8-3.58-8-8s3.58-8 8-8 8 3.58 8 8-3.58 8-8 8zm.5-13H11v6l5.25 3.15.75-1.23-4.5-2.67z"/>
            </svg>
          </div>
          <h3>时间戳 → 日期</h3>
        </div>
        <div class="card-body">
          <div class="input-wrapper">
            <input 
              type="text" 
              v-model="timestampInput" 
              class="modern-input" 
              placeholder="输入时间戳 (秒)" 
              @input="clearResultIfEmpty('timestamp')"
              @keyup.enter="convertToDate"
            >
            <button class="action-btn primary" @click="convertToDate">
              <svg viewBox="0 0 24 24" width="16" height="16">
                <path fill="currentColor" d="M12 4V1L8 5l4 4V6c3.31 0 6 2.69 6 6 0 1.01-.25 1.97-.7 2.8l1.46 1.46C19.54 15.03 20 13.57 20 12c0-4.42-3.58-8-8-8zm0 14c-3.31 0-6-2.69-6-6 0-1.01.25-1.97.7-2.8L5.24 7.74C4.46 8.97 4 10.43 4 12c0 4.42 3.58 8 8 8v3l4-4-4-4v3z"/>
              </svg>
              转换
            </button>
          </div>
          <div class="result-area" v-if="dateResult">
            <div class="result-value">{{ dateResult }}</div>
          </div>
          <div class="result-placeholder" v-else>
            <span>输入时间戳后点击转换</span>
          </div>
        </div>
      </div>

      <!-- 日期转时间戳 -->
      <div class="conversion-card">
        <div class="card-header">
          <div class="card-icon orange">
            <svg viewBox="0 0 24 24" width="24" height="24">
              <path fill="currentColor" d="M19 3h-1V1h-2v2H8V1H6v2H5c-1.11 0-1.99.9-1.99 2L3 19c0 1.1.89 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm0 16H5V8h14v11zM9 10H7v2h2v-2zm4 0h-2v2h2v-2zm4 0h-2v2h2v-2zm-8 4H7v2h2v-2zm4 0h-2v2h2v-2zm4 0h-2v2h2v-2z"/>
            </svg>
          </div>
          <h3>日期 → 时间戳</h3>
        </div>
        <div class="card-body">
          <div class="input-wrapper">
            <input 
              type="text" 
              v-model="dateInput" 
              class="modern-input" 
              placeholder="YYYY-MM-DD HH:mm:ss" 
              @input="clearResultIfEmpty('date')"
              @keyup.enter="convertToTimestamp"
            >
            <button class="action-btn secondary" @click="convertToTimestamp">
              <svg viewBox="0 0 24 24" width="16" height="16">
                <path fill="currentColor" d="M12 4V1L8 5l4 4V6c3.31 0 6 2.69 6 6 0 1.01-.25 1.97-.7 2.8l1.46 1.46C19.54 15.03 20 13.57 20 12c0-4.42-3.58-8-8-8zm0 14c-3.31 0-6-2.69-6-6 0-1.01.25-1.97.7-2.8L5.24 7.74C4.46 8.97 4 10.43 4 12c0 4.42 3.58 8 8 8v3l4-4-4-4v3z"/>
              </svg>
              转换
            </button>
          </div>
          <div class="result-area" v-if="timestampResult">
            <div class="result-value mono">{{ timestampResult }}</div>
          </div>
          <div class="result-placeholder" v-else>
            <span>输入日期后点击转换</span>
          </div>
        </div>
      </div>

      <!-- 日期计算 -->
      <div class="conversion-card">
        <div class="card-header">
          <div class="card-icon purple">
            <svg viewBox="0 0 24 24" width="24" height="24">
              <path fill="currentColor" d="M19 3h-1V1h-2v2H8V1H6v2H5c-1.11 0-1.99.9-1.99 2L3 19c0 1.1.89 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm0 16H5V8h14v11zM9 10H7v2h2v-2zm4 0h-2v2h2v-2zm4 0h-2v2h2v-2z"/>
            </svg>
          </div>
          <h3>日期计算</h3>
        </div>
        <div class="card-body compact">
          <div class="calc-row">
            <span class="calc-label">基准日期</span>
            <input 
              type="text" 
              v-model="baseDateInput" 
              class="modern-input small" 
              placeholder="YYYY-MM-DD" 
              @input="clearResultIfEmpty('dateCalc')"
            >
          </div>
          <div class="calc-row">
            <span class="calc-label">相差天数</span>
            <input 
              type="number" 
              v-model="daysOffset" 
              class="modern-input small days-input" 
              placeholder="0"
              @keyup.enter="calculateDate"
            >
            <span class="calc-hint">(负数向前，正数向后)</span>
          </div>
          <button class="action-btn purple full-width" @click="calculateDate">
            <svg viewBox="0 0 24 24" width="16" height="16">
              <path fill="currentColor" d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-5 14H7v-2h7v2zm3-4H7v-2h10v2zm0-4H7V7h10v2z"/>
            </svg>
            计算日期
          </button>
          <div class="result-area compact" v-if="dateCalcResult">
            <div class="result-value">{{ dateCalcResult }}</div>
          </div>
          <div class="result-placeholder compact" v-else>
            <span>输入日期和天数后计算</span>
          </div>
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
      baseDateInput: "",
      daysOffset: "",
      dateCalcResult: "",
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
    toggleTimeUpdate() {
      this.isPaused = !this.isPaused;
      if (this.isPaused) {
        this.clearTimer();
      } else {
        this.timer = setInterval(this.updateCurrentTime, 1000);
      }
    },
    updateCurrentTime() {
      const now = new Date();
      this.currentTimestamp = Math.floor(now.getTime() / 1000);

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
    clearResultIfEmpty(type) {
      if (type === 'timestamp' && !this.timestampInput.trim()) {
        this.dateResult = '';
      } else if (type === 'date' && !this.dateInput.trim()) {
        this.timestampResult = '';
      } else if (type === 'dateCalc' && !this.baseDateInput.trim() && !this.daysOffset.trim()) {
        this.dateCalcResult = '';
      }
    },
    calculateDate() {
      if (!this.baseDateInput) {
        this.dateCalcResult = '';
        return;
      }
      const days = parseInt(this.daysOffset || 0);
      const baseDate = new Date(this.baseDateInput.replace(/-/g, '/'));
      if (isNaN(baseDate.getTime())) {
        this.dateCalcResult = '日期格式错误';
        return;
      }
      const resultDate = new Date(baseDate);
      resultDate.setDate(resultDate.getDate() + days);
      
      const year = resultDate.getFullYear();
      const month = String(resultDate.getMonth() + 1).padStart(2, '0');
      const day = String(resultDate.getDate()).padStart(2, '0');
      const hours = String(resultDate.getHours()).padStart(2, '0');
      const minutes = String(resultDate.getMinutes()).padStart(2, '0');
      const seconds = String(resultDate.getSeconds()).padStart(2, '0');
      
      this.dateCalcResult = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
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
.time-converter-container {
  max-width: 900px;
  margin: 0 auto;
  padding: 24px;
  height: 100%;
  overflow-y: auto;
}

/* 头部时区选择栏 */
.timezone-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #f8f9fa;
  padding: 12px 20px;
  border-radius: 10px;
  margin-bottom: 20px;
}

.timezone-selector {
  display: flex;
  align-items: center;
  gap: 10px;
}

.selector-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: #666;
  font-weight: 500;
}

.timezone-dropdown {
  padding: 8px 12px;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  background-color: white;
  font-size: 14px;
  color: #333;
  cursor: pointer;
  min-width: 220px;
  outline: none;
  transition: border-color 0.2s;
}

.timezone-dropdown:hover,
.timezone-dropdown:focus {
  border-color: #3498db;
}

.timezone-code {
  display: flex;
  align-items: center;
  gap: 10px;
}

.code-label {
  font-size: 13px;
  color: #888;
}

.code-value {
  font-size: 14px;
  font-weight: 600;
  color: #3498db;
  font-family: 'SF Mono', Monaco, monospace;
  background: white;
  padding: 6px 12px;
  border-radius: 6px;
  border: 1px solid #e0e0e0;
}

/* 当前时间卡片 - 浅色样式 */
.current-time-card {
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  padding: 14px 24px;
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
}

.time-display-section {
  display: flex;
  align-items: center;
  gap: 32px;
  flex: 1;
}

.time-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.time-icon {
  width: 44px;
  height: 44px;
  background: linear-gradient(135deg, #3498db 0%, #2980b9 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.time-icon.timestamp-icon {
  background: linear-gradient(135deg, #42b983 0%, #27ae60 100%);
}

.time-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.time-label {
  font-size: 12px;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.time-value {
  font-size: 20px;
  font-weight: 600;
  color: #1e293b;
  font-family: 'SF Mono', Monaco, monospace;
}

.time-divider {
  width: 1px;
  height: 48px;
  background: #cbd5e1;
}

.pause-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 18px;
  background: #64748b;
  border: none;
  border-radius: 8px;
  color: white;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.pause-btn:hover {
  background: #475569;
}

.pause-btn.paused {
  background: #42b983;
}

/* 转换区域 */
.conversion-section {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

/* 日期计算卡片独占一行 */
.conversion-card:last-child {
  grid-column: 1 / -1;
  max-width: 50%;
}

.conversion-card {
  background: white;
  border-radius: 10px;
  border: 1px solid #f0f0f0;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
  overflow: hidden;
  transition: transform 0.2s, box-shadow 0.2s;
}

.conversion-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 20px;
  background: #f8f9fa;
  border-bottom: 1px solid #f0f0f0;
}

.card-icon {
  width: 36px;
  height: 36px;
  background: #e8f4f8;
  border: 1px solid #3498db;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #3498db;
}

.card-icon.orange {
  background: #fef3e8;
  border-color: #e67e22;
  color: #e67e22;
}

.card-icon.purple {
  background: #f3e8fd;
  border-color: #9b59b6;
  color: #9b59b6;
}

.card-header h3 {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: #333;
}

.card-body {
  padding: 16px;
}

.card-body.compact {
  padding: 10px;
}

.calc-row {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 8px;
}

.calc-label {
  font-size: 11px;
  color: #666;
  min-width: 48px;
  white-space: nowrap;
}

.calc-hint {
  font-size: 10px;
  color: #999;
  white-space: nowrap;
}

.modern-input.small {
  padding: 10px 8px;
  font-size: 12px;
  flex: 1;
}

.modern-input.days-input {
  width: 50px;
  flex: none;
  text-align: center;
}

.action-btn.purple {
  border-color: #9b59b6;
  color: #9b59b6;
  width: 100%;
  justify-content: center;
  margin-top: 4px;
  padding: 8px;
  font-size: 12px;
}

.action-btn.purple:hover {
  background: #9b59b6;
  color: white;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(155, 89, 182, 0.3);
}

.action-btn.full-width {
  width: 100%;
  justify-content: center;
}

.result-area.compact {
  padding: 8px 10px;
  margin-top: 8px;
  border-left-color: #9b59b6;
}

.result-placeholder.compact {
  padding: 12px;
  margin-top: 8px;
}

.input-wrapper {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
}

.modern-input {
  flex: 1;
  padding: 10px 12px;
  border: 2px solid #e8e8e8;
  border-radius: 8px;
  font-size: 13px;
  outline: none;
  transition: all 0.2s;
  background: #fafafa;
}

.modern-input:focus {
  border-color: #3498db;
  background: white;
}

.modern-input::placeholder {
  color: #aaa;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 10px 16px;
  border: 1px solid;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
  background: white;
  align-self: center;
  height: fit-content;
}

.action-btn.primary {
  border-color: #3498db;
  color: #3498db;
}

.action-btn.primary:hover {
  background: #3498db;
  color: white;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(52, 152, 219, 0.3);
}

.action-btn.secondary {
  border-color: #e67e22;
  color: #e67e22;
}

.action-btn.secondary:hover {
  background: #e67e22;
  color: white;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(230, 126, 34, 0.3);
}

.result-area {
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8ec 100%);
  border-radius: 8px;
  padding: 12px;
  border-left: 4px solid #3498db;
}

.result-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #666;
  margin-bottom: 8px;
  font-weight: 500;
}

.result-value {
  font-size: 14px;
  color: #333;
  font-weight: 600;
  word-break: break-all;
}

.result-value.mono {
  font-family: 'SF Mono', Monaco, monospace;
}

.result-placeholder {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 12px;
  text-align: center;
  color: #999;
  font-size: 13px;
  border: 2px dashed #e0e0e0;
}

/* 响应式 */
@media (max-width: 1024px) {
  .conversion-section {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .conversion-section {
    grid-template-columns: 1fr;
  }
  
  .current-time-card {
    flex-direction: column;
    gap: 16px;
  }
  
  .time-display-section {
    flex-direction: column;
    gap: 16px;
    width: 100%;
  }
  
  .time-divider {
    width: 100%;
    height: 1px;
  }
}
</style>
