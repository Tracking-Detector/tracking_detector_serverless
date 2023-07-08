export default defineNuxtConfig({
  css: ['vuetify/lib/styles/main.sass'],
  build: {
    transpile: ['vuetify'],
  },
  modules: [
    'nuxt-highcharts',
  ],
  components: true,
  vite: {
    define: {
      'process.env.DEBUG': false,
    },
  },
  
})
