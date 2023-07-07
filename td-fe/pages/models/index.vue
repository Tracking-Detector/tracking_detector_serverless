<template>
    <div>
        <h3 class="mb-4">Welcome on Models</h3>
        <p style="width: 800px;" class="mb-4">
            Here you can explore the different models and trigger the training process for different datasets. Additionally you can see great insights
            on the models performance and other metrics of each training run for each dataset.
        </p>
        <v-card>
            <v-progress-linear v-if="isLoading" indeterminate></v-progress-linear>
            <div class="pa-4">
                <v-select density="compact" style="width: 300px;" v-model="selectedModel" label="Select Model"
                    :items="models.map(x => x.name)"></v-select>
                   
                <div v-if="selectedModel != 'No model selected'">
                    <h4 class="mb-2">Run Training on DataSet</h4>
                    <div class="d-flex">
                        <v-select density="compact" style="max-width: 300px;" v-model="selectedDataSet" label="Dataset to train the model on"
                    :items="dataSets"></v-select>
             
                    <v-btn icon="mdi-play" variant="text" @click="triggerTraining()"></v-btn>
                    </div>
                    
                    <training-data-view  v-for="dataset in getUniqueDataNames()" :data="runs" :dataSet="dataset"></training-data-view>
                </div>
               
            </div>

        </v-card>

    </div>
</template>
<script setup>


const models = ref([])
const dataSets = ref([])
const runs = ref([])
const isLoading = ref(false)
const selectedModel = ref("No model selected")
const selectedDataSet = ref("")

const triggerTraining = () => {
    fetch(`api/train/models/${selectedModel.value}/run/${selectedDataSet.value}`, {
        method: "POST"
    })
}

const getUniqueDataNames = () => {
    return new Set(runs.value.map(x=>x.dataSet))
}

const loadModelRuns = () => {
    isLoading.value = true
    fetch("/api/training-runs/" + selectedModel.value).then(response => {
        return response.json()
    }).then(body => {
        runs.value = body.data
        isLoading.value = false
    })
}
const loadAvailableDataSets = () => {
    isLoading.value = true
    fetch("/api/export").then(response => {
        return response.json()
    }).then(body => {
        dataSets.value = body.data.map(x => x.name)
        isLoading.value = false
    })
}

const loadAvailableModels = () => {
    isLoading.value = true
    fetch("/api/train/models").then(response => {
        return response.json()
    }).then(body => {
        models.value = body.data
        isLoading.value = false
    })
}
watch(selectedModel, () => {
    loadModelRuns()
})
onMounted(() => {
    loadAvailableModels()
    loadAvailableDataSets()
})
</script>