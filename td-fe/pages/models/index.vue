<template>
  <div>
    <h3 class="mb-4 text-h5">Welcome on Models</h3>
    <p style="width: 600px" class="mb-4 text-body-1">
      Here you can explore the different models and trigger the training process
      for different datasets. Additionally you can see great insights on the
      models performance and other metrics of each training run for each
      dataset.
    </p>
    <v-card>
      <v-progress-linear v-if="isLoading" indeterminate></v-progress-linear>

      <div v-if="models.length > 0" class="pa-4">
        <v-select
          density="compact"
          style="width: 300px"
          v-model="selectedModel"
          label="Select Model"
          :items="models.map((x) => x.name)"
        ></v-select>
        <div v-if="selectedModel != 'No model selected'">
          <h4 class="mb-2 text-h6">Run Training on DataSet</h4>
          <div class="d-flex">
            <v-select
              density="compact"
              style="max-width: 300px"
              v-model="selectedDataSet"
              label="Dataset to train the model on"
              :items="dataSets"
            ></v-select>

            <v-btn
              icon="mdi-play"
              variant="text"
              @click="triggerTraining()"
            ></v-btn>
          </div>
          <div v-if="runs != undefined && runs.length > 0">
            <training-data-view
              v-for="dataset in getUniqueDataNames()"
              :key="dataset"
              :data="runs"
              :dataSet="dataset"
            ></training-data-view>
          </div>
        </div>
      </div>
    </v-card>
  </div>
</template>
<script setup>
useMeta({
  title: "Tracking Detector - Models",
});

const models = ref([]);
const dataSets = ref([]);
const runs = ref([]);
const isLoading = ref(false);
const config = useRuntimeConfig();
const selectedModel = ref("No model selected");
const selectedDataSet = ref("");

const triggerTraining = () => {
  fetch(
    `/api/dispatch/train/${selectedModel.value}/run/${selectedDataSet.value}`,
    {
      headers: {
        "Content-Type": "application/json",
        "X-API-Key": "Bearer " + config.public.apiBase,
      },
      method: "POST",
      body: JSON.stringify({}),
    }
  );
};

const getUniqueDataNames = () => {
  return new Set(runs.value.map((x) => x.dataSet));
};

const loadModelRuns = () => {
  isLoading.value = true;
  fetch("/api/training-runs/" + selectedModel.value, {
    headers: {
      "X-API-Key": "Bearer " + config.public.apiBase,
    },
  })
    .then((response) => {
      if (response.status != 200) {
        return undefined;
      }
      return response.json();
    })
    .then((body) => {
      if (body != undefined) {
        runs.value = body.data;
      }
      isLoading.value = false;
    });
};
const loadAvailableDataSets = () => {
  isLoading.value = true;
  fetch("/api/dispatch/export", {
    headers: {
      "X-API-Key": "Bearer " + config.public.apiBase,
    },
  })
    .then((response) => {
      return response.json();
    })
    .then((body) => {
      dataSets.value = body.data.map((x) => x.name);
      isLoading.value = false;
    });
};

const loadAvailableModels = () => {
  isLoading.value = true;
  fetch("/api/dispatch/model", {
    headers: {
      "X-API-Key": "Bearer " + config.public.apiBase,
    },
  })
    .then((response) => {
      return response.json();
    })
    .then((body) => {
      models.value = body.data;
      isLoading.value = false;
    });
};
watch(selectedModel, () => {
  loadModelRuns();
});
onMounted(() => {
  loadAvailableModels();
  loadAvailableDataSets();
});
</script>
