<template>
  <svg 
    :width="size" 
    :height="size" 
    viewBox="0 0 200 200" 
    xmlns="http://www.w3.org/2000/svg"
    class="gotux-logo"
  >
    <!-- 背景渐变 -->
    <defs>
      <linearGradient id="logoGradient" x1="0%" y1="0%" x2="100%" y2="100%">
        <stop offset="0%" style="stop-color:#667eea;stop-opacity:1" />
        <stop offset="100%" style="stop-color:#764ba2;stop-opacity:1" />
      </linearGradient>
      <linearGradient id="accentGradient" x1="0%" y1="0%" x2="100%" y2="100%">
        <stop offset="0%" style="stop-color:#f093fb;stop-opacity:1" />
        <stop offset="100%" style="stop-color:#f5576c;stop-opacity:1" />
      </linearGradient>
      
      <!-- 阴影 -->
      <filter id="shadow">
        <feGaussianBlur in="SourceAlpha" stdDeviation="3"/>
        <feOffset dx="0" dy="2" result="offsetblur"/>
        <feComponentTransfer>
          <feFuncA type="linear" slope="0.3"/>
        </feComponentTransfer>
        <feMerge>
          <feMergeNode/>
          <feMergeNode in="SourceGraphic"/>
        </feMerge>
      </filter>
    </defs>
    
    <!-- 圆形背景 -->
    <circle 
      cx="100" 
      cy="100" 
      r="90" 
      fill="url(#logoGradient)"
      filter="url(#shadow)"
    />
    
    <!-- 图片框架 - 外层 -->
    <rect 
      x="50" 
      y="60" 
      width="100" 
      height="80" 
      rx="8" 
      fill="none" 
      stroke="#ffffff" 
      stroke-width="6"
      opacity="0.9"
    />
    
    <!-- 山峰图标 (代表图片) -->
    <path 
      d="M 70 110 L 90 85 L 110 100 L 130 75 L 130 120 L 70 120 Z" 
      fill="#ffffff"
      opacity="0.85"
    />
    
    <!-- 太阳/光点 -->
    <circle 
      cx="125" 
      cy="80" 
      r="8" 
      fill="url(#accentGradient)"
    />
    
    <!-- Go 字母装饰 -->
    <text 
      x="100" 
      y="165" 
      font-family="'Segoe UI', sans-serif" 
      font-size="36" 
      font-weight="bold" 
      fill="#ffffff" 
      text-anchor="middle"
    >
      G
    </text>
  </svg>
</template>

<script setup>
import { defineProps } from 'vue'

defineProps({
  size: {
    type: [String, Number],
    default: 200
  }
})
</script>

<style scoped>
.gotux-logo {
  display: inline-block;
  vertical-align: middle;
}
</style>
