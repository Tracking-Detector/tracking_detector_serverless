<template>
    <div class="pa-2">
        <v-divider></v-divider>
        <h4 class="mb-3 mt-3">Trained on {{ dataSet }}</h4>
        <v-row>
            <v-col cols="9">
                <v-row>
                    <v-col cols="6">
                        <highchart :options="createAccLossConfig('accuracy', 'Accuracy against Train Set for: ')">
                        </highchart>
                    </v-col>
                    <v-col cols="6">
                        <highchart :options="createAccLossConfig('loss', 'Loss against Train Set for: ')"></highchart>
                    </v-col>
                    <v-col cols="6">
                        <highchart :options="createAccLossConfig('val_accuracy', 'Accuracy against Test Set for: ')">
                        </highchart>
                    </v-col>
                    <v-col cols="6">
                        <highchart :options="createAccLossConfig('val_loss', 'Loss against Test Set for: ')"></highchart>
                    </v-col>
                </v-row>
            </v-col>
            <v-col cols="3">
                <v-list>
                    <v-list-item>
                        Trigger Training for {{ dataSet }}
                        <template v-slot:append>
                            <v-btn icon="mdi-play" variant="text" @click="triggerTraining(dataSet[0].name, dataSet[0].dataSet)"></v-btn>
                        </template>
                    </v-list-item>
                    <v-list-subheader>Run Metrics</v-list-subheader>
                    <v-list-item v-for="item in genRunMetrics()">
                        <v-list-item-title>Run on: {{ item.time }}</v-list-item-title>
                        <v-list-item-subtitle>
                            - F1Train: {{ item.f1Train }}
                        </v-list-item-subtitle>
                        <v-list-item-subtitle>
                            - F1Test: {{ item.f1Test }}
                        </v-list-item-subtitle>
                        <v-list-item-subtitle>
                            - BatchSize: {{ item.batchSize }}
                        </v-list-item-subtitle>
                        <v-list-item-subtitle>
                            - Epochs: {{ item.epochs }}
                        </v-list-item-subtitle>

                    </v-list-item>
                </v-list>
            </v-col>
        </v-row>


    </div>
</template>
<script setup>
const props = defineProps({
    data: {
        required: true
    },
    dataSet: {
        required: true,
        type: String
    }
})
const { dataSet, data } = toRefs(props)

const genRunMetrics = () => {
    return data.value.filter(x => x.dataSet == dataSet.value)
}

const triggerTraining = (modelName, dataSet) => {
    fetch(`/api/train/models/${modelName}/run/${dataSet}`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({})
    })
}

const createAccLossConfig = (type, title) => {
    const series = data.value.filter(x => x.dataSet == dataSet.value).map(x => {
        return {
            name: `${type} from ${x.time}`,
            data: x.trainingHistory[type]
        }
    })
    return {
        xAxis: {
            text: "Epoch"
        },
        yAxis: {
            text: "Loss / Accuracy"
        },
        title: {
            text: `${title} ${dataSet.value}`
        },
        series: series
    }
}
</script>