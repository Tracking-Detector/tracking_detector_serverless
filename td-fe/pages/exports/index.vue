<template>
    <div class="pa-4">
        <v-alert v-model="alert.isShowing" color="success" icon="$success" :title="alert.title" :text="alert.message"
            closable></v-alert>
        <h3 class="mb-4">Welcome on Exports</h3>
        <p style="width: 800px;" class="mb-4">Here you can see all the available exports for the data. Each export exports
            into a .csv.gz file and stores
            it into the MinIO export bucket. You can download each export on the download page. This page can be useful in
            case you want to manually start a export
            because you added new data into the database. Normally all exports are triggered once every two weeks
            automatically.
        </p>
        <v-card class="pa-4">
            <v-list>
                <v-list-subheader>Available Exports</v-list-subheader>
                <v-progress-linear v-if="isLoading" indeterminate></v-progress-linear>
                <div v-for="avex in exports">
                    <v-list-item @click="startExport(avex.name)">
                        <template v-slot:prepend>
                            <v-icon>mdi-export</v-icon>
                        </template>
                        <v-list-item-title>{{ avex.name }}</v-list-item-title>
                        <v-list-item-subtitle>{{ avex.description }}</v-list-item-subtitle>
                        <template v-slot:append>
                            <v-btn icon="mdi-play" variant="text" @click=""></v-btn>
                        </template>
                    </v-list-item>
                    <v-divider></v-divider>
                </div>
            </v-list>
        </v-card>
    </div>
</template>
<script setup>
const exports = ref([])
const isLoading = ref(true)
const alert = ref({
    isShowing: false,
    title: "Export triggered",
    message: ""
})

const loadAvailableExports = () => {
    isLoading.value = true
    fetch("http://localhost/export").then(response => {
        return response.json()
    }).then(body => {
        exports.value = body.data
        isLoading.value = false
    })
}

const startExport = (name) => {
    isLoading.value = true
    fetch("http://localhost/export/" + name + "/run", {
        method: "POST"
    }).then(response => {
        return response.json()
    }).then(body => {
        alert.value.message = body.message
        isLoading.value = false
        alert.value.isShowing = true
    })
}

onMounted(() => {

    loadAvailableExports()
})

</script>