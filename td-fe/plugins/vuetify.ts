import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

// make sure to also import the coresponding css
import '@mdi/font/css/materialdesignicons.css' // Ensure you are using css-loader
import '@fortawesome/fontawesome-free/css/all.css'

export default defineNuxtPlugin(nuxtApp => {
  const vuetify = createVuetify({
    components,
    icons: {
      defaultSet: 'mdi'
    },
    directives,
  })

  nuxtApp.vueApp.use(vuetify)
})